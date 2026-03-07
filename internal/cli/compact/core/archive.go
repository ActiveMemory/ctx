//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// WriteArchive writes content to a dated archive file in .context/archive/.
//
// Creates the archive directory if needed. If a file for today already exists,
// the new content is appended. Otherwise, a new file is created with a header.
//
// Parameters:
//   - prefix: File name prefix (e.g., "tasks", "decisions", "learnings")
//   - heading: Markdown heading for new archive files (e.g., config.HeadingArchivedTasks)
//   - content: The content to archive
//
// Returns the path to the written archive file.
func WriteArchive(prefix, heading, content string) (string, error) {
	archiveDir := filepath.Join(rc.ContextDir(), config.DirArchive)
	if mkErr := os.MkdirAll(archiveDir, config.PermExec); mkErr != nil {
		return "", ctxerr.CreateArchiveDir(mkErr)
	}

	now := time.Now()
	dateStr := now.Format(config.DateFormat)
	archiveFile := filepath.Join(
		archiveDir,
		fmt.Sprintf(config.TplArchiveFilename, prefix, dateStr),
	)

	nl := config.NewlineLF
	var finalContent string
	if existing, readErr := os.ReadFile(filepath.Clean(archiveFile)); readErr == nil {
		finalContent = string(existing) + nl + content
	} else {
		finalContent = heading + config.ArchiveDateSep +
			dateStr + nl + nl + content
	}

	if writeErr := os.WriteFile(archiveFile, []byte(finalContent), config.PermFile); writeErr != nil {
		return "", ctxerr.WriteArchive(writeErr)
	}

	return archiveFile, nil
}
