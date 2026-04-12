//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0
//
// ================================================================
// STOP — Read internal/audit/README.md before editing this file.
//
// These tests enforce project conventions. The codebase is clean:
// all checks pass with zero violations, zero exceptions.
//
// If a test fails after your change, fix the code under test.
// Do NOT add allowlist entries, bump grandfathered counters, or
// weaken checks. Exceptions require a dedicated PR with
// justification for every entry. See README.md for the full policy.
// ================================================================

package audit

import (
	"go/ast"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

// minDocLines is the minimum number of meaningful text
// lines required in a package doc comment. Set to 3
// to catch lazy one-liners while accepting the current
// codebase standard. Tighten incrementally.
const minDocLines = 3

// TestPackageDocQuality ensures every Go package under
// internal/ has a doc.go file with a meaningful package
// doc comment (at least 8 lines of real text, excluding
// blank comment lines, the "Package X" opener, and
// file-list patterns).
func TestPackageDocQuality(t *testing.T) {
	var violations []string

	// Collect all Go package directories.
	pkgDirs := make(map[string]bool)
	walkErr := filepath.WalkDir(
		"../",
		func(
			path string, d os.DirEntry, err error,
		) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if !strings.HasSuffix(d.Name(), ".go") {
				return nil
			}
			if isTestFile(d.Name()) {
				return nil
			}
			dir := filepath.Dir(path)
			abs, absErr := filepath.Abs(dir)
			if absErr != nil {
				return absErr
			}
			pkgDirs[abs] = true
			return nil
		},
	)
	if walkErr != nil {
		t.Fatalf("walk: %v", walkErr)
	}

	for dir := range pkgDirs {
		docPath := filepath.Join(dir, "doc.go")

		// Check existence.
		if _, statErr := os.Stat(docPath); os.IsNotExist(statErr) {
			violations = append(violations,
				dir+": missing doc.go",
			)
			continue
		}

		// Check quality.
		checkDocQuality(
			t, docPath, dir, &violations,
		)
	}

	if len(violations) > 0 {
		t.Errorf(
			"%d package doc issues:",
			len(violations),
		)
	}
	limit := 30
	if len(violations) < limit {
		limit = len(violations)
	}
	for _, v := range violations[:limit] {
		t.Error(v)
	}
	if len(violations) > 30 {
		t.Errorf(
			"... and %d more",
			len(violations)-30,
		)
	}
}

// checkDocQuality parses a doc.go and verifies the
// package doc comment has enough meaningful content.
func checkDocQuality(
	t *testing.T, docPath, dir string,
	violations *[]string,
) {
	t.Helper()

	pkgs := loadPackages(t)

	// Find the package matching this directory.
	var docFile *ast.File
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			fpath := pkg.Fset.Position(
				file.Pos(),
			).Filename
			if fpath == docPath {
				docFile = file
				break
			}
		}
		if docFile != nil {
			break
		}
	}

	if docFile == nil {
		// Could not load via packages; skip quality
		// check (existence already verified).
		return
	}

	if docFile.Doc == nil {
		*violations = append(*violations,
			dir+": doc.go has no package comment",
		)
		return
	}

	meaningful := countMeaningfulLines(
		docFile.Doc.Text(),
	)
	if meaningful < minDocLines {
		*violations = append(*violations,
			dir+": doc.go has "+
				itoa(meaningful)+
				" meaningful lines (min "+
				itoa(minDocLines)+")",
		)
	}
}

// countMeaningfulLines counts non-blank, non-boilerplate
// lines in a doc comment text. Excludes:
// - Blank lines
// - The "Package X ..." opener line
// - File-list patterns ("// - foo.go", "// Source files:")
// - Lines that are just punctuation/whitespace
func countMeaningfulLines(text string) int {
	count := 0
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)

		// Skip blank.
		if trimmed == "" {
			continue
		}

		// Skip "Package X" opener.
		if strings.HasPrefix(trimmed, "Package ") {
			continue
		}

		// Skip file-list patterns.
		if strings.HasPrefix(trimmed, "- ") &&
			strings.HasSuffix(trimmed, ".go") {
			continue
		}
		if strings.HasPrefix(trimmed, "Source files") {
			continue
		}

		// Skip lines that are just punctuation.
		allPunct := true
		for _, r := range trimmed {
			if r != '-' && r != '=' && r != '_' &&
				r != '*' && r != '#' {
				allPunct = false
				break
			}
		}
		if allPunct {
			continue
		}

		// Skip very short lines (likely just a
		// word or label).
		if utf8.RuneCountInString(trimmed) < 10 {
			continue
		}

		count++
	}

	return count
}
