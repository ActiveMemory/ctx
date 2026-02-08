//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

const summaryEmpty = "empty"

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
		content, []byte(config.PrefixTaskUndone),
	) +
		bytes.Count(
			content, []byte(config.PrefixTaskDone),
		)
	if count == 0 {
		return "loaded"
	}
	return fmt.Sprintf("%d invariants", count)
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
	active := bytes.Count(content, []byte(config.PrefixTaskUndone))
	completed := bytes.Count(content, []byte(config.PrefixTaskDone))

	if active == 0 && completed == 0 {
		return summaryEmpty
	}

	var parts []string
	if active > 0 {
		parts = append(parts, fmt.Sprintf("%d active", active))
	}
	if completed > 0 {
		parts = append(parts, fmt.Sprintf("%d completed", completed))
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
	matches := config.RegExEntryHeading.FindAll(content, -1)
	count := len(matches)

	if count == 0 {
		return summaryEmpty
	}
	if count == 1 {
		return "1 decision"
	}
	return fmt.Sprintf("%d decisions", count)
}

// summarizeGlossary counts term definitions (**term**) in GLOSSARY.md.
//
// Parameters:
//   - content: Raw file content to analyze
//
// Returns:
//   - string: Summary like "5 terms" or "empty" if none
func summarizeGlossary(content []byte) string {
	matches := config.RegExGlossary.FindAll(content, -1)
	count := len(matches)

	if count == 0 {
		return summaryEmpty
	}
	if count == 1 {
		return "1 term"
	}
	return fmt.Sprintf("%d terms", count)
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
	case config.FileConstitution:
		return summarizeConstitution(content)
	case config.FileTask:
		return summarizeTasks(content)
	case config.FileDecision:
		return summarizeDecisions(content)
	case config.FileGlossary:
		return summarizeGlossary(content)
	default:
		if len(content) == 0 || effectivelyEmpty(content) {
			return summaryEmpty
		}
		return "loaded"
	}
}
