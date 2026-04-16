//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx serve" command
// for running a local static site server via
// zensical.
//
// # Behavior
//
// The command starts a local development server for
// a zensical-powered static site. With no arguments,
// it serves the journal site located at
// .context/journal-site. With a directory argument,
// it serves that directory instead.
//
// Before starting, the command validates that the
// target directory exists, is actually a directory
// (not a file), and contains a zensical.toml
// configuration file. Missing or invalid targets
// produce clear error messages.
//
// This command does NOT start a ctx Hub server. For
// Hub functionality, use "ctx hub start" instead.
//
// # Flags
//
// None. This command takes no flags.
//
// # Output
//
// The command delegates to the zensical binary, which
// handles its own output (server URL, request logs).
// On startup failure, returns an error identifying
// the problem (missing directory, missing config,
// or zensical not found).
//
// # Delegation
//
// Directory validation is done inline. Site serving
// is delegated to [execZensical.Run] which locates
// and executes the zensical binary. The default
// journal site path is constructed from
// [rc.ContextDir] and [dir.JournalSite].
package root
