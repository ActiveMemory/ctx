//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
)

// New returns an empty [Registry] with no factories or
// configs bound.
//
// Returns:
//   - *Registry: ready for [Registry.Register] /
//     [Registry.Configure] calls.
func New() *Registry {
	return &Registry{
		factories: make(map[string]Factory),
		configs:   make(map[string]Config),
	}
}

// Register binds a backend type label to its [Factory].
// Returns a wrapped [errBackend.ErrDuplicateRegistration]
// if name is already bound.
//
// Parameters:
//   - name: backend type label, e.g., "vllm".
//   - f: constructor for the backend.
//
// Returns:
//   - error: duplicate-registration error or nil.
func (r *Registry) Register(name string, f Factory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.factories[name]; exists {
		return errBackend.DuplicateRegistration(name)
	}
	r.factories[name] = f
	return nil
}

// Configure stores a per-project [Config] addressed by
// its Name field. The Name does not have to match a
// previously-registered Factory at Configure time;
// resolution-time failure is preferred so .ctxrc parse
// order is independent of Factory registration order.
//
// Parameters:
//   - cfg: per-backend settings loaded from `.ctxrc`.
func (r *Registry) Configure(cfg Config) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.configs[cfg.Name] = cfg
}

// SetDefault names which configured backend
// [Registry.Default] returns when no explicit name is
// passed. The empty string clears the default. The named
// backend is not validated at SetDefault time; an
// invalid name surfaces at Default() as
// [errBackend.ErrBackendNotFound].
//
// Parameters:
//   - name: configured backend name, or empty to clear.
func (r *Registry) SetDefault(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.deflt = name
}

// Resolve returns the [Backend] instance for the named
// type.
//
// Parameters:
//   - name: configured backend name.
//
// Returns:
//   - Backend: constructed via the registered Factory.
//   - error: [errBackend.ErrNoBackends] when no
//     backends are configured at all;
//     [errBackend.ErrBackendNotFound] when name has no
//     Config or no Factory.
func (r *Registry) Resolve(name string) (Backend, error) {
	r.mu.RLock()
	cfg, hasCfg := r.configs[name]
	f, hasFactory := r.factories[name]
	n := len(r.configs)
	r.mu.RUnlock()

	if n == 0 {
		return nil, errBackend.ErrNoBackends
	}
	if !hasCfg || !hasFactory {
		return nil, errBackend.NotFound(name)
	}
	return f(cfg)
}

// Default returns the [Backend] selected by
// [Registry.SetDefault]. If no default is set and
// exactly one backend is configured, that one is
// returned. Otherwise the call fails closed so the user
// is forced to pass `--backend <name>` rather than have
// a wrong default chosen for them.
//
// Returns:
//   - Backend: the default-selected instance.
//   - error: [errBackend.ErrNoBackends] when nothing
//     is configured; [errBackend.ErrAmbiguousDefault]
//     when multiple backends are configured and no
//     explicit default is set.
func (r *Registry) Default() (Backend, error) {
	r.mu.RLock()
	name := r.deflt
	n := len(r.configs)
	var only string
	if n == 1 {
		for k := range r.configs {
			only = k
		}
	}
	r.mu.RUnlock()

	switch {
	case name != "":
		return r.Resolve(name)
	case n == 0:
		return nil, errBackend.ErrNoBackends
	case n == 1:
		return r.Resolve(only)
	default:
		return nil, errBackend.ErrAmbiguousDefault
	}
}
