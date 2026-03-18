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

	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/loop"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/err/backup"
	fsErr "github.com/ActiveMemory/ctx/internal/err/fs"
	initErr "github.com/ActiveMemory/ctx/internal/err/initialize"
	promptErr "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/write/initialize"

	"github.com/ActiveMemory/ctx/internal/assets"
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
			return initErr.ReadTemplate("ralph PROMPT.md", err)
		}
	} else {
		templateContent, err = assets.Template(loop.PromptMd)
		if err != nil {
			return initErr.ReadTemplate("PROMPT.md", err)
		}
	}
	existingContent, err := os.ReadFile(loop.PromptMd)
	fileExists := err == nil
	if !fileExists {
		if err := os.WriteFile(loop.PromptMd, templateContent, fs.PermFile); err != nil {
			return fsErr.FileWrite(loop.PromptMd, err)
		}
		mode := ""
		if ralph {
			mode = " (ralph mode)"
		}
		initialize.CreatedWith(cmd, loop.PromptMd, mode)
		return nil
	}
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, marker.PromptMarkerStart)
	if hasCtxMarkers {
		if !force {
			initialize.CtxContentExists(cmd, loop.PromptMd)
			return nil
		}
		return UpdatePromptSection(cmd, existingStr, templateContent)
	}
	if !autoMerge {
		initialize.FileExistsNoCtx(cmd, loop.PromptMd)
		cmd.Println("Would you like to merge ctx prompt instructions?")
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fsErr.ReadInput(err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != cli.ConfirmShort && response != cli.ConfirmLong {
			initialize.SkippedPlain(cmd, loop.PromptMd)
			return nil
		}
	}
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", loop.PromptMd, timestamp)
	if err := os.WriteFile(backupName, existingContent, fs.PermFile); err != nil {
		return backup.Create(backupName, err)
	}
	initialize.Backup(cmd, backupName)
	insertPos := FindInsertionPoint(existingStr)
	var mergedContent string
	if insertPos == 0 {
		mergedContent = string(templateContent) + token.NewlineLF + existingStr
	} else {
		mergedContent = existingStr[:insertPos] + token.NewlineLF + string(templateContent) + token.NewlineLF + existingStr[insertPos:]
	}
	if err := os.WriteFile(loop.PromptMd, []byte(mergedContent), fs.PermFile); err != nil {
		return fsErr.WriteMerged(loop.PromptMd, err)
	}
	initialize.Merged(cmd, loop.PromptMd)
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
		return promptErr.MarkerNotFound("prompt")
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
		return promptErr.TemplateMissingMarkers("prompt")
	}
	promptContent := templateStr[templateStart : templateEnd+len(marker.PromptMarkerEnd)]
	newContent := existing[:startIdx] + promptContent + existing[endIdx:]
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", loop.PromptMd, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), fs.PermFile); err != nil {
		return backup.CreateGeneric(err)
	}
	initialize.Backup(cmd, backupName)
	if err := os.WriteFile(loop.PromptMd, []byte(newContent), fs.PermFile); err != nil {
		return fsErr.FileUpdate(loop.PromptMd, err)
	}
	initialize.UpdatedPromptSection(cmd, loop.PromptMd)
	return nil
}
