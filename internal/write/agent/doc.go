//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package agent provides terminal output for the agent
// command (ctx agent).
//
// # Exported Functions
//
// [Packet] prints a pre-rendered markdown context packet
// to stdout. The packet is assembled by the agent core
// package and contains a budget-constrained snapshot of
// the project's context files, formatted for consumption
// by AI coding assistants.
//
// # Message Categories
//
//   - Info: rendered markdown content printed verbatim
//     to stdout without additional formatting
//
// # Nil Safety
//
// A nil *cobra.Command is treated as a no-op, making
// Packet safe to call from paths where a command may
// not be available.
//
// # Usage
//
//	content := core.Assemble(files, budget)
//	agent.Packet(cmd, content)
package agent
