//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"
)

// MigrateKeyFile renames the legacy .scratchpad.key to .context.key if needed.
func MigrateKeyFile(contextDir string) {
	old := filepath.Join(contextDir, ".scratchpad.key")
	nw := filepath.Join(contextDir, FileContextKey)
	if _, err := os.Stat(nw); err == nil {
		return // already migrated
	}
	if _, err := os.Stat(old); err == nil {
		_ = os.Rename(old, nw)
	}
}
