//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides shared helpers for config subcommands.
package core

import (
	"fmt"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/err/config"
	execGit "github.com/ActiveMemory/ctx/internal/exec/git"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
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
		return config.ReadProfile(srcFile, readErr)
	}

	dst := filepath.Join(root, file.CtxRC)
	return io.SafeWriteFile(dst, data, fs.PermFile)
}

// SwitchTo copies the requested profile to .ctxrc and returns a status message.
//
// If the requested profile is already active, returns a no-op message.
// If .ctxrc did not previously exist, returns a "created" message.
//
// Parameters:
//   - root: Git repository root directory
//   - profile: Target profile name (file.ProfileDev or file.ProfileBase)
//
// Returns:
//   - string: Status message for the user
//   - error: Non-nil if the profile file copy fails
func SwitchTo(root, profile string) (string, error) {
	current := DetectProfile()
	if current == profile {
		return fmt.Sprintf(
			desc.Text(text.DescKeyConfigAlreadyOn), profile), nil
	}

	srcFile := file.CtxRCBase
	if profile == file.ProfileDev {
		srcFile = file.CtxRCDev
	}

	if copyErr := CopyProfile(root, srcFile); copyErr != nil {
		return "", copyErr
	}

	if current == "" {
		return fmt.Sprintf(
			desc.Text(text.DescKeyConfigCreated), file.CtxRC, profile), nil
	}
	return fmt.Sprintf(
		desc.Text(text.DescKeyConfigSwitched), profile), nil
}

// GitRoot returns the git repository root directory.
//
// Returns an error if git is not installed or the current directory is
// not inside a git repository. Features that depend on git should
// degrade gracefully when this returns an error.
//
// Returns:
//   - string: Absolute path to the repository root
//   - error: If git is not found or CWD is not in a repo
func GitRoot() (string, error) {
	return execGit.Root()
}
