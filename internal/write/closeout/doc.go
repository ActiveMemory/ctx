//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package closeout writes and reads per-pass closeout artifacts
// for the ctx knowledge-base editorial pipeline (Phase KB).
//
// A closeout lives at
// `.context/ingest/closeouts/<TIMESTAMP>-<mode>-closeout.md`
// and carries required frontmatter:
//
//	---
//	sha: <short>
//	branch: <name>
//	mode: <ingest|ask|site-review|ground|note>
//	pass-mode: <topic-page|triage|evidence-only>  (only for ingest)
//	life-stage: <bootstrap|maintenance>
//	generated-at: <RFC-3339 with timezone>
//	---
//
// The body sections (Inputs, Pass-mode block, Topic(s) touched,
// What changed, etc.) are mode-aware; this package treats the
// body as opaque bytes and only asserts on frontmatter.
//
// Closeouts are append-never-rewrite. Once written they are
// immutable; the only legal state change is archival, which
// physically moves the file to .context/archive/closeouts/
// without modifying its contents.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/write/handover]
//     consumes closeouts via [List] + [PostdatedBy] + [Archive]
//     during handover fold.
//   - [github.com/ActiveMemory/ctx/internal/git_meta] supplies
//     the (sha, branch) pair stamped into the frontmatter.
//   - [github.com/ActiveMemory/ctx/internal/config/kb] supplies
//     the mode + closeout-suffix constants.
package closeout
