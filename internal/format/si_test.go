//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
)

func init() { lookup.Init() }

func TestTokens(t *testing.T) {
	tests := []struct {
		name   string
		tokens int
		want   string
	}{
		{"zero", 0, "0"},
		{"small", 500, "500"},
		{"below-K", 999, "999"},
		{"exactly-1K", 1000, "1.0K"},
		{"mid-K", 1500, "1.5K"},
		{"large-K", 50000, "50.0K"},
		{"below-M", 999999, "1000.0K"},
		{"exactly-1M", 1000000, "1.0M"},
		{"mid-M", 2300000, "2.3M"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Tokens(tt.tokens)
			if got != tt.want {
				t.Errorf("Tokens(%d) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}
