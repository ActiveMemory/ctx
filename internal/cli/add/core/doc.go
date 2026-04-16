//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains business logic for the add command.
//
// This package is an umbrella for the subpackages that power
// ctx add. It does not export functions itself; instead it
// coordinates five concerns through its child packages:
//
// # Entry Classification
//
// The entry subpackage classifies user input into context
// file types (task, decision, learning, convention) using
// predicate functions such as FileTypeIsTask. It also
// detects when a task description is complex enough to
// warrant a spec nudge via NeedsSpec.
//
// # Content Extraction
//
// The extract subpackage resolves the entry body from one
// of three sources in priority order: the --file flag, CLI
// positional arguments, or piped stdin. It returns an error
// when no source provides content.
//
// # Markdown Formatting
//
// The format subpackage renders each entry type into its
// Markdown representation. Tasks become checkbox items with
// provenance tags (session, branch, commit, timestamp).
// Decisions and learnings become structured ADR-style
// sections with context, rationale, and consequence fields.
// Conventions become simple list items.
//
// # Section-Aware Insertion
//
// The insert subpackage places formatted entries into the
// correct position within existing context files. Tasks
// land before the first pending item or after an explicit
// section header. Decisions and learnings insert in
// reverse-chronological order before existing entries.
// Conventions append at the end.
//
// # Section Normalization
//
// The normalize subpackage ensures user-provided section
// names carry the correct Markdown heading prefix before
// insertion.
//
// # Example Text
//
// The example subpackage loads type-specific usage examples
// from embedded YAML assets for display in cobra help text.
//
// # Data Flow
//
// The cmd/ layer calls extract.Content to obtain text, then
// entry predicates to choose a formatter from the format
// package, and finally insert.AppendEntry to merge the
// result into the target file. The write/ layer persists
// the bytes to disk.
package core
