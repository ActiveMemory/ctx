//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/initialize/core"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	backupPkg "github.com/ActiveMemory/ctx/internal/err/backup"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	promptErr "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// backupFile creates a timestamped .bak copy and reports it.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: Original file path
//   - content: Content to back up
//
// Returns:
//   - error: Non-nil if the backup write fails
func backupFile(cmd *cobra.Command, filename string, content []byte) error {
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf(file.BackupFormat, filename, timestamp)
	if writeErr := os.WriteFile(backupName, content, fs.PermFile); writeErr != nil {
		return backupPkg.Create(backupName, writeErr)
	}
	initialize.Backup(cmd, backupName)
	return nil
}

// CreateOrMerge handles the common pattern of creating a new file or
// merging ctx content into an existing one.
//
// Parameters:
//   - cmd: Cobra command for output and input
//   - p: Merge parameters
//
// Returns:
//   - created: True if the file was created fresh (no existing file)
//   - error: Non-nil if file operations fail
func CreateOrMerge(cmd *cobra.Command, p core.MergeParams) (bool, error) {
	existingContent, readErr := os.ReadFile(p.Filename)
	fileExists := readErr == nil

	if !fileExists {
		if writeErr := os.WriteFile(
			p.Filename, p.TemplateContent, fs.PermFile,
		); writeErr != nil {
			return false, errFs.FileWrite(p.Filename, writeErr)
		}
		return true, nil
	}

	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, p.MarkerStart)

	if hasCtxMarkers {
		if !p.Force {
			initialize.CtxContentExists(cmd, p.Filename)
			return false, nil
		}
		return false, p.UpdateFn(cmd, existingStr, p.TemplateContent)
	}

	if !p.AutoMerge {
		initialize.FileExistsNoCtx(cmd, p.Filename)
		initialize.MergePrompt(cmd, p.ConfirmPrompt)
		reader := bufio.NewReader(os.Stdin)
		response, inputErr := reader.ReadString(token.NewlineLF[0])
		if inputErr != nil {
			return false, errFs.ReadInput(inputErr)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != cli.ConfirmShort && response != cli.ConfirmLong {
			initialize.SkippedPlain(cmd, p.Filename)
			return false, nil
		}
	}

	if bkErr := backupFile(cmd, p.Filename, existingContent); bkErr != nil {
		return false, bkErr
	}

	insertPos := FindInsertionPoint(existingStr)
	var mergedContent string
	if insertPos == 0 {
		mergedContent = string(p.TemplateContent) + token.NewlineLF + existingStr
	} else {
		mergedContent = existingStr[:insertPos] + token.NewlineLF +
			string(p.TemplateContent) + token.NewlineLF + existingStr[insertPos:]
	}

	if writeErr := os.WriteFile(
		p.Filename, []byte(mergedContent), fs.PermFile,
	); writeErr != nil {
		return false, errFs.WriteMerged(p.Filename, writeErr)
	}
	initialize.Merged(cmd, p.Filename)
	return false, nil
}

// UpdateMarkedSection replaces content between start/end markers in a file.
//
// Creates a timestamped backup before writing. If the end marker is missing,
// replaces it from the start marker to the end of the file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: Path to the file being updated
//   - existing: Current file content
//   - newTemplate: New template content (must contain both markers)
//   - markerStart: Opening marker string
//   - markerEnd: Closing marker string
//   - updateInfoFn: Called after a successful write to report the update
//
// Returns:
//   - error: Non-nil if markers are missing or file operations fail
func UpdateMarkedSection(
	cmd *cobra.Command,
	filename, existing string,
	newTemplate []byte,
	markerStart, markerEnd string,
	updateInfoFn func(cmd *cobra.Command, filename string),
) error {
	startIdx := strings.Index(existing, markerStart)
	if startIdx == -1 {
		return promptErr.MarkerNotFound(filename)
	}

	endIdx := strings.Index(existing, markerEnd)
	if endIdx == -1 {
		endIdx = len(existing)
	} else {
		endIdx += len(markerEnd)
	}

	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, markerStart)
	templateEnd := strings.Index(templateStr, markerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return promptErr.TemplateMissingMarkers(filename)
	}

	sectionContent := templateStr[templateStart : templateEnd+len(markerEnd)]
	newContent := existing[:startIdx] + sectionContent + existing[endIdx:]

	if bkErr := backupFile(cmd, filename, []byte(existing)); bkErr != nil {
		return bkErr
	}

	if writeErr := os.WriteFile(
		filename, []byte(newContent), fs.PermFile,
	); writeErr != nil {
		return errFs.FileUpdate(filename, writeErr)
	}

	updateInfoFn(cmd, filename)
	return nil
}
