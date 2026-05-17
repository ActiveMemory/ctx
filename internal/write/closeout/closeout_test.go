//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	cfgGitmeta "github.com/ActiveMemory/ctx/internal/config/git_meta"
	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	errCloseout "github.com/ActiveMemory/ctx/internal/err/closeout"
	"github.com/ActiveMemory/ctx/internal/write/closeout"
)

// gitInit creates a real git repo at root so gitmeta.ResolveHead
// works without env overrides.
func gitInit(t *testing.T, root string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skipf("git not on PATH: %v", err)
	}
	for _, args := range [][]string{
		{"init", "-q"},
		{"config", "user.email", "test@example.com"},
		{"config", "user.name", "Test User"},
		{"commit", "--allow-empty", "-m", "init", "-q"},
	} {
		//nolint:gosec // G204: test fixture, args are hardcoded above
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
}

func TestWrite_RoundTripWithGitRepo(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)
	closeoutsDir := filepath.Join(root, "ingest", "closeouts")

	f, err := closeout.Write(
		closeoutsDir,
		root,
		cfgKB.CloseoutModeIngest,
		"topic-page",
		"bootstrap",
		"## Inputs\n\nfoo\n",
	)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if f.Path == "" {
		t.Fatal("Write returned empty path")
	}
	if !strings.HasSuffix(f.Path, "-ingest-closeout.md") {
		t.Errorf("filename suffix: got %s", f.Path)
	}
	if f.Frontmatter.Mode != cfgKB.CloseoutModeIngest {
		t.Errorf("mode: want %q; got %q", cfgKB.CloseoutModeIngest, f.Frontmatter.Mode)
	}
	if f.Frontmatter.PassMode != "topic-page" {
		t.Errorf("pass-mode: got %q", f.Frontmatter.PassMode)
	}
	if f.Frontmatter.LifeStage != "bootstrap" {
		t.Errorf("life-stage: got %q", f.Frontmatter.LifeStage)
	}
	if f.Frontmatter.SHA == "" {
		t.Error("SHA empty after Write")
	}
	if f.Frontmatter.GeneratedAt.IsZero() {
		t.Error("GeneratedAt zero after Write")
	}

	// Read it back from disk.
	got, err := closeout.Read(f.Path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got.Frontmatter.SHA != f.Frontmatter.SHA {
		t.Errorf("SHA round-trip mismatch")
	}
	if got.Frontmatter.Mode != f.Frontmatter.Mode {
		t.Errorf("Mode round-trip mismatch")
	}
	if !strings.Contains(got.Body, "foo") {
		t.Errorf("body lost: %q", got.Body)
	}
}

func TestWrite_RejectsEmptyMode(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)
	_, err := closeout.Write(
		filepath.Join(root, "closeouts"),
		root, "", "", "", "body",
	)
	if err == nil {
		t.Fatal("want error for empty mode; got nil")
	}
}

func TestWrite_UsesCtxTaskCommitOverride(t *testing.T) {
	root := t.TempDir()
	// No git init; rely entirely on CTX_TASK_COMMIT override.
	t.Setenv(cfgGitmeta.EnvCtxTaskCommit, "abc1234")
	t.Setenv(cfgGitmeta.EnvGithubActions, "")
	t.Setenv(cfgGitmeta.EnvGithubSHA, "")
	closeoutsDir := filepath.Join(root, "closeouts")

	f, err := closeout.Write(
		closeoutsDir, root,
		cfgKB.CloseoutModeAsk, "", "", "body",
	)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if f.Frontmatter.SHA != "abc1234" {
		t.Errorf("SHA: want abc1234; got %q", f.Frontmatter.SHA)
	}
}

func TestRead_MissingFrontmatter(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "broken.md")
	if err := os.WriteFile(path, []byte("just text, no delim\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := closeout.Read(path)
	if !errors.Is(err, errCloseout.ErrMissingFrontmatter) {
		t.Fatalf("want ErrMissingFrontmatter; got %v", err)
	}
}

func TestRead_MissingFields(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "incomplete.md")
	body := "---\nsha: abc\nbranch: main\n---\n\nbody\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := closeout.Read(path)
	if !errors.Is(err, errCloseout.ErrMissingFields) {
		t.Fatalf("want ErrMissingFields; got %v", err)
	}
}

func TestList_SortedAscending(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)
	closeoutsDir := filepath.Join(root, "closeouts")

	for i, mode := range []string{
		cfgKB.CloseoutModeIngest,
		cfgKB.CloseoutModeAsk,
		cfgKB.CloseoutModeGround,
	} {
		// Force ordering by writing sequentially with a small
		// real-time gap; the test does not depend on exact
		// timing, only on monotonic order at the granularity
		// of the timestamps we use.
		if _, err := closeout.Write(
			closeoutsDir, root, mode, "", "", "body",
		); err != nil {
			t.Fatalf("Write %d: %v", i, err)
		}
		time.Sleep(1100 * time.Millisecond)
	}

	files, bad, err := closeout.List(closeoutsDir)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(bad) != 0 {
		t.Errorf("unexpected bad files: %v", bad)
	}
	if len(files) != 3 {
		t.Fatalf("count: want 3; got %d", len(files))
	}
	for i := 1; i < len(files); i++ {
		prev := files[i-1].Frontmatter.GeneratedAt
		curr := files[i].Frontmatter.GeneratedAt
		if curr.Before(prev) {
			t.Errorf("List not ascending at index %d", i)
		}
	}
}

func TestPostdatedBy(t *testing.T) {
	now := time.Now()
	files := []closeout.File{
		{Frontmatter: closeout.Frontmatter{GeneratedAt: now.Add(-2 * time.Hour)}},
		{Frontmatter: closeout.Frontmatter{GeneratedAt: now.Add(-1 * time.Hour)}},
		{Frontmatter: closeout.Frontmatter{GeneratedAt: now}},
	}
	cursor := now.Add(-90 * time.Minute)
	got := closeout.PostdatedBy(files, cursor)
	if len(got) != 2 {
		t.Errorf("want 2 postdated; got %d", len(got))
	}
}

func TestArchive_MovesFiles(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)
	closeoutsDir := filepath.Join(root, "closeouts")
	archiveDir := filepath.Join(root, "archive")

	f, err := closeout.Write(
		closeoutsDir, root,
		cfgKB.CloseoutModeIngest, "", "", "body",
	)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := closeout.Archive(archiveDir, []closeout.File{f}); err != nil {
		t.Fatalf("Archive: %v", err)
	}
	if _, err := os.Stat(f.Path); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("source still exists: %v", err)
	}
	moved := filepath.Join(archiveDir, filepath.Base(f.Path))
	if _, err := os.Stat(moved); err != nil {
		t.Errorf("archived file missing: %v", err)
	}
}
