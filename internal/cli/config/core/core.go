//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides shared helpers for config subcommands.
package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	internalConfig "github.com/ActiveMemory/ctx/internal/config"
)

// Profile file names and identifiers.
const (
	FileCtxRC     = ".ctxrc"
	FileCtxRCBase = ".ctxrc.base"
	FileCtxRCDev  = ".ctxrc.dev"

	ProfileDev  = "dev"
	ProfileBase = "base"
)

// DetectProfile reads .ctxrc and returns "dev" or "base" based on the
// presence of an uncommented "notify:" line. Returns "" if the file is missing.
//
// Parameters:
//   - root: Git repository root directory
//
// Returns:
//   - string: Profile name ("dev", "base", or "" if missing)
func DetectProfile(root string) string {
	data, readErr := os.ReadFile(filepath.Join(root, FileCtxRC)) //nolint:gosec // project-local config file
	if readErr != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), internalConfig.NewlineLF) {
		if strings.HasPrefix(strings.TrimSpace(line), "notify:") {
			return ProfileDev
		}
	}
	return ProfileBase
}

// CopyProfile copies a source profile file to .ctxrc.
//
// Parameters:
//   - root: Git repository root directory
//   - srcFile: Source profile filename (e.g., ".ctxrc.dev")
//
// Returns:
//   - error: Non-nil on read or write failure
func CopyProfile(root, srcFile string) error {
	src := filepath.Join(root, srcFile)
	data, readErr := os.ReadFile(src) //nolint:gosec // project-local file
	if readErr != nil {
		return fmt.Errorf("read %s: %w", srcFile, readErr)
	}

	dst := filepath.Join(root, FileCtxRC)
	return os.WriteFile(dst, data, internalConfig.PermFile)
}

// GitRoot returns the git repository root directory.
//
// Returns:
//   - string: Absolute path to the git root
//   - error: Non-nil when not inside a git repository
func GitRoot() (string, error) {
	out, execErr := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if execErr != nil {
		return "", fmt.Errorf("not in a git repository: %w", execErr)
	}
	return strings.TrimSpace(string(out)), nil
}
