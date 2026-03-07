//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides shared helpers for config subcommands.
package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	internalConfig "github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// Profile file names and identifiers — aliased from internal/config.
const (
	FileCtxRC     = internalConfig.FileCtxRC
	FileCtxRCBase = internalConfig.FileCtxRCBase
	FileCtxRCDev  = internalConfig.FileCtxRCDev
	ProfileDev    = internalConfig.ProfileDev
	ProfileBase   = internalConfig.ProfileBase
	ProfileProd   = internalConfig.ProfileProd
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
	data, readErr := validation.SafeReadFile(root, FileCtxRC)
	if readErr != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), internalConfig.NewlineLF) {
		if strings.HasPrefix(strings.TrimSpace(line), internalConfig.ProfileDetectKey) {
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
	data, readErr := validation.SafeReadFile(root, srcFile)
	if readErr != nil {
		return ctxerr.ReadProfile(srcFile, readErr)
	}

	dst := filepath.Join(root, FileCtxRC)
	return os.WriteFile(dst, data, internalConfig.PermFile)
}

// GitRoot returns the git repository root directory.
//
// Returns an error if git is not installed or the current directory is
// not inside a git repository. Features that depend on git should
// degrade gracefully when this returns an error.
//
// Returns:
//   - string: Absolute path to the git root
//   - error: Non-nil when git is missing or not inside a repository
//
// SwitchTo copies the requested profile to .ctxrc and returns a status message.
//
// If the requested profile is already active, returns a no-op message.
// If .ctxrc did not previously exist, returns a "created" message.
//
// Parameters:
//   - root: Git repository root directory
//   - profile: Target profile name (ProfileDev or ProfileBase)
//
// Returns:
//   - string: Status message for the user
//   - error: Non-nil if the profile file copy fails
func SwitchTo(root, profile string) (string, error) {
	current := DetectProfile(root)
	if current == profile {
		return "already on " + profile + " profile", nil
	}

	srcFile := FileCtxRCBase
	if profile == ProfileDev {
		srcFile = FileCtxRCDev
	}

	if copyErr := CopyProfile(root, srcFile); copyErr != nil {
		return "", copyErr
	}

	if current == "" {
		return "created " + FileCtxRC + " from " + profile + " profile", nil
	}
	return "switched to " + profile + " profile", nil
}

// GitRoot returns the git repository root directory.
//
// Returns an error if git is not installed or the current directory is
// not inside a git repository. Features that depend on git should
// degrade gracefully when this returns an error.
func GitRoot() (string, error) {
	if _, lookErr := exec.LookPath("git"); lookErr != nil {
		return "", ctxerr.GitNotFound()
	}

	out, execErr := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if execErr != nil {
		return "", ctxerr.NotInGitRepo(execErr)
	}
	return strings.TrimSpace(string(out)), nil
}
