//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// HandleClaudeMd creates or merges CLAUDE.md with ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - force: If true, overwrite existing ctx section
//   - autoMerge: If true, skip interactive confirmation
//
// Returns:
//   - error: Non-nil if file operations fail
func HandleClaudeMd(cmd *cobra.Command, force, autoMerge bool) error {
	templateContent, err := assets.ClaudeMd()
	if err != nil {
		return fmt.Errorf("failed to read CLAUDE.md template: %w", err)
	}
	existingContent, err := os.ReadFile(config.FileClaudeMd)
	fileExists := err == nil
	if !fileExists {
		if err := os.WriteFile(config.FileClaudeMd, templateContent, config.PermFile); err != nil {
			return fmt.Errorf("failed to write %s: %w", config.FileClaudeMd, err)
		}
		cmd.Println(fmt.Sprintf("  ✓ %s", config.FileClaudeMd))
		return nil
	}
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, config.CtxMarkerStart)
	if hasCtxMarkers {
		if !force {
			cmd.Println(fmt.Sprintf("  ○ %s (ctx content exists, skipped)\n", config.FileClaudeMd))
			return nil
		}
		return UpdateCtxSection(cmd, existingStr, templateContent)
	}
	if !autoMerge {
		cmd.Println(fmt.Sprintf("\n%s exists but has no ctx content.\n", config.FileClaudeMd))
		cmd.Println("Would you like to append ctx context management instructions?")
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != config.ConfirmShort && response != config.ConfirmLong {
			cmd.Println(fmt.Sprintf("  ○ %s (skipped)", config.FileClaudeMd))
			return nil
		}
	}
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FileClaudeMd, timestamp)
	if err := os.WriteFile(backupName, existingContent, config.PermFile); err != nil {
		return fmt.Errorf("failed to create backup %s: %w", backupName, err)
	}
	cmd.Println(fmt.Sprintf("  ✓ %s (backup)", backupName))
	insertPos := FindInsertionPoint(existingStr)
	var mergedContent string
	if insertPos == 0 {
		mergedContent = string(templateContent) + config.NewlineLF + existingStr
	} else {
		mergedContent = existingStr[:insertPos] + config.NewlineLF + string(templateContent) + config.NewlineLF + existingStr[insertPos:]
	}
	if err := os.WriteFile(config.FileClaudeMd, []byte(mergedContent), config.PermFile); err != nil {
		return fmt.Errorf("failed to write merged %s: %w", config.FileClaudeMd, err)
	}
	cmd.Println(fmt.Sprintf("  ✓ %s (merged)", config.FileClaudeMd))
	return nil
}
