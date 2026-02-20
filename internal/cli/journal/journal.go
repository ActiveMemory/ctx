//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the journal command with subcommands.
//
// The journal system provides LLM-powered analysis and synthesis of
// exported session files in .context/journal/.
//
// Returns:
//   - *cobra.Command: The journal command with subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "journal",
		Short: "Analyze and synthesize exported sessions",
		Long: `Work with exported session files in .context/journal/.

The journal system provides tools for analyzing, enriching, and
publishing your AI session history.

Subcommands:
  site      Generate a static site from journal entries
  obsidian  Generate an Obsidian vault from journal entries
  mark      Update processing state for a journal entry

Examples:
  ctx journal site                    # Generate site in .context/journal-site/
  ctx journal site --output ~/public  # Custom output directory
  ctx journal site --serve            # Generate and serve locally
  ctx journal obsidian                # Generate Obsidian vault
  ctx journal mark session.md enriched`,
	}

	cmd.AddCommand(journalSiteCmd())
	cmd.AddCommand(journalObsidianCmd())
	cmd.AddCommand(journalMarkCmd())

	return cmd
}

// journalMarkCmd returns the "ctx journal mark" subcommand.
func journalMarkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mark <filename> <stage>",
		Short: "Update journal processing state",
		Long: fmt.Sprintf(`Mark a journal entry as having completed a processing stage.

Valid stages: %s

The state is recorded in .context/journal/.state.json with today's date.

Examples:
  ctx journal mark 2026-01-21-session-abc12345.md exported
  ctx journal mark 2026-01-21-session-abc12345.md enriched
  ctx journal mark 2026-01-21-session-abc12345.md normalized
  ctx journal mark 2026-01-21-session-abc12345.md fences_verified`, strings.Join(state.ValidStages, ", ")),
		Args: cobra.ExactArgs(2), //nolint:mnd // 2 positional args: filename, stage
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJournalMark(cmd, args[0], args[1])
		},
	}

	cmd.Flags().Bool("check", false, "Check if stage is set (exit 1 if not)")

	return cmd
}

// runJournalMark handles the journal mark command.
func runJournalMark(cmd *cobra.Command, filename, stage string) error {
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
