//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// markJournalCmd returns the "ctx system mark-journal" subcommand.
//
// Updates the processing state for a journal entry in
// .context/journal/.state.json. Used by journal skills to record
// pipeline progress (exported → enriched → normalized → fences_verified).
//
// Hidden because it is a plumbing command called by skills, not a
// user-facing workflow.
func markJournalCmd() *cobra.Command {
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

// runMarkJournal handles the mark-journal command.
func runMarkJournal(cmd *cobra.Command, filename, stage string) error {
	journalDir := filepath.Join(rc.ContextDir(), config.DirJournal)

	jstate, err := state.Load(journalDir)
	if err != nil {
		return fmt.Errorf("load journal state: %w", err)
	}

	check, _ := cmd.Flags().GetBool("check")
	if check {
		fs := jstate.Entries[filename]
		var val string
		switch stage {
		case "exported":
			val = fs.Exported
		case "enriched":
			val = fs.Enriched
		case "normalized":
			val = fs.Normalized
		case "fences_verified":
			val = fs.FencesVerified
		case "locked":
			val = fs.Locked
		default:
			return fmt.Errorf("unknown stage %q; valid: %s", stage, strings.Join(state.ValidStages, ", "))
		}
		if val == "" {
			return fmt.Errorf("%s: %s not set", filename, stage)
		}
		cmd.Printf("%s: %s = %s\n", filename, stage, val)
		return nil
	}

	if ok := jstate.Mark(filename, stage); !ok {
		return fmt.Errorf("unknown stage %q; valid: %s", stage, strings.Join(state.ValidStages, ", "))
	}

	if err := jstate.Save(journalDir); err != nil {
		return fmt.Errorf("save journal state: %w", err)
	}

	cmd.Printf("%s: marked %s\n", filename, stage)
	return nil
}
