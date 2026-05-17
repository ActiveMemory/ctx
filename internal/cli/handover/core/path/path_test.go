//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package path_test

import (
	"os"
	"path/filepath"
	"testing"

	handoverPath "github.com/ActiveMemory/ctx/internal/cli/handover/core/path"
	cfgHandover "github.com/ActiveMemory/ctx/internal/config/handover"
)

// canonicalCtxDir builds a CTX_DIR honoring the rc-required
// `.context` basename and points the env var at it.
func canonicalCtxDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	ctxDir := filepath.Join(root, ".context")
	if err := os.MkdirAll(ctxDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	t.Setenv("CTX_DIR", ctxDir)
	return ctxDir
}

func TestDir(t *testing.T) {
	ctxDir := canonicalCtxDir(t)

	got, err := handoverPath.Dir()
	if err != nil {
		t.Fatalf("Dir: %v", err)
	}
	want := filepath.Join(ctxDir, cfgHandover.Subdir)
	if got != want {
		t.Errorf("Dir: want %q; got %q", want, got)
	}
}
