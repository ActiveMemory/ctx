//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package event formats and renders event log entries
// for the system events command.
//
// Events are notifications emitted by hooks during
// agent sessions. This package transforms raw event
// payloads into displayable output in two formats.
//
// # Timestamp Formatting
//
// [FormatTimestamp] converts an RFC3339 timestamp to
// local time using a precise date-time layout. Returns
// the original string on parse failure, ensuring
// display never fails.
//
// # Hook Name Extraction
//
// [ExtractHookName] gets the hook name from an event
// payload. It first checks the Detail.Hook field, then
// falls back to extracting a prefix before the first
// colon in the message. Returns a fallback constant
// when neither source yields a name.
//
// # Message Truncation
//
// [TruncateMessage] limits message length for columnar
// display, appending a truncation suffix when the
// message exceeds the maximum length.
//
// # Output Formats
//
// [FormatJSON] serializes events as JSONL lines, one
// per event. Marshal errors are silently skipped.
//
// [FormatHuman] renders events in aligned columns
// showing timestamp, event type, hook name, and
// truncated message. The column format string is loaded
// from embedded assets.
package event
