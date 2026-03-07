//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package markjournal

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/journal/state"
)

// Cmd returns the "ctx system mark-journal" subcommand.
//
// Returns:
//   - *cobra.Command: Configured mark-journal subcommand
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mark-journal <filename> <stage>",
		Short: "Update journal processing state",
		Long: fmt.Sprintf(`Mark a journal entry as having completed a processing stage.

Valid stages: %s

The state is recorded in .context/journal/.state.json with today's date.

Examples:
  ctx system mark-journal 2026-01-21-session-abc12345.md exported
  ctx system mark-journal 2026-01-21-session-abc12345.md enriched
  ctx system mark-journal 2026-01-21-session-abc12345.md normalized
  ctx system mark-journal 2026-01-21-session-abc12345.md fences_verified`, strings.Join(state.ValidStages, ", ")),
		Hidden: true,
		Args:   cobra.ExactArgs(2), //nolint:mnd // 2 positional args: filename, stage
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMarkJournal(cmd, args[0], args[1])
		},
	}

	cmd.Flags().Bool("check", false, "Check if stage is set (exit 1 if not)")

	return cmd
}
