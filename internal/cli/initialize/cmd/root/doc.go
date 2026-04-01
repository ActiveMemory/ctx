//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the ctx init command for creating and
// updating .context/ directories.
//
// [Cmd] builds the cobra.Command with --force and --auto-merge
// flags. [Run] orchestrates the full init workflow: create
// directories, deploy templates, generate encryption key, deploy
// hooks and skills, merge settings, and write CLAUDE.md.
package root
