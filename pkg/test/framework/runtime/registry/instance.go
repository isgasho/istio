//  Copyright 2018 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package registry

import (
	"fmt"
	"testing"

	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/runtime/api"
)

var _ component.Defaults = &Instance{}

// Instance of a component registry
type Instance struct {
	defaults  map[component.ID]component.Descriptor
	factories map[component.Key]api.ComponentFactory
}

// New component registry
func New() *Instance {
	return &Instance{
		defaults:  make(map[component.ID]component.Descriptor),
		factories: make(map[component.Key]api.ComponentFactory),
	}
}

// Register a component
func (r *Instance) Register(desc component.Descriptor, isDefault bool, factory api.ComponentFactory) {
	k := desc.Key
	if k.ID == "" {
		panic("attempting to register framework component without an ID")
	}

	if r.factories[k] != nil {
		panic(fmt.Sprintf("duplicate components registered `%s`", desc.FriendlyName()))
	}
	r.factories[k] = factory

	if isDefault {
		if _, ok := r.defaults[k.ID]; ok {
			panic(fmt.Sprintf("default already set for component `%s`", desc.FriendlyName()))
		}
		r.defaults[k.ID] = desc
	}
}

// GetFactory for a component
func (r *Instance) GetFactory(desc component.Descriptor) (api.ComponentFactory, error) {
	k := desc.Key
	f := r.factories[k]
	// If the key was a Variant and there was no factory for the variant, try the default factory.
	if f == nil && desc.Key.Variant != "" {
		f = r.factories[desc.Key.ID.GetKey()]
	}
	if f == nil {
		return nil, fmt.Errorf("unknown component `%s`", desc.FriendlyName())
	}
	return f, nil
}

// GetDefaultDescriptor implements Defaults interface
func (r *Instance) GetDefaultDescriptor(id component.ID) (component.Descriptor, error) {
	d, ok := r.defaults[id]
	if !ok {
		return component.Descriptor{}, fmt.Errorf("unknown component `%s`", id)
	}
	return d, nil
}

// GetDefaultDescriptorOrFail implements Defaults interface
func (r *Instance) GetDefaultDescriptorOrFail(id component.ID, t testing.TB) component.Descriptor {
	d, err := r.GetDefaultDescriptor(id)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
