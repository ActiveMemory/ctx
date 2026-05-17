//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	"github.com/ActiveMemory/ctx/internal/write/closeout"
	"github.com/ActiveMemory/ctx/internal/write/handover"
)

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

func TestWrite_HappyPathWithFold(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)
	handoversDir := filepath.Join(root, "handovers")
	closeoutsDir := filepath.Join(root, "closeouts")
	archiveDir := filepath.Join(root, "archive")

	// Pre-create a closeout that will be folded.
	_, err := closeout.Write(
		closeoutsDir, root, cfgKB.CloseoutModeIngest,
		"topic-page", "bootstrap", "## What changed\n\nstuff\n",
	)
	if err != nil {
		t.Fatalf("seed closeout: %v", err)
	}

	res, err := handover.Write(
		handoversDir, closeoutsDir, archiveDir, root,
		handover.Entry{
			Title:   "First Session",
			Summary: "did the thing",
			Next:    "do the other thing",
		},
	)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if res.File.Path == "" {
		t.Fatal("Write returned empty path")
	}
	if !strings.HasSuffix(res.File.Path, "-first-session.md") {
		t.Errorf("filename: got %s", res.File.Path)
	}
	if len(res.FoldedCloseouts) != 1 {
		t.Errorf("folded count: want 1; got %d", len(res.FoldedCloseouts))
	}
	if !strings.Contains(res.File.Body, "## Folded closeouts") {
		t.Errorf("folded section missing from body")
	}
	// Source closeout should have been archived.
	if _, statErr := os.Stat(res.FoldedCloseouts[0].Path); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("folded closeout still in source: %v", statErr)
	}
}

func TestWrite_NoFoldKeepsCloseouts(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)
	handoversDir := filepath.Join(root, "handovers")
	closeoutsDir := filepath.Join(root, "closeouts")
	archiveDir := filepath.Join(root, "archive")

	f, err := closeout.Write(
		closeoutsDir, root, cfgKB.CloseoutModeIngest, "", "", "body",
	)
	if err != nil {
		t.Fatalf("seed closeout: %v", err)
	}

	res, err := handover.Write(
		handoversDir, closeoutsDir, archiveDir, root,
		handover.Entry{
			Title:   "Mid-session checkpoint",
			Summary: "checkpoint",
			Next:    "resume",
			NoFold:  true,
		},
	)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	if len(res.FoldedCloseouts) != 0 {
		t.Errorf("expected no folds; got %d", len(res.FoldedCloseouts))
	}
	if _, statErr := os.Stat(f.Path); statErr != nil {
		t.Errorf("closeout should still exist after --no-fold: %v", statErr)
	}
}

func TestWrite_RejectsEmptyFields(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)
	dirs := struct{ h, c, a string }{
		filepath.Join(root, "handovers"),
		filepath.Join(root, "closeouts"),
		filepath.Join(root, "archive"),
	}
	cases := []handover.Entry{
		{Title: "", Summary: "s", Next: "n"},
		{Title: "t", Summary: "", Next: "n"},
		{Title: "t", Summary: "s", Next: ""},
	}
	for i, e := range cases {
		_, err := handover.Write(dirs.h, dirs.c, dirs.a, root, e)
		if err == nil {
			t.Errorf("case %d: expected error", i)
		}
	}
}

func TestLatest_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	latestAt, path, err := handover.Latest(dir)
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if !latestAt.IsZero() {
		t.Errorf("latest: want zero; got %v", latestAt)
	}
	if path != "" {
		t.Errorf("path: want empty; got %s", path)
	}
}

func TestLatest_ReturnsMostRecent(t *testing.T) {
	root := t.TempDir()
	gitInit(t, root)
	handoversDir := filepath.Join(root, "handovers")
	closeoutsDir := filepath.Join(root, "closeouts")
	archiveDir := filepath.Join(root, "archive")

	for i, title := range []string{"first", "second"} {
		_, err := handover.Write(
			handoversDir, closeoutsDir, archiveDir, root,
			handover.Entry{
				Title:   title,
				Summary: "s",
				Next:    "n",
				NoFold:  true,
			},
		)
		if err != nil {
			t.Fatalf("Write %d: %v", i, err)
		}
		time.Sleep(1100 * time.Millisecond)
	}

	latestAt, path, err := handover.Latest(handoversDir)
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if latestAt.IsZero() {
		t.Error("latest zero after writes")
	}
	if !strings.HasSuffix(path, "-second.md") {
		t.Errorf("want second handover; got %s", path)
	}
}
