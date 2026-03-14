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
	"github.com/ActiveMemory/ctx/internal/config/loop"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
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
		templateContent, err = assets.RalphTemplate(loop.PromptMd)
		if err != nil {
			return ctxerr.ReadInitTemplate("ralph PROMPT.md", err)
		}
	} else {
		templateContent, err = assets.Template(loop.PromptMd)
		if err != nil {
			return ctxerr.ReadInitTemplate("PROMPT.md", err)
		}
	}
	existingContent, err := os.ReadFile(loop.PromptMd)
	fileExists := err == nil
	if !fileExists {
		if err := os.WriteFile(loop.PromptMd, templateContent, fs.PermFile); err != nil {
			return ctxerr.FileWrite(loop.PromptMd, err)
		}
		mode := ""
		if ralph {
			mode = " (ralph mode)"
		}
		write.InitCreatedWith(cmd, loop.PromptMd, mode)
		return nil
	}
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, marker.PromptMarkerStart)
	if hasCtxMarkers {
		if !force {
			write.InitCtxContentExists(cmd, loop.PromptMd)
			return nil
		}
		return UpdatePromptSection(cmd, existingStr, templateContent)
	}
	if !autoMerge {
		write.InitFileExistsNoCtx(cmd, loop.PromptMd)
		cmd.Println("Would you like to merge ctx prompt instructions?")
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return ctxerr.ReadInput(err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != cli.ConfirmShort && response != cli.ConfirmLong {
			write.InitSkippedPlain(cmd, loop.PromptMd)
			return nil
		}
	}
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", loop.PromptMd, timestamp)
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
	if err := os.WriteFile(loop.PromptMd, []byte(mergedContent), fs.PermFile); err != nil {
		return ctxerr.WriteMerged(loop.PromptMd, err)
	}
	write.InitMerged(cmd, loop.PromptMd)
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
	startIdx := strings.Index(existing, marker.PromptMarkerStart)
	if startIdx == -1 {
		return ctxerr.MarkerNotFound("prompt")
	}
	endIdx := strings.Index(existing, marker.PromptMarkerEnd)
	if endIdx == -1 {
		endIdx = len(existing)
	} else {
		endIdx += len(marker.PromptMarkerEnd)
	}
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, marker.PromptMarkerStart)
	templateEnd := strings.Index(templateStr, marker.PromptMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return ctxerr.TemplateMissingMarkers("prompt")
	}
	promptContent := templateStr[templateStart : templateEnd+len(marker.PromptMarkerEnd)]
	newContent := existing[:startIdx] + promptContent + existing[endIdx:]
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", loop.PromptMd, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), fs.PermFile); err != nil {
		return ctxerr.CreateBackupGeneric(err)
	}
	write.InitBackup(cmd, backupName)
	if err := os.WriteFile(loop.PromptMd, []byte(newContent), fs.PermFile); err != nil {
		return ctxerr.FileUpdate(loop.PromptMd, err)
	}
	write.InitUpdatedPromptSection(cmd, loop.PromptMd)
	return nil
}
