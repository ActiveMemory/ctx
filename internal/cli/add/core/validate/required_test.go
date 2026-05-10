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

func TestRequireBodyFlagsRejectsMissingFlag(t *testing.T) {
	c := &cobra.Command{Use: "test", Run: func(_ *cobra.Command, _ []string) {}}
	c.Flags().String("context", "", "")
	err := RequireBodyFlags(c, "context", "nonexistent")
	if err == nil {
		t.Fatal("expected error for unregistered flag")
	}
}

func TestRequireBodyFlagsMarksRequired(t *testing.T) {
	c := &cobra.Command{Use: "test", Run: func(_ *cobra.Command, _ []string) {}}
	c.Flags().String("context", "", "")
	c.Flags().String("rationale", "", "")
	if err := RequireBodyFlags(c, "context", "rationale"); err != nil {
		t.Fatalf("RequireBodyFlags error: %v", err)
	}
	// Missing --context and --rationale: PreRunE fires first and
	// reports the empty value before cobra's required-flag check.
	c.SetArgs([]string{})
	c.SetOut(&strings.Builder{})
	c.SetErr(&strings.Builder{})
	err := c.Execute()
	if err == nil {
		t.Fatal("expected error for missing required body flags")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Errorf("error should name the empty flag: %v", err)
	}
}

// TestRequireBodyFlagsCobraRequiredAnnotation verifies that
// MarkFlagRequired's annotation is set on each flag. The
// PreRunE shadow normally fires first; with PreRunE bypassed,
// cobra still rejects the missing flag at validateRequiredFlags.
func TestRequireBodyFlagsCobraRequiredAnnotation(t *testing.T) {
	c := &cobra.Command{Use: "test", Run: func(_ *cobra.Command, _ []string) {}}
	c.Flags().String("context", "", "")
	if err := RequireBodyFlags(c, "context"); err != nil {
		t.Fatalf("RequireBodyFlags error: %v", err)
	}
	flag := c.Flag("context")
	if flag == nil {
		t.Fatal("context flag missing")
	}
	if _, ok := flag.Annotations[cobra.BashCompOneRequiredFlag]; !ok {
		t.Error("expected cobra required annotation on --context")
	}
}

func TestRequireBodyFlagsRejectsPlaceholderViaPreRunE(t *testing.T) {
	ran := false
	c := &cobra.Command{
		Use: "test",
		RunE: func(_ *cobra.Command, _ []string) error {
			ran = true
			return nil
		},
	}
	c.Flags().String("context", "", "")
	c.Flags().String("rationale", "", "")
	if err := RequireBodyFlags(c, "context", "rationale"); err != nil {
		t.Fatalf("RequireBodyFlags error: %v", err)
	}
	c.SetArgs([]string{
		"--context", "TBD",
		"--rationale", "good reason",
	})
	c.SetOut(&strings.Builder{})
	c.SetErr(&strings.Builder{})
	err := c.Execute()
	if err == nil {
		t.Fatal("expected placeholder rejection")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Errorf("error should name the offending flag: %v", err)
	}
	if ran {
		t.Error("RunE should not have executed after PreRunE rejected input")
	}
}

func TestRequireBodyFlagsAllowsLegitimateInput(t *testing.T) {
	ran := false
	c := &cobra.Command{
		Use: "test",
		RunE: func(_ *cobra.Command, _ []string) error {
			ran = true
			return nil
		},
	}
	c.Flags().String("context", "", "")
	c.Flags().String("rationale", "", "")
	if err := RequireBodyFlags(c, "context", "rationale"); err != nil {
		t.Fatalf("RequireBodyFlags error: %v", err)
	}
	c.SetArgs([]string{
		"--context", "real context",
		"--rationale", "real rationale",
	})
	c.SetOut(&strings.Builder{})
	c.SetErr(&strings.Builder{})
	if err := c.Execute(); err != nil {
		t.Fatalf("legitimate input rejected: %v", err)
	}
	if !ran {
		t.Error("RunE should have executed")
	}
}

func TestRequireBodyFlagsPreservesExistingPreRunE(t *testing.T) {
	preRan := false
	c := &cobra.Command{
		Use: "test",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			preRan = true
			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	}
	c.Flags().String("context", "", "")
	if err := RequireBodyFlags(c, "context"); err != nil {
		t.Fatalf("RequireBodyFlags error: %v", err)
	}
	c.SetArgs([]string{"--context", "legitimate"})
	c.SetOut(&strings.Builder{})
	c.SetErr(&strings.Builder{})
	if err := c.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !preRan {
		t.Error("existing PreRunE should still execute")
	}
}
