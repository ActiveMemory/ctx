//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx hook" command for generating AI tool integrations.
//
// The command outputs configuration snippets and instructions for integrating
// Context with various AI coding tools like Claude Code, Cursor, Aider, etc.
//
// Flags:
//   - --write, -w: Write the configuration file instead of printing
//
// Returns:
//   - *cobra.Command: Configured hook command that accepts a tool name argument
func Cmd() *cobra.Command {
	var write bool

	cmd := &cobra.Command{
		Use:   "hook <tool>",
		Short: "Generate AI tool integration configs",
		Long: `Generate configuration and instructions
for integrating Context with AI tools.

Supported tools:
  claude-code  - Anthropic's Claude Code CLI (use plugin instead)
  cursor       - Cursor IDE
  aider        - Aider AI coding assistant
  copilot      - GitHub Copilot
  windsurf     - Windsurf IDE

Use --write to generate the configuration file directly:
  ctx hook copilot --write    # Creates .github/copilot-instructions.md

Example:
  ctx hook cursor`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHook(cmd, args, write)
		},
	}

	cmd.Flags().BoolVarP(
		&write, "write", "w", false,
		"Write the configuration file instead of printing",
	)

	return cmd
}
