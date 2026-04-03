//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package file provides the "ctx trace file" CLI subcommand.
//
// It runs git log for a given file path and displays context refs
// attached to each commit that touched the file. Supports an optional
// :line-range suffix on the path argument (e.g. "src/auth.go:42-60").
//
// Key exports: [Cmd], [Run].
// Used by the trace command tree to register the file subcommand.
package file
