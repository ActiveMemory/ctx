//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package skill implements the "ctx skill" command group
// for managing reusable instruction bundles that can be
// installed, listed, and removed from the project context.
//
// Skills are markdown files that define slash-command
// behaviors for AI coding assistants. They live in
// .claude/skills/ and are discovered by tools like
// Claude Code at session start. The skill command group
// provides CLI tooling to manage the installed set.
//
// # Subcommands
//
//   - install: copy a skill file from a source path or
//     URL into the project's .claude/skills/ directory
//   - list: display all installed skills with their
//     trigger descriptions
//   - remove: delete a skill by name from the project
//
// # Subpackages
//
//	cmd/install -- skill installation logic
//	cmd/list -- skill enumeration and display
//	cmd/remove -- skill deletion
package skill
