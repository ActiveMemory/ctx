//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/zensical"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run handles the serve command.
//
// Parameters:
//   - args: Optional directory to serve
//
// Returns:
//   - error: Non-nil if directory is invalid, config is missing,
//     or zensical is not found
func Run(args []string) error {
	var d string

	if len(args) > 0 {
		d = args[0]
	} else {
		d = filepath.Join(rc.ContextDir(), dir.JournalSite)
	}

	// Verify directory exists
	info, statErr := os.Stat(d)
	if statErr != nil {
		return ctxerr.DirNotFound(d)
	}
	if !info.IsDir() {
		return ctxerr.NotDirectory(d)
	}

	// Check zensical.toml exists
	tomlPath := filepath.Join(d, zensical.Toml)
	if _, statErr = os.Stat(tomlPath); os.IsNotExist(statErr) {
		return ctxerr.NoSiteConfig(d)
	}

	// Check if zensical is available
	_, lookErr := exec.LookPath(zensical.Bin)
	if lookErr != nil {
		return ctxerr.ZensicalNotFound()
	}

	return runZensical(d)
}

// runZensical launches zensical serve in the given directory.
//
// Parameters:
//   - dir: Working directory for the zensical process
//
// Returns:
//   - error: Non-nil if the process fails
func runZensical(dir string) error {
	z := exec.Command(zensical.Bin, "serve") //nolint:gosec // G204: args are constants
	z.Dir = dir
	z.Stdout = os.Stdout
	z.Stderr = os.Stderr
	z.Stdin = os.Stdin

	return z.Run()
}
