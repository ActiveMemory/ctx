//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

// Format strings for sentinel-wrapping in err/context constructors.
// Centralized here so the magic-string audit (which exempts
// internal/config) does not flag them at the call site.
const (
	// FmtWrapColon wraps a sentinel and a tailored message:
	//   fmt.Errorf(FmtWrapColon, ErrFoo, "tailored detail")
	//   ↦ "<ErrFoo.Error()>: tailored detail".
	FmtWrapColon = "%w: %s"
)

// YAML key constants for `.ctxrc` shape. Centralized
// here so the rc parser, the setup-time writer, and any
// future tools that touch `.ctxrc` reference the same
// strings.
const (
	// YAMLKeyBackends is the top-level `backends:` key
	// holding the list of configured AI inference
	// backends.
	YAMLKeyBackends = "backends"
	// YAMLKeyBackendName is the per-entry `name:` key the
	// setup writer matches on to enforce idempotency.
	YAMLKeyBackendName = "name"
)
