//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package path resolves full filesystem paths for every artifact
// the ctx knowledge-base editorial pipeline reads or writes.
//
// All functions compose [rc.ContextDir] with the directory and
// filename constants from
// [github.com/ActiveMemory/ctx/internal/config/kb]; nothing here
// invents new strings. Tests verify the joins are stdlib-clean
// (no string concatenation, per CONSTITUTION.md Quality
// Invariants).
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb] is the
//     name source.
//   - [github.com/ActiveMemory/ctx/internal/rc] resolves the
//     declared .context/ directory.
package path
