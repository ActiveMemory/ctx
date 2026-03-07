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
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
)

// HandlePromptMd creates or merges PROMPT.md with ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - force: If true, overwrite existing prompt section
//   - autoMerge: If true, skip interactive confirmation
//   - ralph: If true, use ralph mode template
//
// Returns:
//   - error: Non-nil if file operations fail
func HandlePromptMd(cmd *cobra.Command, force, autoMerge, ralph bool) error {
	var templateContent []byte
	var err error
	if ralph {
		templateContent, err = assets.RalphTemplate(config.FilePromptMd)
		if err != nil {
			return ctxerr.ReadInitTemplate("ralph PROMPT.md", err)
		}
	} else {
		templateContent, err = assets.Template(config.FilePromptMd)
		if err != nil {
			return ctxerr.ReadInitTemplate("PROMPT.md", err)
		}
	}
	existingContent, err := os.ReadFile(config.FilePromptMd)
	fileExists := err == nil
	if !fileExists {
		if err := os.WriteFile(config.FilePromptMd, templateContent, config.PermFile); err != nil {
			return ctxerr.FileWrite(config.FilePromptMd, err)
		}
		mode := ""
		if ralph {
			mode = " (ralph mode)"
		}
		write.InitCreatedWith(cmd, config.FilePromptMd, mode)
		return nil
	}
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, config.PromptMarkerStart)
	if hasCtxMarkers {
		if !force {
			write.InitCtxContentExists(cmd, config.FilePromptMd)
			return nil
		}
		return UpdatePromptSection(cmd, existingStr, templateContent)
	}
	if !autoMerge {
		write.InitFileExistsNoCtx(cmd, config.FilePromptMd)
		cmd.Println("Would you like to merge ctx prompt instructions?")
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return ctxerr.ReadInput(err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != config.ConfirmShort && response != config.ConfirmLong {
			write.InitSkippedPlain(cmd, config.FilePromptMd)
			return nil
		}
	}
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FilePromptMd, timestamp)
	if err := os.WriteFile(backupName, existingContent, config.PermFile); err != nil {
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
	if err := os.WriteFile(config.FilePromptMd, []byte(mergedContent), config.PermFile); err != nil {
		return ctxerr.WriteMerged(config.FilePromptMd, err)
	}
	write.InitMerged(cmd, config.FilePromptMd)
	return nil
}

// UpdatePromptSection replaces the existing prompt section between markers with
// new content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - existing: Existing file content
//   - newTemplate: New template content
//
// Returns:
//   - error: Non-nil if markers are missing or file operations fail
func UpdatePromptSection(cmd *cobra.Command, existing string, newTemplate []byte) error {
	startIdx := strings.Index(existing, config.PromptMarkerStart)
	if startIdx == -1 {
		return ctxerr.MarkerNotFound("prompt")
	}
	endIdx := strings.Index(existing, config.PromptMarkerEnd)
	if endIdx == -1 {
		endIdx = len(existing)
	} else {
		endIdx += len(config.PromptMarkerEnd)
	}
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, config.PromptMarkerStart)
	templateEnd := strings.Index(templateStr, config.PromptMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return ctxerr.TemplateMissingMarkers("prompt")
	}
	promptContent := templateStr[templateStart : templateEnd+len(config.PromptMarkerEnd)]
	newContent := existing[:startIdx] + promptContent + existing[endIdx:]
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FilePromptMd, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), config.PermFile); err != nil {
		return ctxerr.CreateBackupGeneric(err)
	}
	write.InitBackup(cmd, backupName)
	if err := os.WriteFile(config.FilePromptMd, []byte(newContent), config.PermFile); err != nil {
		return ctxerr.FileUpdate(config.FilePromptMd, err)
	}
	write.InitUpdatedPromptSection(cmd, config.FilePromptMd)
	return nil
}
