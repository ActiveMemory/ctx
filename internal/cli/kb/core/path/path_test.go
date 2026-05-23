//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package path_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	kbPath "github.com/ActiveMemory/ctx/internal/cli/kb/core/path"
	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
)

// canonicalCtxDir creates $TMP/.context, chdirs into $TMP so the
// cwd-anchored resolver finds it, and returns the absolute ctx-dir
// path.
func canonicalCtxDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	ctxDir := filepath.Join(root, ".context")
	if err := os.MkdirAll(ctxDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	t.Chdir(root)
	return ctxDir
}

func TestKBDir(t *testing.T) {
	ctxDir := canonicalCtxDir(t)

	got, err := kbPath.KBDir()
	if err != nil {
		t.Fatalf("KBDir: %v", err)
	}
	want := filepath.Join(ctxDir, cfgKB.KBSubdir)
	if got != want {
		t.Errorf("KBDir: want %q; got %q", want, got)
	}
}

func TestKBTopicIndexFile(t *testing.T) {
	ctxDir := canonicalCtxDir(t)

	got, err := kbPath.KBTopicIndexFile("cursor-hooks")
	if err != nil {
		t.Fatalf("KBTopicIndexFile: %v", err)
	}
	want := filepath.Join(
		ctxDir, cfgKB.KBSubdir, cfgKB.TopicsSubdir,
		"cursor-hooks", cfgKB.TopicIndex,
	)
	if got != want {
		t.Errorf("KBTopicIndexFile: want %q; got %q", want, got)
	}
	// Sanity: index.md must end the path; topic slug must be
	// the parent.
	if !strings.HasSuffix(got, "/cursor-hooks/index.md") {
		t.Errorf("path suffix unexpected: %s", got)
	}
}

func TestIngestArtifactFile(t *testing.T) {
	ctxDir := canonicalCtxDir(t)

	got, err := kbPath.IngestArtifactFile(cfgKB.Rules)
	if err != nil {
		t.Fatalf("IngestArtifactFile: %v", err)
	}
	want := filepath.Join(ctxDir, cfgKB.IngestSubdir, cfgKB.Rules)
	if got != want {
		t.Errorf("IngestArtifactFile: want %q; got %q", want, got)
	}
}

func TestCloseoutsDir(t *testing.T) {
	ctxDir := canonicalCtxDir(t)

	got, err := kbPath.CloseoutsDir()
	if err != nil {
		t.Fatalf("CloseoutsDir: %v", err)
	}
	want := filepath.Join(
		ctxDir, cfgKB.IngestSubdir, cfgKB.CloseoutsSubdir,
	)
	if got != want {
		t.Errorf("CloseoutsDir: want %q; got %q", want, got)
	}
}

func TestArchiveCloseoutsDir(t *testing.T) {
	_ = canonicalCtxDir(t)

	got, err := kbPath.ArchiveCloseoutsDir()
	if err != nil {
		t.Fatalf("ArchiveCloseoutsDir: %v", err)
	}
	if !strings.HasSuffix(got, "/archive/closeouts") {
		t.Errorf("path suffix unexpected: %s", got)
	}
}

func TestSiteKBDir(t *testing.T) {
	ctxDir := canonicalCtxDir(t)

	got, err := kbPath.SiteKBDir()
	if err != nil {
		t.Fatalf("SiteKBDir: %v", err)
	}
	want := filepath.Join(
		ctxDir, cfgKB.SiteSubdir, cfgKB.SiteKBSubdir,
	)
	if got != want {
		t.Errorf("SiteKBDir: want %q; got %q", want, got)
	}
}
