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

// FakeSidecarSets implements SidecarSetInterface
type FakeSidecarSets struct {
	Fake *FakeAppsV1alpha1
}

var sidecarsetsResource = schema.GroupVersionResource{Group: "apps.kruise.io", Version: "v1alpha1", Resource: "sidecarsets"}

var sidecarsetsKind = schema.GroupVersionKind{Group: "apps.kruise.io", Version: "v1alpha1", Kind: "SidecarSet"}

// Get takes name of the sidecarSet, and returns the corresponding sidecarSet object, and an error if there is any.
func (c *FakeSidecarSets) Get(name string, options v1.GetOptions) (result *v1alpha1.SidecarSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(sidecarsetsResource, name), &v1alpha1.SidecarSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SidecarSet), err
}

// List takes label and field selectors, and returns the list of SidecarSets that match those selectors.
func (c *FakeSidecarSets) List(opts v1.ListOptions) (result *v1alpha1.SidecarSetList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(sidecarsetsResource, sidecarsetsKind, opts), &v1alpha1.SidecarSetList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SidecarSetList{ListMeta: obj.(*v1alpha1.SidecarSetList).ListMeta}
	for _, item := range obj.(*v1alpha1.SidecarSetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested sidecarSets.
func (c *FakeSidecarSets) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(sidecarsetsResource, opts))
}

// Create takes the representation of a sidecarSet and creates it.  Returns the server's representation of the sidecarSet, and an error, if there is any.
func (c *FakeSidecarSets) Create(sidecarSet *v1alpha1.SidecarSet) (result *v1alpha1.SidecarSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(sidecarsetsResource, sidecarSet), &v1alpha1.SidecarSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SidecarSet), err
}

// Update takes the representation of a sidecarSet and updates it. Returns the server's representation of the sidecarSet, and an error, if there is any.
func (c *FakeSidecarSets) Update(sidecarSet *v1alpha1.SidecarSet) (result *v1alpha1.SidecarSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(sidecarsetsResource, sidecarSet), &v1alpha1.SidecarSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SidecarSet), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSidecarSets) UpdateStatus(sidecarSet *v1alpha1.SidecarSet) (*v1alpha1.SidecarSet, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(sidecarsetsResource, "status", sidecarSet), &v1alpha1.SidecarSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SidecarSet), err
}

// Delete takes name of the sidecarSet and deletes it. Returns an error if one occurs.
func (c *FakeSidecarSets) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(sidecarsetsResource, name), &v1alpha1.SidecarSet{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSidecarSets) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(sidecarsetsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.SidecarSetList{})
	return err
}

// Patch applies the patch and returns the patched sidecarSet.
func (c *FakeSidecarSets) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SidecarSet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(sidecarsetsResource, name, pt, data, subresources...), &v1alpha1.SidecarSet{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SidecarSet), err
}
