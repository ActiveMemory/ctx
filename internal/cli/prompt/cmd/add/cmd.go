//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"
)

// Cmd returns the prompt add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func Cmd() *cobra.Command {
	var fromStdin bool

	cmd := &cobra.Command{
		Use:   "add NAME",
		Short: "Create a new prompt from embedded template or stdin",
		Long: `Create a new prompt template in .context/prompts/.

By default, creates from an embedded starter template if one exists
with the given name. Use --stdin to read content from standard input.

Examples:
  ctx prompt add code-review
  echo "# My Prompt" | ctx prompt add my-prompt --stdin`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd, args[0], fromStdin)
		},
	}

	cmd.Flags().BoolVar(&fromStdin, "stdin", false, "read prompt content from stdin")

	return cmd
}
