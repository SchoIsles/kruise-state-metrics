/*
Copyright 2020 The SchoIsles Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/openkruise/kruise-api/apps/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeUnitedDeployments implements UnitedDeploymentInterface
type FakeUnitedDeployments struct {
	Fake *FakeAppsV1alpha1
	ns   string
}

var uniteddeploymentsResource = schema.GroupVersionResource{Group: "apps.kruise.io", Version: "v1alpha1", Resource: "uniteddeployments"}

var uniteddeploymentsKind = schema.GroupVersionKind{Group: "apps.kruise.io", Version: "v1alpha1", Kind: "UnitedDeployment"}

// Get takes name of the unitedDeployment, and returns the corresponding unitedDeployment object, and an error if there is any.
func (c *FakeUnitedDeployments) Get(name string, options v1.GetOptions) (result *v1alpha1.UnitedDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(uniteddeploymentsResource, c.ns, name), &v1alpha1.UnitedDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.UnitedDeployment), err
}

// List takes label and field selectors, and returns the list of UnitedDeployments that match those selectors.
func (c *FakeUnitedDeployments) List(opts v1.ListOptions) (result *v1alpha1.UnitedDeploymentList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(uniteddeploymentsResource, uniteddeploymentsKind, c.ns, opts), &v1alpha1.UnitedDeploymentList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.UnitedDeploymentList{ListMeta: obj.(*v1alpha1.UnitedDeploymentList).ListMeta}
	for _, item := range obj.(*v1alpha1.UnitedDeploymentList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested unitedDeployments.
func (c *FakeUnitedDeployments) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(uniteddeploymentsResource, c.ns, opts))

}

// Create takes the representation of a unitedDeployment and creates it.  Returns the server's representation of the unitedDeployment, and an error, if there is any.
func (c *FakeUnitedDeployments) Create(unitedDeployment *v1alpha1.UnitedDeployment) (result *v1alpha1.UnitedDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(uniteddeploymentsResource, c.ns, unitedDeployment), &v1alpha1.UnitedDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.UnitedDeployment), err
}

// Update takes the representation of a unitedDeployment and updates it. Returns the server's representation of the unitedDeployment, and an error, if there is any.
func (c *FakeUnitedDeployments) Update(unitedDeployment *v1alpha1.UnitedDeployment) (result *v1alpha1.UnitedDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(uniteddeploymentsResource, c.ns, unitedDeployment), &v1alpha1.UnitedDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.UnitedDeployment), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeUnitedDeployments) UpdateStatus(unitedDeployment *v1alpha1.UnitedDeployment) (*v1alpha1.UnitedDeployment, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(uniteddeploymentsResource, "status", c.ns, unitedDeployment), &v1alpha1.UnitedDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.UnitedDeployment), err
}

// Delete takes name of the unitedDeployment and deletes it. Returns an error if one occurs.
func (c *FakeUnitedDeployments) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(uniteddeploymentsResource, c.ns, name), &v1alpha1.UnitedDeployment{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeUnitedDeployments) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(uniteddeploymentsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.UnitedDeploymentList{})
	return err
}

// Patch applies the patch and returns the patched unitedDeployment.
func (c *FakeUnitedDeployments) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.UnitedDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(uniteddeploymentsResource, c.ns, name, pt, data, subresources...), &v1alpha1.UnitedDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.UnitedDeployment), err
}
