//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"errors"
	"testing"

	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"

	"github.com/ActiveMemory/ctx/internal/disclosure"
)

// T01: the entry prefix is per-kind — timestamped kinds open an entry
// with "## [", a convention root with a bare "## ". This is the whole
// structural difference the unified parse path is parametrized on.
func TestEntryPrefix(t *testing.T) {
	tests := []struct {
		name string
		kind disclosure.Kind
		want string
	}{
		{"learning is timestamped", disclosure.KindLearning, "## ["},
		{"decision is timestamped", disclosure.KindDecision, "## ["},
		{"convention is a prose section", disclosure.KindConvention, "## "},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := disclosure.EntryPrefix(tc.kind); got != tc.want {
				t.Errorf("EntryPrefix(%v) = %q, want %q",
					tc.kind, got, tc.want)
			}
		})
	}
}

// T03: the duplicate-title sentinel resolves to embedded text. A
// sentinel whose message comes back empty is the "uninitialized
// desc.Text" failure mode, and it hides the reason a digest refused.
func TestErrDuplicateStagedTitle_Resolves(t *testing.T) {
	if errDisc.ErrDuplicateStagedTitle.Error() == "" {
		t.Fatal("ErrDuplicateStagedTitle resolves to empty text")
	}
	wrapped := errors.New("wrapped: " + errDisc.ErrDuplicateStagedTitle.Error())
	if wrapped.Error() == "wrapped: " {
		t.Error("sentinel text did not survive wrapping")
	}
}
