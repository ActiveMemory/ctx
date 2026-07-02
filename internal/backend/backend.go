//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import errBackend "github.com/ActiveMemory/ctx/internal/err/backend"

// SetDefault records the configured default backend name.
//
// Parameters:
//   - name: configured backend name used by Default
func (registry *Registry) SetDefault(name string) {
	registry.defaultName = name
}

// RegisterBuiltin adds a built-in backend factory by type.
//
// Parameters:
//   - name: configured backend name
//   - config: backend configuration passed to the factory
//
// Returns:
//   - error: duplicate registration or missing built-in backend
func (registry *Registry) RegisterBuiltin(name string, config Config) error {
	factory, ok := builtinFactory(config.Type)
	if !ok {
		factory, ok = builtinFactory(name)
	}
	if !ok {
		return errBackend.MissingBackend{Name: name}
	}
	return registry.Register(name, config, factory)
}

// Register adds a backend factory and its config to the registry.
//
// Parameters:
//   - name: configured backend name
//   - config: backend configuration passed to the factory
//   - factory: backend constructor
//
// Returns:
//   - error: duplicate registration error, when name already exists
func (registry *Registry) Register(
	name string,
	config Config,
	factory Factory,
) error {
	registry.ensure()
	if _, ok := registry.factories[name]; ok {
		return errBackend.DuplicateRegistration{Name: name}
	}
	registry.factories[name] = factory
	registry.configs[name] = config
	return nil
}

// Resolve returns the backend selected by name.
//
// Parameters:
//   - name: configured backend name
//
// Returns:
//   - Backend: resolved backend instance
//   - error: missing backend or factory failure
func (registry *Registry) Resolve(name string) (Backend, error) {
	registry.ensure()
	factory, ok := registry.factories[name]
	if !ok {
		return nil, errBackend.MissingBackend{Name: name}
	}
	resolved, resolveErr := factory(registry.configs[name])
	if resolveErr != nil {
		return nil, errBackend.Factory{Name: name, Cause: resolveErr}
	}
	return resolved, nil
}

// Default resolves the registry default backend.
//
// Returns:
//   - Backend: resolved default backend instance
//   - error: empty registry, ambiguous registry, missing default,
//     or factory failure
func (registry *Registry) Default() (Backend, error) {
	registry.ensure()
	if registry.defaultName != "" {
		return registry.Resolve(registry.defaultName)
	}
	count := len(registry.factories)
	if count == 0 {
		return nil, errBackend.NoBackendConfigured{}
	}
	if count > 1 {
		return nil, errBackend.MultipleBackends{}
	}
	for name := range registry.factories {
		return registry.Resolve(name)
	}
	return nil, errBackend.NoBackendConfigured{}
}
