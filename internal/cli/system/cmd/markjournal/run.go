//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package markjournal

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// runMarkJournal handles the mark-journal command.
func runMarkJournal(cmd *cobra.Command, filename, stage string) error {
	journalDir := filepath.Join(rc.ContextDir(), config.DirJournal)

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return fmt.Errorf("load journal state: %w", loadErr)
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
		cmd.Println(fmt.Sprintf("%s: %s = %s", filename, stage, val))
		return nil
	}

	if ok := jstate.Mark(filename, stage); !ok {
		return fmt.Errorf("unknown stage %q; valid: %s", stage, strings.Join(state.ValidStages, ", "))
	}

	if saveErr := jstate.Save(journalDir); saveErr != nil {
		return fmt.Errorf("save journal state: %w", saveErr)
	}

	cmd.Println(fmt.Sprintf("%s: marked %s", filename, stage))
	return nil
}
