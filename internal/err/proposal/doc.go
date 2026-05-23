//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package proposal owns the typed-error wrapping
// constructors used by the proposal-queue writer
// (`internal/write/proposal/`). Errors flow from disk
// failures (mkdir, write) and from the JSON-parse step
// the extract command performs before writing.
//
// Per the project's error-handling convention, the
// user-facing strings live in
// `internal/assets/commands/text/errors.yaml` keyed by
// `err.proposal.*`; this package resolves the keys via
// `desc.Text` at error-construction time.
package proposal
