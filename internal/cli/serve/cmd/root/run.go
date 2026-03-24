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
	"github.com/ActiveMemory/ctx/internal/err/fs"
	errSite "github.com/ActiveMemory/ctx/internal/err/site"
	zensicalBin "github.com/ActiveMemory/ctx/internal/exec/zensical"
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
		return fs.DirNotFound(d)
	}
	if !info.IsDir() {
		return fs.NotDirectory(d)
	}

	// Check zensical.toml exists
	tomlPath := filepath.Join(d, zensical.Toml)
	if _, statErr = os.Stat(tomlPath); os.IsNotExist(statErr) {
		return errSite.NoConfig(d)
	}

	// Check if zensical is available
	_, lookErr := exec.LookPath(zensical.Bin)
	if lookErr != nil {
		return errSite.ZensicalNotFound()
	}

	return zensicalBin.Run(d)
}
