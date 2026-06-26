//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package run

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func declareContext(t *testing.T, content string) {
	t.Helper()
	tempDir := t.TempDir()
	ctxDir := filepath.Join(tempDir, dir.Context)
	if mkErr := os.MkdirAll(ctxDir, 0700); mkErr != nil {
		t.Fatalf("mkdir .context: %v", mkErr)
	}
	if content != "" {
		rcPath := filepath.Join(tempDir, ".ctxrc")
		if wrErr := os.WriteFile(rcPath, []byte(content), 0600); wrErr != nil {
			t.Fatalf("write .ctxrc: %v", wrErr)
		}
	}
	t.Chdir(tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)
}

func TestResolveRejectsInvalidBackendConfig(t *testing.T) {
	declareContext(t, `backends:
  vllm:
    endpoint: http://localhost:8000
    api_token: nope
`)

	_, err := resolve("")
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "backends.vllm.api_token") {
		t.Fatalf("error = %v, want full backend field path", err)
	}
}
