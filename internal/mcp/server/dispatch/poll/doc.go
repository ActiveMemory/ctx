//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package poll implements the **resource-change watcher**
// behind the MCP `resources/subscribe` notification. When
// a client subscribes to one or more `.context/` files,
// this package polls their mtimes and emits a
// `notifications/resources/updated` JSON-RPC message
// when any of them changes.
//
// MCP supports change notifications, but the underlying
// substrate (a polling watcher) is opaque to the client.
// This package is that substrate.
//
// # Public Surface
//
//   - **[NewPoller](paths, intervalMs)** — builds a
//     poller for the given file paths with the
//     given polling interval.
//   - **Poller methods** — `Start(ctx, ch)` to
//     begin emitting change events on `ch`,
//     `Stop()` to cease, `Update(paths)` to swap
//     the watch set without restart.
//
// # Why Polling, Not fsnotify
//
// Polling at ~1 Hz is reliable across every
// platform ctx supports (Linux, macOS, Windows)
// without per-platform watcher quirks (file
// renames, fsync timing, cross-FS edge cases).
// MCP's notification cadence does not need
// sub-second precision; "saw a change within a
// second" is enough.
//
// # Concurrency
//
// `Start` spawns a single goroutine that ticks
// on the configured interval; `Stop` signals it
// via context cancellation. Concurrent calls to
// `Update` are serialized through the poller's
// mutex.
package poll
