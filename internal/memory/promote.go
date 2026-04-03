//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/entity"
	ctxEntry "github.com/ActiveMemory/ctx/internal/entry"
)

// Promote writes a classified entry to the appropriate .context/ file.
//
// Uses the entry package for consistent formatting and indexing.
//
// Parameters:
//   - e: the memory entry to promote
//   - classification: target type and confidence from classification
//
// Returns:
//   - error: non-nil if the entry write fails
func Promote(e Entry, classification Classification) error {
	// Extract a title from the entry text
	// (first line, trimmed of Markdown markers)
	title := extractTitle(e.Text)

	params := entity.EntryParams{
		Type:    classification.Target,
		Content: title,
	}

	switch classification.Target {
	case entry.Decision:
		params.Context = desc.Text(text.DescKeyMemoryImportSource)
		params.Rationale = extractBody(e.Text)
		params.Consequence = desc.Text(text.DescKeyMemoryImportReview)

	case entry.Learning:
		params.Context = desc.Text(text.DescKeyMemoryImportSource)
		params.Lesson = extractBody(e.Text)
		params.Application = desc.Text(text.DescKeyMemoryImportReview)

	case entry.Task:
		// Tasks just need content: FormatTask handles the rest

	case entry.Convention:
		// Conventions just need content: FormatConvention handles the rest
	}

	return ctxEntry.Write(params)
}
