//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"strings"

	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// parseRegions splits a root of any kind: preamble, then the staging
// entries, then the ## Themes region (which runs to EOF). It writes the
// Preamble, Staging, and ThemesRaw fields of r.
//
// One path serves all three kinds because the layout is identical; only
// the line prefix that opens an entry differs ([EntryPrefix]). Cutting
// the themes region off first is what makes the convention prefix "## "
// safe to scan for — the structural "## Themes" is no longer in the text
// being searched.
//
// Parameters:
//   - r: the root being assembled (mutated in place)
//   - content: the full root content
//   - themeOffsets: byte offsets of every ## Themes heading line
func parseRegions(r *Root, content string, themeOffsets []int) {
	beforeThemes := content
	if r.HasThemes {
		t := themeOffsets[0]
		beforeThemes = content[:t]
		r.ThemesRaw = content[t:]
	}

	si := firstLinePrefixOffset(beforeThemes, EntryPrefix(r.Kind))
	if si != -1 {
		r.Preamble = beforeThemes[:si]
		r.Staging = beforeThemes[si:]
	} else {
		r.Preamble = beforeThemes
	}
}

// parseThemes reads the "- name — gist → [label](link)" bullet lines of a
// raw ## Themes region into Theme values. It is best-effort: lines it
// does not recognize are skipped (Validate, not Parse, judges shape).
//
// Parameters:
//   - themesRaw: the raw ## Themes region of a root
//
// Returns:
//   - []Theme: one per recognized bullet, in file order; nil if none
func parseThemes(themesRaw string) []Theme {
	if themesRaw == "" {
		return nil
	}
	var themes []Theme
	for _, line := range strings.Split(themesRaw, token.NewlineLF) {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, token.PrefixListDash) {
			continue
		}
		body := strings.TrimPrefix(trimmed, token.PrefixListDash)
		themes = append(themes, parseThemeBullet(body))
	}
	return themes
}

// parseThemeBullet pulls the name, gist, and link out of one theme
// bullet's body. The link is the first "(...)" following a "](".
//
// Parameters:
//   - body: one theme bullet with the leading "- " already stripped
//
// Returns:
//   - Theme: the parsed name/gist/link (fields empty when not present)
func parseThemeBullet(body string) Theme {
	t := Theme{}
	if open := strings.LastIndex(body, cfgDisc.LinkOpen); open != -1 {
		start := open + len(cfgDisc.LinkOpen)
		if end := strings.IndexByte(body[start:], ')'); end != -1 {
			t.Link = body[start : start+end]
		}
	}
	// name is the text before the em-dash metadata separator.
	if dash := strings.Index(body, token.MetaSeparator); dash != -1 {
		t.Name = strings.TrimSpace(body[:dash])
		t.Gist = strings.TrimSpace(body[dash+len(token.MetaSeparator):])
	} else {
		t.Name = strings.TrimSpace(body)
	}
	return t
}

// renderThemeBullet renders one theme bullet line (no trailing newline),
// the inverse of parseThemeBullet and in the exact shape it reads back:
// "- <theme> — <gist> → [<theme>](<noun>/<slug>.md)".
//
// Parameters:
//   - a: the assignment supplying theme name, gist, and slug
//   - noun: the theme-file subdirectory for this kind
//
// Returns:
//   - string: the rendered bullet line
func renderThemeBullet(a Assignment, noun string) string {
	link := noun + token.Slash + a.Slug + cfgFile.ExtMarkdown
	return token.PrefixListDash + a.Theme + token.MetaSeparator + a.Gist +
		cfgDisc.ThemeArrow + cfgDisc.LinkLabelOpen + a.Theme +
		cfgDisc.LinkOpen + link + cfgDisc.LinkClose
}

// headingLineOffsets returns the byte offset of every line whose trimmed
// content equals heading (an exact ATX heading line). Used to find and
// count region-delimiting headings. Headings inside an HTML comment are
// skipped — they are illustrative examples, not structure.
//
// Parameters:
//   - content: the text to scan
//   - heading: the exact heading line to match (e.g. "## Themes")
//
// Returns:
//   - []int: byte offsets of each matching line's start, in order
func headingLineOffsets(content, heading string) []int {
	spans := htmlCommentSpans(content)
	var offsets []int
	for i := 0; i < len(content); {
		line, next := lineAt(content, i)
		if strings.TrimSpace(line) == heading && !insideAnySpan(i, spans) {
			offsets = append(offsets, i)
		}
		if next == -1 {
			break
		}
		i = next
	}
	return offsets
}

