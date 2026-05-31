//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"path"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/asset"
	"github.com/ActiveMemory/ctx/internal/config/file"
)

func TestWhyDoc(t *testing.T) {
	tests := []struct {
		name        string
		doc         string
		wantContain string
		wantErr     bool
	}{
		{"manifesto exists", "manifesto", "Manifesto", false},
		{"about exists", "about", "ctx", false},
		{"design-invariants exists", "design-invariants", "Invariants", false},
		{"nonexistent returns error", "nonexistent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := FS.ReadFile(path.Join(asset.DirWhy, tt.doc+file.ExtMarkdown))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q", tt.doc)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tt.doc, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("content of %q does not contain %q", tt.doc, tt.wantContain)
			}
		})
	}
}

func TestListWhyDocs(t *testing.T) {
	entries, err := FS.ReadDir(asset.DirWhy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"about", "design-invariants", "manifesto"}
	docSet := make(map[string]bool)
	for _, e := range entries {
		name := e.Name()
		if strings.HasSuffix(name, file.ExtMarkdown) {
			docSet[strings.TrimSuffix(name, file.ExtMarkdown)] = true
		}
	}

	for _, exp := range expected {
		if !docSet[exp] {
			t.Errorf("missing expected doc: %s", exp)
		}
	}
}
