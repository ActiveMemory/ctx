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

func TestFold_ASCIIByteIdentical(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"hello", "hello"},
		{"HELLO", "hello"},
		{"Hello World", "hello world"},
		{"123-ABC_xyz", "123-abc_xyz"},
		{"TBD", "tbd"},
		{"n/a", "n/a"},
	}
	for _, tt := range cases {
		got := i18n.Fold(tt.in)
		if got != tt.want {
			t.Errorf("Fold(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// TestFold_UnicodeCorrectness guards the i18n properties
// that strings.ToLower violates. Without this test, a
// future regression to strings.ToLower would still pass
// the ASCII suite above.
//
// Note on Turkish: İ folds to "i" + COMBINING DOT ABOVE,
// which is distinct from plain "i". This is *correct*
// Unicode behavior — Turkish dotted-I and plain-i are
// linguistically distinct letters. Fold makes İPTAL match
// İptal (both produce the same folded form), which is the
// property the placeholder validator needs. It does not
// promise to collapse İ to i; that would require either
// Turkish-locale lowercasing or NFKD + combining-mark
// stripping, neither of which is a comparison primitive.
func TestFold_UnicodeCorrectness(t *testing.T) {
	cases := []struct {
		name     string
		a, b     string
		wantEqal bool
	}{
		{
			name:     "Turkish: İPTAL folds equal to İptal",
			a:        "İPTAL",
			b:        "İptal",
			wantEqal: true,
		},
		{
			name:     "Turkish: İ folded does NOT equal plain i (correct)",
			a:        "İ",
			b:        "i",
			wantEqal: false,
		},
		{
			name:     "German sharp S folds to ss",
			a:        "STRASSE",
			b:        "Straße",
			wantEqal: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gotA, gotB := i18n.Fold(tt.a), i18n.Fold(tt.b)
			equal := gotA == gotB
			if equal != tt.wantEqal {
				t.Errorf(
					"Fold(%q)=%q, Fold(%q)=%q; equal=%v, want %v",
					tt.a, gotA, tt.b, gotB, equal, tt.wantEqal,
				)
			}
		})
	}
}
