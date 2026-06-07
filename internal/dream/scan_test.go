//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream_test

import (
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/dream"
)

// TestScanIdeasMarkdownOnly hashes only markdown under ideas/, skipping
// the done/ archive and non-markdown files.
func TestScanIdeasMarkdownOnly(t *testing.T) {
	root := t.TempDir()
	ideas := filepath.Join(root, "ideas")
	mustWrite(t, filepath.Join(ideas, "a.md"), "alpha")
	mustWrite(t, filepath.Join(ideas, "b.md"), "beta")
	mustWrite(t, filepath.Join(ideas, "binary.png"), "not markdown")
	mustWrite(t, filepath.Join(ideas, "done", "old.md"), "archived")

	got, err := dream.ScanIdeas(root, ideas)
	if err != nil {
		t.Fatalf("ScanIdeas: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("scanned %d files, want 2: %v", len(got), got)
	}
	if _, ok := got[filepath.Join("ideas", "a.md")]; !ok {
		t.Fatal("ideas/a.md missing from scan")
	}
	if _, ok := got[filepath.Join("ideas", "done", "old.md")]; ok {
		t.Fatal("ideas/done/ must be excluded")
	}
}

// TestScanIdeasMissingDir yields an empty map (not an error) when
// ideas/ does not exist.
func TestScanIdeasMissingDir(t *testing.T) {
	root := t.TempDir()
	got, err := dream.ScanIdeas(root, filepath.Join(root, "ideas"))
	if err != nil {
		t.Fatalf("ScanIdeas missing dir: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("scanned %d files, want 0", len(got))
	}
}

// mustWrite writes content to path, creating parents, failing the test
// on error.
func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := writeFixture(path, content); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}
}
