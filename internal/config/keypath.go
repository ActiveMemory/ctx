//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PermKeyDir is the permission for the user-level key directory (owner rwx only).
const PermKeyDir = 0700

// KeyDir returns the user-level directory for encryption keys.
//
// Returns ~/.local/ctx/keys/ using os.UserHomeDir.
// Returns an empty string if the home directory cannot be determined.
func KeyDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "ctx", "keys")
}

// ProjectKeySlug returns the filename (without directory) for a project's key.
//
// Format: <path-slug>--<sha8>.key
// Example: -home-jose-WORKSPACE-ctx--a1b2c3d4.key
//
// Parameters:
//   - projectRoot: Absolute path to the project root directory
//
// Returns:
//   - string: Key filename
func ProjectKeySlug(projectRoot string) string {
	slug := strings.ReplaceAll(projectRoot, string(filepath.Separator), "-")
	hash := sha256.Sum256([]byte(projectRoot))
	return fmt.Sprintf("%s--%x.key", slug, hash[:4])
}

// ProjectKeyPath returns the full path for a project's user-level key.
//
// Returns ~/.local/ctx/keys/<slug>--<sha8>.key.
// Returns an empty string if the home directory cannot be determined.
//
// Parameters:
//   - projectRoot: Absolute path to the project root directory
//
// Returns:
//   - string: Full path to the key file
func ProjectKeyPath(projectRoot string) string {
	dir := KeyDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, ProjectKeySlug(projectRoot))
}

// ResolveKeyPath determines the effective key file path.
//
// Resolution order:
//  1. overridePath if non-empty (explicit .ctxrc key_path)
//  2. User-level path if it exists (~/.local/ctx/keys/<slug>.key)
//  3. Legacy project-local path if it exists (<contextDir>/.ctx.key)
//  4. User-level path as default (for new key generation)
//
// Parameters:
//   - contextDir: The .context/ directory path
//   - projectRoot: Absolute path to the project root
//   - overridePath: Explicit key path from .ctxrc (may be empty)
//
// Returns:
//   - string: The resolved key file path
func ResolveKeyPath(contextDir, projectRoot, overridePath string) string {
	if overridePath != "" {
		return overridePath
	}

	userLevel := ProjectKeyPath(projectRoot)

	// Check user-level path first.
	if userLevel != "" {
		if _, err := os.Stat(userLevel); err == nil {
			return userLevel
		}
	}

	// Check legacy project-local path.
	legacy := filepath.Join(contextDir, FileContextKey)
	if _, err := os.Stat(legacy); err == nil {
		return legacy
	}

	// Default to user-level for new keys.
	if userLevel != "" {
		return userLevel
	}

	// Fallback: project-local (only when home dir unavailable).
	return legacy
}
