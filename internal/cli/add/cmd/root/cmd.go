//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/add/core"
)

// Cmd returns the "ctx add" command for appending entries to context files.
//
// Supported types are defined in [config.FileType] (both singular and plural
// forms accepted, e.g., "decision" or "decisions"). Content can be provided
// via command argument, --file flag, or stdin pipe.
//
// Flags:
//   - --priority, -p: Priority level for tasks (high, medium, low)
//   - --section, -s: Target section within the file
//   - --file, -f: Read content from a file instead of argument
//   - --context, -c: Context for decisions/learnings (required)
//   - --rationale, -r: Rationale for decisions (required for decisions)
//   - --consequences: Consequences for decisions (required for decisions)
//   - --lesson, -l: Lesson for learnings (required for learnings)
//   - --application, -a: Application for learnings (required for learnings)
//
// Returns:
//   - *cobra.Command: Configured add command with flags registered
func Cmd() *cobra.Command {
	var (
		priority     string
		section      string
		fromFile     string
		context      string
		rationale    string
		consequences string
		lesson       string
		application  string
	)

	short, long := assets.CommandDesc("add")

	cmd := &cobra.Command{
		Use:       "add <type> [content]",
		Short:     short,
		Long:      long,
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{"task", "decision", "learning", "convention"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, core.Config{
				Priority:     priority,
				Section:      section,
				FromFile:     fromFile,
				Context:      context,
				Rationale:    rationale,
				Consequences: consequences,
				Lesson:       lesson,
				Application:  application,
			})
		},
	}

	cmd.Flags().StringVarP(
		&priority,
		"priority", "p", "",
		"Priority level for tasks (high, medium, low)",
	)
	_ = cmd.RegisterFlagCompletionFunc("priority", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"high", "medium", "low"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Flags().StringVarP(
		&section,
		"section", "s", "",
		"Target section within file",
	)
	cmd.Flags().StringVarP(
		&fromFile,
		"file", "f", "",
		"Read content from file instead of argument",
	)
	cmd.Flags().StringVarP(
		&context,
		"context", "c", "",
		"Context for decisions: what prompted this decision (required for decisions)",
	)
	cmd.Flags().StringVarP(
		&rationale,
		"rationale", "r", "",
		"Rationale for decisions: why this choice over alternatives (required for decisions)",
	)
	cmd.Flags().StringVar(
		&consequences,
		"consequences", "",
		"Consequences for decisions: what changes as a result (required for decisions)",
	)
	cmd.Flags().StringVarP(
		&lesson,
		"lesson", "l", "",
		"Lesson for learnings: the key insight (required for learnings)",
	)
	cmd.Flags().StringVarP(
		&application,
		"application", "a", "",
		"Application for learnings: how to apply this going forward (required for learnings)",
	)

	return cmd
}
