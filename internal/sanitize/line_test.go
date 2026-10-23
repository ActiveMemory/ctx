//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

import (
	"testing"
)

// TestLine tests the Line helper function.
func TestLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{"ref characters", "feat/x.y_z-1", 0, "feat/x.y_z-1"},
		{"spaces", "Phase 1", 0, "Phase 1"},
		{"line breaks", "x\n- [x] forged\r\n", 0, "x- [x] forged"},
		{"tab and separators", "a\tb c d", 0, "abcd"},
		{"length cap", "abcdef", 3, "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Line(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("Line(%q, %d) = %q, want %q",
					tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}
