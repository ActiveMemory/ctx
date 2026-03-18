//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/err/backup"
	fs2 "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/prompt"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
	"github.com/spf13/cobra"
)

// FindInsertionPoint finds where to insert ctx content in an existing file.
//
// Parameters:
//   - content: Existing file content
//
// Returns:
//   - int: Position to insert at
func FindInsertionPoint(content string) int {
	lines := strings.Split(content, token.NewlineLF)
	pos := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			pos += len(line) + 1
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			level := 0
			for _, ch := range trimmed {
				if ch == '#' {
					level++
				} else {
					break
				}
			}
			if level == 1 {
				pos += len(line) + 1
				for j := i + 1; j < len(lines); j++ {
					if strings.TrimSpace(lines[j]) == "" {
						pos += len(lines[j]) + 1
					} else {
						break
					}
				}
				return pos
			}
			return 0
		}
		return 0
	}
	return 0
}

// UpdateCtxSection replaces the existing ctx section between markers with new content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - existing: Existing file content
//   - newTemplate: New template content
//
// Returns:
//   - error: Non-nil if markers are missing or file operations fail
func UpdateCtxSection(cmd *cobra.Command, existing string, newTemplate []byte) error {
	startIdx := strings.Index(existing, marker.CtxMarkerStart)
	if startIdx == -1 {
		return ctxerr.MarkerNotFound("ctx")
	}
	endIdx := strings.Index(existing, marker.CtxMarkerEnd)
	if endIdx == -1 {
		endIdx = len(existing)
	} else {
		endIdx += len(marker.CtxMarkerEnd)
	}
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, marker.CtxMarkerStart)
	templateEnd := strings.Index(templateStr, marker.CtxMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return ctxerr.TemplateMissingMarkers("ctx")
	}
	ctxContent := templateStr[templateStart : templateEnd+len(marker.CtxMarkerEnd)]
	newContent := existing[:startIdx] + ctxContent + existing[endIdx:]
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", claude.Md, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), fs.PermFile); err != nil {
		return backup.CreateGeneric(err)
	}
	initialize.Backup(cmd, backupName)
	if err := os.WriteFile(claude.Md, []byte(newContent), fs.PermFile); err != nil {
		return fs2.FileUpdate(claude.Md, err)
	}
	initialize.UpdatedCtxSection(cmd, claude.Md)
	return nil
}
