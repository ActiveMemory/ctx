//go:build darwin

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import "testing"

func TestParsePressureLevel(t *testing.T) {
	tests := []struct {
		name          string
		output        string
		wantSev       Severity
		wantSupported bool
	}{
		{"normal", "1\n", SeverityOK, true},
		{"warning", "2\n", SeverityWarning, true},
		{"critical", "4\n", SeverityDanger, true},
		{"normal no newline", "1", SeverityOK, true},
		{"unrecognized value", "3\n", SeverityOK, false},
		{"unrecognized high value", "99", SeverityOK, false},
		{"non-numeric", "n/a\n", SeverityOK, false},
		{"empty", "", SeverityOK, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sev, supported := parsePressureLevel(tt.output)
			if supported != tt.wantSupported {
				t.Errorf("supported = %v, want %v",
					supported, tt.wantSupported)
			}
			if sev != tt.wantSev {
				t.Errorf("severity = %v, want %v", sev, tt.wantSev)
			}
		})
	}
}
