//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"strings"
	"testing"
)

func TestList(t *testing.T) {
	templates, err := List()
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}

	if len(templates) == 0 {
		t.Error("List() returned empty list")
	}

	expected := []string{
		"learning.md",
		"decision.md",
	}

	templateSet := make(map[string]bool)
	for _, name := range templates {
		templateSet[name] = true
	}

	for _, exp := range expected {
		if !templateSet[exp] {
			t.Errorf("List() missing expected template: %s", exp)
		}
	}
}

func TestForName(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		wantContain string
		wantErr     bool
	}{
		{
			name:        "learning.md exists",
			template:    "learning.md",
			wantContain: "Context",
			wantErr:     false,
		},
		{
			name:        "decision.md exists",
			template:    "decision.md",
			wantContain: "Context",
			wantErr:     false,
		},
		{
			name:     "nonexistent entry template returns error",
			template: "nonexistent.md",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ForName(tt.template)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ForName(%q) expected error, got nil", tt.template)
				}
				return
			}
			if err != nil {
				t.Errorf("ForName(%q) unexpected error: %v", tt.template, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("ForName(%q) content does not contain %q",
					tt.template, tt.wantContain)
			}
		})
	}
}
