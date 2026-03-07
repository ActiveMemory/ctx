//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
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

	short, long := assets.CommandDesc("hook")
	cmd := &cobra.Command{
		Use:         "hook <tool>",
		Short:       short,
		Annotations: map[string]string{config.AnnotationSkipInit: "true"},
		Long:        long,
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, write)
		},
	}

	cmd.Flags().BoolVarP(
		&write, "write", "w", false,
		assets.FlagDesc("hook.write"),
	)

	return cmd
}
