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

// MigrateKeyFile consolidates legacy key files and promotes project-local
// keys to the user-level directory.
//
// Migration tiers (executed in order, stopping at first match):
//  1. Rename legacy names (.context.key, .scratchpad.key) → .ctx.key
//  2. Copy project-local .ctx.key → ~/.local/ctx/keys/<slug>.key,
//     then remove the project-local copy
//
// Parameters:
//   - contextDir: The .context/ directory path
//   - projectRoot: Absolute path to the project root (for user-level path)
func MigrateKeyFile(contextDir, projectRoot string) {
	localKey := filepath.Join(contextDir, FileContextKey)

	// Tier 1: rename legacy file names within project dir.
	if _, err := os.Stat(localKey); err != nil {
		for _, legacy := range []string{".context.key", ".scratchpad.key"} {
			old := filepath.Join(contextDir, legacy)
			if _, err := os.Stat(old); err == nil {
				_ = os.Rename(old, localKey)
				break
			}
		}
	}

	// Tier 2: promote project-local key to user-level directory.
	userLevel := ProjectKeyPath(projectRoot)
	if userLevel == "" {
		return // home dir unavailable
	}

	// Already at user level — nothing to do.
	if _, err := os.Stat(userLevel); err == nil {
		// Clean up stale project-local copy if both exist.
		if _, localErr := os.Stat(localKey); localErr == nil {
			_ = os.Remove(localKey)
		}
		return
	}

	// No project-local key to promote.
	if _, err := os.Stat(localKey); err != nil {
		return
	}

	// Copy to user-level, then remove project-local.
	data, readErr := os.ReadFile(localKey) //nolint:gosec // project-local path
	if readErr != nil {
		return
	}

	if mkdirErr := os.MkdirAll(filepath.Dir(userLevel), PermKeyDir); mkdirErr != nil {
		return
	}

	if writeErr := os.WriteFile(userLevel, data, PermSecret); writeErr != nil {
		return
	}

	_ = os.Remove(localKey)
}
