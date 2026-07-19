//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"strings"

	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// WriteThemeBullet folds one assignment's gist into a root's raw ## Themes
// region and returns the updated region. It edits raw bullet lines rather
// than re-rendering from parsed Theme values: parseThemeBullet folds the
// "→ [link]" tail into the parsed gist, so re-rendering would corrupt it,
// and the spec requires untouched themes be left alone — so every bullet
// but the one for a.Theme is byte-preserved.
//
// Behaviour:
//   - empty region (un-migrated root): creates "## Themes" with the bullet.
//   - a.Theme already has a bullet: that line is replaced in place.
//   - a.Theme is new: the bullet is appended after the last existing one
//     (or after the heading when there are none).
//
// The bullet renders as "- <theme> — <gist> → [<theme>](<noun>/<slug>.md)".
//
// Parameters:
//   - themesRaw: the root's raw ## Themes region (may be empty)
//   - a: the assignment whose gist to write
//   - noun: the theme-file subdirectory for this kind (e.g. "learnings")
//
// Returns:
//   - string: the updated ## Themes region
func WriteThemeBullet(themesRaw string, a Assignment, noun string) string {
	bullet := renderThemeBullet(a, noun)

	if strings.TrimSpace(themesRaw) == "" {
		return cfgDisc.HeadingThemes + token.NewlineLF + token.NewlineLF +
			bullet + token.NewlineLF
	}

	lines := strings.Split(themesRaw, token.NewlineLF)
	lastBullet := -1
	for i, ln := range lines {
		trimmed := strings.TrimSpace(ln)
		if !strings.HasPrefix(trimmed, token.PrefixListDash) {
			continue
		}
		lastBullet = i
		body := strings.TrimPrefix(trimmed, token.PrefixListDash)
		if parseThemeBullet(body).Name == a.Theme {
			lines[i] = bullet
			return strings.Join(lines, token.NewlineLF)
		}
	}

	if lastBullet == -1 {
		trimmed := strings.TrimRight(themesRaw, token.NewlineLF)
		return trimmed + token.NewlineLF + token.NewlineLF + bullet + token.NewlineLF
	}

	out := make([]string, 0, len(lines)+1)
	out = append(out, lines[:lastBullet+1]...)
	out = append(out, bullet)
	out = append(out, lines[lastBullet+1:]...)
	return strings.Join(out, token.NewlineLF)
}
