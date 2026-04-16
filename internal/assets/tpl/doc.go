//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tpl holds **Sprintf-based format string constants**
// for output that is too structurally rich to live in the
// flat YAML text store ([internal/config/embed/text]).
//
// The line ctx draws is: simple substitution → YAML; multi-
// line templated content with conditional sections,
// indentation rules, or per-call escaping → here.
//
// # What Lives Here
//
// Each `tpl_*.go` file owns one rendering domain:
//
//   - **`tpl_entry.go`**: the canonical TASKS.md task
//     line and its inline tags (`#priority:`,
//     `#session:`, `#branch:`, `#commit:`, `#added:`).
//     Used by `ctx add task`.
//   - **`tpl_hub_entry.go`**: markdown rendering of one
//     hub entry (date header + origin tag + content
//     body + horizontal rule). Consumed by
//     `ctx connection sync` when materializing entries
//     into `.context/hub/`.
//   - **`tpl_journal.go`**: the journal entry skeleton:
//     YAML frontmatter + body shell that the importer
//     fills in.
//   - **`tpl_loop.go`**: the autonomous-loop shell
//     script template (`ctx loop` output).
//   - **`tpl_obsidian.go`**: the Obsidian vault page
//     templates (note frontmatter + wikilink section).
//   - **`tpl_recall.go`**: the format the legacy
//     `ctx recall` command used; kept here while the
//     journal-merge transition completes.
//   - **`tpl_trigger.go`**: the empty trigger script
//     scaffold installed by `ctx trigger add`.
//
// # Naming Convention
//
// Constants are named for what they render, not how:
// [HubEntryMarkdown], [Task], [TaskPriority], etc.
// Each carries a doc comment listing the Sprintf args
// in order so callers cannot accidentally pass the
// wrong argument order.
//
// # Migration Note
//
// Several templates here are migration candidates for
// Go `text/template`; Sprintf with many positional
// arguments is fragile. The migration is tracked in
// TASKS.md; until then, contributors should add new
// templates here only when the YAML text store cannot
// represent the structure.
//
// # Concurrency
//
// All exports are immutable string constants. Safe
// for any access pattern.
package tpl
