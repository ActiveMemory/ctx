//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package summary

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/context/sanitize"
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
		return desc.TextDesc(text.DescKeySummaryLoaded)
	}
	return fmt.Sprintf(desc.TextDesc(text.DescKeySummaryInvariants), count)
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
		return desc.TextDesc(text.DescKeySummaryEmpty)
	}

	var parts []string
	if active > 0 {
		parts = append(parts, fmt.Sprintf(desc.TextDesc(text.DescKeySummaryActive), active))
	}
	if completed > 0 {
		parts = append(parts, fmt.Sprintf(desc.TextDesc(text.DescKeySummaryCompleted), completed))
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
		return desc.TextDesc(text.DescKeySummaryEmpty)
	}
	if count == 1 {
		return desc.TextDesc(text.DescKeySummaryDecision)
	}
	return fmt.Sprintf(desc.TextDesc(text.DescKeySummaryDecisions), count)
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
		return desc.TextDesc(text.DescKeySummaryEmpty)
	}
	if count == 1 {
		return desc.TextDesc(text.DescKeySummaryTerm)
	}
	return fmt.Sprintf(desc.TextDesc(text.DescKeySummaryTerms), count)
}

// GenerateSummary creates a brief summary for a context file based on its
// name and content.
//
// Parameters:
//   - name: Filename to determine summary strategy
//   - content: Raw file content to analyze
//
// Returns:
//   - string: Summary string (e.g., "3 active, 2 completed" or "empty")
func GenerateSummary(name string, content []byte) string {
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
		if len(content) == 0 || sanitize.EffectivelyEmpty(content) {
			return desc.TextDesc(text.DescKeySummaryEmpty)
		}
		return desc.TextDesc(text.DescKeySummaryLoaded)
	}
}
