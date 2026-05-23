//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"github.com/ActiveMemory/ctx/internal/entity"
)

// Backends returns the configured AI inference backends
// translated into [entity.BackendConfig] values.
//
// Timeout strings are parsed via [time.ParseDuration].
// Unparseable timeouts emit a warning and the resulting
// Config carries a zero Timeout, leaving the backend
// implementation to apply its own default. This matches
// the broader rc parsing posture: parse failures are
// warnings, not hard errors, so a single typo cannot
// prevent ctx from booting.
//
// Returns:
//   - []entity.BackendConfig: configured backends;
//     empty when no `backends:` table is in `.ctxrc`.
func Backends() []entity.BackendConfig {
	raw := RC().Backends
	if len(raw) == 0 {
		return nil
	}
	out := make([]entity.BackendConfig, 0, len(raw))
	for _, b := range raw {
		out = append(out, entity.BackendConfig{
			Name:         b.Name,
			Endpoint:     b.Endpoint,
			APIKeyEnv:    b.APIKeyEnv,
			Timeout:      parseTimeout(b.Name, b.Timeout),
			DefaultModel: b.DefaultModel,
		})
	}
	return out
}

// DefaultBackend returns the name of the backend selected
// when no `--backend` flag is passed to a `ctx ai *`
// command. Empty string means the default is unset; with
// a single backend configured, callers should treat that
// one as implicit.
//
// Returns:
//   - string: configured default backend name, or "".
func DefaultBackend() string {
	return RC().DefaultBackend
}
