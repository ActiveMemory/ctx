//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package stepdown implements the "ctx hub stepdown"
// subcommand that requests leadership transfer from
// the current hub node.
//
// # What It Does
//
// Tells the current Raft leader node to voluntarily
// give up leadership so another node in the cluster
// can take over. This is useful for planned
// maintenance or graceful role rotation in a
// multi-node hub deployment.
//
// # Flags
//
// None. The command accepts no arguments.
//
// # Output
//
// Prints a confirmation line indicating that the
// leadership stepdown was requested.
//
// # Delegation
//
// [Cmd] builds the cobra.Command and delegates
// directly to [coreStep.Run] which sends the
// stepdown request and writes the confirmation
// message.
package stepdown
