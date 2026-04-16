//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory implements the "ctx memory" command for
// bridging Claude Code's auto memory into the .context/
// directory.
//
// The memory bridge discovers MEMORY.md from Claude
// Code's project-scoped auto memory, mirrors it locally
// for drift detection, and supports importing classified
// entries into structured context files. This ensures
// knowledge captured by Claude Code's auto-memory system
// feeds back into the persistent project context.
//
// # Subcommands
//
//   - sync: copy MEMORY.md to .context/memory/mirror.md
//     with archival of the previous mirror
//   - status: show drift status and line counts between
//     source and mirror
//   - diff: show line-level differences between mirror
//     and source
//   - import: classify and promote entries into the
//     appropriate context files (decisions, learnings,
//     conventions, tasks)
//   - publish: push curated context back into MEMORY.md
//   - unpublish: remove published sections from
//     MEMORY.md
//
// # Subpackages
//
//	cmd/sync, cmd/status, cmd/diff -- mirror operations
//	cmd/importer -- entry classification and promotion
//	cmd/publish, cmd/unpublish -- bidirectional sync
//	core -- mirror discovery, diff, and classification
package memory
