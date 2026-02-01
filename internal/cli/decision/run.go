//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package decision

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// runReindex regenerates the DECISIONS.md index.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - args: Command arguments (unused)
//
// Returns:
//   - error: Non-nil if file read/write fails
func runReindex(cmd *cobra.Command, _ []string) error {
	filePath := filepath.Join(rc.GetContextDir(), config.FileDecision)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("%s not found. Run 'ctx init' first", config.FileDecision)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	updated := index.UpdateDecisions(string(content))

	if err := os.WriteFile(filePath, []byte(updated), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}

	entries := index.ParseHeaders(string(content))
	green := color.New(color.FgGreen).SprintFunc()
	if len(entries) == 0 {
		cmd.Printf("%s Index cleared (no decisions found)\n", green("✓"))
	} else {
		cmd.Printf(
			"%s Index regenerated with %d entries\n", green("✓"),
			len(entries),
		)
	}

	return nil
}
