//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover

import (
	"strings"
	"time"

	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	cfgHandover "github.com/ActiveMemory/ctx/internal/config/handover"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// buildFilename derives a handover's on-disk name. Shape:
// `<TS>-<slug>.md` where `<TS>` is the UTC compact RFC-3339
// form (`20260517T021837Z`; colons stripped) and `<slug>` is
// the kebab-case normalisation of the caller's title.
//
// Parameters:
//   - now: GeneratedAt timestamp.
//   - title: caller-supplied title.
//
// Returns:
//   - string: filename portion only (no directory).
func buildFilename(now time.Time, title string) string {
	slug := titleToSlug(title)
	if slug == "" {
		slug = cfgHandover.DefaultSlug
	}
	var sb strings.Builder
	sb.WriteString(now.UTC().Format(cfgTime.RFC3339Compact))
	sb.WriteString(token.Dash)
	sb.WriteString(slug)
	sb.WriteString(cfgFile.ExtMarkdown)
	return sb.String()
}

// titleToSlug normalises a free-text title into a kebab-case
// slug-safe form.
//
// Parameters:
//   - s: free text.
//
// Returns:
//   - string: lowercase, hyphen-separated, with non-alnum
//     characters stripped and leading / trailing hyphens
//     trimmed.
func titleToSlug(s string) string {
	low := strings.ToLower(strings.TrimSpace(s))
	low = strings.ReplaceAll(low, token.Space, token.Dash)
	low = regex.Slug.ReplaceAllString(low, "")
	low = strings.Trim(low, token.Dash)
	return low
}
