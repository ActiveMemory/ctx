//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/cli/add/core/format"
	coreAppend "github.com/ActiveMemory/ctx/internal/cli/add/core/insert"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errAdd "github.com/ActiveMemory/ctx/internal/err/add"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Write formats and writes an entry to the appropriate context file.
//
// Handles the complete write cycle: read existing content,
// format the entry,
// append it, write back, and update the index if needed.
//
// Parameters:
//   - params: Params containing type, content, and optional fields
//
// Returns:
//   - error: Non-nil if the type is unknown, the file
//     doesn't exist, or write fails
func Write(params Params) error {
	fType := strings.ToLower(params.Type)

	fileName, ok := entry.CtxFile(fType)
	if !ok {
		return errAdd.UnknownType(fType)
	}

	contextDir := params.ContextDir
	if contextDir == "" {
		contextDir = rc.ContextDir()
	}
	filePath := filepath.Join(contextDir, fileName)

	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return errAdd.FileNotFound(filePath)
	}

	existing, readErr := os.ReadFile(filepath.Clean(filePath))
	if readErr != nil {
		return errFs.FileRead(filePath, readErr)
	}

	var formatted string
	switch fType {
	case entry.Decision:
		formatted = format.Decision(
			params.Content, params.Context, params.Rationale, params.Consequence,
		)
	case entry.Task:
		formatted = format.Task(params.Content, params.Priority)
	case entry.Learning:
		formatted = format.Learning(
			params.Content, params.Context, params.Lesson, params.Application,
		)
	case entry.Convention:
		formatted = format.Convention(params.Content)
	default:
		return errAdd.UnknownType(fType)
	}

	newContent := coreAppend.AppendEntry(
		existing, formatted, fType, params.Section,
	)

	if writeErr := os.WriteFile( //nolint:gosec // path from rc.ContextDir, trusted
		filePath, newContent, fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(filePath, writeErr)
	}

	switch fType {
	case entry.Decision:
		indexed := index.UpdateDecisions(string(newContent))
		if indexErr := os.WriteFile( //nolint:gosec // path from rc.ContextDir, trusted
			filePath, []byte(indexed), fs.PermFile,
		); indexErr != nil {
			return errAdd.IndexUpdate(filePath, indexErr)
		}
	case entry.Learning:
		indexed := index.UpdateLearnings(string(newContent))
		if indexErr := os.WriteFile( //nolint:gosec // path from rc.ContextDir, trusted
			filePath, []byte(indexed), fs.PermFile,
		); indexErr != nil {
			return errAdd.IndexUpdate(filePath, indexErr)
		}
		// case entry.Task, entry.Convention:
		// No index to update for these types
	}

	return nil
}

// ValidateAndWrite validates the entry params and writes the entry.
//
// Parameters:
//   - params: entry parameters with type, content, and optional fields
//
// Returns:
//   - error: validation or write error
func ValidateAndWrite(params Params) error {
	if vErr := Validate(params, nil); vErr != nil {
		return vErr
	}
	return Write(params)
}
