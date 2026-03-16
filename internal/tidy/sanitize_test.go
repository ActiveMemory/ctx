//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tidy

import (
	"testing"
)

// TestTruncateString tests the TruncateString helper function.
func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a longer string", 10, "this is..."},
		{"", 10, ""},
		{"hello — world of things", 10, "hello —..."},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := TruncateString(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
			}
		})
	}
}

// TestRemoveEmptySections tests the RemoveEmptySections helper function.
func TestRemoveEmptySections(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		removed  int
	}{
		{
			name:     "no empty sections",
			input:    "# Title\n\n## Section\n\nContent here\n",
			expected: "# Title\n\n## Section\n\nContent here\n",
			removed:  0,
		},
		{
			name:     "single empty section",
			input:    "# Title\n\n## Empty\n\n## HasContent\n\nSome content\n",
			expected: "# Title\n\n## HasContent\n\nSome content\n",
			removed:  1,
		},
		{
			name:     "multiple empty sections",
			input:    "# Title\n\n## Empty1\n\n## Empty2\n\n## HasContent\n\nContent\n",
			expected: "# Title\n\n## HasContent\n\nContent\n",
			removed:  2,
		},
		{
			name:     "empty section at end",
			input:    "# Title\n\n## Content\n\nText\n\n## EmptyAtEnd\n",
			expected: "# Title\n\n## Content\n\nText\n",
			removed:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, count := RemoveEmptySections(tt.input)
			if count != tt.removed {
				t.Errorf("RemoveEmptySections() removed %d sections, want %d", count, tt.removed)
			}
			if result != tt.expected {
				t.Errorf("RemoveEmptySections() result:\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}
