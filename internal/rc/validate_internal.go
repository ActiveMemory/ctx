//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	cfgRC "github.com/ActiveMemory/ctx/internal/config/rc"
)

// validateBackends checks semantic requirements for .ctxrc backends.
//
// Returns:
//   - error: validation failure, or nil when configuration is valid
func (cfg CtxRC) validateBackends() error {
	backends := cfg.Backends
	if len(backends.Configs) == 0 {
		return nil
	}

	if backends.Default != "" {
		if _, ok := backends.Configs[backends.Default]; !ok {
			return &yaml.TypeError{Errors: []string{
				cfgRC.ErrBackendsDefaultMissing + backends.Default,
			}}
		}
	}

	for name, backend := range backends.Configs {
		if backend.Endpoint == "" {
			return &yaml.TypeError{Errors: []string{
				fmt.Sprintf(cfgRC.ErrBackendsEndpointRequired, name),
			}}
		}
	}

	return nil
}

// backendsShapeError reports whether type errors are malformed backends shape.
//
// Parameters:
//   - errs: YAML type error strings
//
// Returns:
//   - bool: true when the error should fail validation, not warn
func backendsShapeError(errs []string) bool {
	for _, msg := range errs {
		if strings.Contains(msg, cfgRC.ErrBackendsMapping) ||
			strings.Contains(msg, cfgRC.ErrBackendsDefaultScalar) {
			return true
		}
	}

	return false
}
