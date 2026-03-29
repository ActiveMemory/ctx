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
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
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
		return desc.Text(text.DescKeySummaryLoaded)
	}
	return fmt.Sprintf(desc.Text(text.DescKeySummaryInvariants), count)
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
		return desc.Text(text.DescKeySummaryEmpty)
	}

	var parts []string
	if active > 0 {
		activeFmt := desc.Text(text.DescKeySummaryActive)
		parts = append(parts, fmt.Sprintf(activeFmt, active))
	}
	if completed > 0 {
		completedFmt := desc.Text(text.DescKeySummaryCompleted)
		parts = append(parts, fmt.Sprintf(completedFmt, completed))
	}
	return strings.Join(parts, token.CommaSpace)
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
		return desc.Text(text.DescKeySummaryEmpty)
	}
	if count == 1 {
		return desc.Text(text.DescKeySummaryDecision)
	}
	return fmt.Sprintf(desc.Text(text.DescKeySummaryDecisions), count)
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
		return desc.Text(text.DescKeySummaryEmpty)
	}
	if count == 1 {
		return desc.Text(text.DescKeySummaryTerm)
	}
	return fmt.Sprintf(desc.Text(text.DescKeySummaryTerms), count)
}
