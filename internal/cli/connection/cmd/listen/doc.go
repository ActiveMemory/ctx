//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package listen implements the "ctx connection listen"
// subcommand that streams context entries from a
// connected ctx Hub in real time.
//
// # What It Does
//
// Opens a persistent gRPC stream to the hub using
// the Listen RPC. As new entries arrive they are
// written to .context/hub/ as markdown files and a
// receipt line is printed to stdout. The stream
// runs until the user presses Ctrl-C.
//
// # Flags
//
// None. The command accepts no arguments. Connection
// settings (hub address, token, subscribed types)
// are read from the encrypted config file at
// .context/.connect.enc.
//
// # Output
//
// Prints "Listening..." on startup, then one line
// per received entry showing the entry type. Entries
// are filtered by the types configured via
// "ctx connection subscribe".
//
// # Delegation
//
// [Cmd] builds the cobra.Command and delegates
// directly to [coreListen.Run] which handles config
// loading, gRPC client setup, signal handling, and
// the streaming receive loop.
package listen
