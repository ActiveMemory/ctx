//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package journal implements the "ctx journal" command for
// analyzing and publishing exported AI session files.
//
// The journal system ingests session transcripts from
// .context/journal/, enriches them with YAML frontmatter
// metadata (topics, type, outcome, technologies,
// key_files), and publishes them in browsable formats.
//
// # Subcommands
//
//   - source: list or inspect raw journal entries
//   - import: ingest exported session files into the
//     journal directory
//   - schema: output the journal entry JSON Schema
//   - lock: mark a journal entry as finalized
//   - unlock: revert a locked entry to editable state
//   - sync: synchronize journal state with the source
//   - site: generate a zensical-compatible static site
//     with browsable history, indices, and search
//   - obsidian: generate an Obsidian vault with
//     wikilinks, MOC pages, and graph cross-linking
//
// Both site and obsidian formats reuse the same
// scan/parse/index infrastructure and consume identical
// enriched journal entries.
//
// # Subpackages
//
//	cmd/source: entry listing and inspection
//	cmd/importer: session file ingestion
//	cmd/schema: JSON Schema output
//	cmd/lock, cmd/unlock: entry finalization
//	cmd/sync: state synchronization
//	cmd/site: static site generation
//	cmd/obsidian: Obsidian vault generation
//	core: scan, parse, index, and enrichment logic
package journal
