//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status implements the ctx config status subcommand.
//
// [Cmd] builds the cobra.Command. [Run] reads the active .ctxrc
// file, detects the current profile, and displays resolved
// configuration values including context directory, token budget,
// and notification settings.
package status
