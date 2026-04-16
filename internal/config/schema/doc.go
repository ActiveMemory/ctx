//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema defines the vocabulary for the JSONL
// schema validation system that checks Claude Code
// session transcripts for structural drift.
//
// When Claude Code writes a session, every line is a
// JSON object with a type field and a nested message
// envelope. This package declares every field name,
// record type, content block type, and format string
// the validator and drift reporter need.
//
// # Field Lists
//
//   - [RequiredFields]  — fields that must appear on
//     every user/assistant record (uuid, sessionId,
//     timestamp, type, message, ...).
//   - [OptionalFields]  — fields that may appear
//     (gitBranch, todos, agentId, ...).
//
// # Record and Block Types
//
//   - [RecordUser], [RecordAssistant] — message
//     record types that carry conversation content.
//   - [MetadataRecordTypes] — metadata-only types
//     (no field validation): last-prompt,
//     custom-title, attachment, tag, etc.
//   - [InfraRecordTypes] — infrastructure types
//     skipped by the parser: progress, summary, etc.
//   - [ParsedBlockTypes] — content blocks the
//     parser extracts (text, thinking, tool_use,
//     tool_result).
//   - [KnownBlockTypes] — recognized but unparsed
//     (server_tool_use, mcp_tool_use, ...).
//
// # Report Format Strings
//
// Format strings for both terminal summaries and
// Markdown drift reports live here: headings, table
// rows, suggestion text, and scan statistics.
//
// # Schema Version
//
// [Version] tracks the current schema version and
// [CCVersionRange] records the Claude Code versions
// the schema has been tested against.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package schema
