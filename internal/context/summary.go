//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
)

// summarizeConstitution counts checkbox items (invariants) in CONSTITUTION.md.
//
// Parameters:
//   - content: Raw file content to analyze
//
// Returns:
//   - string: Summary like "5 invariants" or "loaded" if none found
func summarizeConstitution(content []byte) string {
	// Count checkbox items (invariants)
	count := bytes.Count(
		content, []byte(marker.PrefixTaskUndone),
	) +
		bytes.Count(
			content, []byte(marker.PrefixTaskDone),
		)
	if count == 0 {
		return assets.TextDesc(assets.TextDescKeySummaryLoaded)
	}
	return fmt.Sprintf(assets.TextDesc(assets.TextDescKeySummaryInvariants), count)
}

// summarizeTasks counts active and completed tasks in TASKS.md.
//
// Parameters:
//   - content: Raw file content to analyze
//
// Returns:
//   - string: Summary like "3 active, 2 completed" or "empty" if none
func summarizeTasks(content []byte) string {
	// Count active (unchecked) and completed (checked) tasks
	active := bytes.Count(content, []byte(marker.PrefixTaskUndone))
	completed := bytes.Count(content, []byte(marker.PrefixTaskDone))

	if active == 0 && completed == 0 {
		return assets.TextDesc(assets.TextDescKeySummaryEmpty)
	}

	var parts []string
	if active > 0 {
		parts = append(parts, fmt.Sprintf(assets.TextDesc(assets.TextDescKeySummaryActive), active))
	}
	if completed > 0 {
		parts = append(parts, fmt.Sprintf(assets.TextDesc(assets.TextDescKeySummaryCompleted), completed))
	}
	return strings.Join(parts, ", ")
}

// summarizeDecisions counts decision headers (## sections) in DECISIONS.md.
//
// Parameters:
//   - content: Raw file content to analyze
//
// Returns:
//   - string: Summary like "3 decisions" or "empty" if none
func summarizeDecisions(content []byte) string {
	// Count decision headers (## [date] or ## Decision)
	matches := regex.EntryHeading.FindAll(content, -1)
	count := len(matches)

	if count == 0 {
		return assets.TextDesc(assets.TextDescKeySummaryEmpty)
	}
	if count == 1 {
		return assets.TextDesc(assets.TextDescKeySummaryDecision)
	}
	return fmt.Sprintf(assets.TextDesc(assets.TextDescKeySummaryDecisions), count)
}

// summarizeGlossary counts term definitions (**term**) in GLOSSARY.md.
//
// Parameters:
//   - content: Raw file content to analyze
//
// Returns:
//   - string: Summary like "5 terms" or "empty" if none
func summarizeGlossary(content []byte) string {
	matches := regex.Glossary.FindAll(content, -1)
	count := len(matches)

	if count == 0 {
		return assets.TextDesc(assets.TextDescKeySummaryEmpty)
	}
	if count == 1 {
		return assets.TextDesc(assets.TextDescKeySummaryTerm)
	}
	return fmt.Sprintf(assets.TextDesc(assets.TextDescKeySummaryTerms), count)
}

// generateSummary creates a brief summary for a context file based on its
// name and content.
//
// Parameters:
//   - name: Filename to determine summary strategy
//   - content: Raw file content to analyze
//
// Returns:
//   - string: Summary string (e.g., "3 active, 2 completed" or "empty")
func generateSummary(name string, content []byte) string {
	switch name {
	case ctx.Constitution:
		return summarizeConstitution(content)
	case ctx.Task:
		return summarizeTasks(content)
	case ctx.Decision:
		return summarizeDecisions(content)
	case ctx.Glossary:
		return summarizeGlossary(content)
	default:
		if len(content) == 0 || effectivelyEmpty(content) {
			return assets.TextDesc(assets.TextDescKeySummaryEmpty)
		}
		return assets.TextDesc(assets.TextDescKeySummaryLoaded)
	}
}
