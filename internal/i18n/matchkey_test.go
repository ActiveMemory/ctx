//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package i18n_test

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/i18n"
)

func TestMatchKey_ASCIIIsByteIdentical(t *testing.T) {
	for _, in := range []string{
		"", "hello", "HELLO", "Hello World", "123-abc_XYZ", "tbd", "n/a",
	} {
		got := i18n.MatchKey(in)
		want := i18n.Fold(in) // Fold and MatchKey agree on ASCII.
		if got != want {
			t.Errorf("MatchKey(%q)=%q, want %q (== Fold)", in, got, want)
		}
	}
}

// TestMatchKey_CollapsesLatinDiacritics enumerates the
// pairs MatchKey is *supposed* to collapse. These are the
// ergonomic wins for cross-keyboard placeholder matching:
// a casual user typing the ASCII form should hit the
// vocabulary entry written in the diacritic form, and
// vice versa.
func TestMatchKey_CollapsesLatinDiacritics(t *testing.T) {
	cases := []struct {
		name, a, b string
	}{
		{"Turkish: İPTAL matches iptal", "İPTAL", "iptal"},
		{"Turkish: İptal matches iptal", "İptal", "iptal"},
		{"German: Straße matches strasse", "Straße", "strasse"},
		{"German: über matches uber", "über", "uber"},
		{"French: café matches cafe", "café", "cafe"},
		{"French: naïve matches naive", "naïve", "naive"},
		{"Spanish: Niño matches nino", "Niño", "nino"},
		{"Catalan: façana matches facana", "façana", "facana"},
		{"Czech: škola matches skola", "škola", "skola"},
		{"Vietnamese: tươi matches tuoi", "tươi", "tuoi"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			a, b := i18n.MatchKey(tt.a), i18n.MatchKey(tt.b)
			if a != b {
				t.Errorf("MatchKey(%q)=%q, MatchKey(%q)=%q; should match", tt.a, a, tt.b, b)
			}
		})
	}
}

// TestMatchKey_PreservesScriptEssentialMarks guards the
// other half of the contract: non-Latin script marks
// (Arabic, Indic, Hebrew) live outside U+0300–U+036F and
// must stay distinct so the comparison doesn't silently
// erase script-essential meaning.
func TestMatchKey_PreservesScriptEssentialMarks(t *testing.T) {
	cases := []struct {
		name, a, b string
	}{
		{"Arabic: hamza below preserved", "إلغاء", "الغاء"},
		{"Bengali: vowel sign preserved", "কা", "ক"},
		{"Devanagari: vowel sign preserved", "का", "क"},
		{"Hindi: chandra-bindu preserved", "हाँ", "हा"},
		{"Hebrew: niqqud preserved", "בִּטּוּל", "בטול"},
		{"Chinese: no false match", "取消", "取"},
		{"Korean: no false match", "취소", "취"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			a, b := i18n.MatchKey(tt.a), i18n.MatchKey(tt.b)
			if a == b {
				t.Errorf(
					"MatchKey(%q)=%q == MatchKey(%q)=%q; should stay distinct",
					tt.a, a, tt.b, b,
				)
			}
		})
	}
}

// TestMatchKey_FoldStaysStrict guards the API split:
// Fold remains a strict Unicode primitive (does NOT
// collapse İ→i or ü→u). Callers that need the strict
// behavior keep getting it. This test pairs with
// TestFold_UnicodeCorrectness in fold_test.go.
func TestMatchKey_FoldStaysStrict(t *testing.T) {
	if i18n.Fold("İ") == i18n.Fold("i") {
		t.Error("Fold should keep İ and i distinct; use MatchKey for casual matching")
	}
	if i18n.Fold("über") == i18n.Fold("uber") {
		t.Error("Fold should keep über and uber distinct; use MatchKey for casual matching")
	}
	// MatchKey collapses, Fold does not.
	if i18n.MatchKey("İ") != i18n.MatchKey("i") {
		t.Error("MatchKey should collapse İ and i")
	}
}
