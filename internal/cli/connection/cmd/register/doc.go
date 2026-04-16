//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package register implements the "ctx connection register"
// subcommand that registers this project with a ctx Hub
// instance.
//
// # What It Does
//
// Connects to the hub at the given address, sends the
// admin token and project name, receives a client
// token in return, and stores the encrypted connection
// config in .context/.connect.enc. After registration
// the project can publish, sync, and listen for
// entries.
//
// # Arguments
//
// Requires exactly one positional argument:
//
//   - args[0]: hub gRPC address (host:port)
//
// # Flags
//
//   - --token (required): the admin token printed
//     by "ctx hub start" at server startup. Used to
//     authenticate the registration request.
//
// # Output
//
// Prints a confirmation line showing the assigned
// client ID from the hub.
//
// # Delegation
//
// [Cmd] builds the cobra.Command, binds the --token
// flag, marks it required, and delegates to
// [coreReg.Run] which handles the gRPC handshake
// and config persistence.
package register
