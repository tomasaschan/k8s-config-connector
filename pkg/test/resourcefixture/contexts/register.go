// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package contexts

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/core/v1alpha1"
	dclcontroller "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/dcl"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl"
	dclconversion "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/conversion"
	dclextension "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/extension"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/kcclite"
	dcllivestate "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/livestate"
	dclmetadata "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/metadata"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/dcl/schema/dclschemaloader"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/k8s"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/krmtotf"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/resourceoverrides"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/servicemapping/servicemappingloader"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/test/resourcefixture"
	"github.com/nasa9084/go-openapi"

	mmdcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	dclunstruct "github.com/GoogleCloudPlatform/declarative-resource-client-library/unstructured"
	tfschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourceContext struct {
	ResourceGVK        schema.GroupVersionKind
	ResourceKind       string
	SkipNoChange       bool
	SkipUpdate         bool
	SkipDriftDetection bool
	SkipDelete         bool

	// fields related to DCL-based resources
	DCLSchema *openapi.Schema
	DCLBased  bool
}

var (
	resourceContextMap = map[string]ResourceContext{}
	emptyGVK           = schema.GroupVersionKind{}
)

func GetResourceContext(fixture resourcefixture.ResourceFixture, serviceMetadataLoader dclmetadata.ServiceMetadataLoader, dclSchemaLoader dclschemaloader.DCLSchemaLoader) ResourceContext {
	rc, ok := resourceContextMap[fixture.Name]
	if !ok {
		rc = ResourceContext{
			ResourceGVK:  fixture.GVK,
			ResourceKind: fixture.GVK.Kind,
		}
	}
	if rc.ResourceGVK == emptyGVK {
		rc.ResourceGVK = fixture.GVK
	}
	if dclmetadata.IsDCLBasedResourceKind(rc.ResourceGVK, serviceMetadataLoader) {
		s, err := dclschemaloader.GetDCLSchemaForGVK(rc.ResourceGVK, serviceMetadataLoader, dclSchemaLoader)
		if err != nil {
			panic(fmt.Sprintf("error getting the DCL schema for GVK %v: %v", rc.ResourceGVK, err))
		}
		rc.DCLSchema = s
		rc.DCLBased = true
	}
	return rc
}

func (rc ResourceContext) SupportsLabels(smLoader *servicemappingloader.ServiceMappingLoader) bool {
	if rc.DCLBased {
		_, _, found, err := dclextension.GetLabelsFieldSchema(rc.DCLSchema)
		if err != nil {
			panic(fmt.Sprintf("error getting the DCL schema for labels field: %v", err))
		}
		return found
	}
	// For tf based resources, resolve the label info from ResourceConfig
	resourceConfig := rc.getResourceConfig(smLoader)
	return resourceConfig.MetadataMapping.Labels != ""
}

func (rc ResourceContext) getResourceConfig(smLoader *servicemappingloader.ServiceMappingLoader) v1alpha1.ResourceConfig {
	for _, sm := range smLoader.GetServiceMappings() {
		for _, resourceConfig := range sm.Spec.Resources {
			if resourceConfig.Kind == rc.ResourceKind {
				return resourceConfig
			}
		}
	}
	panic(fmt.Errorf("no resource config found for kind: %v", rc.ResourceKind))
}

func (rc ResourceContext) Create(t *testing.T, u *unstructured.Unstructured, provider *tfschema.Provider, c client.Client, smLoader *servicemappingloader.ServiceMappingLoader, config *mmdcl.Config, dclConverter *dclconversion.Converter) (*unstructured.Unstructured, error) {
	if rc.DCLBased {
		return dclCreate(u, config, c, dclConverter, smLoader)
	}
	return terraformCreate(u, provider, c, smLoader)
}

func (rc ResourceContext) Get(t *testing.T, u *unstructured.Unstructured, provider *tfschema.Provider, c client.Client, smLoader *servicemappingloader.ServiceMappingLoader, config *mmdcl.Config, dclConverter *dclconversion.Converter) (*unstructured.Unstructured, error) {
	if rc.DCLBased {
		return dclGet(u, config, c, dclConverter, smLoader)
	}
	return terraformGet(u, provider, c, smLoader)
}

func (rc ResourceContext) Delete(t *testing.T, u *unstructured.Unstructured, provider *tfschema.Provider, c client.Client, smLoader *servicemappingloader.ServiceMappingLoader, config *mmdcl.Config, dclConverter *dclconversion.Converter) error {
	if rc.DCLBased {
		return dclDelete(u, config, c, dclConverter, smLoader)
	}
	return terraformDelete(u, provider, c, smLoader)
}

