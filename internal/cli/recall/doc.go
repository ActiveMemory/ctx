//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
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
//   - ctx recall export: Export sessions to editable journal files
package recall
