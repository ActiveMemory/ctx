//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import (
	"testing"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0B"},
		{100, "100B"},
		{1023, "1023B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{10240, "10.0KB"},
		{1048576, "1.0MB"},
		{1572864, "1.5MB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := Size(tt.bytes)
			if got != tt.want {
				t.Errorf("Size(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestKeyFileSlug(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"internal/config/limit.go", "internal_config_limit_go"},
		{"cmd/*.go", "cmd_x_go"},
		{"simple.txt", "simple_txt"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := KeyFileSlug(tt.path)
			if got != tt.want {
				t.Errorf("KeyFileSlug(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
