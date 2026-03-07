//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	ctxentry "github.com/ActiveMemory/ctx/internal/entry"
)

const importSource = "auto-memory import"

// Promote writes a classified entry to the appropriate .context/ file.
// Uses the add package's WriteEntry for consistent formatting and indexing.
func Promote(entry Entry, classification Classification) error {
	// Extract a title from the entry text (first line, trimmed of Markdown markers)
	title := extractTitle(entry.Text)

	params := ctxentry.Params{
		Type:    classification.Target,
		Content: title,
	}

	switch classification.Target {
	case config.EntryDecision:
		params.Context = importSource
		params.Rationale = extractBody(entry.Text)
		params.Consequences = "Imported from MEMORY.md — review and update as needed"

	case config.EntryLearning:
		params.Context = importSource
		params.Lesson = extractBody(entry.Text)
		params.Application = "Imported from MEMORY.md — review and update as needed"

	case config.EntryTask:
		// Tasks just need content — FormatTask handles the rest

	case config.EntryConvention:
		// Conventions just need content — FormatConvention handles the rest
	}

	return ctxentry.Write(params)
}

// extractTitle returns the first meaningful line of an entry, cleaned of
// Markdown heading markers and list item prefixes.
func extractTitle(text string) string {
	line := strings.SplitN(text, config.NewlineLF, 2)[0]
	line = strings.TrimSpace(line)
	// Strip heading markers
	line = strings.TrimLeft(line, "#")
	line = strings.TrimSpace(line)
	// Strip list item markers
	if strings.HasPrefix(line, "- ") {
		line = line[2:]
	} else if strings.HasPrefix(line, "* ") {
		line = line[2:]
	}
	return strings.TrimSpace(line)
}

// extractBody returns everything after the first line, or the first line
// itself if there's only one line.
func extractBody(text string) string {
	parts := strings.SplitN(text, config.NewlineLF, 2)
	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return extractTitle(text)
	}
	return strings.TrimSpace(parts[1])
}
