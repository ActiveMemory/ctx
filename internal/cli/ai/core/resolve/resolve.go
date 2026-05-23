//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"github.com/ActiveMemory/ctx/internal/backend"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Build constructs a [backend.Registry] populated with
// the six built-in factories ([backend.RegisterAll]) and
// every per-project backend configuration loaded from
// `.ctxrc` ([rc.Backends]). The `default_backend:` rc
// key is forwarded via [backend.Registry.SetDefault].
//
// Returns:
//   - *backend.Registry: ready for Resolve / Default
//     calls.
//   - error: any error from [backend.RegisterAll]
//     (duplicate factory name; under normal use returns
//     nil).
func Build() (*backend.Registry, error) {
	r := backend.New()
	if err := backend.RegisterAll(r); err != nil {
		return nil, err
	}
	for _, cfg := range rc.Backends() {
		r.Configure(backend.Config(cfg))
	}
	if def := rc.DefaultBackend(); def != "" {
		r.SetDefault(def)
	}
	return r, nil
}

// Pick returns either the user-named backend (when
// `--backend <name>` was passed) or the registry's
// default. Helper used by every `ctx ai *` subcommand so
// the flag-vs-default branching lives in one place.
//
// Parameters:
//   - r: a Registry built via [Build].
//   - name: --backend value, or "" for default.
//
// Returns:
//   - backend.Backend: the resolved backend.
//   - error: typed err/backend sentinel on failure.
func Pick(r *backend.Registry, name string) (backend.Backend, error) {
	if name != "" {
		return r.Resolve(name)
	}
	return r.Default()
}
