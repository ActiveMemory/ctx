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
