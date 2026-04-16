//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package convert transforms context file names into
// human-readable titles for display in agent context
// packets and CLI output.
//
// [FileNameToTitle] converts SCREAMING_SNAKE_CASE.md
// filenames into Title Case strings. The transformation
// strips the .md extension, replaces underscores with
// spaces, and capitalizes the first letter of each word
// while lowercasing the rest. For example:
//
//   - "TASKS.md"           -> "Tasks"
//   - "AGENT_PLAYBOOK.md"  -> "Agent Playbook"
//   - "CONSTITUTION.md"    -> "Constitution"
//
// This is used by the context load pipeline to produce
// section headers in the agent context packet, making
// machine-named files readable without requiring a
// separate display-name mapping.
package convert
