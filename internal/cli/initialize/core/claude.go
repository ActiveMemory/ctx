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

	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/err/backup"
	fs2 "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
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
		return ctxerr.ReadTemplate("CLAUDE.md", err)
	}
	existingContent, err := os.ReadFile(claude.Md)
	fileExists := err == nil
	if !fileExists {
		if err := os.WriteFile(claude.Md, templateContent, fs.PermFile); err != nil {
			return fs2.FileWrite(claude.Md, err)
		}
		initialize.Created(cmd, claude.Md)
		return nil
	}
	existingStr := string(existingContent)
	hasCtxMarkers := strings.Contains(existingStr, marker.CtxMarkerStart)
	if hasCtxMarkers {
		if !force {
			initialize.CtxContentExists(cmd, claude.Md)
			return nil
		}
		return UpdateCtxSection(cmd, existingStr, templateContent)
	}
	if !autoMerge {
		initialize.FileExistsNoCtx(cmd, claude.Md)
		cmd.Println("Would you like to append ctx context management instructions?")
		cmd.Print("[y/N] ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fs2.ReadInput(err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != cli.ConfirmShort && response != cli.ConfirmLong {
			initialize.SkippedPlain(cmd, claude.Md)
			return nil
		}
	}
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", claude.Md, timestamp)
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
	if err := os.WriteFile(claude.Md, []byte(mergedContent), fs.PermFile); err != nil {
		return fs2.WriteMerged(claude.Md, err)
	}
	initialize.Merged(cmd, claude.Md)
	return nil
}
