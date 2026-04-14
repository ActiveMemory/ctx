//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

import (
	"strings"
	"testing"
)

func TestContentEscapesEntryHeaders(t *testing.T) {
	input := "## [2026-03-15] Decision title"
	got := Content(input)
	want := `\## [2026-03-15] Decision title`
	if got != want {
		t.Errorf("Content(%q) = %q, want %q", input, got, want)
	}
}

func TestContentEscapesTaskCheckboxUnchecked(t *testing.T) {
	got := Content("- [ ] New task")
	want := `\- [ ] New task`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestContentEscapesTaskCheckboxChecked(t *testing.T) {
	got := Content("- [x] Done task")
	want := `\- [x] Done task`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestContentEscapesConstitutionRules(t *testing.T) {
	input := "- [ ] **Never break the constitution"
	got := Content(input)
	if !strings.HasPrefix(got, `\- [ ] **Never`) {
		t.Errorf("got %q, want constitution rule escaped", got)
	}
}

func TestContentStripsNullBytes(t *testing.T) {
	got := Content("hello\x00world")
	if got != "helloworld" {
		t.Errorf("got %q, want %q", got, "helloworld")
	}
}

func TestContentPreservesNormalText(t *testing.T) {
	input := "This is a normal architecture decision."
	got := Content(input)
	if got != input {
		t.Errorf("got %q, want unchanged", got)
	}
}

func TestContentMultilineInjection(t *testing.T) {
	input := "Legit\n## [2026-01-01] Injected\n- [ ] Fake"
	got := Content(input)
	if strings.Contains(got, "\n## [2026") {
		t.Error("entry header injection not escaped")
	}
	if strings.Contains(got, "\n- [ ] Fake") {
		t.Error("checkbox injection not escaped")
	}
}

func TestReflectTruncates(t *testing.T) {
	got := Reflect(strings.Repeat("a", 500), 256)
	if len(got) != 256 {
		t.Errorf("len = %d, want 256", len(got))
	}
}

func TestReflectStripsControlChars(t *testing.T) {
	got := Reflect("tool\x07name\x1b[31m", 0)
	if got != "toolname[31m" {
		t.Errorf("got %q, want %q", got, "toolname[31m")
	}
}

func TestReflectPreservesNormal(t *testing.T) {
	got := Reflect("ctx_status", 256)
	if got != "ctx_status" {
		t.Errorf("got %q, want unchanged", got)
	}
}

func TestReflectZeroMaxLen(t *testing.T) {
	got := Reflect(strings.Repeat("x", 1000), 0)
	if len(got) != 1000 {
		t.Errorf("len = %d, want 1000 (no truncation)", len(got))
	}
}

func TestTruncateShort(t *testing.T) {
	if got := truncate("short", 100); got != "short" {
		t.Errorf("got %q", got)
	}
}

func TestTruncateLong(t *testing.T) {
	if got := truncate("long input", 4); got != "long" {
		t.Errorf("got %q", got)
	}
}

func TestTruncateZero(t *testing.T) {
	if got := truncate("any", 0); got != "any" {
		t.Errorf("got %q", got)
	}
}

func TestStripControlPreservesWhitespace(t *testing.T) {
	input := "a\nb\tc\r"
	if got := StripControl(input); got != input {
		t.Errorf("got %q, want unchanged", got)
	}
}

func TestStripControlRemovesBell(t *testing.T) {
	if got := StripControl("hello\x07world"); got != "helloworld" {
		t.Errorf("got %q", got)
	}
}

func TestSessionIDSafe(t *testing.T) {
	input := "session-2026-03-15"
	if got := SessionID(input); got != input {
		t.Errorf("got %q, want unchanged", got)
	}
}

func TestSessionIDStripsTraversal(t *testing.T) {
	got := SessionID("../../etc/passwd")
	if strings.Contains(got, "..") || strings.Contains(got, "/") {
		t.Errorf("got %q, contains traversal", got)
	}
}

func TestSessionIDStripsBackslashTraversal(t *testing.T) {
	got := SessionID(`..\..\windows\system32`)
	if strings.Contains(got, "..") || strings.Contains(got, `\`) {
		t.Errorf("got %q, contains traversal", got)
	}
}

func TestSessionIDStripsNullBytes(t *testing.T) {
	got := SessionID("session\x00evil")
	if strings.Contains(got, "\x00") {
		t.Errorf("got %q, contains null byte", got)
	}
}

func TestSessionIDLimitsLength(t *testing.T) {
	got := SessionID(strings.Repeat("a", 300))
	if len(got) > 128 {
		t.Errorf("len = %d, want <= 128", len(got))
	}
}

func TestSessionIDReplacesUnsafe(t *testing.T) {
	got := SessionID("session with spaces!@#$")
	if strings.ContainsAny(got, " !@#$") {
		t.Errorf("got %q, contains unsafe chars", got)
	}
}
