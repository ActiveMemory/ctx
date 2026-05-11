//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"strings"
	"testing"
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
