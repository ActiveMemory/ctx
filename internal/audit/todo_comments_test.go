//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"
)

// todoPattern matches actionable TODO, FIXME, HACK,
// and XXX markers in comments. Requires the marker
// to be followed by a colon, parenthesis, or dash
// (conventional action-item forms), not just prose
// mentions of the word.
var todoPattern = regexp.MustCompile(
	`//.*\b(TODO|FIXME|HACK|XXX)\s*[:(/-]` +
		`|/\*.*\b(TODO|FIXME|HACK|XXX)\s*[:(/-]`,
)

// TestNoTODOComments ensures no TODO, FIXME, HACK, or
// XXX comments exist in non-test Go source files.
//
// Deferred work must be tracked in TASKS.md, not
// hidden in source comments. Test files are exempt.
func TestNoTODOComments(t *testing.T) {
	var violations []string

	walkErr := filepath.WalkDir(
		"../",
		func(path string, d os.DirEntry, err error) error {
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

			abs, absErr := filepath.Abs(path)
			if absErr != nil {
				return absErr
			}

			data, readErr := os.ReadFile(path) //nolint:gosec
			if readErr != nil {
				return readErr
			}

			scanner := bufio.NewScanner(
				strings.NewReader(string(data)),
			)
			lineNum := 0
			for scanner.Scan() {
				lineNum++
				line := scanner.Text()
				if todoPattern.MatchString(line) {
					// Trim for display.
					display := line
					if utf8.RuneCountInString(display) > 60 {
						display = display[:60] + "..."
					}
					violations = append(violations,
						abs+":"+
							itoa(lineNum)+
							": "+
							strings.TrimSpace(display),
					)
				}
			}

			return scanner.Err()
		},
	)
	if walkErr != nil {
		t.Fatalf("filepath.WalkDir: %v", walkErr)
	}

	for _, v := range violations {
		t.Error(v)
	}
}
