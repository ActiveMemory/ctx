//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxentry "github.com/ActiveMemory/ctx/internal/entry"
)

// Promote writes a classified entry to the appropriate .context/ file.
// Uses the add package's WriteEntry for consistent formatting and indexing.
func Promote(e Entry, classification Classification) error {
	// Extract a title from the entry text (first line, trimmed of Markdown markers)
	title := extractTitle(e.Text)

	params := ctxentry.Params{
		Type:    classification.Target,
		Content: title,
	}

	switch classification.Target {
	case entry.Decision:
		params.Context = assets.TextDesc(assets.TextDescKeyMemoryImportSource)
		params.Rationale = extractBody(e.Text)
		params.Consequences = assets.TextDesc(assets.TextDescKeyMemoryImportReview)

	case entry.Learning:
		params.Context = assets.TextDesc(assets.TextDescKeyMemoryImportSource)
		params.Lesson = extractBody(e.Text)
		params.Application = assets.TextDesc(assets.TextDescKeyMemoryImportReview)

	case entry.Task:
		// Tasks just need content — FormatTask handles the rest

	case entry.Convention:
		// Conventions just need content — FormatConvention handles the rest
	}

	return ctxentry.Write(params)
}

// extractTitle returns the first meaningful line of an entry, cleaned of
// Markdown heading markers and list item prefixes.
func extractTitle(text string) string {
	line := strings.SplitN(text, token.NewlineLF, 2)[0]
	line = strings.TrimSpace(line)
	// Strip heading markers
	line = strings.TrimLeft(line, token.PrefixHeading)
	line = strings.TrimSpace(line)
	// Strip list item markers
	if strings.HasPrefix(line, token.PrefixListDash) {
		line = line[len(token.PrefixListDash):]
	} else if strings.HasPrefix(line, token.PrefixListStar) {
		line = line[len(token.PrefixListStar):]
	}
	return strings.TrimSpace(line)
}

// extractBody returns everything after the first line, or the first line
// itself if there's only one line.
func extractBody(text string) string {
	parts := strings.SplitN(text, token.NewlineLF, 2)
	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return extractTitle(text)
	}
	return strings.TrimSpace(parts[1])
}
