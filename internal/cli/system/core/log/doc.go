//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package log appends timestamped messages to hook log
// files with automatic rotation. Hooks use this to record
// diagnostic output without blocking the session.
//
// # Logging
//
// [Message] appends a single timestamped log line to the
// given file. Each line contains a precise timestamp, a
// truncated session ID (first 8 characters), and the
// message text. The parent directory is created if it
// does not exist.
//
// Before appending, Message calls [Rotate] to check
// whether the log file needs rotation.
//
// # Rotation
//
// [Rotate] implements a simple one-generation rotation
// policy. When the log file exceeds HookLogMaxBytes, the
// current file is renamed with a ".1" suffix, replacing
// any previous generation. This keeps at most two files
// on disk: the active log and one rotated copy.
//
// The rotation strategy matches the eventlog pattern
// used elsewhere in the codebase.
//
// # Data Flow
//
//  1. Hook calls Message with file path and message
//  2. Message ensures parent directory exists
//  3. Rotate checks file size against threshold
//  4. If oversized, current file becomes .1 backup
//  5. New line is appended with timestamp and session ID
package log
