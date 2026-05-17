//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gitmeta_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	errGitmeta "github.com/ActiveMemory/ctx/internal/err/git_meta"
	"github.com/ActiveMemory/ctx/internal/git_meta"
)

func TestRequireGitTree_DirAccepted(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := gitmeta.RequireGitTree(root); err != nil {
		t.Fatalf("want nil; got %v", err)
	}
}

func TestRequireGitTree_WorktreePointerFileAccepted(t *testing.T) {
	root := t.TempDir()
	// Worktrees write a regular file at .git containing
	// "gitdir: <path>" instead of a directory.
	if err := os.WriteFile(
		filepath.Join(root, ".git"),
		[]byte("gitdir: /main/.git/worktrees/x\n"),
		0o600,
	); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := gitmeta.RequireGitTree(root); err != nil {
		t.Fatalf("want nil; got %v", err)
	}
}

func TestRequireGitTree_MissingReturnsSentinel(t *testing.T) {
	root := t.TempDir()
	err := gitmeta.RequireGitTree(root)
	if err == nil {
		t.Fatal("want error; got nil")
	}
	if !errors.Is(err, errGitmeta.ErrMissingGitTree) {
		t.Fatalf("want ErrMissingGitTree; got %v", err)
	}
	if !strings.Contains(err.Error(), root) {
		t.Errorf("want project root in error message: %q", err.Error())
	}
}

func TestMissingGitTreeForCmd_FormatsCommandPrefix(t *testing.T) {
	err := errGitmeta.MissingGitTreeForCmd("init", "/tmp/x")
	if !errors.Is(err, errGitmeta.ErrMissingGitTree) {
		t.Fatalf("want errors.Is match against ErrMissingGitTree; got %v", err)
	}
	msg := err.Error()
	if !strings.Contains(msg, "ctx init") {
		t.Errorf("want subcommand name in message: %q", msg)
	}
	if !strings.Contains(msg, "/tmp/x") {
		t.Errorf("want project root in message: %q", msg)
	}
	if !strings.Contains(msg, "git init") {
		t.Errorf("want recovery hint in message: %q", msg)
	}
}
