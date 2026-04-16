//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package main is the entry point for the ctx CLI.
//
// The binary performs three steps at startup:
//
//  1. Initializes the embedded asset lookup table via
//     [lookup.Init] so that YAML-backed descriptions,
//     templates, and text constants are available to
//     every subcommand.
//  2. Builds the root cobra.Command through
//     [bootstrap.Initialize], which registers all
//     subcommand packages, wires persistent flags,
//     and injects build-time version metadata via
//     ldflags.
//  3. Calls cmd.Execute and translates any returned
//     error into a formatted stderr message via
//     [writeErr.With] before exiting with code 1.
//
// No business logic lives in this package. All
// domain behavior is delegated to packages under
// internal/.
//
// # Build-Time Injection
//
// Version, commit hash, and build date are injected
// via ldflags at compile time. The Makefile and
// goreleaser config set these values; the bootstrap
// package reads them.
package main
