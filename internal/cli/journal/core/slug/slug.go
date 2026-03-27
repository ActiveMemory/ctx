//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package slug

import (
	"strings"
	"unicode/utf8"

	"github.com/ActiveMemory/ctx/internal/config/journal"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// SlugifyTitle converts a human-readable title into a URL-friendly slug.
//
// Lowercases the input, replaces non-alphanumeric characters with hyphens,
// collapses consecutive hyphens, trims leading/trailing hyphens, and
// truncates on a word boundary at journal.TitleSlugMaxLen characters.
//
// Parameters:
//   - title: Human-readable title string
//
// Returns:
//   - string: Slugified string (may be empty if input is empty or all punctuation)
func SlugifyTitle(title string) string {
	// Strip the "..." truncation suffix from FirstUserMsg if present.
	title = strings.TrimSuffix(title, token.Ellipsis)

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
				sb.WriteString(token.Dash)
				prevHyphen = true
			}
		}
	}

	slug := strings.TrimRight(sb.String(), token.Dash)

	if len(slug) <= journal.TitleSlugMaxLen {
		return slug
	}

	// Truncate on a word (hyphen) boundary.
	truncated := slug[:journal.TitleSlugMaxLen]
	if idx := strings.LastIndex(truncated, token.Dash); idx > 0 {
		truncated = truncated[:idx]
	}
	return truncated
}

// CleanTitle normalises a title for storage in YAML frontmatter.
//
// Replaces newlines, tabs, and consecutive whitespace with single spaces,
// trims the result, and strips the "..." truncation suffix that
// entity.Session.FirstUserMsg may carry.
//
// Parameters:
//   - s: Raw title string
//
// Returns:
//   - string: Cleaned title string
func CleanTitle(s string) string {
	s = strings.TrimSuffix(s, token.Ellipsis)
	s = regex.SystemClaudeTag.ReplaceAllString(s, "")
	var sb strings.Builder
	prevSpace := false
	for _, r := range s {
		if r == rune(token.NewlineLF[0]) ||
			r == rune(token.NewlineCRLF[0]) || r == rune(token.Tab[0]) {
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

	// Truncate to MaxTitleLen on a word boundary.
	if utf8.RuneCountInString(out) > journal.MaxTitleLen {
		runes := []rune(out)
		truncated := string(runes[:journal.MaxTitleLen])
		if idx := strings.LastIndex(truncated, " "); idx > 0 {
			truncated = truncated[:idx]
		}
		out = truncated
	}

	return out
}

// TitleSlug returns the best available slug for a session, following a
// fallback hierarchy:
//
//  1. existingTitle: enriched title from previously exported frontmatter
//  2. s.FirstUserMsg: first user message text
//  3. s.Slug: Claude Code's random slug
//  4. s.ID[:8]: short ID prefix
//
// The chosen source (except s.Slug and s.ID[:8], which are already slugs)
// is passed through SlugifyTitle.
//
// Parameters:
//   - s: Session to derive the slug from
//   - existingTitle: Title from enriched YAML frontmatter (may be empty)
//
// Returns:
//   - slug: URL-friendly slug for the filename
//   - title: Human-readable title for the H1 heading (empty when falling
//     back to s.Slug or s.ID)
func TitleSlug(s *entity.Session, existingTitle string) (slug, title string) {
	if existingTitle != "" {
		clean := CleanTitle(existingTitle)
		sl := SlugifyTitle(clean)
		if sl != "" {
			return sl, clean
		}
	}

	if s.FirstUserMsg != "" {
		clean := CleanTitle(s.FirstUserMsg)
		sl := SlugifyTitle(clean)
		if sl != "" {
			return sl, clean
		}
	}

	if s.Slug != "" {
		return s.Slug, ""
	}

	short := s.ID
	if len(short) > journal.ShortIDLen {
		short = short[:journal.ShortIDLen]
	}
	return short, ""
}
