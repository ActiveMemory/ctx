//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package counter reads and writes integer counters
// persisted as plain-text state files.
//
// Counters are used throughout the system to track
// invocation counts, nudge frequencies, and other
// numeric state that must survive across process
// invocations. Each counter is stored as a single
// integer in a file at a known path.
//
// # Read
//
// [Read] loads an integer from a file. The file
// content is trimmed of whitespace and parsed as a
// base-10 integer. Returns 0 if the file does not
// exist or cannot be parsed, making it safe to call
// without checking file existence first.
//
// # Write
//
// [Write] persists an integer to a file using
// restrictive permissions (0600). Write errors are
// logged as warnings but not returned, treating
// counter persistence as best-effort. This prevents
// counter write failures from interrupting the main
// operation.
package counter
