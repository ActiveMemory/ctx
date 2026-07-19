//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"os"
	"path/filepath"
	"strings"

	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// appendTheme appends an assignment's moved entry spans to its theme
// file, creating the file (with an H1 header) when absent and separating
// with a blank line when appending. Spans are written verbatim so the
// verify step's byte-presence check holds.
//
// Parameters:
//   - path: the theme file to write
//   - a: the assignment whose entries are appended
//   - moved: id -> verbatim span, from SplitStaging
//
// Returns:
//   - error: an IO error, or nil
func appendTheme(path string, a Assignment, moved map[string]string) error {
	var spans strings.Builder
	for _, id := range a.Entries {
		spans.WriteString(moved[id])
	}

	var payload string
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		payload = token.HeadingLevelOneStart + a.Theme +
			token.NewlineLF + token.NewlineLF + spans.String()
	} else {
		payload = token.NewlineLF + spans.String()
	}
	return internalIo.AppendBytes(path, []byte(payload), cfgFs.PermFile)
}

// verifyThemes re-reads each written theme file and confirms every moved
// entry body is byte-present. It runs after all appends and before the
// root rewrite, so a miss aborts with the root untouched.
//
// Parameters:
//   - themeDir: the directory holding this kind's theme files
//   - plan: the digest plan whose assignments were appended
//   - moved: id -> verbatim span, from SplitStaging
//
// Returns:
//   - error: a read error, ErrVerifyFailed, or nil
func verifyThemes(themeDir string, plan Plan, moved map[string]string) error {
	for _, a := range plan.Assignments {
		path := filepath.Join(themeDir, a.Slug+cfgFile.ExtMarkdown)
		content, readErr := internalIo.SafeReadUserFile(path)
		if readErr != nil {
			return readErr
		}
		for _, id := range a.Entries {
			if vErr := verifyContains(string(content), moved[id]); vErr != nil {
				return vErr
			}
		}
	}
	return nil
}

// verifyContains is the byte-presence guard: it confirms a moved entry
// body is present in its theme file after the append. Absence means the
// append did not land intact, so the pass must abort with the root
// untouched rather than remove the entry from staging.
//
// Parameters:
//   - fileContent: the theme file re-read after appending
//   - span: the verbatim entry span that should be present
//
// Returns:
//   - error: ErrVerifyFailed when span is absent, else nil
func verifyContains(fileContent, span string) error {
	if !strings.Contains(fileContent, span) {
		return errDisc.ErrVerifyFailed
	}
	return nil
}

// rewriteRoot assembles the post-move root: the preamble, the remaining
// staging, and the ## Themes region with every assignment's gist folded
// in. On first-run migration it inserts a blank-line boundary before the
// freshly created ## Themes section.
//
// Parameters:
//   - root: the parsed root before the move
//   - remaining: the staging with moved entries cut out
//   - plan: the digest plan whose gists to fold
//   - noun: the theme-file subdirectory for this kind
//
// Returns:
//   - string: the rewritten root content
func rewriteRoot(root Root, remaining string, plan Plan, noun string) string {
	newThemes := root.ThemesRaw
	for _, a := range plan.Assignments {
		newThemes = WriteThemeBullet(newThemes, a, noun)
	}
	sep := ""
	if !root.HasThemes {
		sep = freshThemesBoundary(remaining)
	}
	return root.Preamble + remaining + sep + newThemes
}

// freshThemesBoundary returns the separator to place between the
// remaining staging and a newly created ## Themes section, so the heading
// sits on its own line after a blank line. It is used only on first-run
// migration (an un-migrated root has no existing boundary to preserve).
//
// Parameters:
//   - remaining: the staging that will precede the new ## Themes section
//
// Returns:
//   - string: empty, one newline, or two, so the boundary ends in a
//     blank line before the heading
func freshThemesBoundary(remaining string) string {
	blank := token.NewlineLF + token.NewlineLF
	switch {
	case remaining == "" || strings.HasSuffix(remaining, blank):
		return ""
	case strings.HasSuffix(remaining, token.NewlineLF):
		return token.NewlineLF
	default:
		return blank
	}
}

// renderThemeBullet renders one theme bullet line (no trailing newline),
// in the exact shape parseThemeBullet reads back:
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

// lineByteOffsets returns the byte offset of the start of each line in
// content, with a trailing sentinel offset. Index i is the byte offset
// where line i begins; the value can overshoot len(content) for the
// synthetic final element, so callers clamp with clampOffset.
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
