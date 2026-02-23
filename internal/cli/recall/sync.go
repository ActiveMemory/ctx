//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// recallSyncCmd returns the "ctx recall sync" subcommand.
//
// Scans journal markdowns and syncs their frontmatter lock state into
// .state.json. This is the inverse of "ctx recall lock": the frontmatter
// is treated as the source of truth, and state is updated to match.
//
// Returns:
//   - *cobra.Command: Command for syncing lock state from frontmatter
func recallSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync lock state from journal frontmatter to state file",
		Long: `Scan journal markdowns and sync their lock state to .state.json.

This is the sister command to "ctx recall lock". Instead of marking files
locked in state and updating frontmatter, it reads "locked: true" from
each file's YAML frontmatter and updates .state.json to match.

Typical workflow:
  1. Enrich journal entries (add "locked: true" to frontmatter)
  2. Run "ctx recall sync" to propagate lock state to .state.json

Files with "locked: true" in frontmatter will be marked locked in state.
Files without a "locked:" line (or with "locked: false") will have their
lock cleared if one exists in state.

Examples:
  ctx recall sync`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSync(cmd)
		},
	}

	return cmd
}

// runSync scans all journal markdowns and syncs frontmatter lock state
// to .state.json.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil on I/O failure
func runSync(cmd *cobra.Command) error {
	journalDir := filepath.Join(rc.ContextDir(), config.DirJournal)

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return fmt.Errorf("load journal state: %w", loadErr)
	}

	files, matchErr := matchJournalFiles(journalDir, nil, true)
	if matchErr != nil {
		return matchErr
	}
	if len(files) == 0 {
		cmd.Println("No journal entries found.")
		return nil
	}

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	locked, unlocked := 0, 0

	for _, filename := range files {
		path := filepath.Join(journalDir, filename)
		fmLocked := frontmatterHasLocked(path)
		stateLocked := jstate.Locked(filename)

		switch {
		case fmLocked && !stateLocked:
			jstate.Mark(filename, "locked")
			cmd.Printf("  %s %s (locked)\n", green("✓"), filename)
			locked++
		case !fmLocked && stateLocked:
			jstate.Clear(filename, "locked")
			cmd.Printf("  %s %s (unlocked)\n", yellow("✓"), filename)
			unlocked++
		}
	}

	if saveErr := jstate.Save(journalDir); saveErr != nil {
		return fmt.Errorf("save journal state: %w", saveErr)
	}

	if locked == 0 && unlocked == 0 {
		cmd.Println("No changes — state already matches frontmatter.")
	} else {
		if locked > 0 {
			cmd.Printf("\nLocked %d entry(s).\n", locked)
		}
		if unlocked > 0 {
			cmd.Printf("\nUnlocked %d entry(s).\n", unlocked)
		}
	}

	return nil
}

// frontmatterHasLocked reads a journal file and returns true if its
// YAML frontmatter contains a "locked:" line with a truthy value.
//
// Parameters:
//   - path: Absolute path to the journal .md file
//
// Returns:
//   - bool: True if frontmatter contains "locked: true"
func frontmatterHasLocked(path string) bool {
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return false
	}
	content := string(data)

	nl := config.NewlineLF
	fmOpen := config.Separator + nl

	if !strings.HasPrefix(content, fmOpen) {
		return false
	}

	closeIdx := strings.Index(content[len(fmOpen):], nl+config.Separator+nl)
	if closeIdx < 0 {
		return false
	}

	fmBlock := content[len(fmOpen) : len(fmOpen)+closeIdx]

	for _, line := range strings.Split(fmBlock, nl) {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "locked:") {
			continue
		}
		val := strings.TrimSpace(strings.TrimPrefix(trimmed, "locked:"))
		// Strip inline comment (e.g. "true  # managed by ctx").
		if idx := strings.Index(val, "#"); idx >= 0 {
			val = strings.TrimSpace(val[:idx])
		}
		return val == "true"
	}

	return false
}
