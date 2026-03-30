//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package source implements the ctx journal source subcommand.
//
// [Cmd] builds the cobra.Command with --limit, --project, --tool,
// --latest, --full, and date range flags. [Run] routes to list
// mode (tabular session overview) or show mode (detailed session
// display) based on flags.
package source
