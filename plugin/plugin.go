/*
   Copyright The containerd Authors.

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

package plugin

import (
	"errors"
	"fmt"
	"sync"

	"github.com/containerd/plugin"
)

var (
	// ErrNoType is returned when no type is specified
	ErrNoType = errors.New("plugin: no type")
	// ErrNoPluginID is returned when no id is specified
	ErrNoPluginID = errors.New("plugin: no id")
	// ErrIDRegistered is returned when a duplicate id is already registered
	ErrIDRegistered = errors.New("plugin: id already registered")
	// ErrSkipPlugin is used when a plugin is not initialized and should not be loaded,
	// this allows the plugin loader differentiate between a plugin which is configured
	// not to load and one that fails to load.
	ErrSkipPlugin = errors.New("skip plugin")

	// ErrInvalidRequires will be thrown if the requirements for a plugin are
	// defined in an invalid manner.
	ErrInvalidRequires = errors.New("invalid requires")
)

// IsSkipPlugin returns true if the error is skipping the plugin
func IsSkipPlugin(err error) bool {
	return errors.Is(err, ErrSkipPlugin)
}

// Registration contains information for registering a plugin
type Registration struct {
	// Type of the plugin
	Type plugin.Type
	// ID of the plugin
	ID string
	// Config specific to the plugin
	Config interface{}
	// Requires is a list of plugins that the registered plugin requires to be available
	Requires []plugin.Type

	// InitFn is called when initializing a plugin. The registration and
	// context are passed in. The init function may modify the registration to
	// add exports, capabilities and platform support declarations.
	InitFn func(*InitContext) (interface{}, error)
	// Disable the plugin from loading
	Disable bool
}

// Init the registered plugin
func (r *Registration) Init(ic *InitContext) *Plugin {
	p, err := r.InitFn(ic)
	return &Plugin{
		Registration: r,
		Config:       ic.Config,
		Meta:         ic.Meta,
		instance:     p,
		err:          err,
	}
}

// URI returns the full plugin URI
func (r *Registration) URI() string {
	return fmt.Sprintf("%s.%s", r.Type, r.ID)
}

var register = struct {
	sync.RWMutex
	r []*Registration
}{}

// Register allows plugins to register
func Register(r *Registration) {
	register.Lock()
	defer register.Unlock()

	if r.Type == "" {
		panic(ErrNoType)
	}
	if r.ID == "" {
		panic(ErrNoPluginID)
	}
	if err := checkUnique(r); err != nil {
		panic(err)
	}

	for _, requires := range r.Requires {
		if requires == "*" && len(r.Requires) != 1 {
			panic(ErrInvalidRequires)
		}
	}

	register.r = append(register.r, r)
}

func checkUnique(r *Registration) error {
	for _, registered := range register.r {
		if r.URI() == registered.URI() {
			return fmt.Errorf("%s: %w", r.URI(), ErrIDRegistered)
		}
	}
	return nil
}

// DisableFilter filters out disabled plugins
type DisableFilter func(r *Registration) bool

// Graph returns an ordered list of registered plugins for initialization.
// Plugins in disableList specified by id will be disabled.
func Graph(filter DisableFilter) (ordered []*Registration) {
	register.RLock()
	defer register.RUnlock()

	for _, r := range register.r {
		if filter(r) {
			r.Disable = true
		}
	}

	added := map[*Registration]bool{}
	for _, r := range register.r {
		if r.Disable {
			continue
		}
		children(r, added, &ordered)
		if !added[r] {
			ordered = append(ordered, r)
			added[r] = true
		}
	}
	return ordered
}

func children(reg *Registration, added map[*Registration]bool, ordered *[]*Registration) {
	for _, t := range reg.Requires {
		for _, r := range register.r {
			if !r.Disable &&
				r.URI() != reg.URI() &&
				(t == "*" || r.Type == t) {
				children(r, added, ordered)
				if !added[r] {
					*ordered = append(*ordered, r)
					added[r] = true
				}
			}
		}
	}
}
