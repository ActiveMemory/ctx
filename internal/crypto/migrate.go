//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package crypto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cryptocfg "github.com/ActiveMemory/ctx/internal/config/crypto"
)

// MigrateKeyFile warns about legacy key files that should be moved
// to the global path (~/.ctx/.ctx.key).
//
// If the global key exists, no action is taken. Otherwise, legacy
// locations are checked and a warning is printed to stderr.
//
// Parameters:
//   - contextDir: The .context/ directory path
func MigrateKeyFile(contextDir string) {
	global := GlobalKeyPath()
	if global == "" {
		return
	}

	// Global key exists — nothing to do.
	if _, err := os.Stat(global); err == nil {
		return
	}

	// Check legacy locations and warn.
	var found string

	// Legacy project-local names.
	for _, name := range []string{cryptocfg.ContextKey, ".context.key", ".scratchpad.key"} {
		candidate := filepath.Join(contextDir, name)
		if _, err := os.Stat(candidate); err == nil {
			found = candidate
			break
		}
	}

	// Legacy user-level directory (~/.local/ctx/keys/).
	if found == "" {
		home, homeErr := os.UserHomeDir()
		if homeErr == nil {
			legacyKeyDir := filepath.Join(home, ".local", "ctx", "keys")
			entries, readErr := os.ReadDir(legacyKeyDir)
			if readErr == nil {
				for _, entry := range entries {
					if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".key") {
						found = filepath.Join(legacyKeyDir, entry.Name())
						break
					}
				}
			}
		}
	}

	if found != "" {
		fmt.Fprintf(os.Stderr, "ctx: legacy key found at %s\n"+
			"  Copy it to the new location:\n"+
			"    mkdir -p %s && cp %s %s && chmod 600 %s\n",
			found, filepath.Dir(global), found, global, global)
	}
}
