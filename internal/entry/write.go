//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/cli/add/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Write formats and writes an entry to the appropriate context file.
//
// Handles the complete write cycle: read existing content, format the entry,
// append it, write back, and update the index if needed.
//
// Parameters:
//   - params: Params containing type, content, and optional fields
//
// Returns:
//   - error: Non-nil if the type is unknown, the file doesn't exist, or write fails
func Write(params Params) error {
	fType := strings.ToLower(params.Type)

	fileName, ok := config.FileType[fType]
	if !ok {
		return write.ErrUnknownType(fType)
	}

	contextDir := params.ContextDir
	if contextDir == "" {
		contextDir = rc.ContextDir()
	}
	filePath := filepath.Join(contextDir, fileName)

	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return write.ErrFileNotFound(filePath)
	}

	existing, readErr := os.ReadFile(filepath.Clean(filePath))
	if readErr != nil {
		return write.ErrFileRead(filePath, readErr)
	}

	var formatted string
	switch config.UserInputToEntry(fType) {
	case config.EntryDecision:
		formatted = core.FormatDecision(
			params.Content, params.Context, params.Rationale, params.Consequences,
		)
	case config.EntryTask:
		formatted = core.FormatTask(params.Content, params.Priority)
	case config.EntryLearning:
		formatted = core.FormatLearning(
			params.Content, params.Context, params.Lesson, params.Application,
		)
	case config.EntryConvention:
		formatted = core.FormatConvention(params.Content)
	default:
		return write.ErrUnknownType(fType)
	}

	newContent := core.AppendEntry(existing, formatted, fType, params.Section)

	if writeErr := os.WriteFile(
		filePath, newContent, config.PermFile,
	); writeErr != nil {
		return write.ErrFileWriteAdd(filePath, writeErr)
	}

	switch config.UserInputToEntry(fType) {
	case config.EntryDecision:
		indexed := index.UpdateDecisions(string(newContent))
		if indexErr := os.WriteFile(
			filePath, []byte(indexed), config.PermFile,
		); indexErr != nil {
			return write.ErrIndexUpdate(filePath, indexErr)
		}
	case config.EntryLearning:
		indexed := index.UpdateLearnings(string(newContent))
		if indexErr := os.WriteFile(filePath, []byte(indexed), config.PermFile); indexErr != nil {
			return write.ErrIndexUpdate(filePath, indexErr)
		}
	case config.EntryTask, config.EntryConvention:
		// No index to update for these types
	}

	return nil
}
