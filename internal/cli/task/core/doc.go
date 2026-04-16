//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core is the umbrella for task subcommand
// business logic. It contains no code of its own;
// all functionality lives in its subpackages.
//
// # Subpackages
//
// The task core layer is split into focused packages
// that each handle one aspect of task management:
//
//   - archive -- moves completed tasks from TASKS.md
//     into timestamped archive files
//   - complete -- finds and marks individual tasks as
//     done by number or text search
//   - count -- counts pending top-level tasks, excluding
//     subtasks
//   - path -- resolves the absolute paths to TASKS.md
//     and the archive directory
//
// # Architecture
//
// Each subpackage exports pure business logic functions
// that the cmd/ layer calls. The cmd/ layer handles
// Cobra command wiring, flag parsing, and output
// formatting. This separation keeps the core testable
// without Cobra dependencies.
package core
