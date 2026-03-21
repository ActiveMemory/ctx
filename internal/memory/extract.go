//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/memory"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/io"
)

// readContextFile reads a context file and returns its content as a string.
//
// Parameters:
//   - contextDir: path to the .context/ directory
//   - filename: name of the file to read
//
// Returns:
//   - string: file content, or empty string if the file cannot be read
func readContextFile(contextDir, filename string) string {
	data, readErr := io.SafeReadUserFile(filepath.Join(contextDir, filename))
	if readErr != nil {
		return ""
	}
	return string(data)
}

// writeSection appends a heading and items to buf if items is non-empty.
// When prefix is non-empty, each item is prefixed (e.g., "- ").
//
// Parameters:
//   - buf: string builder to append to
//   - headingKey: assets text description key for the section heading
//   - items: lines to write under the heading
//   - prefix: string prepended to each item (empty for none)
func writeSection(
	buf *strings.Builder, headingKey string, items []string, prefix string,
) {
	if len(items) == 0 {
		return
	}
	buf.WriteString(desc.Text(headingKey) + token.NewlineLF)
	for _, item := range items {
		buf.WriteString(prefix + item + token.NewlineLF)
	}
	buf.WriteString(token.NewlineLF)
}

// extractPendingTasks finds unchecked task items from TASKS.md content.
//
// Parameters:
//   - content: raw Markdown content of TASKS.md
//   - max: maximum number of tasks to return
//
// Returns:
//   - []string: unchecked task lines
func extractPendingTasks(content string, max int) []string {
	var tasks []string
	for _, line := range strings.Split(content, token.NewlineLF) {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, marker.PrefixTaskUndone+token.Space) {
			tasks = append(tasks, trimmed)
			if len(tasks) >= max {
				break
			}
		}
	}
	return tasks
}

// extractRecentEntries returns titles of entries from the last N days.
//
// Parameters:
//   - content: raw Markdown content with timestamped entry headers
//   - max: maximum number of entries to return
//
// Returns:
//   - []string: entry titles within the recency window
func extractRecentEntries(content string, max int) []string {
	blocks := index.ParseEntryBlocks(content)
	cutoff := time.Now().AddDate(
		0, 0, -memory.PublishRecentDays,
	).Format(cfgTime.DateFormat)

	var titles []string
	for _, b := range blocks {
		if b.Entry.Date >= cutoff {
			titles = append(titles, b.Entry.Title)
			if len(titles) >= max {
				break
			}
		}
	}
	return titles
}

// extractConventionItems returns the first N list items from CONVENTIONS.md.
//
// Parameters:
//   - content: raw Markdown content of CONVENTIONS.md
//   - max: maximum number of items to return
//
// Returns:
//   - []string: list items (dash or star prefixed)
func extractConventionItems(content string, max int) []string {
	var items []string
	for _, line := range strings.Split(content, token.NewlineLF) {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, token.PrefixListDash) ||
			strings.HasPrefix(trimmed, token.PrefixListStar) {
			items = append(items, trimmed)
			if len(items) >= max {
				break
			}
		}
	}
	return items
}

// lineCount returns the total number of lines the formatted publish
// block will occupy.
//
// Returns:
//   - int: line count including title, section headers, items, and blanks
func (r *PublishResult) lineCount() int {
	count := 1 // Title line
	if len(r.Tasks) > 0 {
		count += 2 + len(r.Tasks) // header + items + blank
	}
	if len(r.Decisions) > 0 {
		count += 2 + len(r.Decisions)
	}
	if len(r.Conventions) > 0 {
		count += 2 + len(r.Conventions)
	}
	if len(r.Learnings) > 0 {
		count += 2 + len(r.Learnings)
	}
	return count
}

// trimToBudget removes items from lowest-priority sections until
// lineCount fits within the budget.
//
// Trim order: learnings, then conventions, then decisions.
// Tasks are never trimmed.
//
// Parameters:
//   - budget: maximum allowed line count
func (r *PublishResult) trimToBudget(budget int) {
	for r.lineCount() > budget && len(r.Learnings) > 0 {
		r.Learnings = r.Learnings[:len(r.Learnings)-1]
	}
	for r.lineCount() > budget && len(r.Conventions) > 0 {
		r.Conventions = r.Conventions[:len(r.Conventions)-1]
	}
	for r.lineCount() > budget && len(r.Decisions) > 0 {
		r.Decisions = r.Decisions[:len(r.Decisions)-1]
	}
}

// extractTitle returns the first meaningful line of an entry, cleaned
// of Markdown heading markers and list item prefixes.
//
// Parameters:
//   - text: raw entry text
//
// Returns:
//   - string: cleaned first line
func extractTitle(text string) string {
	line := strings.SplitN(text, token.NewlineLF, 2)[0]
	line = strings.TrimSpace(line)
	// Strip heading markers
	line = strings.TrimLeft(line, token.PrefixHeading)
	line = strings.TrimSpace(line)
	// Strip list item markers
	if strings.HasPrefix(line, token.PrefixListDash) {
		line = line[len(token.PrefixListDash):]
	} else if strings.HasPrefix(line, token.PrefixListStar) {
		line = line[len(token.PrefixListStar):]
	}
	return strings.TrimSpace(line)
}

// extractBody returns everything after the first line, or the first
// line itself if there is only one line.
//
// Parameters:
//   - text: raw entry text
//
// Returns:
//   - string: body content after the title
func extractBody(text string) string {
	parts := strings.SplitN(text, token.NewlineLF, 2)
	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return extractTitle(text)
	}
	return strings.TrimSpace(parts[1])
}
