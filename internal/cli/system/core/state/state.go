//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package state

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	ctxContext "github.com/ActiveMemory/ctx/internal/context/validate"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// StateDir returns the project-scoped runtime state directory
// (.context/state/). Ensures the directory exists on each call — MkdirAll
// is a no-op when the directory is already present.
//
// Returns:
//   - string: Absolute path to the state directory
func StateDir() string {
	d := filepath.Join(rc.ContextDir(), dir.State)
	_ = os.MkdirAll(d, 0o750)
	return d
}

// Initialized reports whether the context directory has been properly set up
// via "ctx init". Hooks should no-op when this returns false to avoid
// creating partial state (e.g. logs/) before initialization.
//
// Returns:
//   - bool: True if context directory is initialized
func Initialized() bool {
	return ctxContext.Initialized(rc.ContextDir())
}
