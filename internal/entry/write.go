//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/cli/add/core/format"
	coreAppend "github.com/ActiveMemory/ctx/internal/cli/add/core/insert"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/entity"
	errAdd "github.com/ActiveMemory/ctx/internal/err/add"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/i18n"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Write formats and writes an entry to the appropriate context file.
//
// Handles the complete write cycle: read existing content,
// format the entry, append it, and write back. No index is
// maintained in the file; a table of contents is projected on
// demand by `ctx index <file>`.
//
// Parameters:
//   - params: Params containing type, content, and optional fields
//
// Returns:
//   - error: Non-nil if the type is unknown, the file
//     doesn't exist, or write fails
func Write(params entity.EntryParams) error {
	fType := i18n.Fold(params.Type)

	fileName, ok := entry.CtxFile(fType)
	if !ok {
		return errAdd.UnknownType(fType)
	}

	contextDir := params.ContextDir
	if contextDir == "" {
		declared, ctxErr := rc.ContextDir()
		if ctxErr != nil {
			return ctxErr
		}
		contextDir = declared
	}
	filePath := filepath.Join(contextDir, fileName)

	if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
		return errAdd.FileNotFound(filePath)
	}

	existing, readErr := io.SafeReadUserFile(filepath.Clean(filePath))
	if readErr != nil {
		return errFs.FileRead(filePath, readErr)
	}

	var formatted string
	switch fType {
	case entry.Decision:
		out, fErr := format.Decision(
			params.Content, params.Context, params.Rationale, params.Consequence,
		)
		if fErr != nil {
			return fErr
		}
		formatted = out
	case entry.Task:
		formatted = format.Task(
			params.Content, params.Priority,
			params.SessionID, params.Branch, params.Commit,
		)
	case entry.Learning:
		out, fErr := format.Learning(
			params.Content, params.Context, params.Lesson, params.Application,
		)
		if fErr != nil {
			return fErr
		}
		formatted = out
	case entry.Convention:
		formatted = format.Convention(params.Content)
	default:
		return errAdd.UnknownType(fType)
	}

	newContent := coreAppend.AppendEntry(
		existing, formatted, fType, params.Section,
	)

	if writeErr := io.SafeWriteFile(
		filePath, newContent, fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(filePath, writeErr)
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
func ValidateAndWrite(params entity.EntryParams) error {
	if vErr := Validate(params, nil); vErr != nil {
		return vErr
	}
	return Write(params)
}
