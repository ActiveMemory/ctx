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
	for _, e := range a.Entries {
		spans.WriteString(moved[entryID(e)])
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
		for _, e := range a.Entries {
			if vErr := verifyContains(string(content), moved[entryID(e)]); vErr != nil {
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
