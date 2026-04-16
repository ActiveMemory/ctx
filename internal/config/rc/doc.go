//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package rc defines default values for hooks,
// steering, and timeout configuration read from the
// .ctxrc file.
//
// # What .ctxrc Controls
//
// The .ctxrc file is a YAML configuration file at the
// project root that lets users customize ctx behavior
// without editing source code. Settings include hook
// enable/disable flags, timeout durations, memory
// classification rules, session prefixes, and
// steering overrides.
//
// # Role of This Package
//
// This package provides the compile-time defaults
// that apply when a .ctxrc key is absent or the file
// does not exist. Domain packages read their defaults
// from here rather than hard-coding them, so all
// fallback values are auditable in a single location.
//
// # Why a Separate Package
//
// Keeping defaults in config/rc avoids import cycles:
// the rc loader depends on config packages for
// structure, and config packages depend on rc for
// defaults. Isolating the defaults breaks the cycle.
package rc
