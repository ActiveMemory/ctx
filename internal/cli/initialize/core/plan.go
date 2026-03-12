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

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
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
	templateContent, err := assets.ProjectFile(file.FileImplementationPlan)
	if err != nil {
		return ctxerr.ReadInitTemplate("IMPLEMENTATION_PLAN.md", err)
	}
	existingContent, err := os.ReadFile(file.FileImplementationPlan)
	fileExists := err == nil
	if !fileExists {
		if err := os.WriteFile(file.FileImplementationPlan, templateContent, fs.PermFile); err != nil {
			return ctxerr.FileWrite(file.FileImplementationPlan, err)
		}
		write.InitCreated(cmd, file.FileImplementationPlan)
		return nil
	}
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, config.PlanMarkerStart)
	if hasCtxMarkers {
		if !force {
			write.InitCtxContentExists(cmd, file.FileImplementationPlan)
			return nil
		}
		return UpdatePlanSection(cmd, existingStr, templateContent)
	}
	if !autoMerge {
		write.InitFileExistsNoCtx(cmd, file.FileImplementationPlan)
		cmd.Println("Would you like to merge ctx implementation plan template?")
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return ctxerr.ReadInput(err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != file.ConfirmShort && response != file.ConfirmLong {
			write.InitSkippedPlain(cmd, file.FileImplementationPlan)
			return nil
		}
	}
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", file.FileImplementationPlan, timestamp)
	if err := os.WriteFile(backupName, existingContent, fs.PermFile); err != nil {
		return ctxerr.CreateBackup(backupName, err)
	}
	write.InitBackup(cmd, backupName)
	insertPos := FindInsertionPoint(existingStr)
	var mergedContent string
	if insertPos == 0 {
		mergedContent = string(templateContent) + config.NewlineLF + existingStr
	} else {
		mergedContent = existingStr[:insertPos] + config.NewlineLF + string(templateContent) + config.NewlineLF + existingStr[insertPos:]
	}
	if err := os.WriteFile(file.FileImplementationPlan, []byte(mergedContent), fs.PermFile); err != nil {
		return ctxerr.WriteMerged(file.FileImplementationPlan, err)
	}
	write.InitMerged(cmd, file.FileImplementationPlan)
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
	startIdx := strings.Index(existing, config.PlanMarkerStart)
	if startIdx == -1 {
		return ctxerr.MarkerNotFound("plan")
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
		return ctxerr.TemplateMissingMarkers("plan")
	}
	planContent := templateStr[templateStart : templateEnd+len(config.PlanMarkerEnd)]
	newContent := existing[:startIdx] + planContent + existing[endIdx:]
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", file.FileImplementationPlan, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), fs.PermFile); err != nil {
		return ctxerr.CreateBackupGeneric(err)
	}
	write.InitBackup(cmd, backupName)
	if err := os.WriteFile(file.FileImplementationPlan, []byte(newContent), fs.PermFile); err != nil {
		return ctxerr.FileUpdate(file.FileImplementationPlan, err)
	}
	write.InitUpdatedPlanSection(cmd, file.FileImplementationPlan)
	return nil
}
