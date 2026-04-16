//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package zensical wraps the zensical static site
// generator binary.
//
// # Running Zensical
//
// Run launches zensical with a subcommand in the
// given working directory. It checks LookPath first
// and returns an error if the binary is not found.
//
//	err := zensical.Run("/path/to/site", "build")
//	err := zensical.Run("/path/to/site", "serve")
//
// # I/O Wiring
//
// The zensical process inherits the parent's stdout,
// stderr, and stdin so that build output and
// interactive prompts flow through to the user's
// terminal.
//
// # Error Handling
//
// If the zensical binary is not in PATH, Run returns
// a ZensicalNotFound error from the site error
// package. Other errors propagate from exec.Cmd.Run.
package zensical
