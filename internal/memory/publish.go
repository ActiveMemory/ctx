//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
)

const (
	// MarkerStart is the HTML comment that begins the ctx-published block.
	MarkerStart = "<!-- ctx:published -->"
	// MarkerEnd is the HTML comment that ends the ctx-published block.
	MarkerEnd = "<!-- ctx:end -->"

	// DefaultPublishBudget is the default line budget for published content.
	DefaultPublishBudget = 80

	maxTasks       = 10
	maxDecisions   = 5
	maxConventions = 10
	maxLearnings   = 5
	recentDays     = 7
)

// PublishResult holds what was selected for publishing.
type PublishResult struct {
	Tasks       []string
	Decisions   []string
	Conventions []string
	Learnings   []string
	TotalLines  int
}

// SelectContent reads .context/ files and selects content within the line budget.
//
// Priority order: tasks > decisions > conventions > learnings.
// If over budget, trims from bottom (learnings → conventions → decisions).
func SelectContent(contextDir string, budget int) (PublishResult, error) {
	var result PublishResult

	// Pending tasks
	taskPath := filepath.Join(contextDir, config.FileTask)
	if data, readErr := os.ReadFile(taskPath); readErr == nil { //nolint:gosec // project-local path
		result.Tasks = extractPendingTasks(string(data), maxTasks)
	}

	// Recent decisions
	decPath := filepath.Join(contextDir, config.FileDecision)
	if data, readErr := os.ReadFile(decPath); readErr == nil { //nolint:gosec // project-local path
		result.Decisions = extractRecentEntries(string(data), maxDecisions)
	}

	// Key conventions (first N lines that are list items)
	convPath := filepath.Join(contextDir, config.FileConvention)
	if data, readErr := os.ReadFile(convPath); readErr == nil { //nolint:gosec // project-local path
		result.Conventions = extractConventionItems(string(data), maxConventions)
	}

	// Recent learnings
	lrnPath := filepath.Join(contextDir, config.FileLearning)
	if data, readErr := os.ReadFile(lrnPath); readErr == nil { //nolint:gosec // project-local path
		result.Learnings = extractRecentEntries(string(data), maxLearnings)
	}

	// Trim to budget (tasks always fit, trim from bottom)
	result.trimToBudget(budget)
	result.TotalLines = result.lineCount()

	return result, nil
}

// Format renders the publish result as a Markdown block (without markers).
func (r PublishResult) Format() string {
	var buf strings.Builder
	buf.WriteString("# Project Context (managed by ctx)\n\n")

	if len(r.Tasks) > 0 {
		buf.WriteString("## Pending Tasks\n")
		for _, t := range r.Tasks {
			buf.WriteString(t + "\n")
		}
		buf.WriteString("\n")
	}

	if len(r.Decisions) > 0 {
		buf.WriteString("## Recent Decisions\n")
		for _, d := range r.Decisions {
			buf.WriteString("- " + d + "\n")
		}
		buf.WriteString("\n")
	}

	if len(r.Conventions) > 0 {
		buf.WriteString("## Key Conventions\n")
		for _, c := range r.Conventions {
			buf.WriteString(c + "\n")
		}
		buf.WriteString("\n")
	}

	if len(r.Learnings) > 0 {
		buf.WriteString("## Recent Learnings\n")
		for _, l := range r.Learnings {
			buf.WriteString("- " + l + "\n")
		}
		buf.WriteString("\n")
	}

	return strings.TrimRight(buf.String(), "\n") + "\n"
}

// MergePublished inserts or replaces the marker block in existing MEMORY.md content.
//
// If markers exist, replaces everything between them. If markers are missing,
// appends the block at the end (recovery). Returns (merged content, markers were missing).
func MergePublished(existing, published string) (string, bool) {
	block := MarkerStart + "\n" + published + MarkerEnd + "\n"

	startIdx := strings.Index(existing, MarkerStart)
	endIdx := strings.Index(existing, MarkerEnd)

	if startIdx >= 0 && endIdx > startIdx {
		// Replace existing block
		before := existing[:startIdx]
		after := existing[endIdx+len(MarkerEnd):]
		// Trim trailing newline from after to avoid double blank lines
		after = strings.TrimPrefix(after, "\n")
		return before + block + after, false
	}

	// Markers missing — append
	sep := "\n"
	if !strings.HasSuffix(existing, "\n") {
		sep = "\n\n"
	}
	return existing + sep + block, startIdx < 0
}

// RemovePublished strips the marker block from MEMORY.md content.
// Returns (cleaned content, true if markers were found and removed).
func RemovePublished(content string) (string, bool) {
	startIdx := strings.Index(content, MarkerStart)
	endIdx := strings.Index(content, MarkerEnd)

	if startIdx < 0 || endIdx <= startIdx {
		return content, false
	}

	before := content[:startIdx]
	after := content[endIdx+len(MarkerEnd):]
	after = strings.TrimPrefix(after, "\n")

	result := strings.TrimRight(before, "\n")
	if after != "" {
		result += "\n" + after
	} else {
		result += "\n"
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
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- [ ] ") {
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
	cutoff := time.Now().AddDate(0, 0, -recentDays).Format("2006-01-02")

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
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
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
		return PublishResult{}, fmt.Errorf("selecting content: %w", selectErr)
	}

	formatted := result.Format()

	existing, readErr := os.ReadFile(memoryPath) //nolint:gosec // caller-provided path
	if readErr != nil {
		// MEMORY.md might not exist yet — create with just the block
		existing = []byte{}
	}

	merged, _ := MergePublished(string(existing), formatted)

	if writeErr := os.WriteFile(memoryPath, []byte(merged), config.PermFile); writeErr != nil {
		return PublishResult{}, fmt.Errorf("writing MEMORY.md: %w", writeErr)
	}

	return result, nil
}
