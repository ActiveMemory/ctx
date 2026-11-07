//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write_test

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/handover/cmd/write"
)

func TestCmd_RejectsMissingRequiredFlags(t *testing.T) {
	c := write.Cmd()
	c.SetArgs([]string{
		"Test Title",
		"--summary", "completed work",
		"--next", "next action",
	})
	err := c.Execute()
	if err == nil {
		t.Fatal("expected missing required flag error, got nil")
	}
	for _, flagName := range []string{"facts", "verified", "dead-ends"} {
		if !strings.Contains(err.Error(), flagName) {
			t.Errorf("expected error to mention %q, got: %v", flagName, err)
		}
	}
}

func TestRun_RejectsPlaceholderValues(t *testing.T) {
	cases := []struct {
		name      string
		summary   string
		facts     string
		verified  string
		deadEnds  string
		next      string
		wantField string
	}{
		{
			name:      "summary TBD",
			summary:   "TBD",
			facts:     "verified facts",
			verified:  "verified checks",
			deadEnds:  "no failed approaches",
			next:      "do next",
			wantField: "summary",
		},
		{
			name:      "facts n/a",
			summary:   "summary text",
			facts:     "n/a",
			verified:  "verified checks",
			deadEnds:  "no failed approaches",
			next:      "do next",
			wantField: "facts",
		},
		{
			name:      "verified pending",
			summary:   "summary text",
			facts:     "verified facts",
			verified:  "pending",
			deadEnds:  "no failed approaches",
			next:      "do next",
			wantField: "verified",
		},
		{
			name:      "dead-ends none",
			summary:   "summary text",
			facts:     "verified facts",
			verified:  "verified checks",
			deadEnds:  "none",
			next:      "do next",
			wantField: "dead-ends",
		},
		{
			name:      "next see chat",
			summary:   "summary text",
			facts:     "verified facts",
			verified:  "verified checks",
			deadEnds:  "no failed approaches",
			next:      "see chat",
			wantField: "next",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := write.Cmd()
			err := write.Run(
				c, "Title",
				tc.summary, tc.facts, tc.verified, tc.deadEnds, tc.next,
				"", "", "", true,
			)
			if err == nil {
				t.Fatalf("expected placeholder rejection for %s, got nil", tc.wantField)
			}
			if !strings.Contains(err.Error(), "--"+tc.wantField) {
				t.Errorf("expected error to name --%s, got: %v", tc.wantField, err)
			}
		})
	}
}
