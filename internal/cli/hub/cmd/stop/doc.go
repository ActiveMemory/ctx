//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package stop implements the "ctx hub stop" subcommand
// that shuts down a daemonized ctx Hub server.
//
// # What It Does
//
// Reads the PID file from the hub data directory,
// sends SIGTERM to the running hub process, and
// removes the PID file on success. This is the
// counterpart to "ctx hub start --daemon".
//
// # Flags
//
//   - --data-dir: Directory where the hub stores
//     its PID file and persistent data. Must match
//     the --data-dir used when starting the hub.
//
// # Output
//
// Prints a confirmation line when the hub process
// is successfully terminated.
//
// # Delegation
//
// [Cmd] builds the cobra.Command, binds the
// --data-dir flag, and delegates to [server.Stop]
// which reads the PID file, sends the signal, and
// cleans up.
package stop
