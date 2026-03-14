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

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Profile file names and identifiers — aliased from internal/config.
const (
	FileCtxRC     = file.CtxRC
	FileCtxRCBase = file.CtxRCBase
	FileCtxRCDev  = file.CtxRCDev
	ProfileDev    = file.ProfileDev
	ProfileBase   = file.ProfileBase
	ProfileProd   = file.ProfileProd
)

// DetectProfile returns the active profile name from the parsed .ctxrc.
// Returns "" if .ctxrc is missing or has no profile field.
//
// Returns:
//   - string: Profile name ("dev", "base", or "")
func DetectProfile() string {
	return rc.RC().Profile
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
	data, readErr := io.SafeReadFile(root, srcFile)
	if readErr != nil {
		return ctxerr.ReadProfile(srcFile, readErr)
	}

	dst := filepath.Join(root, FileCtxRC)
	return os.WriteFile(dst, data, fs.PermFile)
}

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
	current := DetectProfile()
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
