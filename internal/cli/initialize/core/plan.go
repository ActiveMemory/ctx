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

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
)

// HandleImplementationPlan creates or merges IMPLEMENTATION_PLAN.md.
//
// Parameters:
//   - cmd: Cobra command for output
//   - force: If true, overwrite existing plan section
//   - autoMerge: If true, skip interactive confirmation
//
// Returns:
//   - error: Non-nil if file operations fail
func HandleImplementationPlan(cmd *cobra.Command, force, autoMerge bool) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	templateContent, err := assets.ProjectFile(config.FileImplementationPlan)
	if err != nil {
		return fmt.Errorf("failed to read IMPLEMENTATION_PLAN.md template: %w", err)
	}
	existingContent, err := os.ReadFile(config.FileImplementationPlan)
	fileExists := err == nil
	if !fileExists {
		if err := os.WriteFile(config.FileImplementationPlan, templateContent, config.PermFile); err != nil {
			return fmt.Errorf("failed to write %s: %w", config.FileImplementationPlan, err)
		}
		cmd.Println(fmt.Sprintf("  %s %s", green("✓"), config.FileImplementationPlan))
		return nil
	}
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, config.PlanMarkerStart)
	if hasCtxMarkers {
		if !force {
			cmd.Println(fmt.Sprintf("  %s %s (ctx content exists, skipped)\n", yellow("○"), config.FileImplementationPlan))
			return nil
		}
		return UpdatePlanSection(cmd, existingStr, templateContent)
	}
	if !autoMerge {
		cmd.Println(fmt.Sprintf("\n%s exists but has no ctx content.\n", config.FileImplementationPlan))
		cmd.Println("Would you like to merge ctx implementation plan template?")
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != config.ConfirmShort && response != config.ConfirmLong {
			cmd.Println(fmt.Sprintf("  %s %s (skipped)\n", yellow("○"), config.FileImplementationPlan))
			return nil
		}
	}
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FileImplementationPlan, timestamp)
	if err := os.WriteFile(backupName, existingContent, config.PermFile); err != nil {
		return fmt.Errorf("failed to create backup %s: %w", backupName, err)
	}
	cmd.Println(fmt.Sprintf("  %s %s (backup)", green("✓"), backupName))
	insertPos := FindInsertionPoint(existingStr)
	var mergedContent string
	if insertPos == 0 {
		mergedContent = string(templateContent) + config.NewlineLF + existingStr
	} else {
		mergedContent = existingStr[:insertPos] + config.NewlineLF + string(templateContent) + config.NewlineLF + existingStr[insertPos:]
	}
	if err := os.WriteFile(config.FileImplementationPlan, []byte(mergedContent), config.PermFile); err != nil {
		return fmt.Errorf("failed to write merged %s: %w", config.FileImplementationPlan, err)
	}
	cmd.Println(fmt.Sprintf("  %s %s (merged)", green("✓"), config.FileImplementationPlan))
	return nil
}

// UpdatePlanSection replaces the existing plan section between markers with
// new content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - existing: Existing file content
//   - newTemplate: New template content
//
// Returns:
//   - error: Non-nil if markers are missing or file operations fail
func UpdatePlanSection(cmd *cobra.Command, existing string, newTemplate []byte) error {
	green := color.New(color.FgGreen).SprintFunc()
	startIdx := strings.Index(existing, config.PlanMarkerStart)
	if startIdx == -1 {
		return fmt.Errorf("plan start marker not found")
	}
	endIdx := strings.Index(existing, config.PlanMarkerEnd)
	if endIdx == -1 {
		endIdx = len(existing)
	} else {
		endIdx += len(config.PlanMarkerEnd)
	}
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, config.PlanMarkerStart)
	templateEnd := strings.Index(templateStr, config.PlanMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return fmt.Errorf("template missing plan markers")
	}
	planContent := templateStr[templateStart : templateEnd+len(config.PlanMarkerEnd)]
	newContent := existing[:startIdx] + planContent + existing[endIdx:]
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FileImplementationPlan, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), config.PermFile); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	cmd.Println(fmt.Sprintf("  %s %s (backup)", green("✓"), backupName))
	if err := os.WriteFile(config.FileImplementationPlan, []byte(newContent), config.PermFile); err != nil {
		return fmt.Errorf("failed to update %s: %w", config.FileImplementationPlan, err)
	}
	cmd.Println(fmt.Sprintf("  %s %s (updated plan section)\n", green("✓"), config.FileImplementationPlan))
	return nil
}
