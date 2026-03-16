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
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/add"
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

	fileName, ok := entry.ToCtxFile[fType]
	if !ok {
		return add.ErrUnknownType(fType)
	}

	contextDir := params.ContextDir
	if contextDir == "" {
		contextDir = rc.ContextDir()
	}
	filePath := filepath.Join(contextDir, fileName)

	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return add.ErrFileNotFound(filePath)
	}

	existing, readErr := os.ReadFile(filepath.Clean(filePath))
	if readErr != nil {
		return add.ErrFileRead(filePath, readErr)
	}

	var formatted string
	switch fType {
	case entry.Decision:
		formatted = core.FormatDecision(
			params.Content, params.Context, params.Rationale, params.Consequence,
		)
	case entry.Task:
		formatted = core.FormatTask(params.Content, params.Priority)
	case entry.Learning:
		formatted = core.FormatLearning(
			params.Content, params.Context, params.Lesson, params.Application,
		)
	case entry.Convention:
		formatted = core.FormatConvention(params.Content)
	default:
		return add.ErrUnknownType(fType)
	}

	newContent := core.AppendEntry(existing, formatted, fType, params.Section)

	if writeErr := os.WriteFile(
		filePath, newContent, fs.PermFile,
	); writeErr != nil {
		return add.ErrFileWriteAdd(filePath, writeErr)
	}

	switch fType {
	case entry.Decision:
		indexed := index.UpdateDecisions(string(newContent))
		if indexErr := os.WriteFile(
			filePath, []byte(indexed), fs.PermFile,
		); indexErr != nil {
			return add.ErrIndexUpdate(filePath, indexErr)
		}
	case entry.Learning:
		indexed := index.UpdateLearnings(string(newContent))
		if indexErr := os.WriteFile(
			filePath, []byte(indexed), fs.PermFile,
		); indexErr != nil {
			return add.ErrIndexUpdate(filePath, indexErr)
		}
	case entry.Task, entry.Convention:
		// No index to update for these types
	}

	return nil
}

// ValidateAndWrite validates the entry params, writes the entry, and
// returns the target context file name.
//
// Parameters:
//   - params: entry parameters with type, content, and optional fields
//
// Returns:
//   - string: the context file name the entry was written to
//   - error: validation or write error
func ValidateAndWrite(params Params) (string, error) {
	if vErr := Validate(params, nil); vErr != nil {
		return "", vErr
	}

	if wErr := Write(params); wErr != nil {
		return "", wErr
	}

	return entry.ToCtxFile[strings.ToLower(params.Type)], nil
}
