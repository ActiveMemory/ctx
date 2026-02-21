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

// handlePromptMd creates or merges PROMPT.md in the project root.
//
// Behavior:
//   - If PROMPT.md doesn't exist: create it from template
//   - If it exists but has no ctx markers: offer to merge
//     (or auto-merge with --merge)
//   - If it exists with ctx markers: update the ctx section only
//     (or skip if not --force)
//
// The ralph parameter selects between interactive (default) and autonomous
// loop templates. Ralph mode templates include completion signals and
// one-task-per-iteration instructions.
//
// Parameters:
//   - cmd: Cobra command for output and input streams
//   - force: If true, overwrite existing ctx content
//   - autoMerge: If true, merge without prompting user
//   - ralph: If true, use autonomous loop template
//
// Returns:
//   - error: Non-nil if template read or file operations fail
func handlePromptMd(cmd *cobra.Command, force, autoMerge, ralph bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Get template content (ralph or default)
	var templateContent []byte
	var err error
	if ralph {
		templateContent, err = assets.RalphTemplate(config.FilePromptMd)
		if err != nil {
			return fmt.Errorf("failed to read ralph PROMPT.md template: %w", err)
		}
	} else {
		templateContent, err = assets.Template(config.FilePromptMd)
		if err != nil {
			return fmt.Errorf("failed to read PROMPT.md template: %w", err)
		}
	}

	// Check if PROMPT.md exists
	existingContent, err := os.ReadFile(config.FilePromptMd)
	fileExists := err == nil

	if !fileExists {
		// File doesn't exist - create it
		if err := os.WriteFile(
			config.FilePromptMd, templateContent, config.PermFile,
		); err != nil {
			return fmt.Errorf("failed to write %s: %w", config.FilePromptMd, err)
		}
		mode := ""
		if ralph {
			mode = " (ralph mode)"
		}
		cmd.Printf("  %s %s%s\n", green("✓"), config.FilePromptMd, mode)
		return nil
	}

	// File exists - check for ctx markers
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, config.PromptMarkerStart)

	if hasCtxMarkers {
		// Already has ctx content
		if !force {
			cmd.Printf(
				"  %s %s (ctx content exists, skipped)\n", yellow("○"),
				config.FilePromptMd,
			)
			return nil
		}
		// Force update: replace the existing ctx section
		return updatePromptSection(cmd, existingStr, templateContent)
	}

	// No ctx markers: need to merge
	if !autoMerge {
		// Prompt user
		cmd.Printf(
			"\n%s exists but has no ctx content.\n", config.FilePromptMd,
		)
		cmd.Println(
			"Would you like to merge ctx prompt instructions?",
		)
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" { //nolint:goconst // trivial user input check
			cmd.Printf("  %s %s (skipped)\n", yellow("○"), config.FilePromptMd)
			return nil
		}
	}

	// Back up existing file
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FilePromptMd, timestamp)
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
		config.FilePromptMd, []byte(mergedContent), config.PermFile); err != nil {
		return fmt.Errorf(
			"failed to write merged %s: %w", config.FilePromptMd, err)
	}
	cmd.Printf("  %s %s (merged)\n", green("✓"), config.FilePromptMd)

	return nil
}

// updatePromptSection replaces the existing prompt section between markers.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - existing: Current file content containing prompt markers
//   - newTemplate: Template content with updated prompt section
//
// Returns:
//   - error: Non-nil if the markers are not found or file operations fail
func updatePromptSection(
	cmd *cobra.Command, existing string, newTemplate []byte,
) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Find the start marker
	startIdx := strings.Index(existing, config.PromptMarkerStart)
	if startIdx == -1 {
		return fmt.Errorf("prompt start marker not found")
	}

	// Find the end marker
	endIdx := strings.Index(existing, config.PromptMarkerEnd)
	if endIdx == -1 {
		// No end marker - append from start marker to end
		endIdx = len(existing)
	} else {
		endIdx += len(config.PromptMarkerEnd)
	}

	// Extract the prompt content from the template (between markers)
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, config.PromptMarkerStart)
	templateEnd := strings.Index(templateStr, config.PromptMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return fmt.Errorf("template missing prompt markers")
	}
	promptContent := templateStr[templateStart : templateEnd+
		len(config.PromptMarkerEnd)]

	// Build new content: before prompt + new prompt content + after prompt
	newContent := existing[:startIdx] + promptContent + existing[endIdx:]

	// Back up before updating
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FilePromptMd, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), config.PermFile); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	cmd.Printf("  %s %s (backup)\n", green("✓"), backupName)

	if err := os.WriteFile(
		config.FilePromptMd, []byte(newContent), config.PermFile,
	); err != nil {
		return fmt.Errorf("failed to update %s: %w", config.FilePromptMd, err)
	}
	cmd.Printf(
		"  %s %s (updated prompt section)\n", green("✓"), config.FilePromptMd,
	)

	return nil
}
