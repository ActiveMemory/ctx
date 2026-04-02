//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package state

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/validate"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Dir returns the project-scoped runtime state directory
// (.context/state/). Ensures the directory exists on each call. MkdirAll
// is a no-op when the directory is already present.
//
// Returns:
//   - string: Absolute path to the state directory
func Dir() string {
	if dirOverride != "" {
		return dirOverride
	}
	d := filepath.Join(rc.ContextDir(), dir.State)
	if mkdirErr := ctxIo.SafeMkdirAll(d, fs.PermRestrictedDir); mkdirErr != nil {
		ctxLog.Warn(warn.Mkdir, d, mkdirErr)
	}
	return d
}

// dirOverride allows tests to redirect Dir() to a temp directory.
var dirOverride string

// SetDirForTest overrides Dir() for testing. Pass an empty string
// to restore the default behavior. Only call from tests.
//
// Parameters:
//   - d: Directory path to use, or empty string to restore default
func SetDirForTest(d string) {
	dirOverride = d
}

// Initialized reports whether the context directory has been properly set up
// via "ctx init". Hooks should no-op when this returns false to avoid
// creating a partial state (e.g., logs/) before initialization.
//
// Returns:
//   - bool: True if the context directory is initialized
func Initialized() bool {
	return ctxContext.Initialized(rc.ContextDir())
}
