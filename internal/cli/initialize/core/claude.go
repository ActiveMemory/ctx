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
		return ctxerr.ReadInitTemplate("CLAUDE.md", err)
	}
	existingContent, err := os.ReadFile(file.FileClaudeMd)
	fileExists := err == nil
	if !fileExists {
		if err := os.WriteFile(file.FileClaudeMd, templateContent, fs.PermFile); err != nil {
			return ctxerr.FileWrite(file.FileClaudeMd, err)
		}
		write.InitCreated(cmd, file.FileClaudeMd)
		return nil
	}
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, config.CtxMarkerStart)
	if hasCtxMarkers {
		if !force {
			write.InitCtxContentExists(cmd, file.FileClaudeMd)
			return nil
		}
		return UpdateCtxSection(cmd, existingStr, templateContent)
	}
	if !autoMerge {
		write.InitFileExistsNoCtx(cmd, file.FileClaudeMd)
		cmd.Println("Would you like to append ctx context management instructions?")
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return ctxerr.ReadInput(err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != file.ConfirmShort && response != file.ConfirmLong {
			write.InitSkippedPlain(cmd, file.FileClaudeMd)
			return nil
		}
	}
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", file.FileClaudeMd, timestamp)
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
	if err := os.WriteFile(file.FileClaudeMd, []byte(mergedContent), fs.PermFile); err != nil {
		return ctxerr.WriteMerged(file.FileClaudeMd, err)
	}
	write.InitMerged(cmd, file.FileClaudeMd)
	return nil
}
