//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package initialize hosts compile-time constants consumed by
// the ctx init command: backup directory naming for
// ctx init --reset and the canonical reset flag name.
//
// Sentinel error values for ctx init refusal and reset live in
// `internal/err/initialize/`; their user-facing text lives in
// `commands/text/errors.yaml` and is resolved through
// `desc.Text` at error-display time by the sentinels' typed
// `Error()` methods. The wrapping constructors (Populated,
// ResetRequiresInteractive) also flow through desc.Text.
package initialize
