//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
)

func TestInitialized_AllFilesPresent(t *testing.T) {
	tmp := t.TempDir()
	for _, f := range ctx.FilesRequired {
		path := filepath.Join(tmp, f)
		if writeErr := os.WriteFile(path, []byte("# "+f+"\n"), 0o600); writeErr != nil {
			t.Fatalf("setup: %v", writeErr)
		}
	}
	if !Initialized(tmp) {
		t.Error("Initialized() = false, want true when all required files present")
	}
}

func TestInitialized_MissingFile(t *testing.T) {
	tmp := t.TempDir()
	// Create all but the last required file.
	for _, f := range ctx.FilesRequired[:len(ctx.FilesRequired)-1] {
		path := filepath.Join(tmp, f)
		if writeErr := os.WriteFile(path, []byte("# "+f+"\n"), 0o600); writeErr != nil {
			t.Fatalf("setup: %v", writeErr)
		}
	}
	if Initialized(tmp) {
		t.Error("Initialized() = true, want false when a required file is missing")
	}
}

func TestInitialized_EmptyDir(t *testing.T) {
	tmp := t.TempDir()
	if Initialized(tmp) {
		t.Error("Initialized() = true, want false for empty directory")
	}
}
