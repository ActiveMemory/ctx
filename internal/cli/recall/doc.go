//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package recall provides CLI commands for browsing and searching AI session history.
//
// The recall system parses JSONL session files from various AI coding assistants
// (currently Claude Code) and provides commands to list, view, and export sessions.
//
// Commands:
//   - ctx recall list: List all parsed sessions
//   - ctx recall show <id>: Show session details
//   - ctx recall import: Import sessions to editable journal files
//   - ctx recall lock: Protect journal entries from export regeneration
//   - ctx recall unlock: Remove lock protection from journal entries
//   - ctx recall sync: Sync lock state from journal frontmatter to state file
package recall
