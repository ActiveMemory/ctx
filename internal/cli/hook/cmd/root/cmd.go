//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
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

	short, long := desc.CommandDesc(cmd.DescKeyHook)
	cmd := &cobra.Command{
		Use:         cmd.DescKeyHook + " <tool>",
		Short:       short,
		Annotations: map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		Long:        long,
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, write)
		},
	}

	cmd.Flags().BoolVarP(
		&write, "write", "w", false,
		desc.FlagDesc(flag.DescKeyHookWrite),
	)

	return cmd
}
