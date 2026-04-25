//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package activate provides error factories for the `ctx activate`
// subcommand. The factories cover the five failure modes surfaced
// by the command:
//
//   - [NoCandidates]: scan from CWD found zero .context/ dirs.
//   - [Ambiguous]: scan found multiple and refuses to pick.
//   - [InvalidPath]: explicit path cannot be stat()ed.
//   - [NotDirectory]: explicit path exists but is not a directory.
//   - [NotContext]: explicit path is a directory but lacks any
//     canonical context file (CONSTITUTION.md or TASKS.md).
//
// Messages are loaded via desc.Text using the DescKey constants in
// internal/config/embed/text/err_activate.go so they stay editable
// without code changes.
package activate
