//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package state defines the typed error constructors
// for runtime state persistence. These errors fire
// when ctx reads, loads, or saves its internal state
// files (.context/.state/) which track things like
// last-run timestamps, import cursors, and feature
// flags.
//
// # Domain
//
// Three constructors cover the entire surface:
//
//   - [ReadingDir]: the state directory could not
//     be read. Wraps the underlying OS error.
//   - [Load]: a state file could not be loaded
//     (read + unmarshal). Wraps the underlying
//     error.
//   - [Save]: a state file could not be saved
//     (marshal + write). Wraps the underlying
//     error.
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
package state
