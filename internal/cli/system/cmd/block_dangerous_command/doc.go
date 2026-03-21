//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package block_dangerous_commands implements the ctx system
// block-dangerous-commands subcommand.
//
// It provides a regex safety net that catches dangerous command patterns
// such as mid-command sudo, git push, and binary installs that the
// deny-list cannot express.
package block_dangerous_command
