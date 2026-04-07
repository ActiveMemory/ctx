//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema provides write functions for schema CLI output.
//
// All terminal output for the schema check and dump commands is
// routed through this package per the project convention that
// cmd.Print* calls live in internal/write/. Functions accept
// primitive types (strings, ints) rather than domain types to
// avoid cross-package type references that would trigger the
// CrossPackageTypes audit.
package schema
