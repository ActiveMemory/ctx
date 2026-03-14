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
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
)

// WriteSection creates a subdirectory, writes its index page, and calls
// writePages to emit individual pages. All three index sections (topics,
// files, types) follow this identical structure.
//
// Parameters:
//   - docsDir: Parent docs directory
//   - subdir: Subdirectory name (e.g., config.JournalDirTopics)
//   - indexContent: Generated Markdown for the index page
//   - writePages: Callback that writes individual pages into the subdirectory
//
// Returns:
//   - error: Non-nil if directory creation or index write fails
func WriteSection(
	docsDir, subdir, indexContent string,
	writePages func(dir string),
) error {
	dir := filepath.Join(docsDir, subdir)
	if mkErr := os.MkdirAll(dir, fs.PermExec); mkErr != nil {
		return ctxerr.Mkdir(dir, mkErr)
	}

	indexPath := filepath.Join(dir, file.Index)
	if writeErr := os.WriteFile(
		indexPath, []byte(indexContent), fs.PermFile,
	); writeErr != nil {
		return ctxerr.FileWrite(indexPath, writeErr)
	}

	writePages(dir)
	return nil
}

// WriteMonthSections writes month-grouped entry links to a string builder.
//
// Parameters:
//   - sb: String builder to write to
//   - months: Entries keyed by month string (YYYY-MM)
//   - monthOrder: Month strings in display order
//   - linkPrefix: Path prefix for links (e.g., config.LinkPrefixParent for
//     subpages, "" for index)
func WriteMonthSections(
	sb *strings.Builder,
	months map[string][]JournalEntry,
	monthOrder []string, linkPrefix string,
) {
	nl := token.NewlineLF
	for _, month := range monthOrder {
		_, _ = fmt.Fprintf(sb, assets.TplJournalMonthHeading+nl+nl, month)
		for _, e := range months[month] {
			link := strings.TrimSuffix(e.Filename, file.ExtMarkdown)
			timeStr := ""
			if e.Time != "" && len(e.Time) >= journal.TimePrefixLen {
				timeStr = e.Time[:journal.TimePrefixLen] + " "
			}
			_, _ = fmt.Fprintf(sb,
				assets.TplJournalSubpageEntry+nl,
				timeStr, e.Title, linkPrefix, link)
			if e.Summary != "" {
				_, _ = fmt.Fprintf(sb, assets.TplJournalIndexSummary+nl, e.Summary)
			}
		}
		sb.WriteString(nl)
	}
}

// GenerateGroupedPage builds a detail page with a heading, stats line, and
// month-grouped session links. Used by topic, key file, and type pages.
//
// Parameters:
//   - heading: Pre-formatted Markdown heading (e.g., "# refactoring")
//   - stats: Pre-formatted stats line (e.g., "**5 sessions** with this topic.")
//   - entries: Journal entries to group by month
//
// Returns:
//   - string: Complete Markdown page content
func GenerateGroupedPage(heading, stats string, entries []JournalEntry) string {
	var sb strings.Builder
	nl := token.NewlineLF

	sb.WriteString(heading + nl + nl)
	sb.WriteString(stats + nl + nl)

	months, monthOrder := GroupByMonth(entries)
	WriteMonthSections(&sb, months, monthOrder, token.LinkPrefixParent)

	return sb.String()
}

// WritePopularAndLongtail writes the popular and longtail sections of an
// index page. Popular items link to dedicated pages; longtail items show
// inline links to the first matching session.
//
// Parameters:
//   - sb: String builder to write to
//   - popCount: Number of popular items
//   - popHeading: Section heading for popular items
//   - popItem: Returns label, slug, and entry count for popular item at index i
//   - ltCount: Number of longtail items
//   - ltHeading: Section heading for longtail items
//   - ltTpl: Format template for longtail entries
//   - ltItem: Returns label and first entry for longtail item at index i
func WritePopularAndLongtail(
	sb *strings.Builder,
	popCount int, popHeading string,
	popItem func(int) (string, string, int),
	ltCount int, ltHeading, ltTpl string,
	ltItem func(int) (string, JournalEntry),
) {
	nl := token.NewlineLF

	if popCount > 0 {
		sb.WriteString(popHeading + nl + nl)
		for i := range popCount {
			label, slug, count := popItem(i)
			sb.WriteString(FormatSessionLink(label, slug, count))
		}
		sb.WriteString(nl)
	}

	if ltCount > 0 {
		sb.WriteString(ltHeading + nl + nl)
		for i := range ltCount {
			label, e := ltItem(i)
			link := strings.TrimSuffix(e.Filename, file.ExtMarkdown)
			_, _ = fmt.Fprintf(sb, ltTpl+nl, label, e.Title, link)
		}
		sb.WriteString(nl)
	}
}

// ContinuesMultipart reports whether the filename is a continuation part
// (p2, p3, etc.) of a multipart session.
//
// Parameters:
//   - filename: Journal entry filename to check
//
// Returns:
//   - bool: True if the filename matches the multipart continuation pattern
func ContinuesMultipart(filename string) bool {
	return regex.MultiPart.MatchString(filename)
}
