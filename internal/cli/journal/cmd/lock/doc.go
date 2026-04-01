//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package lock implements the ctx journal lock subcommand.
//
// [Cmd] builds the cobra.Command with --all flag. [Run] marks
// journal entries as locked, preventing future import regeneration
// from overwriting enriched content.
package lock
