//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package serve

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// runServe handles the serve command.
//
// Parameters:
//   - args: Optional directory to serve
//
// Returns:
//   - error: Non-nil if directory is invalid, config is missing,
//     or zensical is not found
func runServe(args []string) error {
	var dir string

	if len(args) > 0 {
		dir = args[0]
	} else {
		dir = filepath.Join(rc.ContextDir(), config.DirJournalSite)
	}

	// Verify directory exists
	info, statErr := os.Stat(dir)
	if statErr != nil {
		return errDirNotFound(dir)
	}
	if !info.IsDir() {
		return errNotDir(dir)
	}

	// Check zensical.toml exists
	tomlPath := filepath.Join(dir, config.FileZensicalToml)
	if _, statErr = os.Stat(tomlPath); os.IsNotExist(statErr) {
		return errNoSiteConfig(dir)
	}

	// Check if zensical is available
	_, lookErr := exec.LookPath(config.BinZensical)
	if lookErr != nil {
		return errZensicalNotFound()
	}

	// Run zensical serve
	zensical := exec.Command(config.BinZensical, "serve") //nolint:gosec // G204: args are constants
	zensical.Dir = dir
	zensical.Stdout = os.Stdout
	zensical.Stderr = os.Stderr
	zensical.Stdin = os.Stdin

	return zensical.Run()
}
