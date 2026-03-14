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

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/memory"
	time2 "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Sync copies sourcePath to .context/memory/mirror.md, archiving the
// previous mirror if one exists. Creates directories as needed.
func Sync(contextDir, sourcePath string) (SyncResult, error) {
	mirrorDir := filepath.Join(contextDir, dir.Memory)
	mirrorPath := filepath.Join(mirrorDir, memory.MemoryMirror)

	sourceData, readErr := io.SafeReadUserFile(sourcePath)
	if readErr != nil {
		return SyncResult{}, ctxerr.MemoryReadSource(readErr)
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
			return SyncResult{}, ctxerr.MemoryArchivePrevious(archiveErr)
		}
		result.ArchivedTo = archivePath
	}

	if mkErr := os.MkdirAll(mirrorDir, fs.PermExec); mkErr != nil {
		return SyncResult{}, ctxerr.MemoryCreateDir(mkErr)
	}

	if writeErr := os.WriteFile(mirrorPath, sourceData, fs.PermFile); writeErr != nil {
		return SyncResult{}, ctxerr.MemoryWriteMirror(writeErr)
	}

	return result, nil
}

// Archive copies the current mirror.md to archive/mirror-<timestamp>.md.
// Returns the archive path. Returns an error if no mirror exists.
func Archive(contextDir string) (string, error) {
	mirrorPath := filepath.Join(contextDir, dir.Memory, memory.MemoryMirror)
	archiveDir := filepath.Join(contextDir, dir.MemoryArchive)

	data, readErr := io.SafeReadUserFile(mirrorPath)
	if readErr != nil {
		return "", ctxerr.MemoryReadMirrorArchive(readErr)
	}

	if mkErr := os.MkdirAll(archiveDir, fs.PermExec); mkErr != nil {
		return "", ctxerr.MemoryCreateArchiveDir(mkErr)
	}

	ts := time.Now().Format(time2.TimestampCompact)
	archivePath := filepath.Join(archiveDir, memory.PrefixMirror+ts+file.ExtMarkdown)

	if writeErr := os.WriteFile(archivePath, data, fs.PermFile); writeErr != nil {
		return "", ctxerr.MemoryWriteArchive(writeErr)
	}

	return archivePath, nil
}

// Diff returns a simple line-based diff between the mirror and the source.
// Returns empty string when files are identical.
func Diff(contextDir, sourcePath string) (string, error) {
	mirrorPath := filepath.Join(contextDir, dir.Memory, memory.MemoryMirror)

	mirrorData, mirrorErr := io.SafeReadUserFile(mirrorPath)
	if mirrorErr != nil {
		return "", ctxerr.MemoryReadMirror(mirrorErr)
	}

	sourceData, sourceErr := io.SafeReadUserFile(sourcePath)
	if sourceErr != nil {
		return "", ctxerr.MemoryReadDiffSource(sourceErr)
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
func HasDrift(contextDir, sourcePath string) bool {
	mirrorPath := filepath.Join(contextDir, dir.Memory, memory.MemoryMirror)

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

func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	return bytes.Count(data, []byte(token.NewlineLF))
}

// simpleDiff produces a minimal unified-style diff header with added/removed lines.
func simpleDiff(oldPath, newPath string, oldLines, newLines []string) string {
	var buf strings.Builder
	_, _ = fmt.Fprintf(&buf, assets.TextDesc(assets.TextDescKeyMemoryDiffOldFormat), oldPath)
	_, _ = fmt.Fprintf(&buf, assets.TextDesc(assets.TextDescKeyMemoryDiffNewFormat), newPath)

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
