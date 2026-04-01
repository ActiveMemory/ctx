//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package main is the entry point for the ctx CLI.
//
// The binary delegates immediately to [bootstrap.Execute], which
// builds the root cobra.Command, registers all 24 subcommand
// packages, and calls cmd.Execute. No business logic lives here;
// version injection happens via ldflags at build time.
package main
