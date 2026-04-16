//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package watch implements the "ctx watch" command for
// processing structured context-update commands from AI
// output.
//
// The watch command reads stdin for XML-style update
// commands and applies them to the corresponding context
// files. This enables AI assistants to emit structured
// updates that ctx processes into persistent context,
// bridging free-form AI output with the structured
// .context/ directory.
//
// # Supported Commands
//
//   - <task>: add or update a task in TASKS.md
//   - <decision>: append a decision to DECISIONS.md
//   - <learning>: append a learning to LEARNINGS.md
//   - <convention>: append a convention to CONVENTIONS.md
//   - <complete>: mark a task as completed in TASKS.md
//
// # Dry-Run Mode
//
// The --dry-run flag previews changes without modifying
// any files. This is useful for validating AI output
// before committing changes to context.
//
// # Subpackages
//
//   - cmd/root: cobra command definition, stdin parsing,
//     and command dispatch
//   - core: XML command parsing, file mutation, and
//     dry-run rendering
package watch
