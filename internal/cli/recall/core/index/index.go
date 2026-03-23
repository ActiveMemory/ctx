//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package index

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"
)

// BuildSessionIndex scans journal .md files in journalDir and returns a
// map of session_id → filename.
//
// Two-pass matching:
//  1. Parse YAML frontmatter for a session_id field (authoritative).
//  2. For files without session_id, extract the last 8 characters before
//     config.ExtMarkdown and treat them as a short ID candidate (migration path for
//     legacy exports).
//
// Parameters:
//   - journalDir: Path to the journal directory
//
// Returns:
//   - map[string]string: session ID → filename mapping
func BuildSessionIndex(journalDir string) map[string]string {
	index := make(map[string]string)

	entries, readErr := os.ReadDir(journalDir)
	if readErr != nil {
		return index
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), file.ExtMarkdown) {
			continue
		}

		path := filepath.Join(journalDir, e.Name())
		content, fileErr := os.ReadFile(filepath.Clean(path))
		if fileErr != nil {
			continue
		}

		// Pass 1: look for session_id in YAML frontmatter.
		if sid := ExtractSessionID(string(content)); sid != "" {
			// Only map the base filename (without -p* suffix) or each
			// part file to the same session ID. The caller uses
			// index[sessionID] to find the base filename.
			if _, exists := index[sid]; !exists {
				index[sid] = e.Name()
			}
			continue
		}

		// Pass 2: extract short ID from the filename as fallback.
		// Filename format: YYYY-MM-DD-slug-SHORTID.md or ...-pN.md
		name := e.Name()
		// Strip multipart suffix (e.g., "-p2.md" → config.ExtMarkdown).
		baseName := strings.TrimSuffix(name, file.ExtMarkdown)
		if idx := strings.LastIndex(baseName, journal.MultipartSuffix); idx > 0 {
			suffix := baseName[idx+2:]
			allDigits := true
			for _, r := range suffix {
				if r < '0' || r > '9' {
					allDigits = false
					break
				}
			}
			if allDigits && len(suffix) > 0 {
				// This is a multipart file; skip it — the base file
				// provides the index entry.
				continue
			}
		}

		// Extract the last 8 chars before .md as candidate short ID.
		if len(baseName) >= journal.ShortIDLen {
			shortID := baseName[len(baseName)-journal.ShortIDLen:]
			// Store with the short ID as key (caller matches against
			// session.ID[:8]).
			if _, exists := index[shortID]; !exists {
				index[shortID] = name
			}
		}
	}

	return index
}

// ExtractSessionID parses session_id from YAML frontmatter.
//
// Looks for a line matching `session_id: "..."` or `session_id: ...`
// within the frontmatter block delimited by "---".
//
// Parameters:
//   - content: Full file content
//
// Returns:
//   - string: The session ID, or "" if not found
func ExtractSessionID(content string) string {
	return ExtractFrontmatterField(content, session.FmKeySessionID)
}

// LookupSessionFile finds the existing filename for a session in the index.
//
// Checks the full session ID first (frontmatter-based match), then falls
// back to the short ID (filename-based legacy match).
//
// Parameters:
//   - index: Session index from BuildSessionIndex
//   - sessionID: Full session UUID
//
// Returns:
//   - string: Existing filename, or "" if not found
func LookupSessionFile(index map[string]string, sessionID string) string {
	if name, ok := index[sessionID]; ok {
		return name
	}
	short := sessionID
	if len(short) > journal.ShortIDLen {
		short = short[:journal.ShortIDLen]
	}
	if name, ok := index[short]; ok {
		return name
	}
	return ""
}

// ExtractFrontmatterField extracts a single field value from YAML frontmatter.
//
// Parameters:
//   - content: Full file content
//   - field: Field name to extract (e.g. "title")
//
// Returns:
//   - string: The field value (unquoted), or "" if not found
func ExtractFrontmatterField(content, field string) string {
	nl := token.NewlineLF
	fmOpen := token.Separator + nl

	if !strings.HasPrefix(content, fmOpen) {
		return ""
	}
	end := strings.Index(content[len(fmOpen):], nl+token.Separator+nl)
	if end < 0 {
		return ""
	}
	fmBlock := content[len(fmOpen) : len(fmOpen)+end]

	prefix := field + token.Colon
	for _, line := range strings.Split(fmBlock, nl) {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			val := strings.TrimSpace(strings.TrimPrefix(line, prefix))
			val = strings.Trim(val, `"'`)
			return val
		}
	}
	return ""
}

// RenameJournalFiles renames a journal file (and its multipart siblings)
// from oldBase to newBase within journalDir.
//
// Handles both the base file (oldBase.md → newBase.md) and multipart
// files (oldBase-pN.md → newBase-pN.md). Updates internal navigation
// links in multipart files to reference the new filenames.
//
// Parameters:
//   - journalDir: Path to the journal directory
//   - oldBase: Old base filename without extension
//   - newBase: New base filename without extension
//   - numParts: Expected number of parts (used for nav link updates)
func RenameJournalFiles(journalDir, oldBase, newBase string, numParts int) {
	// Rename base file.
	oldPath := filepath.Join(journalDir, oldBase+file.ExtMarkdown)
	newPath := filepath.Join(journalDir, newBase+file.ExtMarkdown)
	if _, statErr := os.Stat(oldPath); statErr == nil {
		_ = os.Rename(oldPath, newPath)
	}

	// Rename multipart files and update nav links.
	for p := 2; p <= numParts; p++ {
		oldPart := filepath.Join(
			journalDir, fmt.Sprintf(tpl.RecallPartFilename, oldBase, p),
		)
		newPart := filepath.Join(
			journalDir, fmt.Sprintf(tpl.RecallPartFilename, newBase, p),
		)
		if _, statErr := os.Stat(oldPart); statErr == nil {
			_ = os.Rename(oldPart, newPart)
		}
	}

	// Update navigation links inside all parts to reference the new baseName.
	UpdateNavLinks(journalDir, newBase, oldBase, numParts)
}

// UpdateNavLinks replaces references to oldBase with newBase inside
// all part files for a session.
//
// Parameters:
//   - journalDir: Path to the journal directory
//   - newBase: New base filename without extension
//   - oldBase: Old base filename without extension
//   - numParts: Total number of parts
func UpdateNavLinks(journalDir, newBase, oldBase string, numParts int) {
	if numParts <= 1 {
		return
	}

	files := []string{filepath.Join(journalDir, newBase+file.ExtMarkdown)}
	for p := 2; p <= numParts; p++ {
		files = append(files, filepath.Join(journalDir,
			fmt.Sprintf(tpl.RecallPartFilename, newBase, p)))
	}

	for _, f := range files {
		data, readErr := os.ReadFile(filepath.Clean(f))
		if readErr != nil {
			continue
		}
		updated := strings.ReplaceAll(string(data), oldBase, newBase)
		if updated != string(data) {
			_ = io.SafeWriteFile(f, []byte(updated), fs.PermFile)
		}
	}
}
