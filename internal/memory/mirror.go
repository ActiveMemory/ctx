//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/memory"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errMemory "github.com/ActiveMemory/ctx/internal/err/memory"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Sync copies sourcePath to .context/memory/mirror.md, archiving the
// previous mirror if one exists. Creates directories as needed.
//
// Parameters:
//   - contextDir: Path to the project context directory
//   - sourcePath: Path to the source MEMORY.md file
//
// Returns:
//   - SyncResult: Summary of what was copied and archived
//   - error: If reading, writing, or archiving fails
func Sync(contextDir, sourcePath string) (SyncResult, error) {
	mirrorDir := filepath.Join(contextDir, dir.Memory)
	mirrorPath := filepath.Join(mirrorDir, memory.Mirror)

	sourceData, readErr := io.SafeReadUserFile(sourcePath)
	if readErr != nil {
		return SyncResult{}, errMemory.ReadSource(readErr)
	}

	result := SyncResult{
		SourcePath:  sourcePath,
		MirrorPath:  mirrorPath,
		SourceLines: countLines(sourceData),
	}

	// Archive existing mirror before overwrite
	if existingData, statErr := io.SafeReadUserFile(mirrorPath); statErr == nil {
		result.MirrorLines = countLines(existingData)
		archivePath, archiveErr := Archive(contextDir)
		if archiveErr != nil {
			return SyncResult{}, errMemory.ArchivePrevious(archiveErr)
		}
		result.ArchivedTo = archivePath
	}

	if mkErr := os.MkdirAll(mirrorDir, fs.PermExec); mkErr != nil {
		return SyncResult{}, errMemory.CreateDir(mkErr)
	}

	if writeErr := os.WriteFile(
		mirrorPath, sourceData, fs.PermFile,
	); writeErr != nil {
		return SyncResult{}, errMemory.WriteMirror(writeErr)
	}

	return result, nil
}

// Archive copies the current mirror.md to archive/mirror-<timestamp>.md.
// Returns the archive path. Returns an error if no mirror exists.
//
// Parameters:
//   - contextDir: Path to the project context directory
//
// Returns:
//   - string: Path to the written archive file
//   - error: If no mirror exists or writing fails
func Archive(contextDir string) (string, error) {
	mirrorPath := filepath.Join(contextDir, dir.Memory, memory.Mirror)
	archiveDir := filepath.Join(contextDir, dir.MemoryArchive)

	data, readErr := io.SafeReadUserFile(mirrorPath)
	if readErr != nil {
		return "", errMemory.ReadMirrorArchive(readErr)
	}

	if mkErr := os.MkdirAll(archiveDir, fs.PermExec); mkErr != nil {
		return "", errMemory.CreateArchiveDir(mkErr)
	}

	ts := time.Now().Format(cfgTime.CompactTimestamp)
	archiveName := memory.PrefixMirror + ts + file.ExtMarkdown
	archivePath := filepath.Join(archiveDir, archiveName)

	if writeErr := os.WriteFile(archivePath, data, fs.PermFile); writeErr != nil {
		return "", errMemory.WriteArchive(writeErr)
	}

	return archivePath, nil
}

// Diff returns a simple line-based diff between the mirror and the source.
// Returns empty string when files are identical.
//
// Parameters:
//   - contextDir: Path to the project context directory
//   - sourcePath: Path to the source MEMORY.md file
//
// Returns:
//   - string: Line-based diff, or empty if identical
//   - error: If either file cannot be read
func Diff(contextDir, sourcePath string) (string, error) {
	mirrorPath := filepath.Join(contextDir, dir.Memory, memory.Mirror)

	mirrorData, mirrorErr := io.SafeReadUserFile(mirrorPath)
	if mirrorErr != nil {
		return "", errMemory.ReadMirror(mirrorErr)
	}

	sourceData, sourceErr := io.SafeReadUserFile(sourcePath)
	if sourceErr != nil {
		return "", errMemory.ReadDiffSource(sourceErr)
	}

	if bytes.Equal(mirrorData, sourceData) {
		return "", nil
	}

	mirrorLines := strings.Split(string(mirrorData), token.NewlineLF)
	sourceLines := strings.Split(string(sourceData), token.NewlineLF)

	return simpleDiff(mirrorPath, sourcePath, mirrorLines, sourceLines), nil
}

// HasDrift checks whether MEMORY.md has been modified since the last sync.
// Returns false if either file is missing (no drift to report).
//
// Parameters:
//   - contextDir: Path to the project context directory
//   - sourcePath: Path to the source MEMORY.md file
//
// Returns:
//   - bool: True if MEMORY.md has been modified since the last sync
func HasDrift(contextDir, sourcePath string) bool {
	mirrorPath := filepath.Join(contextDir, dir.Memory, memory.Mirror)

	sourceInfo, sourceErr := os.Stat(sourcePath)
	if sourceErr != nil {
		return false
	}

	mirrorInfo, mirrorErr := os.Stat(mirrorPath)
	if mirrorErr != nil {
		return false
	}

	return sourceInfo.ModTime().After(mirrorInfo.ModTime())
}

// ArchiveCount returns the number of archived mirror snapshots.
//
// Parameters:
//   - contextDir: Path to the project context directory
//
// Returns:
//   - int: Number of archived mirror snapshot files
func ArchiveCount(contextDir string) int {
	archiveDir := filepath.Join(contextDir, dir.MemoryArchive)
	entries, readErr := os.ReadDir(archiveDir)
	if readErr != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), memory.PrefixMirror) {
			count++
		}
	}
	return count
}
