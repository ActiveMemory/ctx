//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"strings"
	"unicode/utf8"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// slugMaxLen is the maximum character length for a title-derived slug.
const slugMaxLen = 50

// slugifyTitle converts a human-readable title into a URL-friendly slug.
//
// Lowercases the input, replaces non-alphanumeric characters with hyphens,
// collapses consecutive hyphens, trims leading/trailing hyphens, and
// truncates on a word boundary at slugMaxLen characters.
//
// Parameters:
//   - title: Human-readable title string
//
// Returns:
//   - string: Slugified string (may be empty if input is empty or all punctuation)
func slugifyTitle(title string) string {
	// Strip the "..." truncation suffix from FirstUserMsg if present.
	title = strings.TrimSuffix(title, "...")

	var sb strings.Builder
	prevHyphen := false

	for _, r := range strings.ToLower(title) {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			sb.WriteRune(r)
			prevHyphen = false
		default:
			// Replace any non-alphanumeric character with a single hyphen.
			if !prevHyphen && sb.Len() > 0 {
				sb.WriteByte('-')
				prevHyphen = true
			}
		}
	}

	slug := strings.TrimRight(sb.String(), "-")

	if len(slug) <= slugMaxLen {
		return slug
	}

	// Truncate on a word (hyphen) boundary.
	truncated := slug[:slugMaxLen]
	if idx := strings.LastIndex(truncated, "-"); idx > 0 {
		truncated = truncated[:idx]
	}
	return truncated
}

// cleanTitle normalises a title for storage in YAML frontmatter.
//
// Replaces newlines, tabs and consecutive whitespace with single spaces,
// trims the result, and strips the "..." truncation suffix that
// parser.Session.FirstUserMsg may carry.
func cleanTitle(s string) string {
	s = strings.TrimSuffix(s, "...")
	s = config.RegExClaudeTag.ReplaceAllString(s, "")
	var sb strings.Builder
	prevSpace := false
	for _, r := range s {
		if r == '\n' || r == '\r' || r == '\t' {
			r = ' '
		}
		if r == ' ' {
			if !prevSpace && sb.Len() > 0 {
				sb.WriteRune(r)
			}
			prevSpace = true
			continue
		}
		sb.WriteRune(r)
		prevSpace = false
	}
	out := strings.TrimSpace(sb.String())

	// Truncate to RecallMaxTitleLen on a word boundary.
	if utf8.RuneCountInString(out) > config.RecallMaxTitleLen {
		runes := []rune(out)
		truncated := string(runes[:config.RecallMaxTitleLen])
		if idx := strings.LastIndex(truncated, " "); idx > 0 {
			truncated = truncated[:idx]
		}
		out = truncated
	}

	return out
}

// titleSlug returns the best available slug for a session, following a
// fallback hierarchy:
//
//  1. existingTitle — enriched title from previously exported frontmatter
//  2. s.FirstUserMsg — first user message text
//  3. s.Slug — Claude Code's random slug
//  4. s.ID[:8] — short ID prefix
//
// The chosen source (except s.Slug and s.ID[:8], which are already slugs)
// is passed through slugifyTitle.
//
// Parameters:
//   - s: Session to derive the slug from
//   - existingTitle: Title from enriched YAML frontmatter (may be empty)
//
// Returns:
//   - slug: URL-friendly slug for the filename
//   - title: Human-readable title for the H1 heading (empty when falling
//     back to s.Slug or s.ID)
func titleSlug(s *parser.Session, existingTitle string) (slug, title string) {
	if existingTitle != "" {
		clean := cleanTitle(existingTitle)
		sl := slugifyTitle(clean)
		if sl != "" {
			return sl, clean
		}
	}

	if s.FirstUserMsg != "" {
		clean := cleanTitle(s.FirstUserMsg)
		sl := slugifyTitle(clean)
		if sl != "" {
			return sl, clean
		}
	}

	if s.Slug != "" {
		return s.Slug, ""
	}

	short := s.ID
	if len(short) > config.RecallShortIDLen {
		short = short[:config.RecallShortIDLen]
	}
	return short, ""
}
