//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
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

// newCmd builds a minimal cobra command with two string flags and
// a no-op RunE — the standard test fixture for the BodyFlags
// PreRunE pattern that the noun-level constructors use.
func newCmd() *cobra.Command {
	c := &cobra.Command{
		Use:  "test",
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	}
	c.Flags().String("context", "", "")
	c.Flags().String("rationale", "", "")
	c.SetOut(&strings.Builder{})
	c.SetErr(&strings.Builder{})
	return c
}

func TestBodyFlagsAcceptsLegitimateValues(t *testing.T) {
	c := newCmd()
	c.SetArgs([]string{
		"--context", "real context",
		"--rationale", "real rationale",
	})
	if execErr := c.Execute(); execErr != nil {
		t.Fatalf("parse: %v", execErr)
	}
	if err := BodyFlags(c, "context", "rationale"); err != nil {
		t.Errorf("BodyFlags rejected legitimate input: %v", err)
	}
}

func TestBodyFlagsRejectsPlaceholder(t *testing.T) {
	c := newCmd()
	c.SetArgs([]string{
		"--context", "TBD",
		"--rationale", "good reason",
	})
	if execErr := c.Execute(); execErr != nil {
		t.Fatalf("parse: %v", execErr)
	}
	err := BodyFlags(c, "context", "rationale")
	if err == nil {
		t.Fatal("expected placeholder rejection")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Errorf("error should name the offending flag: %v", err)
	}
}

func TestBodyFlagsRejectsMissingFlag(t *testing.T) {
	// Cobra defaults --context to "" when omitted; the empty-value
	// check catches it through the same path as a placeholder.
	c := newCmd()
	c.SetArgs([]string{"--rationale", "ok"})
	if execErr := c.Execute(); execErr != nil {
		t.Fatalf("parse: %v", execErr)
	}
	err := BodyFlags(c, "context", "rationale")
	if err == nil {
		t.Fatal("expected rejection when --context is missing")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Errorf("error should name --context: %v", err)
	}
}

func TestBodyFlagsStopsAtFirstFailure(t *testing.T) {
	// Flag order in the call determines which failure is reported
	// when multiple flags fail.
	c := newCmd()
	c.SetArgs([]string{
		"--context", "TBD",
		"--rationale", "n/a",
	})
	if execErr := c.Execute(); execErr != nil {
		t.Fatalf("parse: %v", execErr)
	}
	err := BodyFlags(c, "rationale", "context")
	if err == nil {
		t.Fatal("expected rejection")
	}
	if !strings.Contains(err.Error(), "rationale") {
		t.Errorf("expected --rationale to be reported first, got %v", err)
	}
}
