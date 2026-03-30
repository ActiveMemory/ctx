//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package unlock implements the ctx journal unlock subcommand.
//
// [Cmd] builds the cobra.Command with --all flag. [Run] removes
// lock protection from journal entries, allowing future import
// regeneration to update them.
package unlock
