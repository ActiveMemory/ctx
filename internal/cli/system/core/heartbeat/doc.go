//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package heartbeat reads and writes mtime values for
// session activity tracking in state files.
//
// The heartbeat system records the last known
// modification time of context files. Hooks compare
// the stored mtime against the current file mtime to
// detect whether context has changed since the last
// check, triggering nudges or other actions when
// drift is detected.
//
// # ReadMtime
//
// [ReadMtime] loads a stored mtime value from a file.
// The file content is trimmed of whitespace and parsed
// as a base-10 int64. Returns 0 if the file does not
// exist or cannot be parsed, making it safe to call
// without checking file existence first.
//
// # WriteMtime
//
// [WriteMtime] persists an mtime value to a file
// using restrictive permissions (0600). Write errors
// are logged as warnings but not returned, treating
// mtime persistence as best-effort. This prevents
// heartbeat write failures from interrupting hook
// execution.
//
// # Comparison with Counter
//
// The heartbeat package mirrors the counter package
// in structure but stores int64 mtime values instead
// of int counters. Both use the same best-effort
// write pattern with warning-level logging on failure.
package heartbeat
