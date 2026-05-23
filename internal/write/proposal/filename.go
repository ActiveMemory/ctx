//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package proposal

import (
	"strings"
	"time"

	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	cfgProposal "github.com/ActiveMemory/ctx/internal/config/proposal"
	cfgRegex "github.com/ActiveMemory/ctx/internal/config/regex"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	cfgToken "github.com/ActiveMemory/ctx/internal/config/token"
)

// buildFilename derives a proposal's on-disk name.
// Shape: `<TS>-<slug>.md` where `<TS>` is the UTC
// compact RFC-3339 form (`20260523T144530Z`) and
// `<slug>` is the kebab-case normalisation of the
// caller's slug.
//
// Parameters:
//   - now: timestamp.
//   - slug: free-text slug.
//
// Returns:
//   - string: filename portion only (no directory).
func buildFilename(now time.Time, slug string) string {
	clean := slugify(slug)
	if clean == "" {
		clean = cfgProposal.DefaultSlug
	}
	var sb strings.Builder
	sb.WriteString(now.UTC().Format(cfgTime.RFC3339Compact))
	sb.WriteString(cfgToken.Dash)
	sb.WriteString(clean)
	sb.WriteString(cfgFile.ExtMarkdown)
	return sb.String()
}

// slugify normalises a free-text title into a
// kebab-case slug-safe form: lowercased, only ASCII
// alnum + hyphen, runs of non-alnum collapsed to a
// single hyphen, trimmed.
//
// Parameters:
//   - s: free text.
//
// Returns:
//   - string: cleaned slug, possibly empty.
func slugify(s string) string {
	lower := strings.ToLower(s)
	parts := cfgRegex.Slug.Split(lower, -1)
	out := strings.Join(parts, cfgToken.Dash)
	out = strings.Trim(out, cfgToken.Dash)
	return out
}
