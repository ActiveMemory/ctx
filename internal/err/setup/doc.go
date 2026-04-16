//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package setup defines the typed error constructors
// for the tool setup subsystem. These errors fire
// when `ctx setup` creates directories, writes
// configuration files, or syncs steering files for
// a specific AI tool (Claude Code, Aider, etc.).
//
// # Domain
//
// Four constructors cover the entire surface:
//
//   - [CreateDir] -- a setup directory could not be
//     created. Wraps the underlying OS error.
//   - [MarshalConfig] -- the MCP configuration
//     JSON could not be marshaled.
//   - [WriteFile] -- a setup file could not be
//     written to disk.
//   - [SyncSteering] -- steering file sync failed
//     during the setup sequence.
//
// # Wrapping Strategy
//
// All four constructors wrap their cause with
// fmt.Errorf %w so callers can inspect the
// underlying error. All user-facing text is
// resolved through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package setup
