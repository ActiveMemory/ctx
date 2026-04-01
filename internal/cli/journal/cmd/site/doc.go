//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package site implements the ctx journal site subcommand.
//
// [Cmd] builds the cobra.Command with --build and --output flags.
// [Run] generates a static journal site: parses entries, builds
// month-grouped pages, topic indexes, and a zensical configuration.
// With --build, it also invokes zensical to produce HTML.
package site
