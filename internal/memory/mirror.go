//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
)

// SyncResult holds the outcome of a Sync operation.
type SyncResult struct {
	SourcePath  string
	MirrorPath  string
	ArchivedTo  string // empty if no prior mirror existed
	SourceLines int
	MirrorLines int // lines in the previous mirror (0 if first sync)
}

// Sync copies sourcePath to .context/memory/mirror.md, archiving the
// previous mirror if one exists. Creates directories as needed.
func Sync(contextDir, sourcePath string) (SyncResult, error) {
	mirrorDir := filepath.Join(contextDir, config.DirMemory)
	mirrorPath := filepath.Join(mirrorDir, config.FileMemoryMirror)

	sourceData, readErr := os.ReadFile(sourcePath) //nolint:gosec // caller-provided path
	if readErr != nil {
		return SyncResult{}, fmt.Errorf("reading source: %w", readErr)
	}

	result := SyncResult{
		SourcePath:  sourcePath,
		MirrorPath:  mirrorPath,
		SourceLines: countLines(sourceData),
	}

	// Archive existing mirror before overwrite
	if existingData, statErr := os.ReadFile(mirrorPath); statErr == nil { //nolint:gosec // project-local path
		result.MirrorLines = countLines(existingData)
		archivePath, archiveErr := Archive(contextDir)
		if archiveErr != nil {
			return SyncResult{}, fmt.Errorf("archiving previous mirror: %w", archiveErr)
		}
		result.ArchivedTo = archivePath
	}

	if mkErr := os.MkdirAll(mirrorDir, config.PermExec); mkErr != nil {
		return SyncResult{}, fmt.Errorf("creating memory directory: %w", mkErr)
	}

	if writeErr := os.WriteFile(mirrorPath, sourceData, config.PermFile); writeErr != nil {
		return SyncResult{}, fmt.Errorf("writing mirror: %w", writeErr)
	}

	return result, nil
}

// Archive copies the current mirror.md to archive/mirror-<timestamp>.md.
// Returns the archive path. Returns an error if no mirror exists.
func Archive(contextDir string) (string, error) {
	mirrorPath := filepath.Join(contextDir, config.DirMemory, config.FileMemoryMirror)
	archiveDir := filepath.Join(contextDir, config.DirMemoryArchive)

	data, readErr := os.ReadFile(mirrorPath) //nolint:gosec // project-local path
	if readErr != nil {
		return "", fmt.Errorf("reading mirror for archive: %w", readErr)
	}

	if mkErr := os.MkdirAll(archiveDir, config.PermExec); mkErr != nil {
		return "", fmt.Errorf("creating archive directory: %w", mkErr)
	}

	ts := time.Now().Format(config.TimestampCompact)
	archivePath := filepath.Join(archiveDir, "mirror-"+ts+config.ExtMarkdown)

	if writeErr := os.WriteFile(archivePath, data, config.PermFile); writeErr != nil {
		return "", fmt.Errorf("writing archive: %w", writeErr)
	}

	return archivePath, nil
}

// Diff returns a simple line-based diff between the mirror and the source.
// Returns empty string when files are identical.
func Diff(contextDir, sourcePath string) (string, error) {
	mirrorPath := filepath.Join(contextDir, config.DirMemory, config.FileMemoryMirror)

	mirrorData, mirrorErr := os.ReadFile(mirrorPath) //nolint:gosec // project-local path
	if mirrorErr != nil {
		return "", fmt.Errorf("reading mirror: %w", mirrorErr)
	}

	sourceData, sourceErr := os.ReadFile(sourcePath) //nolint:gosec // caller-provided path
	if sourceErr != nil {
		return "", fmt.Errorf("reading source: %w", sourceErr)
	}

	if bytes.Equal(mirrorData, sourceData) {
		return "", nil
	}

	mirrorLines := strings.Split(string(mirrorData), config.NewlineLF)
	sourceLines := strings.Split(string(sourceData), config.NewlineLF)

	return simpleDiff(mirrorPath, sourcePath, mirrorLines, sourceLines), nil
}

// HasDrift checks whether MEMORY.md has been modified since the last sync.
// Returns false if either file is missing (no drift to report).
func HasDrift(contextDir, sourcePath string) bool {
	mirrorPath := filepath.Join(contextDir, config.DirMemory, config.FileMemoryMirror)

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
func ArchiveCount(contextDir string) int {
	archiveDir := filepath.Join(contextDir, config.DirMemoryArchive)
	entries, readErr := os.ReadDir(archiveDir)
	if readErr != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "mirror-") {
			count++
		}
	}
	return count
}

func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	return bytes.Count(data, []byte(config.NewlineLF))
}

// simpleDiff produces a minimal unified-style diff header with added/removed lines.
func simpleDiff(oldPath, newPath string, oldLines, newLines []string) string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("--- %s (mirror)\n", oldPath))
	buf.WriteString(fmt.Sprintf("+++ %s (source)\n", newPath))

	oldSet := make(map[string]bool, len(oldLines))
	for _, l := range oldLines {
		oldSet[l] = true
	}
	newSet := make(map[string]bool, len(newLines))
	for _, l := range newLines {
		newSet[l] = true
	}

	for _, l := range oldLines {
		if !newSet[l] {
			buf.WriteString("-" + l + config.NewlineLF)
		}
	}
	for _, l := range newLines {
		if !oldSet[l] {
			buf.WriteString("+" + l + config.NewlineLF)
		}
	}

	return buf.String()
}
