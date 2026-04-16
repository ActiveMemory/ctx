//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package block_dangerous_command implements the
// **`ctx system block-dangerous-command`** hidden hook,
// which intercepts shell commands that match dangerous
// patterns before they execute.
//
// # What It Does
//
// The hook reads a JSON envelope from stdin containing
// the command the agent is about to run. It tests the
// command against a set of compiled regular expressions
// and, if any match, emits a JSON block response that
// prevents execution. The patterns catch:
//
//   - **mid-command sudo** -- e.g. "echo foo | sudo rm"
//   - **git push** -- any direct push attempt
//   - **cp/mv to bin** -- copying files into system
//     binary directories
//   - **install to /usr/local/bin** -- direct binary
//     installation outside package managers
//
// When a command is blocked, a relay notification is
// also sent to the nudge channel so the agent sees the
// reason for the block.
//
// # Input
//
// A JSON hook envelope on stdin with a ToolInput.Command
// field containing the shell command string.
//
// # Output
//
// On match: a JSON [entity.BlockResponse] with decision
// "block" and a human-readable reason. On no match:
// no output (silent pass-through).
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin, tests each regex in priority order, loads the
// appropriate message template via [core/message.Load],
// and marshals the block response. Relay notifications
// are sent through [core/nudge.Relay].
package block_dangerous_command
