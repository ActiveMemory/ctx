//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgFmt "github.com/ActiveMemory/ctx/internal/config/format"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// Size formats a file size in human-readable form.
//
// Parameters:
//   - bytes: File size in bytes
//
// Returns:
//   - string: Human-readable size (e.g., "512B", "1.5KB", "2.3MB")
func Size(bytes int64) string {
	if bytes < cfgFmt.IECUnit {
		return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatBytes), bytes)
	}
	kb := float64(bytes) / cfgFmt.IECUnit
	if kb < cfgFmt.IECUnit {
		return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatKilobytes), kb)
	}
	mb := kb / cfgFmt.IECUnit
	return fmt.Sprintf(desc.Text(text.DescKeyWriteFormatMegabytes), mb)
}

// KeyFileSlug converts a file path to a safe slug for use as a filename.
//
// Replaces /, . with _ and * with x.
//
// Parameters:
//   - path: Original file path (e.g., "internal/config/*.go")
//
// Returns:
//   - string: Safe slug (e.g., "internal_config_x_go")
func KeyFileSlug(path string) string {
	slug := strings.ReplaceAll(path, token.Slash, token.Underscore)
	slug = strings.ReplaceAll(slug, token.Dot, token.Underscore)
	slug = strings.ReplaceAll(slug, token.GlobStar, token.GlobReplace)
	return slug
}

// SessionLink formats a Markdown list item linking to a page with a
// session count.
//
// Parameters:
//   - label: Display text for the link
//   - slug: Link target without .md extension
//   - count: Number of sessions to display
//
// Returns:
//   - string: Formatted line (e.g., "- [topic](topic.md) (3 sessions)\n")
func SessionLink(label, slug string, count int) string {
	return fmt.Sprintf(desc.Text(text.DescKeyJournalMocSessionLink),
		label, slug, file.ExtMarkdown, count, token.NewlineLF)
}
