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
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

// maxLineLength is the maximum allowed line length in non-test Go
// source files.
const maxLineLength = 80

// TestLineLength ensures all non-test Go files under internal/ have
// lines of at most 80 characters.
//
// Test files are exempt.
func TestLineLength(t *testing.T) {
	var violations []string

	walkErr := filepath.WalkDir("../", func(path string, d os.DirEntry, err error) error {
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

		// config/embed/text/ has DescKey constants
		// whose identifiers exceed 80 chars and
		// cannot be shortened without a full rename.
		// gofmt re-expands any manual wrapping.
		if strings.Contains(path, "config/embed/text") {
			return nil
		}

		abs, absErr := filepath.Abs(path)
		if absErr != nil {
			return absErr
		}

		data, readErr := os.ReadFile(path) //nolint:gosec // audit test reads source files
		if readErr != nil {
			return readErr
		}

		scanner := bufio.NewScanner(
			strings.NewReader(string(data)),
		)
		lineNum := 0
		inRawString := false
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// Track raw string literals (backtick).
			// Toggle on each unescaped backtick.
			for _, r := range line {
				if r == '`' {
					inRawString = !inRawString
				}
			}

			// Lines inside raw strings cannot be
			// reformatted without changing content.
			if inRawString {
				continue
			}

			runeLen := utf8.RuneCountInString(line)
			if runeLen > maxLineLength {
				violations = append(violations,
					abs+":"+
						itoa(lineNum)+
						": "+itoa(runeLen)+
						" chars (max "+
						itoa(maxLineLength)+")",
				)
			}
		}

		return scanner.Err()
	})
	if walkErr != nil {
		t.Fatalf("filepath.WalkDir: %v", walkErr)
	}

	if len(violations) > 0 {
		t.Errorf("%d lines exceed %d characters:",
			len(violations), maxLineLength)
	}
	limit := 30
	if len(violations) < limit {
		limit = len(violations)
	}
	for _, v := range violations[:limit] {
		t.Error(v)
	}
	if len(violations) > 30 {
		t.Errorf("... and %d more", len(violations)-30)
	}
}

// itoa converts an int to a string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	if neg {
		digits = append(digits, '-')
	}
	// reverse
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
