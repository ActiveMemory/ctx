//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// WriteArchive writes task content to a dated archive file in .context/archive/.
//
// Creates the archive directory if needed. If a file for today already exists,
// the new content is appended. Otherwise a new file is created with a header.
//
// Returns the path to the written archive file.
func WriteArchive(content string) (string, error) {
	archiveDir := filepath.Join(rc.ContextDir(), config.DirArchive)
	if err := os.MkdirAll(archiveDir, config.PermExec); err != nil {
		return "", fmt.Errorf("failed to create archive directory: %w", err)
	}

	now := time.Now()
	archiveFile := filepath.Join(
		archiveDir,
		fmt.Sprintf("tasks-%s.md", now.Format("2006-01-02")),
	)

	nl := config.NewlineLF
	var finalContent string
	if existing, err := os.ReadFile(filepath.Clean(archiveFile)); err == nil {
		finalContent = string(existing) + nl + content
	} else {
		finalContent = config.HeadingArchivedTasks + " - " +
			now.Format("2006-01-02") + nl + nl + content
	}

	if err := os.WriteFile(archiveFile, []byte(finalContent), config.PermFile); err != nil {
		return "", fmt.Errorf("failed to write archive: %w", err)
	}

	return archiveFile, nil
}
