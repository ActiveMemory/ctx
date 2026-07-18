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

// parseEntryKind splits an entry-based root: preamble, then the staging
// entries, then the ## Themes region (which runs to EOF). It writes the
// Preamble, Staging, and ThemesRaw fields of r.
//
// Parameters:
//   - r: the root being assembled (mutated in place)
//   - content: the full root content
//   - themeOffsets: byte offsets of every ## Themes heading line
func parseEntryKind(r *Root, content string, themeOffsets []int) {
	beforeThemes := content
	if r.HasThemes {
		t := themeOffsets[0]
		beforeThemes = content[:t]
		r.ThemesRaw = content[t:]
	}

	si := firstLinePrefixOffset(beforeThemes, cfgDisc.EntryLinePrefix)
	if si != -1 {
		r.Preamble = beforeThemes[:si]
		r.Staging = beforeThemes[si:]
	} else {
		r.Preamble = beforeThemes
	}
}

// parseConvention splits a conventions root: preamble, the ## Themes
// region, then the ## Recent staging region. An un-migrated conventions
// file has no ## Themes; its prose sections are the staging directly.
//
// Parameters:
//   - r: the root being assembled (mutated in place)
//   - content: the full root content
//   - themeOffsets: byte offsets of every ## Themes heading line
func parseConvention(r *Root, content string, themeOffsets []int) {
	if !r.HasThemes {
		si := firstLinePrefixOffset(content, cfgDisc.ConventionLinePrefix)
		if si != -1 {
			r.Preamble = content[:si]
			r.Staging = content[si:]
		} else {
			r.Preamble = content
		}
		return
	}

	t := themeOffsets[0]
	r.Preamble = content[:t]
	recentOffsets := headingLineOffsets(content[t:], cfgDisc.HeadingRecent)
	if len(recentOffsets) > 0 {
		rec := t + recentOffsets[0]
		r.ThemesRaw = content[t:rec]
		r.Staging = content[rec:]
	} else {
		r.ThemesRaw = content[t:]
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

// headingLineOffsets returns the byte offset of every line whose trimmed
// content equals heading (an exact ATX heading line). Used to find and
// count region-delimiting headings.
//
// Parameters:
//   - content: the text to scan
//   - heading: the exact heading line to match (e.g. "## Themes")
//
// Returns:
//   - []int: byte offsets of each matching line's start, in order
func headingLineOffsets(content, heading string) []int {
	var offsets []int
	for i := 0; i < len(content); {
		line, next := lineAt(content, i)
		if strings.TrimSpace(line) == heading {
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
// starts with prefix, or -1.
//
// Parameters:
//   - content: the text to scan
//   - prefix: the line-start prefix to match (e.g. "## [")
//
// Returns:
//   - int: byte offset of the first matching line, or -1 if none
func firstLinePrefixOffset(content, prefix string) int {
	for i := 0; i < len(content); {
		line, next := lineAt(content, i)
		if strings.HasPrefix(line, prefix) {
			return i
		}
		if next == -1 {
			break
		}
		i = next
	}
	return -1
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
