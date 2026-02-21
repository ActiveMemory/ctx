//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// handleClaudeMd creates or merges CLAUDE.md in the project root.
//
// Behavior:
//   - If CLAUDE.md doesn't exist: create it from template
//   - If it exists but has no ctx markers: offer to merge
//     (or auto-merge with --merge)
//   - If it exists with ctx markers: update the ctx section only
//     (or skip if not --force)
//
// Parameters:
//   - cmd: Cobra command for output and input streams
//   - force: If true, overwrite existing ctx content
//   - autoMerge: If true, merge without prompting user
//
// Returns:
//   - error: Non-nil if template read or file operations fail
func handleClaudeMd(cmd *cobra.Command, force, autoMerge bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Get template content
	templateContent, err := assets.Template(config.FileClaudeMd)
	if err != nil {
		return fmt.Errorf("failed to read CLAUDE.md template: %w", err)
	}

	// Check if CLAUDE.md exists
	existingContent, err := os.ReadFile(config.FileClaudeMd)
	fileExists := err == nil

	if !fileExists {
		// File doesn't exist - create it
		if err := os.WriteFile(
			config.FileClaudeMd, templateContent, config.PermFile,
		); err != nil {
			return fmt.Errorf("failed to write %s: %w", config.FileClaudeMd, err)
		}
		cmd.Printf("  %s %s\n", green("✓"), config.FileClaudeMd)
		return nil
	}

	// File exists - check for ctx markers
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, config.CtxMarkerStart)

	if hasCtxMarkers {
		// Already has ctx content
		if !force {
			cmd.Printf(
				"  %s %s (ctx content exists, skipped)\n", yellow("○"),
				config.FileClaudeMd,
			)
			return nil
		}
		// Force update: replace the existing ctx section
		return updateCtxSection(cmd, existingStr, templateContent)
	}

	// No ctx markers: need to merge
	if !autoMerge {
		// Prompt user
		cmd.Printf(
			"\n%s exists but has no ctx content.\n", config.FileClaudeMd,
		)
		cmd.Println(
			"Would you like to append ctx context management instructions?",
		)
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" { //nolint:goconst // trivial user input check
			cmd.Printf("  %s %s (skipped)\n", yellow("○"), config.FileClaudeMd)
			return nil
		}
	}

	// Back up existing file
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FileClaudeMd, timestamp)
	if err := os.WriteFile(backupName, existingContent, config.PermFile); err != nil {
		return fmt.Errorf("failed to create backup %s: %w", backupName, err)
	}
	cmd.Printf("  %s %s (backup)\n", green("✓"), backupName)

	// Find the best insertion point (after the H1 title, or at the top)
	insertPos := findInsertionPoint(existingStr)

	// Build merged content: before + ctx content + after
	var mergedContent string
	if insertPos == 0 {
		// Insert at top
		mergedContent = string(templateContent) + "\n" + existingStr
	} else {
		// Insert after H1 heading
		mergedContent = existingStr[:insertPos] + "\n" +
			string(templateContent) + "\n" + existingStr[insertPos:]
	}

	if err := os.WriteFile(
		config.FileClaudeMd, []byte(mergedContent), config.PermFile); err != nil {
		return fmt.Errorf(
			"failed to write merged %s: %w", config.FileClaudeMd, err)
	}
	cmd.Printf("  %s %s (merged)\n", green("✓"), config.FileClaudeMd)

	return nil
}
