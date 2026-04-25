//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package testctx

import (
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Declare wires CTX_DIR to <tempDir>/.context, redirects HOME to
// tempDir so user-home writes (e.g. ~/.claude/settings.json) stay
// inside the temp tree, resets rc state, and returns the absolute
// path that CTX_DIR now points to.
//
// HOME isolation matters because `ctx init` reads and writes
// ~/.claude/settings.json. Without isolation, parallel `go test
// ./...` packages all read-modify-write the same real file and race.
//
// Typical pattern:
//
//	tmpDir := t.TempDir()
//	t.Chdir(tmpDir)
//	ctxPath := testctx.Declare(t, tmpDir)
//	_ = initialize.Cmd().Execute()   // materialize .context/
//	// subsequent ctx commands in the same process resolve to ctxPath
//
// Declare does NOT create the directory; that is the caller's
// responsibility, typically via `ctx init`. Tests that only need the
// environment declared (without materializing .context/) can skip the
// init step.
//
// Parameters:
//   - t:       test handle (required for t.Setenv / t.Cleanup).
//   - tempDir: absolute path to the per-test temp directory, usually
//     the value returned by t.TempDir().
//
// Returns:
//   - string: absolute path `<tempDir>/.context`.
func Declare(t *testing.T, tempDir string) string {
	t.Helper()
	ctxDir := filepath.Join(tempDir, dir.Context)
	t.Setenv(env.CtxDir, ctxDir)
	t.Setenv(env.Home, tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)
	return ctxDir
}
