// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/gardener/gardener/pkg/apis/garden/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSecretBindings implements SecretBindingInterface
type FakeSecretBindings struct {
	Fake *FakeGardenV1beta1
	ns   string
}

var secretbindingsResource = schema.GroupVersionResource{Group: "garden.sapcloud.io", Version: "v1beta1", Resource: "secretbindings"}

var secretbindingsKind = schema.GroupVersionKind{Group: "garden.sapcloud.io", Version: "v1beta1", Kind: "SecretBinding"}

// Get takes name of the secretBinding, and returns the corresponding secretBinding object, and an error if there is any.
func (c *FakeSecretBindings) Get(name string, options v1.GetOptions) (result *v1beta1.SecretBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(secretbindingsResource, c.ns, name), &v1beta1.SecretBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SecretBinding), err
}

// List takes label and field selectors, and returns the list of SecretBindings that match those selectors.
func (c *FakeSecretBindings) List(opts v1.ListOptions) (result *v1beta1.SecretBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(secretbindingsResource, secretbindingsKind, c.ns, opts), &v1beta1.SecretBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.SecretBindingList{ListMeta: obj.(*v1beta1.SecretBindingList).ListMeta}
	for _, item := range obj.(*v1beta1.SecretBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested secretBindings.
func (c *FakeSecretBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(secretbindingsResource, c.ns, opts))

}

// Create takes the representation of a secretBinding and creates it.  Returns the server's representation of the secretBinding, and an error, if there is any.
func (c *FakeSecretBindings) Create(secretBinding *v1beta1.SecretBinding) (result *v1beta1.SecretBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(secretbindingsResource, c.ns, secretBinding), &v1beta1.SecretBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SecretBinding), err
}

// Update takes the representation of a secretBinding and updates it. Returns the server's representation of the secretBinding, and an error, if there is any.
func (c *FakeSecretBindings) Update(secretBinding *v1beta1.SecretBinding) (result *v1beta1.SecretBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(secretbindingsResource, c.ns, secretBinding), &v1beta1.SecretBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SecretBinding), err
}

// Delete takes name of the secretBinding and deletes it. Returns an error if one occurs.
func (c *FakeSecretBindings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(secretbindingsResource, c.ns, name), &v1beta1.SecretBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSecretBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(secretbindingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.SecretBindingList{})
	return err
}

// Patch applies the patch and returns the patched secretBinding.
func (c *FakeSecretBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.SecretBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(secretbindingsResource, c.ns, name, pt, data, subresources...), &v1beta1.SecretBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.SecretBinding), err
}
