//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import "testing"

func TestDate(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"2026-03-05", false},
		{"2026-01-01", false},
		{"", false},
		{"not-a-date", true},
		{"2026/03/05", true},
		{"03-05-2026", true},
	}
	for _, tt := range tests {
		_, err := Date(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("Date(%q): err=%v, wantErr=%v", tt.input, err, tt.wantErr)
		}
	}
}
