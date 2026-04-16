//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package source implements the "ctx journal source"
// command.
//
// # Overview
//
// The source command provides two modes for working with
// raw Claude Code JSONL session files:
//
//   - List mode (default): prints a tabular overview of
//     available sessions with timestamps, slugs, and
//     summary metadata.
//   - Show mode: displays detailed content for a single
//     session, identified by slug, ID, or positional arg.
//
// # Flags
//
//	-s, --show <id>      Show a specific session by slug
//	                      or session ID.
//	    --latest          Show the most recent session.
//	    --full            Include full session content in
//	                      show mode (no truncation).
//	-n, --limit <n>      Maximum sessions to list
//	                      (default from config).
//	-p, --project <name> Filter by project name.
//	-t, --tool <name>    Filter by tool name.
//	    --since <date>    Include sessions on or after
//	                      this date.
//	    --until <date>    Include sessions on or before
//	                      this date.
//	    --all-projects    Scan all project directories.
//
// # Behavior
//
// [Cmd] builds the cobra.Command and registers all flags
// listed above. [Run] inspects the flags to determine the
// mode: if --show, --latest, or a positional argument is
// present, it dispatches to coreSrc.RunShow; otherwise it
// dispatches to coreSrc.RunList.
//
// # Output
//
// List mode outputs a formatted table to stdout. Show
// mode outputs session details including timestamps,
// message counts, and optionally the full conversation
// content.
package source
