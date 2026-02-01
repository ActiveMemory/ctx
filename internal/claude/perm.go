//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

// DefaultPermissions returns the default permissions for ctx commands.
//
// These permissions allow Claude Code to run ctx CLI commands without
// prompting for approval. All ctx subcommands are pre-approved.
//
// Returns:
//   - []string: List of permission patterns for ctx commands
func DefaultPermissions() []string {
	return []string{
		"Bash(ctx status:*)",
		"Bash(ctx agent:*)",
		"Bash(ctx add:*)",
		"Bash(ctx session:*)",
		"Bash(ctx tasks:*)",
		"Bash(ctx loop:*)",
	}
}
