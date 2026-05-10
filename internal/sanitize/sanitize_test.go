//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

import (
	"strings"
	"testing"
	"unicode/utf8"
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

// --- UTF-8-safe truncation (follow-up to PR #76) ---

func TestTruncateMultibyteRuneBoundary(t *testing.T) {
	// "héllo" — 'é' is U+00E9, a 2-byte UTF-8 rune (0xC3 0xA9).
	// Bytes: 'h' 0xC3 0xA9 'l' 'l' 'o' (6 bytes total).
	// Cutting at 2 bytes lands inside 'é'; must back up to 1.
	got := truncate("héllo", 2)
	if got != "h" {
		t.Errorf("truncate at mid-rune = %q (% x), want %q", got, got, "h")
	}
	if !utf8Valid(got) {
		t.Errorf("truncate produced invalid UTF-8: % x", got)
	}
}

func TestTruncateThreeByteRune(t *testing.T) {
	// '€' is U+20AC, 3 bytes (0xE2 0x82 0xAC).
	// Cutting "a€b" at 2 lands inside '€'; back up to 1.
	got := truncate("a€b", 2)
	if got != "a" {
		t.Errorf("got %q, want %q", got, "a")
	}
}

func TestTruncateFourByteRune(t *testing.T) {
	// '🜲' is U+1F732, 4 bytes (0xF0 0x9F 0x9C 0xB2).
	// Cutting "🜲x" at 2 lands inside the emoji; back up to 0.
	got := truncate("🜲x", 2)
	if got != "" {
		t.Errorf("got %q (% x), want empty", got, got)
	}
}

func TestTruncateAtExactRuneBoundary(t *testing.T) {
	// "héllo" cut at 3 bytes lands exactly after 'é' (1+2=3).
	got := truncate("héllo", 3)
	if got != "hé" {
		t.Errorf("got %q, want %q", got, "hé")
	}
}

func TestTruncatePreservesValidUTF8(t *testing.T) {
	// Repeated 2-byte rune; any cut must terminate at an even byte
	// offset relative to the start of the run.
	in := strings.Repeat("é", 100) // 200 bytes
	for cut := 0; cut <= 200; cut++ {
		got := truncate(in, cut)
		if !utf8Valid(got) {
			t.Errorf("cut=%d produced invalid UTF-8: % x", cut, got)
		}
	}
}

func TestReflectMultibyteSafe(t *testing.T) {
	// Reflect should also produce valid UTF-8 when truncating.
	got := Reflect("héllo world", 2)
	if got != "h" {
		t.Errorf("Reflect mid-rune = %q, want %q", got, "h")
	}
}

func TestSessionIDMultibyteSafe(t *testing.T) {
	// Pre-sanitize the string so unsafe chars become hyphens, then
	// run a long input through SessionID; the result must be valid
	// UTF-8 even at the truncation boundary.
	in := strings.Repeat("a", 130) // exceeds MaxSessionIDLen=128
	got := SessionID(in)
	if !utf8Valid(got) {
		t.Errorf("SessionID produced invalid UTF-8: % x", got)
	}
}

// --- Zl/Zp separator stripping (follow-up to PR #76) ---

func TestStripControlRemovesLineSeparator(t *testing.T) {
	// U+2028 — LINE SEPARATOR.
	got := StripControl("a b")
	if got != "ab" {
		t.Errorf("got %q, want %q", got, "ab")
	}
}

func TestStripControlRemovesParagraphSeparator(t *testing.T) {
	// U+2029 — PARAGRAPH SEPARATOR.
	got := StripControl("a b")
	if got != "ab" {
		t.Errorf("got %q, want %q", got, "ab")
	}
}

func TestStripControlRemovesBothSeparators(t *testing.T) {
	got := StripControl("x y z")
	if got != "xyz" {
		t.Errorf("got %q, want %q", got, "xyz")
	}
}

func TestReflectStripsLineSeparator(t *testing.T) {
	// Newline injection via U+2028 must not survive reflection in
	// error messages.
	got := Reflect("tool name", 0)
	if got != "toolname" {
		t.Errorf("got %q, want %q", got, "toolname")
	}
}

// utf8Valid is a thin alias for [utf8.ValidString], used to make the
// "did we cut mid-rune?" assertion read clearly at the call sites.
func utf8Valid(s string) bool {
	return utf8.ValidString(s)
}
