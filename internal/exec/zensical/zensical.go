//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package zensical

import (
	"os"
	"os/exec"

	cfgZensical "github.com/ActiveMemory/ctx/internal/config/zensical"
	errSite "github.com/ActiveMemory/ctx/internal/err/site"
)

// Run launches zensical with the given subcommand in the given directory.
//
// Parameters:
//   - dir: Working directory for the zensical process
//   - command: Zensical subcommand (e.g. "build", "serve")
//
// Returns:
//   - error: Non-nil if zensical is not found or the process fails
func Run(dir, command string) error {
	if _, lookErr := exec.LookPath(cfgZensical.Bin); lookErr != nil {
		return errSite.ZensicalNotFound()
	}

	//nolint:gosec // G204: binary is a constant, command is from caller
	z := exec.Command(cfgZensical.Bin, command)
	z.Dir = dir
	z.Stdout = os.Stdout
	z.Stderr = os.Stderr
	z.Stdin = os.Stdin

	return z.Run()
}
