// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/gardener/gardener/pkg/apis/security/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CredentialsBindingLister helps list CredentialsBindings.
// All objects returned here must be treated as read-only.
type CredentialsBindingLister interface {
	// List lists all CredentialsBindings in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.CredentialsBinding, err error)
	// CredentialsBindings returns an object that can list and get CredentialsBindings.
	CredentialsBindings(namespace string) CredentialsBindingNamespaceLister
	CredentialsBindingListerExpansion
}

// credentialsBindingLister implements the CredentialsBindingLister interface.
type credentialsBindingLister struct {
	indexer cache.Indexer
}

// NewCredentialsBindingLister returns a new CredentialsBindingLister.
func NewCredentialsBindingLister(indexer cache.Indexer) CredentialsBindingLister {
	return &credentialsBindingLister{indexer: indexer}
}

// List lists all CredentialsBindings in the indexer.
func (s *credentialsBindingLister) List(selector labels.Selector) (ret []*v1alpha1.CredentialsBinding, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CredentialsBinding))
	})
	return ret, err
}

// CredentialsBindings returns an object that can list and get CredentialsBindings.
func (s *credentialsBindingLister) CredentialsBindings(namespace string) CredentialsBindingNamespaceLister {
	return credentialsBindingNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CredentialsBindingNamespaceLister helps list and get CredentialsBindings.
// All objects returned here must be treated as read-only.
type CredentialsBindingNamespaceLister interface {
	// List lists all CredentialsBindings in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.CredentialsBinding, err error)
	// Get retrieves the CredentialsBinding from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.CredentialsBinding, error)
	CredentialsBindingNamespaceListerExpansion
}

// credentialsBindingNamespaceLister implements the CredentialsBindingNamespaceLister
// interface.
type credentialsBindingNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CredentialsBindings in the indexer for a given namespace.
func (s credentialsBindingNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.CredentialsBinding, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.CredentialsBinding))
	})
	return ret, err
}

// Get retrieves the CredentialsBinding from the indexer for a given namespace and name.
func (s credentialsBindingNamespaceLister) Get(name string) (*v1alpha1.CredentialsBinding, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("credentialsbinding"), name)
	}
	return obj.(*v1alpha1.CredentialsBinding), nil
}
