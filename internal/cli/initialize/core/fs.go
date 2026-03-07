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

	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
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
	lines := strings.Split(content, config.NewlineLF)
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
	startIdx := strings.Index(existing, config.CtxMarkerStart)
	if startIdx == -1 {
		return ctxerr.MarkerNotFound("ctx")
	}
	endIdx := strings.Index(existing, config.CtxMarkerEnd)
	if endIdx == -1 {
		endIdx = len(existing)
	} else {
		endIdx += len(config.CtxMarkerEnd)
	}
	templateStr := string(newTemplate)
	templateStart := strings.Index(templateStr, config.CtxMarkerStart)
	templateEnd := strings.Index(templateStr, config.CtxMarkerEnd)
	if templateStart == -1 || templateEnd == -1 {
		return ctxerr.TemplateMissingMarkers("ctx")
	}
	ctxContent := templateStr[templateStart : templateEnd+len(config.CtxMarkerEnd)]
	newContent := existing[:startIdx] + ctxContent + existing[endIdx:]
	timestamp := time.Now().Unix()
	backupName := fmt.Sprintf("%s.%d.bak", config.FileClaudeMd, timestamp)
	if err := os.WriteFile(backupName, []byte(existing), config.PermFile); err != nil {
		return ctxerr.CreateBackupGeneric(err)
	}
	write.InitBackup(cmd, backupName)
	if err := os.WriteFile(config.FileClaudeMd, []byte(newContent), config.PermFile); err != nil {
		return ctxerr.FileUpdate(config.FileClaudeMd, err)
	}
	write.InitUpdatedCtxSection(cmd, config.FileClaudeMd)
	return nil
}
