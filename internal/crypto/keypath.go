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
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// GlobalKeyPath returns the global encryption key path.
//
// Returns ~/.ctx/.ctx.key using os.UserHomeDir.
// Returns an empty string if the home directory cannot be determined.
//
// Returns:
//   - string: Absolute path to the global encryption key,
//     or empty string on failure
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
	if !strings.HasPrefix(path, token.PrefixHomeDir) {
		return path
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return path
	}
	return filepath.Join(home, path[len(token.PrefixHomeDir):])
}

// ResolveKeyPath determines the effective key file path.
//
// Resolution order:
//  1. overridePath if non-empty (explicit .ctxrc key_path, with
//     tilde expansion) — the supported per-project isolation knob
//  2. Global default (~/.ctx/.ctx.key)
//  3. Project-local path (<contextDir>/.ctx.key) as a degenerate
//     fallback ONLY when the home directory is unavailable
//
// A project-local <contextDir>/.ctx.key is never auto-detected or
// preferred over the global key. That implicit tier was removed: it
// stored the key next to the ciphertext (a security antipattern) and
// was the sole cause of key divergence in git worktrees, where the
// gitignored key is absent from the checkout. Per
// specs/notify-resolution-hardening.md.
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

	// Tier 2: global default.
	if global := GlobalKeyPath(); global != "" {
		return global
	}

	// Degenerate fallback: the project-local path, used only when the
	// home directory is unavailable so no global location can be
	// computed. This is NOT project-local key auto-detection — a stray
	// <contextDir>/.ctx.key is never preferred over the global key.
	return filepath.Join(contextDir, cfgCrypto.ContextKey)
}
