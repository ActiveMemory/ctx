//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package start implements the "ctx hub start" subcommand
// that launches the ctx Hub gRPC server.
//
// # What It Does
//
// Starts the hub server either in the foreground or
// as a detached daemon process. The hub provides a
// gRPC API for publishing, syncing, and streaming
// context entries across projects. When --peers is
// set, the server joins a Raft cluster for leader
// election and entry replication.
//
// # Flags
//
//   - --port: gRPC listen port. Defaults to the
//     value from [server.DefaultPort].
//   - --data-dir: Directory for persistent storage
//     (entries, Raft state, PID file).
//   - --daemon: Run as a background daemon. Writes
//     a PID file for later stop/status commands.
//   - --peers: Comma-separated list of peer
//     addresses (host:port) to form a Raft cluster.
//
// # Output
//
// In foreground mode, prints the admin token and
// listen address, then blocks until interrupted.
// In daemon mode, prints the admin token and PID,
// then detaches.
//
// # Delegation
//
// [Cmd] builds the cobra.Command and binds all
// flags. When --daemon is set it calls
// [server.RunDaemon]; otherwise it parses the
// peers string via [server.ParsePeers] and calls
// [server.Run] for foreground operation.
package start
