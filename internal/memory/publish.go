//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/memory"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errMemory "github.com/ActiveMemory/ctx/internal/err/memory"
	"github.com/ActiveMemory/ctx/internal/io"
)

// SelectContent reads .context/ files and selects content within the
// line budget.
//
// Priority order: tasks > decisions > conventions > learnings.
// If over budget, trims from bottom (learnings, conventions, decisions).
//
// Parameters:
//   - contextDir: path to the .context/ directory
//   - budget: maximum number of lines in the published block
//
// Returns:
//   - PublishResult: selected content with per-section slices
//   - error: non-nil if content selection fails
func SelectContent(contextDir string, budget int) (PublishResult, error) {
	var result PublishResult

	result.Tasks = extractPendingTasks(
		readContextFile(contextDir, ctx.Task), memory.PublishMaxTasks,
	)
	result.Decisions = extractRecentEntries(
		readContextFile(contextDir, ctx.Decision), memory.PublishMaxDecisions,
	)
	result.Conventions = extractConventionItems(
		readContextFile(contextDir, ctx.Convention), memory.PublishMaxConventions,
	)
	result.Learnings = extractRecentEntries(
		readContextFile(contextDir, ctx.Learning), memory.PublishMaxLearnings,
	)

	// Trim to budget (tasks always fit, trim from bottom)
	result.trimToBudget(budget)
	result.TotalLines = result.lineCount()

	return result, nil
}

// Format renders the publish result as a Markdown block (without markers).
//
// Returns:
//   - string: formatted Markdown with section headings and items
func (r *PublishResult) Format() string {
	var buf strings.Builder
	buf.WriteString(desc.Text(text.DescKeyMemoryPublishTitle))

	writeSection(&buf, text.DescKeyMemoryPublishTasks, r.Tasks, "")
	writeSection(&buf, text.DescKeyMemoryPublishDec,
		r.Decisions, token.PrefixListDash,
	)
	writeSection(&buf, text.DescKeyMemoryPublishConv, r.Conventions, "")
	writeSection(&buf, text.DescKeyMemoryPublishLrn,
		r.Learnings, token.PrefixListDash,
	)

	return strings.TrimRight(buf.String(), token.NewlineLF) + token.NewlineLF
}

// MergePublished inserts or replaces the marker block in existing
// MEMORY.md content.
//
// If markers exist, replaces everything between them. If markers are
// missing, appends the block at the end (recovery).
//
// Parameters:
//   - existing: current MEMORY.md content
//   - published: formatted publish block to insert
//
// Returns:
//   - string: merged content
//   - bool: true if markers were missing (appended instead of replaced)
func MergePublished(existing, published string) (string, bool) {
	block := marker.PublishMarkerStart + token.NewlineLF +
		published + marker.PublishMarkerEnd + token.NewlineLF

	startIdx := strings.Index(existing, marker.PublishMarkerStart)
	endIdx := strings.Index(existing, marker.PublishMarkerEnd)

	if startIdx >= 0 && endIdx > startIdx {
		// Replace the existing block
		before := existing[:startIdx]
		after := existing[endIdx+len(marker.PublishMarkerEnd):]
		// Trim trailing newline from after to avoid double blank lines
		after = strings.TrimPrefix(after, token.NewlineLF)
		return before + block + after, false
	}

	// Markers missing: append
	sep := token.NewlineLF
	if !strings.HasSuffix(existing, token.NewlineLF) {
		sep = token.NewlineLF + token.NewlineLF
	}
	return existing + sep + block, startIdx < 0
}

// RemovePublished strips the marker block from MEMORY.md content.
//
// Parameters:
//   - content: current MEMORY.md content
//
// Returns:
//   - string: content with the publish block removed
//   - bool: true if markers were found and removed
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

// Publish writes selected content to MEMORY.md with marker-based merge.
//
// Parameters:
//   - contextDir: path to the .context/ directory
//   - memoryPath: path to the MEMORY.md file
//   - budget: maximum number of lines in the published block
//
// Returns:
//   - PublishResult: the content that was published
//   - error: non-nil if selection or file write fails
func Publish(contextDir, memoryPath string, budget int) (PublishResult, error) {
	result, selectErr := SelectContent(contextDir, budget)
	if selectErr != nil {
		return PublishResult{}, errMemory.SelectContent(selectErr)
	}

	formatted := result.Format()

	existing, readErr := io.SafeReadUserFile(memoryPath)
	if readErr != nil {
		// MEMORY.md might not exist yet: create with just the block
		existing = []byte{}
	}

	merged, _ := MergePublished(string(existing), formatted)

	if writeErr := os.WriteFile(
		memoryPath, []byte(merged), fs.PermFile,
	); writeErr != nil {
		return PublishResult{}, errMemory.WriteMemory(writeErr)
	}

	return result, nil
}
