// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSeeds implements SeedInterface
type FakeSeeds struct {
	Fake *FakeCoreV1beta1
}

var seedsResource = v1beta1.SchemeGroupVersion.WithResource("seeds")

var seedsKind = v1beta1.SchemeGroupVersion.WithKind("Seed")

// Get takes name of the seed, and returns the corresponding seed object, and an error if there is any.
func (c *FakeSeeds) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.Seed, err error) {
	emptyResult := &v1beta1.Seed{}
	obj, err := c.Fake.
		Invokes(testing.NewRootGetActionWithOptions(seedsResource, name, options), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1beta1.Seed), err
}

// List takes label and field selectors, and returns the list of Seeds that match those selectors.
func (c *FakeSeeds) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.SeedList, err error) {
	emptyResult := &v1beta1.SeedList{}
	obj, err := c.Fake.
		Invokes(testing.NewRootListActionWithOptions(seedsResource, seedsKind, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.SeedList{ListMeta: obj.(*v1beta1.SeedList).ListMeta}
	for _, item := range obj.(*v1beta1.SeedList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested seeds.
func (c *FakeSeeds) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchActionWithOptions(seedsResource, opts))
}

// Create takes the representation of a seed and creates it.  Returns the server's representation of the seed, and an error, if there is any.
func (c *FakeSeeds) Create(ctx context.Context, seed *v1beta1.Seed, opts v1.CreateOptions) (result *v1beta1.Seed, err error) {
	emptyResult := &v1beta1.Seed{}
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateActionWithOptions(seedsResource, seed, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1beta1.Seed), err
}

// Update takes the representation of a seed and updates it. Returns the server's representation of the seed, and an error, if there is any.
func (c *FakeSeeds) Update(ctx context.Context, seed *v1beta1.Seed, opts v1.UpdateOptions) (result *v1beta1.Seed, err error) {
	emptyResult := &v1beta1.Seed{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateActionWithOptions(seedsResource, seed, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1beta1.Seed), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSeeds) UpdateStatus(ctx context.Context, seed *v1beta1.Seed, opts v1.UpdateOptions) (result *v1beta1.Seed, err error) {
	emptyResult := &v1beta1.Seed{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceActionWithOptions(seedsResource, "status", seed, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1beta1.Seed), err
}

// Delete takes name of the seed and deletes it. Returns an error if one occurs.
func (c *FakeSeeds) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(seedsResource, name, opts), &v1beta1.Seed{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSeeds) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionActionWithOptions(seedsResource, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.SeedList{})
	return err
}

// Patch applies the patch and returns the patched seed.
func (c *FakeSeeds) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Seed, err error) {
	emptyResult := &v1beta1.Seed{}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceActionWithOptions(seedsResource, name, pt, data, opts, subresources...), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1beta1.Seed), err
}
