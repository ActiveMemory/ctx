//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// FormatSize formats a file size in human-readable form.
//
// Parameters:
//   - bytes: File size in bytes
//
// Returns:
//   - string: Human-readable size (e.g., "512B", "1.5KB", "2.3MB")
func FormatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	}
	kb := float64(bytes) / 1024
	if kb < 1024 {
		return fmt.Sprintf("%.1fKB", kb)
	}
	mb := kb / 1024
	return fmt.Sprintf("%.1fMB", mb)
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
	slug := strings.ReplaceAll(path, "/", "_")
	slug = strings.ReplaceAll(slug, ".", "_")
	slug = strings.ReplaceAll(slug, "*", "x")
	return slug
}

// FormatSessionLink formats a Markdown list item linking to a page with a
// session count.
//
// Parameters:
//   - label: Display text for the link
//   - slug: Link target without .md extension
//   - count: Number of sessions to display
//
// Returns:
//   - string: Formatted line (e.g., "- [topic](topic.md) (3 sessions)\n")
func FormatSessionLink(label, slug string, count int) string {
	return fmt.Sprintf(assets.TextDesc(assets.TextDescKeyJournalMocSessionLink),
		label, slug, file.ExtMarkdown, count, token.NewlineLF)
}
