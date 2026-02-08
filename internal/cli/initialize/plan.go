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

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/tpl"
)

// handleImplementationPlan creates or merges IMPLEMENTATION_PLAN.md in the
// project root.
//
// Behavior:
//   - If IMPLEMENTATION_PLAN.md doesn't exist: create it from template
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
func handleImplementationPlan(cmd *cobra.Command, force, autoMerge bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Get template content
	templateContent, err := tpl.Template(config.FileImplementationPlan)
	if err != nil {
		return fmt.Errorf(
			"failed to read IMPLEMENTATION_PLAN.md template: %w", err)
	}

	// Check if file exists
	existingContent, err := os.ReadFile(config.FileImplementationPlan)
	fileExists := err == nil

	if !fileExists {
		// File doesn't exist - create it
		if err := os.WriteFile(
			config.FileImplementationPlan, templateContent, config.PermFile,
		); err != nil {
			return fmt.Errorf(
				"failed to write %s: %w", config.FileImplementationPlan, err)
		}
		cmd.Printf("  %s %s\n", green("✓"), config.FileImplementationPlan)
		return nil
	}

	// File exists - check for ctx markers
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, config.PlanMarkerStart)

	if hasCtxMarkers {
		// Already has ctx content
		if !force {
			cmd.Printf(
				"  %s %s (ctx content exists, skipped)\n", yellow("○"),
				config.FileImplementationPlan,
			)
			return nil
		}
		// Force update: replace the existing ctx section
		return updatePlanSection(cmd, existingStr, templateContent)
	}

	// No ctx markers: need to merge
	if !autoMerge {
		// Prompt user
		cmd.Printf(
			"\n%s exists but has no ctx content.\n",
			config.FileImplementationPlan,
		)
		cmd.Println(
			"Would you like to merge ctx implementation plan template?",
		)
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" { //nolint:goconst // trivial user input check
			cmd.Printf(
				"  %s %s (skipped)\n", yellow("○"), config.FileImplementationPlan)
			return nil
		}
	}

	// Back up existing file
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf(
		"%s.%d.bak", config.FileImplementationPlan, timestamp)
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
		config.FileImplementationPlan, []byte(mergedContent), config.PermFile); err != nil {
		return fmt.Errorf(
			"failed to write merged %s: %w", config.FileImplementationPlan, err)
	}
	cmd.Printf("  %s %s (merged)\n", green("✓"), config.FileImplementationPlan)

	return nil
}

// updatePlanSection replaces the existing plan section between markers.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - existing: Current file content containing plan markers
//   - newTemplate: Template content with updated plan section
//
// Returns:
//   - error: Non-nil if the markers are not found or file operations fail
func updatePlanSection(
	cmd *cobra.Command, existing string, newTemplate []byte,
) error {
	green := color.New(color.FgGreen).SprintFunc()

	// Find the start marker
	startIdx := strings.Index(existing, config.PlanMarkerStart)
	if startIdx == -1 {
		return fmt.Errorf("plan start marker not found")
	}

	// Find the end marker
	endIdx := strings.Index(existing, config.PlanMarkerEnd)
	if endIdx == -1 {
		// No end marker - append from start marker to end
		endIdx = len(existing)
	} else {
		endIdx += len(config.PlanMarkerEnd)
	}

	// Extract the plan content from the template (between markers)
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, config.PlanMarkerStart)
	templateEnd := strings.Index(templateStr, config.PlanMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return fmt.Errorf("template missing plan markers")
	}
	planContent := templateStr[templateStart : templateEnd+
		len(config.PlanMarkerEnd)]

	// Build new content: before plan + new plan content + after plan
	newContent := existing[:startIdx] + planContent + existing[endIdx:]

	// Back up before updating
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf(
		"%s.%d.bak", config.FileImplementationPlan, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), config.PermFile); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	cmd.Printf("  %s %s (backup)\n", green("✓"), backupName)

	if err := os.WriteFile(
		config.FileImplementationPlan, []byte(newContent), config.PermFile,
	); err != nil {
		return fmt.Errorf(
			"failed to update %s: %w", config.FileImplementationPlan, err)
	}
	cmd.Printf(
		"  %s %s (updated plan section)\n",
		green("✓"), config.FileImplementationPlan,
	)

	return nil
}
