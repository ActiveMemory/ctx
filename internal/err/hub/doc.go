//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hub defines the typed error constructors
// for the hub subsystem -- the background daemon
// that coordinates multi-project context sharing
// and peer synchronization.
//
// # Domain
//
// Errors fall into three categories:
//
//   - **Token generation** -- the hub failed to
//     generate a cryptographic token for peer
//     authentication. Constructor: [GenerateToken].
//   - **Internal errors** -- a catch-all wrapper
//     for unexpected failures inside the hub
//     server. Constructor: [InternalErr].
//   - **Registration** -- a project is already
//     registered with the hub, or a peer action
//     is unrecognized. Constructors:
//     [DuplicateProject], [InvalidPeerAction].
//
// # Wrapping Strategy
//
// [GenerateToken] and [InternalErr] wrap their
// cause with fmt.Errorf %w so callers can inspect
// the underlying crypto/rand or server error.
// [DuplicateProject] and [InvalidPeerAction]
// return plain formatted errors. All user-facing
// text is resolved through
// [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package hub
