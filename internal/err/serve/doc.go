//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package serve defines the typed error constructors
// for the serve subsystem, which manages the hub
// daemon lifecycle: PID file reading, process
// lookups, and graceful shutdown.
//
// # Domain
//
// Three constructors cover the entire surface:
//
//   - [NoRunningHub]: the PID file could not be
//     read, meaning no hub daemon is running.
//     Wraps the underlying read error.
//   - [InvalidPID]: the PID file contents could
//     not be parsed as a valid process ID.
//     Wraps the underlying parse error.
//   - [Kill]: sending a signal to the daemon
//     process failed. Wraps the underlying
//     os.Process.Kill error and includes the PID.
//
// # Wrapping Strategy
//
// All three constructors wrap their cause with
// fmt.Errorf %w so callers can inspect the
// underlying error. All user-facing text is
// resolved through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package serve
