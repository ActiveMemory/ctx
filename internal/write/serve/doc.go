//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package serve provides terminal output for the
// ctx Hub server command (ctx serve).
//
// The Hub is a local HTTP server that exposes
// context operations over a REST API. Output
// functions cover the server lifecycle from
// startup through shutdown.
//
// # Startup
//
// [HubStarted] prints the network address the
// server is listening on. [AdminToken] prints the
// generated admin token for authenticating API
// requests. Both are emitted at launch before the
// server begins accepting connections.
//
// # Background Mode
//
// [Daemonized] confirms the hub started as a
// background process, printing the daemon PID.
// [Stopped] confirms a running daemon was killed,
// also printing the PID that was terminated.
package serve
