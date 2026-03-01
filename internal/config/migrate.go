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

// MigrateKeyFile renames legacy key files (.context.key, .scratchpad.key)
// to .ctx.key if needed.
func MigrateKeyFile(contextDir string) {
	nw := filepath.Join(contextDir, FileContextKey)
	if _, err := os.Stat(nw); err == nil {
		return // already migrated
	}
	for _, legacy := range []string{".context.key", ".scratchpad.key"} {
		old := filepath.Join(contextDir, legacy)
		if _, err := os.Stat(old); err == nil {
			_ = os.Rename(old, nw)
			return
		}
	}
}
