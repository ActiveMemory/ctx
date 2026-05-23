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

// Declare positions the test at tempDir as if the user had `cd`'d
// there. Under the cwd-anchored resolution model
// (spec: specs/cwd-anchored-context.md), this is what `ctx` needs
// to resolve `<tempDir>/.context` from any subsequent call.
//
// Side effects:
//
//   - [t.Chdir](tempDir): subsequent ctx calls resolve via $PWD.
//   - [t.Setenv]("HOME", tempDir): redirects user-home writes
//     (e.g. ~/.claude/settings.json) into the temp tree so
//     parallel-package `go test ./...` runs do not race on the
//     real config file.
//   - [rc.Reset]: clears any cached rc state from prior tests in
//     the process; registered as a [t.Cleanup] for symmetry on
//     test exit.
//
// Declare does NOT create the directory; that is the caller's
// responsibility, typically via `ctx init`. Tests that only need
// the cwd positioned (without materializing .context/) can skip
// the init step.
//
// Note: [t.Chdir] is incompatible with [t.Parallel] (it changes
// process-global state). Tests using Declare must not call
// t.Parallel.
//
// Parameters:
//   - t:       test handle (required for t.Chdir / t.Setenv / t.Cleanup).
//   - tempDir: absolute path to the per-test temp directory, usually
//     the value returned by t.TempDir().
//
// Returns:
//   - string: absolute path `<tempDir>/.context` (whether or not
//     it has been materialized).
func Declare(t *testing.T, tempDir string) string {
	t.Helper()
	t.Chdir(tempDir)
	t.Setenv(env.Home, tempDir)
	rc.Reset()
	t.Cleanup(rc.Reset)
	return filepath.Join(tempDir, dir.Context)
}
