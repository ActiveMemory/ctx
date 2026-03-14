//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/memory"
	time2 "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/io"
)

// SelectContent reads .context/ files and selects content within the line budget.
//
// Priority order: tasks > decisions > conventions > learnings.
// If over budget, trims from bottom (learnings → conventions → decisions).
func SelectContent(contextDir string, budget int) (PublishResult, error) {
	var result PublishResult

	// Pending tasks
	taskPath := filepath.Join(contextDir, ctx.Task)
	if data, readErr := io.SafeReadUserFile(taskPath); readErr == nil {
		result.Tasks = extractPendingTasks(string(data), memory.PublishMaxTasks)
	}

	// Recent decisions
	decPath := filepath.Join(contextDir, ctx.Decision)
	if data, readErr := io.SafeReadUserFile(decPath); readErr == nil {
		result.Decisions = extractRecentEntries(string(data), memory.PublishMaxDecisions)
	}

	// Key conventions (first N lines that are list items)
	convPath := filepath.Join(contextDir, ctx.Convention)
	if data, readErr := io.SafeReadUserFile(convPath); readErr == nil {
		result.Conventions = extractConventionItems(string(data), memory.PublishMaxConventions)
	}

	// Recent learnings
	lrnPath := filepath.Join(contextDir, ctx.Learning)
	if data, readErr := io.SafeReadUserFile(lrnPath); readErr == nil {
		result.Learnings = extractRecentEntries(string(data), memory.PublishMaxLearnings)
	}

	// Trim to budget (tasks always fit, trim from bottom)
	result.trimToBudget(budget)
	result.TotalLines = result.lineCount()

	return result, nil
}

// Format renders the publish result as a Markdown block (without markers).
func (r PublishResult) Format() string {
	var buf strings.Builder
	buf.WriteString(assets.TextDesc(assets.TextDescKeyMemoryPublishTitle))

	if len(r.Tasks) > 0 {
		buf.WriteString(assets.TextDesc(assets.TextDescKeyMemoryPublishTasks) + token.NewlineLF)
		for _, t := range r.Tasks {
			buf.WriteString(t + token.NewlineLF)
		}
		buf.WriteString(token.NewlineLF)
	}

	if len(r.Decisions) > 0 {
		buf.WriteString(assets.TextDesc(assets.TextDescKeyMemoryPublishDec) + token.NewlineLF)
		for _, d := range r.Decisions {
			buf.WriteString(token.PrefixListDash + d + token.NewlineLF)
		}
		buf.WriteString(token.NewlineLF)
	}

	if len(r.Conventions) > 0 {
		buf.WriteString(assets.TextDesc(assets.TextDescKeyMemoryPublishConv) + token.NewlineLF)
		for _, c := range r.Conventions {
			buf.WriteString(c + token.NewlineLF)
		}
		buf.WriteString(token.NewlineLF)
	}

	if len(r.Learnings) > 0 {
		buf.WriteString(assets.TextDesc(assets.TextDescKeyMemoryPublishLrn) + token.NewlineLF)
		for _, l := range r.Learnings {
			buf.WriteString(token.PrefixListDash + l + token.NewlineLF)
		}
		buf.WriteString(token.NewlineLF)
	}

	return strings.TrimRight(buf.String(), token.NewlineLF) + token.NewlineLF
}

// MergePublished inserts or replaces the marker block in existing MEMORY.md content.
//
// If markers exist, replaces everything between them. If markers are missing,
// appends the block at the end (recovery). Returns (merged content, markers were missing).
func MergePublished(existing, published string) (string, bool) {
	block := marker.PublishMarkerStart + token.NewlineLF + published + marker.PublishMarkerEnd + token.NewlineLF

	startIdx := strings.Index(existing, marker.PublishMarkerStart)
	endIdx := strings.Index(existing, marker.PublishMarkerEnd)

	if startIdx >= 0 && endIdx > startIdx {
		// Replace existing block
		before := existing[:startIdx]
		after := existing[endIdx+len(marker.PublishMarkerEnd):]
		// Trim trailing newline from after to avoid double blank lines
		after = strings.TrimPrefix(after, token.NewlineLF)
		return before + block + after, false
	}

	// Markers missing — append
	sep := token.NewlineLF
	if !strings.HasSuffix(existing, token.NewlineLF) {
		sep = token.NewlineLF + token.NewlineLF
	}
	return existing + sep + block, startIdx < 0
}

// RemovePublished strips the marker block from MEMORY.md content.
// Returns (cleaned content, true if markers were found and removed).
func RemovePublished(content string) (string, bool) {
	startIdx := strings.Index(content, marker.PublishMarkerStart)
	endIdx := strings.Index(content, marker.PublishMarkerEnd)

	if startIdx < 0 || endIdx <= startIdx {
		return content, false
	}

	before := content[:startIdx]
	after := content[endIdx+len(marker.PublishMarkerEnd):]
	after = strings.TrimPrefix(after, token.NewlineLF)

	result := strings.TrimRight(before, token.NewlineLF)
	if after != "" {
		result += token.NewlineLF + after
	} else {
		result += token.NewlineLF
	}

	return result, true
}

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

// extractPendingTasks finds unchecked task items from TASKS.md.
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
func extractRecentEntries(content string, max int) []string {
	blocks := index.ParseEntryBlocks(content)
	cutoff := time.Now().AddDate(0, 0, -memory.PublishRecentDays).Format(time2.DateFormat)

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
func extractConventionItems(content string, max int) []string {
	var items []string
	for _, line := range strings.Split(content, token.NewlineLF) {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, token.PrefixListDash) || strings.HasPrefix(trimmed, token.PrefixListStar) {
			items = append(items, trimmed)
			if len(items) >= max {
				break
			}
		}
	}
	return items
}

// Publish writes selected content to MEMORY.md with marker-based merge.
func Publish(contextDir, memoryPath string, budget int) (PublishResult, error) {
	result, selectErr := SelectContent(contextDir, budget)
	if selectErr != nil {
		return PublishResult{}, ctxerr.MemorySelectContent(selectErr)
	}

	formatted := result.Format()

	existing, readErr := io.SafeReadUserFile(memoryPath)
	if readErr != nil {
		// MEMORY.md might not exist yet — create with just the block
		existing = []byte{}
	}

	merged, _ := MergePublished(string(existing), formatted)

	if writeErr := os.WriteFile(memoryPath, []byte(merged), fs.PermFile); writeErr != nil {
		return PublishResult{}, ctxerr.MemoryWriteMemory(writeErr)
	}

	return result, nil
}
