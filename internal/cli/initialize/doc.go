//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package initialize implements the "ctx init" command for initializing a
// .context/ directory with template files.
//
// The init command creates the foundation for persistent AI context by
// generating template files for constitution rules, tasks, decisions,
// learnings, conventions, and architecture documentation. It also sets
// up Claude Code integration with hooks and slash commands.
//
// # File Organization
//
//   - init.go: Command definition and flag registration
//   - run.go: Main execution logic and orchestration
//   - validate.go: PATH validation for ctx executable
//   - fs.go: File system operations and marker handling
//   - claude.go: CLAUDE.md creation and merge logic
//   - tpl.go: Entry template creation
//   - cmd.go: Claude Code slash command creation
//   - hook.go: Claude Code hook and settings creation
package initialize
