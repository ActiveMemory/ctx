//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package inspect implements `ctx disclosure inspect <file>`: a
// read-only report of a canonical knowledge file's staged entries and
// current themes, for the progressive-disclosure dry-run pass.
//
// # Domain
//
// [Cmd] wires the cobra command; [Run] resolves the file's kind (via
// internal/disclosure.KindFor), reads it, builds an
// internal/disclosure.Inspection, and renders it through
// internal/write/disclosure. It writes nothing to the file.
//
// # Related packages
//
//   - internal/disclosure — Parse/Inspect the root.
//   - internal/write/disclosure — render the Inspection.
//   - internal/err/disclosure — the not-a-knowledge-file error.
package inspect
