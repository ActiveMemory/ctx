//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tidy

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errBackup "github.com/ActiveMemory/ctx/internal/err/backup"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// WriteArchive writes content to a dated archive file in .context/archive/.
//
// Creates the archive directory if needed. If a file for today already exists,
// the new content is appended. Otherwise, a new file is created with a header.
//
// Parameters:
//   - prefix: File name prefix (e.g., "tasks", "decisions", "learnings")
//   - heading: Markdown heading for new archive files
//     (e.g., config.HeadingArchivedTasks)
//   - content: The content to archive
//
// Returns:
//   - string: Path to the written archive file
//   - error: If creating the archive directory or writing fails
func WriteArchive(prefix, heading, content string) (string, error) {
	ctxDir, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		return "", ctxErr
	}
	archiveDir := filepath.Join(ctxDir, dir.Archive)
	if mkErr := io.SafeMkdirAll(archiveDir, fs.PermExec); mkErr != nil {
		return "", errBackup.CreateArchiveDir(mkErr)
	}

	now := time.Now()
	dateStr := now.Format(cfgTime.DateFormat)
	archiveFile := filepath.Join(
		archiveDir,
		fmt.Sprintf(archive.TplFilename, prefix, dateStr),
	)

	nl := token.NewlineLF
	var finalContent string
	cleanPath := filepath.Clean(archiveFile)
	if existing, readErr := io.SafeReadUserFile(cleanPath); readErr == nil {
		finalContent = string(existing) + nl + content
	} else {
		finalContent = heading + archive.DateSep +
			dateStr + nl + nl + content
	}

	if writeErr := io.SafeWriteFile(
		archiveFile, []byte(finalContent), fs.PermFile,
	); writeErr != nil {
		return "", errBackup.WriteArchive(writeErr)
	}

	return archiveFile, nil
}
