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

	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/project"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
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
	templateContent, err := assets.ProjectFile(project.ImplementationPlan)
	if err != nil {
		return ctxerr.ReadInitTemplate("IMPLEMENTATION_PLAN.md", err)
	}
	existingContent, err := os.ReadFile(project.ImplementationPlan)
	fileExists := err == nil
	if !fileExists {
		if err := os.WriteFile(project.ImplementationPlan, templateContent, fs.PermFile); err != nil {
			return ctxerr.FileWrite(project.ImplementationPlan, err)
		}
		write.InitCreated(cmd, project.ImplementationPlan)
		return nil
	}
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, marker.PlanMarkerStart)
	if hasCtxMarkers {
		if !force {
			write.InitCtxContentExists(cmd, project.ImplementationPlan)
			return nil
		}
		return UpdatePlanSection(cmd, existingStr, templateContent)
	}
	if !autoMerge {
		write.InitFileExistsNoCtx(cmd, project.ImplementationPlan)
		cmd.Println("Would you like to merge ctx implementation plan template?")
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return ctxerr.ReadInput(err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != cli.ConfirmShort && response != cli.ConfirmLong {
			write.InitSkippedPlain(cmd, project.ImplementationPlan)
			return nil
		}
	}
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", project.ImplementationPlan, timestamp)
	if err := os.WriteFile(backupName, existingContent, fs.PermFile); err != nil {
		return ctxerr.CreateBackup(backupName, err)
	}
	write.InitBackup(cmd, backupName)
	insertPos := FindInsertionPoint(existingStr)
	var mergedContent string
	if insertPos == 0 {
		mergedContent = string(templateContent) + token.NewlineLF + existingStr
	} else {
		mergedContent = existingStr[:insertPos] + token.NewlineLF + string(templateContent) + token.NewlineLF + existingStr[insertPos:]
	}
	if err := os.WriteFile(project.ImplementationPlan, []byte(mergedContent), fs.PermFile); err != nil {
		return ctxerr.WriteMerged(project.ImplementationPlan, err)
	}
	write.InitMerged(cmd, project.ImplementationPlan)
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
	startIdx := strings.Index(existing, marker.PlanMarkerStart)
	if startIdx == -1 {
		return ctxerr.MarkerNotFound("plan")
	}
	endIdx := strings.Index(existing, marker.PlanMarkerEnd)
	if endIdx == -1 {
		endIdx = len(existing)
	} else {
		endIdx += len(marker.PlanMarkerEnd)
	}
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, marker.PlanMarkerStart)
	templateEnd := strings.Index(templateStr, marker.PlanMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return ctxerr.TemplateMissingMarkers("plan")
	}
	planContent := templateStr[templateStart : templateEnd+len(marker.PlanMarkerEnd)]
	newContent := existing[:startIdx] + planContent + existing[endIdx:]
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", project.ImplementationPlan, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), fs.PermFile); err != nil {
		return ctxerr.CreateBackupGeneric(err)
	}
	write.InitBackup(cmd, backupName)
	if err := os.WriteFile(project.ImplementationPlan, []byte(newContent), fs.PermFile); err != nil {
		return ctxerr.FileUpdate(project.ImplementationPlan, err)
	}
	write.InitUpdatedPlanSection(cmd, project.ImplementationPlan)
	return nil
}
