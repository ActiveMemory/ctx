//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package crypto

import (
	"os"
	"path/filepath"
	"strings"

	cfgCrypto "github.com/ActiveMemory/ctx/internal/config/crypto"
	"github.com/ActiveMemory/ctx/internal/config/dir"
)

// GlobalKeyPath returns the global encryption key path.
//
// Returns ~/.ctx/.ctx.key using os.UserHomeDir.
// Returns an empty string if the home directory cannot be determined.
//
// Returns:
//   - string: Absolute path to the global encryption key, or empty string on failure
func GlobalKeyPath() string {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return ""
	}
	return filepath.Join(home, dir.CtxData, cfgCrypto.ContextKey)
}

// ExpandHome expands a leading ~/ prefix to the user's home directory.
//
// If the path does not start with "~/", it is returned unchanged.
// If the home directory cannot be determined, the path is returned unchanged.
//
// Parameters:
//   - path: File path that may contain a leading ~/
//
// Returns:
//   - string: Path with ~/ expanded to the home directory
func ExpandHome(path string) string {
	if !strings.HasPrefix(path, "~/") {
		return path
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return path
	}
	return filepath.Join(home, path[2:])
}

// ResolveKeyPath determines the effective key file path.
//
// Resolution order:
//  1. overridePath if non-empty (explicit .ctxrc key_path, with tilde expansion)
//  2. Project-local path if it exists (<contextDir>/.ctx.key)
//  3. Global default (~/.ctx/.ctx.key)
//  4. Project-local path as fallback (when home dir unavailable)
//
// Parameters:
//   - contextDir: The .context/ directory path
//   - overridePath: Explicit key path from .ctxrc (may be empty)
//
// Returns:
//   - string: The resolved key file path
func ResolveKeyPath(contextDir, overridePath string) string {
	// Tier 1: explicit override from .ctxrc key_path.
	if overridePath != "" {
		return ExpandHome(overridePath)
	}

	// Tier 2: project-local key.
	local := filepath.Join(contextDir, cfgCrypto.ContextKey)
	if _, statErr := os.Stat(local); statErr == nil {
		return local
	}

	// Tier 3: global default.
	global := GlobalKeyPath()
	if global != "" {
		return global
	}

	// Fallback: project-local (home dir unavailable).
	return local
}
