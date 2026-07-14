//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package index implements the `ctx index` command: a computed table of
// contents over any Markdown knowledge file.
//
// `ctx index <file>` projects the file's ATX headings (`##`, and `###` with
// --depth 3) in file order, replacing the stored `<!-- INDEX -->` blocks that
// DECISIONS.md / LEARNINGS.md once carried. The index is never written back —
// it is recomputed on demand, so it cannot drift from the entries it
// summarizes. One generic command serves every knowledge file: DECISIONS,
// LEARNINGS, CONVENTIONS, and TASKS (`## Phase …`) alike.
//
// The heading recognizer lives in [internal/heading]; this package is the CLI
// surface (argument + flags) and delegates rendering to [internal/write/index].
//
// # Subpackages
//
//	cmd/root: cobra command wiring (positional FILE arg, --depth, --json)
package index
