//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package wrap

import (
	"strings"
	"testing"
)

func TestSoftWrap(t *testing.T) {
	line := "This is a test line that is much longer" +
		" than eighty characters and should be" +
		" wrapped at word boundaries."
	result := Soft(line, 40)

	if len(result) < 2 {
		t.Errorf("expected multiple lines, got %d", len(result))
	}

	for _, l := range result {
		if len(l) > 45 {
			t.Errorf("line too long: %q (%d chars)", l, len(l))
		}
	}

	joined := strings.Join(result, " ")
	if joined != line {
		t.Errorf("content changed after wrap:\n  got:  %q\n  want: %q", joined, line)
	}
}

func TestSoftWrapContent(t *testing.T) {
	content := "---\ntitle: test with a very long value" +
		" that should not be wrapped because it" +
		" is inside frontmatter block\n---\n" +
		"Short line\nThis is a very long line" +
		" that exceeds eighty characters and" +
		" should be wrapped at word boundaries" +
		" properly.\n"
	got := Content(content)

	if !strings.Contains(got, "title: test with a very long value") {
		t.Error("frontmatter should not be wrapped")
	}
	if !strings.Contains(got, "Short line") {
		t.Error("short lines should be preserved")
	}
}
