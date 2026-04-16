//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package server is the **server-side runtime** for
// `ctx hub start` — daemon lifecycle, PID file management,
// and the wire-up between the [internal/hub] package and
// the user-facing CLI flags (`--port`, `--peers`,
// `--daemon`).
//
// The package is the bridge: [internal/hub] knows how to
// be a hub, this package knows how to *run* one as a
// daemon process.
//
// # Public Surface
//
//   - **[Run](opts)** — foreground server boot. Binds
//     the listener, instantiates the [hub.Server],
//     wires the optional [hub.Cluster] when `--peers`
//     is passed, blocks on serve. Honors signals
//     (SIGINT, SIGTERM) for graceful shutdown.
//   - **[DefaultPort]** — the canonical port (9900)
//     used by docs, examples, and the recipes.
//
// # Daemon Mode
//
// When the user passes `--daemon`, the parent forks a
// detached child, writes `<dataDir>/hub.pid` with the
// child's PID, and exits. The PID file is what
// `ctx hub stop` consumes to send SIGTERM.
//
// # PID File Lifecycle
//
//   - **Created** atomically on daemon start.
//   - **Removed** by the child on graceful shutdown.
//   - **Stale-detected** by `ctx hub status` (PID does
//     not refer to a running process) so a crashed
//     hub does not block a fresh start.
//
// # Concurrency
//
// The server runs in the same process as gRPC
// dispatch; this package starts it and waits. No
// in-process concurrency primitives beyond what
// [internal/hub] and the gRPC runtime already provide.
package server
