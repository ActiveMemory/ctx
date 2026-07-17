//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"os"
	"path/filepath"

	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// CheckPairing verifies the root's gists and the theme files are 1:1: no
// theme file the root cannot reach, and no gist pointing at a file that
// does not exist.
//
// An un-migrated root (no gists) with no theme directory is vacuously
// paired.
//
// Parameters:
//   - root: the parsed root (its Themes supply the gist links)
//   - themeDir: directory holding this kind's theme files
//
// Returns:
//   - error: ErrOrphanThemeFile, ErrMissingThemeFile, a read error, or nil
func CheckPairing(root Root, themeDir string) error {
	gists := gistBasenames(root)
	files, err := themeFiles(themeDir)
	if err != nil {
		return err
	}

	fileSet := map[string]bool{}
	for _, f := range files {
		name := filepath.Base(f)
		fileSet[name] = true
		if !gists[name] {
			return errDisc.ErrOrphanThemeFile
		}
	}
	for g := range gists {
		if !fileSet[g] {
			return errDisc.ErrMissingThemeFile
		}
	}
	return nil
}

// CheckUniqueness verifies every entry lives in exactly one place: the
// staging zone or a single theme file, never two. Entry identity is its
// timestamp+title.
//
// Parameters:
//   - root: the parsed root (its Staging supplies staged entries)
//   - themeDir: directory holding this kind's theme files
//
// Returns:
//   - error: ErrDuplicateEntry, a read error, or nil
func CheckUniqueness(root Root, themeDir string) error {
	seen := map[string]bool{}
	for _, id := range entryIDs(root.Staging) {
		if seen[id] {
			return errDisc.ErrDuplicateEntry
		}
		seen[id] = true
	}

	files, err := themeFiles(themeDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		data, readErr := internalIo.SafeReadUserFile(f)
		if readErr != nil {
			return readErr
		}
		for _, id := range entryIDs(string(data)) {
			if seen[id] {
				return errDisc.ErrDuplicateEntry
			}
			seen[id] = true
		}
	}
	return nil
}

// CheckLinks verifies every theme link in the root resolves to a path
// that exists, relative to the context directory.
//
// Parameters:
//   - root: the parsed root (its Themes supply the links)
//   - ctxDir: the context directory the links are relative to
//
// Returns:
//   - error: ErrBrokenThemeLink, or nil
func CheckLinks(root Root, ctxDir string) error {
	for _, t := range root.Themes {
		if t.Link == "" {
			continue
		}
		if _, statErr := os.Stat(filepath.Join(ctxDir, t.Link)); statErr != nil {
			return errDisc.ErrBrokenThemeLink
		}
	}
	return nil
}
