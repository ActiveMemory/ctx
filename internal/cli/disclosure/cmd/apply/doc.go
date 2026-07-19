//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package apply implements `ctx disclosure apply <file> --plan <path|->`:
// the write half of progressive disclosure.
//
// # Domain
//
// [Cmd] wires the cobra command; [Run] resolves the file's kind (via
// internal/disclosure.KindFor), reads a digest plan (JSON, from --plan or
// stdin), and hands both to internal/disclosure.Apply, which moves the
// named staged entries into per-theme files and folds their gists into
// the root under its append→verify→remove guards. The result is rendered
// through internal/write/disclosure.
//
// # Related packages
//
//   - internal/disclosure — Apply: the guarded mover.
//   - internal/write/disclosure — render the ApplyResult.
//   - internal/err/disclosure — the not-a-knowledge-file error.
package apply
