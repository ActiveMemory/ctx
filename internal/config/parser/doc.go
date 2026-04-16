//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parser defines buffer sizes, directory
// names, and prefix lists used by the session file
// parser that converts Claude Code JSONL transcripts
// into structured journal entries.
//
// # Scanner Buffers
//
// Claude Code session files can contain very long
// lines (tool results with full file contents). Three
// buffer tiers handle different parsing contexts:
//
//   - BufInitSize (64 KB) -- initial scanner buffer.
//   - BufMaxSize (1 MB) -- maximum for normal session
//     parsing.
//   - BufMaxSizeSchema (4 MB) -- maximum for schema
//     validation, where tool result lines can exceed
//     the normal cap.
//
// # Format Detection
//
// LinesToPeek (50) is the number of lines the parser
// reads when auto-detecting whether a file is JSONL
// or Markdown format.
//
// # Subagent Filtering
//
// DirSubagents ("subagents") names the directory
// that holds sidechain sessions sharing the parent
// session ID. The parser skips this directory to
// avoid importing duplicate content.
//
// # Session Header Prefixes
//
// DefaultSessionPrefixes (["Session:"]) lists the
// Markdown heading prefixes that mark the start of
// a new session in hand-written transcript files.
// Users can extend this list via the
// session_prefixes key in .ctxrc.
package parser
