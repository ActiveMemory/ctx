package task

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// tasksFilePath returns the path to TASKS.md.
//
// Returns:
//   - string: Full path to .context/TASKS.md
func tasksFilePath() string {
	return filepath.Join(rc.ContextDir(), config.FileTask)
}

// archiveDirPath returns the path to the archive directory.
//
// Returns:
//   - string: Full path to .context/archive/
func archiveDirPath() string {
	return filepath.Join(rc.ContextDir(), config.DirArchive)
}
