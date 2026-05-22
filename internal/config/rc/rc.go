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