// firstLinePrefixOffset returns the byte offset of the first line that
// starts with prefix, or -1. A match inside an HTML comment is skipped: a
// "## [" example in a knowledge file's format guide is not a staging entry.
//
// Parameters:
//   - content: the text to scan
//   - prefix: the line-start prefix to match (e.g. "## [")
//
// Returns:
//   - int: byte offset of the first matching line, or -1 if none
func firstLinePrefixOffset(content, prefix string) int {
	spans := htmlCommentSpans(content)
	for i := 0; i < len(content); {
		line, next := lineAt(content, i)
		if strings.HasPrefix(line, prefix) && !insideAnySpan(i, spans) {
			return i
		}
		if next == -1 {
			break
		}
		i = next
	}
	return -1
}

// entryBelowThemes reports whether an entry heading of kind k appears
// inside the themes region — below "## Themes", where the digesting pass
// cannot reach it.
//
// The region opens with its own "## Themes" heading, which shares the
// convention entry prefix ("## "), so that line is skipped explicitly:
// without it every migrated convention root would report a violation
// against itself. Headings inside an HTML comment are skipped for the
// same reason they are elsewhere — they are examples, not structure.
//
// Parameters:
//   - themesRaw: the root's raw themes region
//   - k: the root kind, selecting the entry prefix
//
// Returns:
//   - bool: true when an entry heading sits in the themes region
func entryBelowThemes(themesRaw string, k Kind) bool {
	if themesRaw == "" {
		return false
	}
	prefix := EntryPrefix(k)
	spans := htmlCommentSpans(themesRaw)
	for i := 0; i < len(themesRaw); {
		line, next := lineAt(themesRaw, i)
		if strings.HasPrefix(line, prefix) &&
			strings.TrimSpace(line) != cfgDisc.HeadingThemes &&
			!insideAnySpan(i, spans) {
			return true
		}
		if next == -1 {
			break
		}
		i = next
	}
	return false
}

// htmlCommentSpans returns the [start, end) byte ranges of every HTML
// comment (token.HTMLCommentOpen … token.HTMLCommentClose) in content; an
// unterminated open runs to EOF. The structural heading scans use it to
// ignore headings that are only illustrative examples inside a comment —
// e.g. the "## [YYYY-MM-DD] Title" line in a knowledge file's format guide,
// which must not be mistaken for a staging entry.
//
// Parameters:
//   - content: the text to scan
//
// Returns:
//   - [][2]int: {start, end} byte ranges in order; nil when none
func htmlCommentSpans(content string) [][2]int {
	var spans [][2]int
	for i := 0; i < len(content); {
		open := strings.Index(content[i:], token.HTMLCommentOpen)
		if open == -1 {
			break
		}
		start := i + open
		afterOpen := start + len(token.HTMLCommentOpen)
		rel := strings.Index(content[afterOpen:], token.HTMLCommentClose)
		if rel == -1 {
			spans = append(spans, [2]int{start, len(content)})
			break
		}
		end := afterOpen + rel + len(token.HTMLCommentClose)
		spans = append(spans, [2]int{start, end})
		i = end
	}
	return spans
}

// insideAnySpan reports whether byte offset off falls within any of the
// [start, end) ranges (from htmlCommentSpans).
//
// Parameters:
//   - off: the byte offset to test
//   - spans: byte ranges to test against
//
// Returns:
//   - bool: true when off is inside a span
func insideAnySpan(off int, spans [][2]int) bool {
	for _, s := range spans {
		if off >= s[0] && off < s[1] {
			return true
		}
	}
	return false
}

// lineAt returns the line beginning at byte offset i (without its
// trailing newline) and the offset of the next line's start, or -1 when
// this is the last line.
//
// Parameters:
//   - content: the text being scanned
//   - i: byte offset of the start of the line to read
//
// Returns:
//   - line: the line's content without its trailing newline
//   - next: byte offset of the next line's start, or -1 if last
func lineAt(content string, i int) (line string, next int) {
	rel := strings.Index(content[i:], token.NewlineLF)
	if rel == -1 {
		return content[i:], -1
	}
	return content[i : i+rel], i + rel + 1
}

// lineByteOffsets returns the byte offset of the start of each line in
// content, with a trailing sentinel offset. Index i is the byte offset
// where line i begins; the value can overshoot len(content) for the
// synthetic final element, so callers clamp with clampOffset. SplitStaging
// uses it to cut entry spans on raw byte boundaries.
//
// Parameters:
//   - content: the text to index
//
// Returns:
//   - []int: byte offsets per line, plus one trailing sentinel
func lineByteOffsets(content string) []int {
	lines := strings.Split(content, token.NewlineLF)
	offs := make([]int, len(lines)+1)
	for i, ln := range lines {
		offs[i+1] = offs[i] + len(ln) + len(token.NewlineLF)
	}
	return offs
}

// clampOffset bounds a byte offset to [0, n]. The final line offset from
// lineByteOffsets overshoots by one newline when content has no trailing
// newline; clamping keeps slice bounds valid without special-casing.
//
// Parameters:
//   - off: a candidate byte offset
//   - n: the length to clamp to
//
// Returns:
//   - int: off bounded to [0, n]
func clampOffset(off, n int) int {
	if off > n {
		return n
	}
	return off
}
