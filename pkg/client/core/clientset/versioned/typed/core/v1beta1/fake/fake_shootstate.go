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

// FakeShootStates implements ShootStateInterface
type FakeShootStates struct {
	Fake *FakeCoreV1beta1
	ns   string
}

var shootstatesResource = v1beta1.SchemeGroupVersion.WithResource("shootstates")

var shootstatesKind = v1beta1.SchemeGroupVersion.WithKind("ShootState")

// Get takes name of the shootState, and returns the corresponding shootState object, and an error if there is any.
func (c *FakeShootStates) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.ShootState, err error) {
	emptyResult := &v1beta1.ShootState{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(shootstatesResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1beta1.ShootState), err
}

// List takes label and field selectors, and returns the list of ShootStates that match those selectors.
func (c *FakeShootStates) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.ShootStateList, err error) {
	emptyResult := &v1beta1.ShootStateList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(shootstatesResource, shootstatesKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.ShootStateList{ListMeta: obj.(*v1beta1.ShootStateList).ListMeta}
	for _, item := range obj.(*v1beta1.ShootStateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested shootStates.
func (c *FakeShootStates) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(shootstatesResource, c.ns, opts))

}

// Create takes the representation of a shootState and creates it.  Returns the server's representation of the shootState, and an error, if there is any.
func (c *FakeShootStates) Create(ctx context.Context, shootState *v1beta1.ShootState, opts v1.CreateOptions) (result *v1beta1.ShootState, err error) {
	emptyResult := &v1beta1.ShootState{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(shootstatesResource, c.ns, shootState, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1beta1.ShootState), err
}

// Update takes the representation of a shootState and updates it. Returns the server's representation of the shootState, and an error, if there is any.
func (c *FakeShootStates) Update(ctx context.Context, shootState *v1beta1.ShootState, opts v1.UpdateOptions) (result *v1beta1.ShootState, err error) {
	emptyResult := &v1beta1.ShootState{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(shootstatesResource, c.ns, shootState, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1beta1.ShootState), err
}

// Delete takes name of the shootState and deletes it. Returns an error if one occurs.
func (c *FakeShootStates) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(shootstatesResource, c.ns, name, opts), &v1beta1.ShootState{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeShootStates) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(shootstatesResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.ShootStateList{})
	return err
}

// Patch applies the patch and returns the patched shootState.
func (c *FakeShootStates) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.ShootState, err error) {
	emptyResult := &v1beta1.ShootState{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(shootstatesResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1beta1.ShootState), err
}