func (rc ResourceContext) DoPreActuationTransformFor(u *unstructured.Unstructured, provider *tfschema.Provider, smLoader *servicemappingloader.ServiceMappingLoader, dclConverter *dclconversion.Converter) (*unstructured.Unstructured, error) {
	resource, err := unstructuredToKRMResource(rc.DCLBased, u, provider, smLoader, dclConverter)
	if err != nil {
		return nil, err
	}
	if err := resourceoverrides.Handler.PreActuationTransform(resource); err != nil {
		return nil, fmt.Errorf("could not run pre-acutuation transform on resource %s: %v", u.GetName(), err)
	}
	return resource.MarshalAsUnstructured()
}

func terraformDelete(u *unstructured.Unstructured, provider *tfschema.Provider, c client.Client, smLoader *servicemappingloader.ServiceMappingLoader) error {
	ctx := context.Background()
	resource, liveState, err := getTerraformResourceAndLiveState(ctx, u, provider, c, smLoader)
	if err != nil {
		return err
	}
	if liveState.Empty() {
		return fmt.Errorf("resource '%v' of type '%v' cannot be deleted as it does not exist", u.GetName(), u.GroupVersionKind())
	}
	_, diagnostics := resource.TFResource.Apply(ctx, liveState, &terraform.InstanceDiff{Destroy: true}, provider.Meta())
	if err := krmtotf.NewErrorFromDiagnostics(diagnostics); err != nil {
		return fmt.Errorf("error deleting resource: %v", err)
	}
	return err
}

func terraformCreate(u *unstructured.Unstructured, provider *tfschema.Provider, c client.Client, smLoader *servicemappingloader.ServiceMappingLoader) (*unstructured.Unstructured, error) {
	ctx := context.Background()
	resource, liveState, err := getTerraformResourceAndLiveState(ctx, u, provider, c, smLoader)
	if err != nil {
		return nil, err
	}
	if !liveState.Empty() {
		return nil, fmt.Errorf("resource '%v' of type '%v' cannot be created as it already exists", u.GetName(), u.GroupVersionKind())
	}
	config, _, err := krmtotf.KRMResourceToTFResourceConfig(resource, c, smLoader)
	if err != nil {
		return nil, fmt.Errorf("error expanding resource configuration: %v", err)
	}
	diff, err := resource.TFResource.Diff(ctx, liveState, config, provider.Meta())
	if err != nil {
		return nil, fmt.Errorf("error calculating diff: %v", err)
	}
	newState, diagnostics := resource.TFResource.Apply(ctx, liveState, diff, provider.Meta())
	if err := krmtotf.NewErrorFromDiagnostics(diagnostics); err != nil {
		return nil, fmt.Errorf("error applying resource change: %v", err)
	}
	return resourceToKRM(resource, newState)
}

func terraformGet(u *unstructured.Unstructured, provider *tfschema.Provider, c client.Client, smLoader *servicemappingloader.ServiceMappingLoader) (*unstructured.Unstructured, error) {
	resource, liveState, err := getTerraformResourceAndLiveState(context.Background(), u, provider, c, smLoader)
	if err != nil {
		return nil, err
	}
	if liveState.Empty() {
		return nil, fmt.Errorf("resource '%v' of type '%v' is not found", u.GetName(), u.GroupVersionKind())
	}
	return resourceToKRM(resource, liveState)
}

func dclCreate(u *unstructured.Unstructured, config *mmdcl.Config, kubeClient client.Client, converter *dclconversion.Converter, serviceMappingLoader *servicemappingloader.ServiceMappingLoader) (*unstructured.Unstructured, error) {
	ctx := context.Background()
	resource, err := newDCLResource(u, converter)
	if err != nil {
		return nil, err
	}
	liveLite, err := dcllivestate.FetchLiveState(ctx, resource, config, converter, serviceMappingLoader, kubeClient)
	if err != nil {
		return nil, err
	}
	if liveLite != nil {
		return nil, fmt.Errorf("resource '%v' of type '%v' cannot be created as it already exists", u.GetName(), u.GroupVersionKind())
	}
	lite, err := kcclite.ToKCCLite(resource, converter.MetadataLoader, converter.SchemaLoader, serviceMappingLoader, kubeClient)
	if err != nil {
		return nil, fmt.Errorf("error converting KCC full to KCC lite: %w", err)
	}
	dclObj, err := converter.KRMObjectToDCLObject(lite)
	if err != nil {
		return nil, fmt.Errorf("error converting KCC lite to DCL resource: %w", err)
	}
	createdDCLObj, err := dclunstruct.Apply(ctx, config, dclObj, dclcontroller.LifecycleParams...)
	if err != nil {
		return nil, fmt.Errorf("error applying the desired resource: %v", err)
	}
	// get the new state in KCC lite format
	newStateLite, err := converter.DCLObjectToKRMObject(createdDCLObj)
	if err != nil {
		return nil, fmt.Errorf("error converting DCL resource to KCC lite: %v", err)
	}
	return dclStateToKRM(resource, newStateLite, converter.MetadataLoader)
}

