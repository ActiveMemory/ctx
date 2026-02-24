//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add provides the "ctx add" command for appending entries to context
// files.
//
// It supports adding decisions, tasks, learnings, and conventions to their
// respective files in the .context/ directory. Content can be provided via
// command argument, --file flag, or stdin pipe.
//
// Supported entry types (defined in [config.FileType]):
//   - decision/decisions: Appends to DECISIONS.md
//   - task/tasks: Inserts into TASKS.md before first unchecked task,
//     or under a named section when --section is provided
//   - learning/learnings: Appends to LEARNINGS.md
//   - convention/conventions: Appends to CONVENTIONS.md
//
// Example usage:
//
//	ctx add decision "Use PostgreSQL for primary database"
//	ctx add task "Implement auth" --priority high --section "Phase 1"
//	ctx add learning --file notes.md
//	echo "Use camelCase" | ctx add convention
package add
