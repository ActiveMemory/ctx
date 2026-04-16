//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core is the umbrella for the journal command's
// processing pipeline.
//
// # Overview
//
// The journal system converts Claude Code session JSONL
// transcripts into structured markdown files with YAML
// frontmatter, then generates navigable indexes and
// optional Obsidian vault output. This package groups
// the sub-packages that implement each stage.
//
// # Pipeline Stages
//
// The journal pipeline flows through these stages:
//
//  1. plan: scans session files and builds an import
//     plan of file actions (new, regenerate, skip).
//  2. confirm: presents the plan summary and prompts
//     for user confirmation.
//  3. execute: writes journal entries to disk,
//     preserving existing frontmatter on regeneration.
//  4. normalize: cleans up content boundaries and
//     formatting.
//  5. generate: builds site indexes, maps of content,
//     and navigation pages.
//
// Each stage is idempotent and tracks progress via a
// shared .state.json file.
//
// # Sub-packages
//
//   - collapse: collapses verbose tool output blocks
//   - confirm: user confirmation prompts
//   - consolidate: merges repeated tool runs
//   - execute: import plan execution
//   - extract: YAML frontmatter extraction
//   - format: size, slug, and link formatting
//   - frontmatter: frontmatter parsing and types
//   - generate: site and index generation
//   - group: entry grouping by topic
//   - index: journal index management
//   - lock: concurrent access locking
//   - moc: map of content generation
//   - normalize: content boundary normalisation
//   - obsidian: Obsidian vault builder
//   - parse: session JSONL parsing
//   - plan: import planning
//   - query: journal entry querying
//   - reduce: content reduction
//   - schema: schema validation
//   - section: section indexing
//   - session: session metadata
//   - slug: filename slug generation
//   - source: source file listing and formatting
//   - turn: turn header and body extraction
//   - validate: entry validation
//   - wikilink: wikilink processing
package core
