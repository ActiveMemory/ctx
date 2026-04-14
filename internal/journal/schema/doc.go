//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema is ctx's defense against the **silent
// drift** of Claude Code's session-file format.
//
// Claude Code stores sessions as JSONL files under
// `~/.claude/projects/<slug>/` with an **undocumented,
// unversioned** record format that changes between
// releases. There is no schema URL, no version tag, no
// compatibility commitment. The only way to know whether
// a new Claude Code release added a field, removed a type,
// or quietly renamed a property is to compare empirical
// reality to a frozen reference shape — which is what this
// package does.
//
// # The Reference Shape
//
// [schema.go] declares the **expected** record shape
// derived from analysis of real session files: the set of
// known top-level fields, the set of known record `type`
// values, the set of known content-block types within
// `assistant` records, and the per-type required-field
// list.
//
// [build.go] / [check.go] / [validate.go] walk an actual
// JSONL file and accumulate findings into a [Collector]:
//
//   - **Unknown fields** — a key the reference shape does
//     not list (Claude added a property).
//   - **Missing required fields** — a key the reference
//     shape requires but the record omits (Claude
//     removed a property; we may now silently drop data
//     downstream).
//   - **Unknown record types** — a `type` value not in
//     the reference set (a new record kind appeared).
//   - **Unknown content block types** — same, but for
//     content blocks inside `assistant` records.
//
// # Strictly Informational
//
// Validation **never blocks** an import or other
// operation. Findings flow into [Report] which formats a
// markdown drift report consumed by `ctx doctor` and the
// release-prep runbook. The intent is "tell me when the
// upstream shape moved so I can update the parser", not
// "refuse to ingest anything we have not pre-blessed".
//
// # Updating the Reference
//
// When a new Claude Code release introduces fields the
// drift report flags:
//
//  1. Inspect the new records to confirm semantics.
//  2. Update the reference declarations in [schema.go].
//  3. Update [internal/journal/parser] if the new
//     fields carry session-relevant data.
//  4. Add a learning to LEARNINGS.md so the change is
//     not repeated when reviewing the next release.
//
// # Concurrency
//
// All exported functions are pure data transformations
// over byte slices and `[]Finding`. Concurrent callers
// never race.
//
// # Related Packages
//
//   - [internal/journal/parser]   — the actual session
//     parser; this package is its **canary**.
//   - [internal/cli/journal/cmd/schema]  — the
//     `ctx journal schema` CLI surface.
//   - [internal/entity]           — finding and report
//     types.
package schema
