//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ceremony_test

import (
	"testing"

	coreCeremony "github.com/ActiveMemory/ctx/internal/cli/system/core/ceremony"
)

// TestInvokedByPrompt verifies that a live ceremony invocation is
// recognized from the UserPromptSubmit prompt in both the bare and
// plugin-scoped forms, with trailing arguments and leading whitespace,
// while unrelated prompts (including near-miss tokens) are not matched.
func TestInvokedByPrompt(t *testing.T) {
	cases := []struct {
		name   string
		prompt string
		want   bool
	}{
		{"bare remember", "/ctx-remember", true},
		{"bare wrapup", "/ctx-wrap-up", true},
		{"plugin remember", "/ctx:ctx-remember", true},
		{"plugin wrapup", "/ctx:ctx-wrap-up", true},
		{"remember with args", "/ctx-remember please", true},
		{"leading whitespace", "  /ctx-remember", true},
		{"non-ceremony command", "/ctx-status", false},
		{"prose mentioning command", "run /ctx-remember later", false},
		{"near-miss suffix", "/ctx-remembering", false},
		{"empty", "", false},
		{"whitespace only", "   ", false},
		{"plain prompt", "what were we working on?", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := coreCeremony.InvokedByPrompt(tc.prompt); got != tc.want {
				t.Errorf("InvokedByPrompt(%q) = %v, want %v",
					tc.prompt, got, tc.want)
			}
		})
	}
}
