//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hub implements the shared context hub server and client.
//
// The hub is a gRPC service that aggregates published entries
// (decisions, learnings, conventions) from multiple ctx instances
// and streams them to subscribers in real-time.
//
// Storage is append-only JSONL. Auth is token-based (admin token
// for registration, per-client tokens for RPCs). Connection config
// is encrypted locally using AES-256-GCM via [internal/crypto].
//
// Key exports: [Store], [Entry], [Auth], [Server], [Client].
// See source files for implementation details.
// Part of the internal subsystem.
package hub
