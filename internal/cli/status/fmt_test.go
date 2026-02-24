//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"testing"
	"time"
)

func TestFormatTimeAgo(t *testing.T) {
	tests := []struct {
		name   string
		offset time.Duration
		want   string
	}{
		{"just now", 5 * time.Second, "just now"},
		{"1 minute", 90 * time.Second, "1 minute ago"},
		{"several minutes", 5*time.Minute + 10*time.Second, "5 minutes ago"},
		{"59 minutes", 59*time.Minute + 30*time.Second, "59 minutes ago"},
		{"1 hour", 90 * time.Minute, "1 hour ago"},
		{"several hours", 5*time.Hour + 10*time.Minute, "5 hours ago"},
		{"1 day", 36 * time.Hour, "1 day ago"},
		{"several days", 4 * 24 * time.Hour, "4 days ago"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTimeAgo(time.Now().Add(-tt.offset))
			if got != tt.want {
				t.Errorf("formatTimeAgo() = %q, want %q", got, tt.want)
			}
		})
	}

	// Old date (>7 days) returns formatted date
	t.Run("old date", func(t *testing.T) {
		old := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
		got := formatTimeAgo(old)
		if got != "Jun 15, 2025" {
			t.Errorf("formatTimeAgo() = %q, want %q", got, "Jun 15, 2025")
		}
	})
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"zero", 0, "0"},
		{"small", 500, "500"},
		{"below-threshold", 999, "999"},
		{"exactly-1000", 1000, "1,000"},
		{"mid-range", 1500, "1,500"},
		{"large", 12345, "12,345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatNumber(tt.n)
			if got != tt.want {
				t.Errorf("formatNumber(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name string
		b    int64
		want string
	}{
		{"zero", 0, "0 B"},
		{"small bytes", 500, "500 B"},
		{"below-KB", 1023, "1023 B"},
		{"1 KB", 1024, "1.0 KB"},
		{"1.5 KB", 1536, "1.5 KB"},
		{"1 MB", 1024 * 1024, "1.0 MB"},
		{"1 GB", 1024 * 1024 * 1024, "1.0 GB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatBytes(tt.b)
			if got != tt.want {
				t.Errorf("formatBytes(%d) = %q, want %q", tt.b, got, tt.want)
			}
		})
	}
}
