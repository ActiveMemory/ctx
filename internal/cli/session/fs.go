//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// sessionsDirPath returns the path to the `sessions` directory.
//
// Returns:
//   - string: Full path to .context/sessions/
func sessionsDirPath() string {
	return filepath.Join(rc.ContextDir(), config.DirSessions)
}
