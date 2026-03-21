//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/err/journal"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/session"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/recall"
)

// LockedFrontmatterLine is the YAML line inserted into frontmatter when
// a journal entry is locked.
var LockedFrontmatterLine = session.FrontmatterLocked + token.Colon + " " +
	cli.AnnotationTrue + "  # managed by ctx"

// MatchJournalFiles returns journal .md filenames matching the given
// patterns. If all is true, returns every .md file in the directory.
// Multi-part files (base + -pN parts) are included when the base matches.
//
// Parameters:
//   - journalDir: Path to the journal directory
//   - patterns: Slug, date, or short-ID substrings to match
//   - all: If true, return all .md files
//
// Returns:
//   - []string: Matching filenames (basename only)
//   - error: Non-nil on I/O failure
func MatchJournalFiles(
	journalDir string,
	patterns []string,
	all bool,
) ([]string, error) {
	entries, readErr := os.ReadDir(journalDir)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, nil
		}
		return nil, journal.ReadDir(readErr)
	}

	// Collect all .md filenames.
	var mdFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), file.ExtMarkdown) {
			mdFiles = append(mdFiles, e.Name())
		}
	}

	if all {
		return mdFiles, nil
	}

	// Build a set of matching base names, then expand to include parts.
	matchedBases := make(map[string]bool)
	for _, f := range mdFiles {
		lower := strings.ToLower(f)
		for _, pat := range patterns {
			if strings.Contains(lower, strings.ToLower(pat)) {
				base := MultipartBase(f)
				matchedBases[base] = true
			}
		}
	}

	// Expand: include all files sharing a matched base.
	var result []string
	for _, f := range mdFiles {
		base := MultipartBase(f)
		if matchedBases[base] {
			result = append(result, f)
		}
	}

	return result, nil
}

// MultipartBase returns the base name for a potentially multi-part file.
// For "2026-01-21-slug-abc12345-p2.md" it returns
// "2026-01-21-slug-abc12345.md". For non-multipart files, returns the
// filename as-is.
//
// Parameters:
//   - filename: Journal entry filename
//
// Returns:
//   - string: Base filename (without -pN suffix)
func MultipartBase(filename string) string {
	base := strings.TrimSuffix(filename, file.ExtMarkdown)
	if idx := strings.LastIndex(base, "-p"); idx > 0 {
		suffix := base[idx+2:]
		allDigits := true
		for _, r := range suffix {
			if r < '0' || r > '9' {
				allDigits = false
				break
			}
		}
		if allDigits && len(suffix) > 0 {
			return base[:idx] + file.ExtMarkdown
		}
	}
	return filename
}

// lockedPrefix is the frontmatter key prefix for locked lines.
var lockedPrefix = session.FrontmatterLocked + token.Colon

// UpdateLockFrontmatter inserts or removes the "locked: true" line in
// a journal file's YAML frontmatter. The state file is the source of
// truth; this is for human visibility only.
//
// Parameters:
//   - path: Absolute path to the journal .md file
//   - lock: True to insert, false to remove
func UpdateLockFrontmatter(path string, lock bool) {
	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		return
	}
	content := string(data)

	nl := token.NewlineLF
	fmOpen := token.Separator + nl

	if !strings.HasPrefix(content, fmOpen) {
		// No frontmatter — nothing to modify.
		return
	}

	closeIdx := strings.Index(content[len(fmOpen):], nl+token.Separator+nl)
	if closeIdx < 0 {
		return
	}

	fmEnd := len(fmOpen) + closeIdx // index of the newline before closing ---
	fmBlock := content[len(fmOpen):fmEnd]

	if lock {
		// Already has locked line?
		if strings.Contains(fmBlock, lockedPrefix) {
			return
		}
		// Insert before closing ---.
		updated := content[:fmEnd] + nl + LockedFrontmatterLine +
			content[fmEnd:]
		_ = io.SafeWriteFile(path, []byte(updated), fs.PermFile)
	} else {
		// Remove the locked line.
		lines := strings.Split(fmBlock, nl)
		var filtered []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, lockedPrefix) {
				continue
			}
			filtered = append(filtered, line)
		}
		newFM := strings.Join(filtered, nl)
		updated := content[:len(fmOpen)] + newFM + content[fmEnd:]
		_ = io.SafeWriteFile(path, []byte(updated), fs.PermFile)
	}
}

// FrontmatterHasLocked reads a journal file and returns true if its
// YAML frontmatter contains a "locked:" line with a truthy value.
//
// Parameters:
//   - path: Absolute path to the journal .md file
//
// Returns:
//   - bool: True if frontmatter contains "locked: true"
func FrontmatterHasLocked(path string) bool {
	data, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		return false
	}
	content := string(data)

	nl := token.NewlineLF
	fmOpen := token.Separator + nl

	if !strings.HasPrefix(content, fmOpen) {
		return false
	}

	closeIdx := strings.Index(content[len(fmOpen):], nl+token.Separator+nl)
	if closeIdx < 0 {
		return false
	}

	fmBlock := content[len(fmOpen) : len(fmOpen)+closeIdx]

	for _, line := range strings.Split(fmBlock, nl) {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, lockedPrefix) {
			continue
		}
		val := strings.TrimSpace(strings.TrimPrefix(trimmed, lockedPrefix))
		// Strip inline comment (e.g. "true  # managed by ctx").
		if idx := strings.Index(val, "#"); idx >= 0 {
			val = strings.TrimSpace(val[:idx])
		}
		return val == cli.AnnotationTrue
	}

	return false
}

// RunLockUnlock handles both lock and unlock commands.
//
// Parameters:
//   - cmd: Cobra command for output
//   - args: Patterns to match against journal filenames
//   - all: If true, apply to all journal entries
//   - lock: True for lock, false for unlock
//
// Returns:
//   - error: Non-nil on validation or I/O failure
func RunLockUnlock(
	cmd *cobra.Command,
	args []string,
	all, lock bool,
) error {
	if len(args) == 0 && !all {
		return cmd.Help()
	}
	if len(args) > 0 && all {
		return ctxErr.AllWithPattern()
	}

	journalDir := filepath.Join(rc.ContextDir(), dir.Journal)

	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return journal.LoadState(loadErr)
	}

	// Collect matching .md files.
	files, matchErr := MatchJournalFiles(journalDir, args, all)
	if matchErr != nil {
		return matchErr
	}
	if len(files) == 0 {
		if all {
			recall.LockUnlockNone(cmd)
		} else {
			return journal.NoEntriesMatch(strings.Join(args, ", "))
		}
		return nil
	}

	verb := session.FrontmatterLocked
	if !lock {
		verb = session.Unlocked
	}

	count := 0
	for _, filename := range files {
		alreadyLocked := jstate.Locked(filename)
		if lock && alreadyLocked {
			continue
		}
		if !lock && !alreadyLocked {
			continue
		}

		// Update state.
		if lock {
			jstate.Mark(filename, session.FrontmatterLocked)
		} else {
			jstate.Clear(filename, session.FrontmatterLocked)
		}

		// Update frontmatter for human visibility.
		path := filepath.Join(journalDir, filename)
		UpdateLockFrontmatter(path, lock)

		recall.LockUnlockEntry(cmd, filename, verb)
		count++
	}

	if saveErr := jstate.Save(journalDir); saveErr != nil {
		return journal.SaveState(saveErr)
	}

	recall.LockUnlockSummary(cmd, verb, count)

	return nil
}
