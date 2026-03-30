//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sanitize

import (
	"testing"
)

// TestFilename tests the Filename helper function.
func TestFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple topic", "simple-topic"},
		{"Uppercase Topic", "uppercase-topic"},
		{"topic with   multiple   spaces", "topic-with-multiple-spaces"},
		{"special!@#$%chars", "special-chars"},
		{"already-valid", "already-valid"},
		{"", "session"},
		{"   ", "session"},
		{"---", "session"},
		{
			"a very long topic name that exceeds" +
				" the maximum allowed length" +
				" of fifty characters",
			"a-very-long-topic-name-that-exceeds-the-maximum-al",
		},
		{"trailing---", "trailing"},
		{"---leading", "leading"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Filename(tt.input)
			if result != tt.expected {
				t.Errorf("Filename(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}
