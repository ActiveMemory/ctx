//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// ScanJournalEntries reads all journal Markdown files and extracts metadata.
//
// Parameters:
//   - journalDir: Path to .context/journal/
//
// Returns:
//   - []JournalEntry: Parsed entries sorted by date (newest first)
//   - error: Non-nil if directory scanning fails
func ScanJournalEntries(journalDir string) ([]JournalEntry, error) {
	files, readErr := os.ReadDir(journalDir)
	if readErr != nil {
		return nil, readErr
	}

	var entries []JournalEntry
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), file.ExtMarkdown) {
			continue
		}

		path := filepath.Join(journalDir, f.Name())
		entry := ParseJournalEntry(path, f.Name())
		entries = append(entries, entry)
	}

	// Sort by datetime (newest first) - combine Date and Time
	sort.Slice(entries, func(i, j int) bool {
		// Compare Date+Time strings (YYYY-MM-DD + HH:MM:SS)
		di := entries[i].Date + " " + entries[i].Time
		dj := entries[j].Date + " " + entries[j].Time
		return di > dj
	})

	return entries, nil
}

// ParseJournalEntry extracts metadata from a journal file.
//
// Parameters:
//   - path: Full path to the journal file
//   - filename: Filename (e.g., "2026-01-21-async-roaming-allen-af7cba21.md")
//
// Returns:
//   - JournalEntry: Parsed entry with title, date, project extracted
func ParseJournalEntry(path, filename string) JournalEntry {
	entry := JournalEntry{
		Filename: filename,
		Path:     path,
	}

	// Extract date from the filename (YYYY-MM-DD-slug-id.md)
	if len(filename) >= journal.DatePrefixLen {
		entry.Date = filename[:journal.DatePrefixLen]
	}

	// Read the file to extract metadata
	content, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		entry.Title = strings.TrimSuffix(filename, file.ExtMarkdown)
		return entry
	}

	// File size
	entry.Size = int64(len(content))

	contentStr := string(content)

	// Parse YAML frontmatter if present
	nl := token.NewlineLF
	fmOpen := len(token.Separator + nl)
	if strings.HasPrefix(contentStr, token.Separator+nl) {
		if end := strings.Index(
			contentStr[fmOpen:], nl+token.Separator+nl,
		); end >= 0 {
			fmRaw := contentStr[fmOpen : fmOpen+end]
			var fm JournalFrontmatter
			if yaml.Unmarshal([]byte(fmRaw), &fm) == nil {
				if fm.Title != "" {
					entry.Title = fm.Title
				}
				if fm.Time != "" {
					entry.Time = fm.Time
				}
				if fm.Project != "" {
					entry.Project = fm.Project
				}
				if fm.SessionID != "" {
					entry.SessionID = fm.SessionID
				}
				if fm.Model != "" {
					entry.Model = fm.Model
				}
				entry.TokensIn = fm.TokensIn
				entry.TokensOut = fm.TokensOut
				entry.Topics = fm.Topics
				entry.Type = fm.Type
				entry.Outcome = fm.Outcome
				entry.KeyFiles = fm.KeyFiles
				entry.Summary = fm.Summary
			}
		}
	}

	// Check for suggestion mode sessions
	if strings.Contains(contentStr, desc.Text(text.DescKeyLabelSuggestionMode)) {
		entry.Suggestive = true
	}

	// Line-by-line parsing as fallback for fields not in frontmatter
	lines := strings.Split(contentStr, nl)
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Title from first H1 (only if frontmatter didn't set it)
		if strings.HasPrefix(
			line, token.HeadingLevelOneStart,
		) && entry.Title == "" {
			entry.Title = strings.TrimPrefix(line, token.HeadingLevelOneStart)
		}

		// Time from metadata
		if strings.HasPrefix(line, desc.Text(text.DescKeyLabelMetadataTime)) {
			entry.Time = strings.TrimSpace(
				strings.TrimPrefix(line, desc.Text(text.DescKeyLabelMetadataTime)),
			)
		}

		// Project from metadata
		if strings.HasPrefix(line, desc.Text(text.DescKeyLabelMetadataProject)) {
			entry.Project = strings.TrimSpace(
				strings.TrimPrefix(line, desc.Text(text.DescKeyLabelMetadataProject)),
			)
		}

		// Stop after we have all three
		if entry.Title != "" && entry.Time != "" && entry.Project != "" {
			break
		}
	}

	if entry.Title == "" {
		entry.Title = strings.TrimSuffix(filename, file.ExtMarkdown)
	}

	// Strip Claude Code internal markup tags from titles
	entry.Title = strings.TrimSpace(regex.SystemClaudeTag.ReplaceAllString(entry.Title, ""))

	// Sanitize characters that break Markdown link text: angle brackets
	// become HTML entities; backticks and # are stripped (they add no
	// meaning inside [...] link labels).
	entry.Title = strings.NewReplacer(
		"<", "&lt;", ">", "&gt;",
		"`", "", "#", "",
	).Replace(entry.Title)
	entry.Title = strings.TrimSpace(entry.Title)

	return entry
}
