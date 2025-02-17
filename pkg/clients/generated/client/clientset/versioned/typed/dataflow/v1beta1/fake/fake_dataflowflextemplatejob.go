// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// *** DISCLAIMER ***
// Config Connector's go-client for CRDs is currently in ALPHA, which means
// that future versions of the go-client may include breaking changes.
// Please try it out and give us feedback!

// Code generated by main. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/clients/generated/apis/dataflow/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeDataflowFlexTemplateJobs implements DataflowFlexTemplateJobInterface
type FakeDataflowFlexTemplateJobs struct {
	Fake *FakeDataflowV1beta1
	ns   string
}

var dataflowflextemplatejobsResource = schema.GroupVersionResource{Group: "dataflow.cnrm.cloud.google.com", Version: "v1beta1", Resource: "dataflowflextemplatejobs"}

var dataflowflextemplatejobsKind = schema.GroupVersionKind{Group: "dataflow.cnrm.cloud.google.com", Version: "v1beta1", Kind: "DataflowFlexTemplateJob"}

// Get takes name of the dataflowFlexTemplateJob, and returns the corresponding dataflowFlexTemplateJob object, and an error if there is any.
func (c *FakeDataflowFlexTemplateJobs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.DataflowFlexTemplateJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(dataflowflextemplatejobsResource, c.ns, name), &v1beta1.DataflowFlexTemplateJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.DataflowFlexTemplateJob), err
}

// List takes label and field selectors, and returns the list of DataflowFlexTemplateJobs that match those selectors.
func (c *FakeDataflowFlexTemplateJobs) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.DataflowFlexTemplateJobList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(dataflowflextemplatejobsResource, dataflowflextemplatejobsKind, c.ns, opts), &v1beta1.DataflowFlexTemplateJobList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.DataflowFlexTemplateJobList{ListMeta: obj.(*v1beta1.DataflowFlexTemplateJobList).ListMeta}
	for _, item := range obj.(*v1beta1.DataflowFlexTemplateJobList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested dataflowFlexTemplateJobs.
func (c *FakeDataflowFlexTemplateJobs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(dataflowflextemplatejobsResource, c.ns, opts))

}

// Create takes the representation of a dataflowFlexTemplateJob and creates it.  Returns the server's representation of the dataflowFlexTemplateJob, and an error, if there is any.
func (c *FakeDataflowFlexTemplateJobs) Create(ctx context.Context, dataflowFlexTemplateJob *v1beta1.DataflowFlexTemplateJob, opts v1.CreateOptions) (result *v1beta1.DataflowFlexTemplateJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(dataflowflextemplatejobsResource, c.ns, dataflowFlexTemplateJob), &v1beta1.DataflowFlexTemplateJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.DataflowFlexTemplateJob), err
}

// Update takes the representation of a dataflowFlexTemplateJob and updates it. Returns the server's representation of the dataflowFlexTemplateJob, and an error, if there is any.
func (c *FakeDataflowFlexTemplateJobs) Update(ctx context.Context, dataflowFlexTemplateJob *v1beta1.DataflowFlexTemplateJob, opts v1.UpdateOptions) (result *v1beta1.DataflowFlexTemplateJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(dataflowflextemplatejobsResource, c.ns, dataflowFlexTemplateJob), &v1beta1.DataflowFlexTemplateJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.DataflowFlexTemplateJob), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeDataflowFlexTemplateJobs) UpdateStatus(ctx context.Context, dataflowFlexTemplateJob *v1beta1.DataflowFlexTemplateJob, opts v1.UpdateOptions) (*v1beta1.DataflowFlexTemplateJob, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(dataflowflextemplatejobsResource, "status", c.ns, dataflowFlexTemplateJob), &v1beta1.DataflowFlexTemplateJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.DataflowFlexTemplateJob), err
}

// Delete takes name of the dataflowFlexTemplateJob and deletes it. Returns an error if one occurs.
func (c *FakeDataflowFlexTemplateJobs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(dataflowflextemplatejobsResource, c.ns, name, opts), &v1beta1.DataflowFlexTemplateJob{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeDataflowFlexTemplateJobs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(dataflowflextemplatejobsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.DataflowFlexTemplateJobList{})
	return err
}

// Patch applies the patch and returns the patched dataflowFlexTemplateJob.
func (c *FakeDataflowFlexTemplateJobs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.DataflowFlexTemplateJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(dataflowflextemplatejobsResource, c.ns, name, pt, data, subresources...), &v1beta1.DataflowFlexTemplateJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.DataflowFlexTemplateJob), err
}
