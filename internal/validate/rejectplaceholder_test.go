//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

func TestRejectPlaceholderAcceptsLegitimate(t *testing.T) {
	for _, v := range []string{
		"a real rationale",
		"because PostgreSQL is well-supported",
		"we need TBD-style markers in the body but the field is real",
		"see above the line break",
	} {
		if err := RejectPlaceholder("context", v); err != nil {
			t.Errorf("RejectPlaceholder(%q) = %v, want nil", v, err)
		}
	}
}

func TestRejectPlaceholderRejectsExactMatches(t *testing.T) {
	for _, v := range []string{
		"TBD", "tbd", "Tbd",
		"N/A", "n/a", "na",
		"see chat", "See Chat",
		"see above", "see below",
		"pending", "PENDING",
		"none", "to be done",
	} {
		if err := RejectPlaceholder("context", v); err == nil {
			t.Errorf("RejectPlaceholder(%q) = nil, want error", v)
		}
	}
}

func TestRejectPlaceholderRejectsWhitespace(t *testing.T) {
	for _, v := range []string{
		"",
		" ",
		"\t",
		"\n",
		"   \t  \n  ",
	} {
		err := RejectPlaceholder("rationale", v)
		if err == nil {
			t.Errorf("RejectPlaceholder(%q) = nil, want error", v)
		}
		if !strings.Contains(err.Error(), "rationale") {
			t.Errorf("error should name flag 'rationale': %v", err)
		}
	}
}

func TestRejectPlaceholderTrimsBeforeMatching(t *testing.T) {
	// "  TBD  " is still a placeholder after trim.
	err := RejectPlaceholder("consequence", "  TBD  ")
	if err == nil {
		t.Error("padded placeholder should still be rejected")
	}
}

// TestRejectPlaceholderHonorsCtxrcExtensions wires the
// whole flow end-to-end: seed a .ctxrc with a custom
// `placeholders:` list, verify both the shipped default
// (tbd) and the user-supplied entry (iptal) reject. Guards
// against the validator silently dropping the merge step.
func TestRejectPlaceholderHonorsCtxrcExtensions(t *testing.T) {
	tmpDir := t.TempDir()
	ctxDir := filepath.Join(tmpDir, ".context")
	if err := os.MkdirAll(ctxDir, 0o755); err != nil {
		t.Fatal(err)
	}
	rcContent := "placeholders:\n  - iptal\n  - yapılacak\n"
	if err := os.WriteFile(
		filepath.Join(tmpDir, ".ctxrc"), []byte(rcContent), 0o644,
	); err != nil {
		t.Fatal(err)
	}
	testctx.Declare(t, tmpDir)

	for _, v := range []string{"tbd", "TBD", "iptal", "IPTAL", "Iptal", "yapılacak"} {
		if err := RejectPlaceholder("context", v); err == nil {
			t.Errorf("RejectPlaceholder(%q) = nil, want error (default+user merge)", v)
		}
	}
	// A non-placeholder still passes through.
	if err := RejectPlaceholder("context", "a real reason"); err != nil {
		t.Errorf("RejectPlaceholder(%q) = %v, want nil", "a real reason", err)
	}
	// Turkish dotted-I sanity check: the user list contains plain
	// "iptal", so "İPTAL" (with dotted İ) does NOT match — i18n.Fold
	// preserves the linguistic distinction between İ and i, by
	// design (see specs/i18n-fold-helper-and-ban.md). A user who
	// wants both rejected adds both spellings to .ctxrc.
	if err := RejectPlaceholder("context", "İPTAL"); err != nil {
		t.Errorf("RejectPlaceholder(\"İPTAL\") = %v, want nil "+
			"(distinct from user-supplied \"iptal\" under Unicode fold)", err)
	}
}
