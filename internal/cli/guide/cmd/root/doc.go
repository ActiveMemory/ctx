//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the ctx guide command that lists
// available skills and CLI commands.
//
// [Cmd] builds the cobra.Command with --skills and --commands
// flags. [Run] reads skills from the embedded plugin and commands
// from the bootstrap registry, then renders both lists.
package root