func dclGet(u *unstructured.Unstructured, config *mmdcl.Config, kubeClient client.Client, converter *dclconversion.Converter, serviceMappingLoader *servicemappingloader.ServiceMappingLoader) (*unstructured.Unstructured, error) {
	resource, err := newDCLResource(u, converter)
	if err != nil {
		return nil, err
	}
	liveLite, err := dcllivestate.FetchLiveState(context.Background(), resource, config, converter, serviceMappingLoader, kubeClient)
	if err != nil {
		return nil, err
	}
	if liveLite == nil {
		return nil, fmt.Errorf("resource '%v' of type '%v' is not found", u.GetName(), u.GroupVersionKind())
	}
	return dclStateToKRM(resource, liveLite, converter.MetadataLoader)
}

func dclDelete(u *unstructured.Unstructured, config *mmdcl.Config, kubeClient client.Client, converter *dclconversion.Converter, serviceMappingLoader *servicemappingloader.ServiceMappingLoader) error {
	ctx := context.Background()
	resource, err := newDCLResource(u, converter)
	if err != nil {
		return err
	}
	lite, err := kcclite.ToKCCLiteBestEffort(resource, converter.MetadataLoader, converter.SchemaLoader, serviceMappingLoader, kubeClient)
	if err != nil {
		return fmt.Errorf("error converting KCC full to KCC lite: %w", err)
	}
	dclObj, err := converter.KRMObjectToDCLObject(lite)
	if err != nil {
		return fmt.Errorf("error converting KCC lite to DCL resource: %w", err)
	}
	if err := dclunstruct.Delete(ctx, config, dclObj); err != nil {
		return fmt.Errorf("error deleting the resource %v: %w", resource.GetNamespacedName(), err)
	}
	return nil
}

func newDCLResource(u *unstructured.Unstructured, converter *dclconversion.Converter) (*dcl.Resource, error) {
	s, err := dclschemaloader.GetDCLSchemaForGVK(u.GroupVersionKind(), converter.MetadataLoader, converter.SchemaLoader)
	if err != nil {
		return nil, err
	}
	resource, err := dcl.NewResource(u, s)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func dclStateToKRM(resource *dcl.Resource, liveState *unstructured.Unstructured, smLoader dclmetadata.ServiceMetadataLoader) (*unstructured.Unstructured, error) {
	spec, status, err := kcclite.ResolveSpecAndStatus(liveState, resource, smLoader)
	if err != nil {
		return nil, err
	}
	resource.Spec = spec
	resource.Status = status
	resource.Labels = liveState.GetLabels()
	return resource.MarshalAsUnstructured()
}

func resourceToKRM(resource *krmtotf.Resource, state *terraform.InstanceState) (*unstructured.Unstructured, error) {
	resource.Spec, resource.Status = krmtotf.ResolveSpecAndStatusWithResourceID(resource, state)
	resource.Labels = krmtotf.GetLabelsFromState(resource, state)
	return resource.MarshalAsUnstructured()
}

func getTerraformResourceAndLiveState(ctx context.Context, u *unstructured.Unstructured, provider *tfschema.Provider, c client.Client, smLoader *servicemappingloader.ServiceMappingLoader) (*krmtotf.Resource,
	*terraform.InstanceState, error) {
	resource, err := newTerraformResource(u, provider, smLoader)
	if err != nil {
		return nil, nil, err
	}
	liveState, err := krmtotf.FetchLiveState(ctx, resource, provider, c, smLoader)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching live state: %v", err)
	}
	return resource, liveState, nil
}

func newTerraformResource(u *unstructured.Unstructured, provider *tfschema.Provider, smLoader *servicemappingloader.ServiceMappingLoader) (*krmtotf.Resource, error) {
	sm, err := smLoader.GetServiceMapping(u.GroupVersionKind().Group)
	if err != nil {
		return nil, err
	}
	resource, err := krmtotf.NewResource(u, sm, provider)
	if err != nil {
		return nil, fmt.Errorf("could not parse resource %s: %v", u.GetName(), err)
	}
	return resource, nil
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "is not found")
}

func unstructuredToKRMResource(isDCLBasedResource bool, u *unstructured.Unstructured, provider *tfschema.Provider, smLoader *servicemappingloader.ServiceMappingLoader, converter *dclconversion.Converter) (*k8s.Resource, error) {
	if isDCLBasedResource {
		return dclUnstructuredToKRMResource(u, converter)
	}
	return terraformUnstructuredToKRMResource(u, provider, smLoader)
}

func dclUnstructuredToKRMResource(u *unstructured.Unstructured, converter *dclconversion.Converter) (*k8s.Resource, error) {
	resource, err := newDCLResource(u, converter)
	if err != nil {
		return nil, err
	}
	return &resource.Resource, nil
}

func terraformUnstructuredToKRMResource(u *unstructured.Unstructured, provider *tfschema.Provider, smLoader *servicemappingloader.ServiceMappingLoader) (*k8s.Resource, error) {
	resource, err := newTerraformResource(u, provider, smLoader)
	if err != nil {
		return nil, err
	}
	return &resource.Resource, nil
}
