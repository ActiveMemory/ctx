//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package summary

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/context/sanitize"
)

// Generate creates a brief summary for a context file based on its
// name and content.
//
// Parameters:
//   - name: Filename to determine summary strategy
//   - content: Raw file content to analyze
//
// Returns:
//   - string: Summary string (e.g., "3 active, 2 completed" or "empty")
func Generate(name string, content []byte) string {
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
			return desc.Text(text.DescKeySummaryEmpty)
		}
		return desc.Text(text.DescKeySummaryLoaded)
	}
}
