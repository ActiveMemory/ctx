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

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
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

	if writeErr := os.WriteFile(mirrorPath, sourceData, fs.PermFile); writeErr != nil {
		return SyncResult{}, errMemory.WriteMirror(writeErr)
	}

	return result, nil
}

// Archive copies the current mirror.md to archive/mirror-<timestamp>.md.
// Returns the archive path. Returns an error if no mirror exists.
//
// Parameters:
//   - contextDir: Path to the project context directory
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
	archivePath := filepath.Join(archiveDir, memory.PrefixMirror+ts+file.ExtMarkdown)

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

// countLines returns the number of newline characters in data.
//
// Parameters:
//   - data: Raw byte content to scan
//
// Returns:
//   - int: Count of newline characters; zero for empty input
func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	return bytes.Count(data, []byte(token.NewlineLF))
}

// simpleDiff produces a minimal unified-style diff header with added/removed lines.
//
// Parameters:
//   - oldPath: Label for the old file in the diff header
//   - newPath: Label for the new file in the diff header
//   - oldLines: Lines from the previous version
//   - newLines: Lines from the current version
//
// Returns:
//   - string: Formatted diff showing added and removed lines
func simpleDiff(oldPath, newPath string, oldLines, newLines []string) string {
	var buf strings.Builder
	_, _ = fmt.Fprintf(&buf, desc.Text(text.DescKeyMemoryDiffOldFormat), oldPath)
	_, _ = fmt.Fprintf(&buf, desc.Text(text.DescKeyMemoryDiffNewFormat), newPath)

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
			buf.WriteString("-" + l + token.NewlineLF)
		}
	}
	for _, l := range newLines {
		if !oldSet[l] {
			buf.WriteString("+" + l + token.NewlineLF)
		}
	}

	return buf.String()
}
