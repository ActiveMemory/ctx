//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package zensical

import (
	"os"
	"os/exec"

	"github.com/ActiveMemory/ctx/internal/config/zensical"
)

// Run launches zensical serve in the given directory.
//
// Parameters:
//   - dir: Working directory for the zensical process
//
// Returns:
//   - error: Non-nil if the process fails
func Run(dir string) error {
	z := exec.Command(zensical.Bin, "serve") //nolint:gosec // G204: args are constants
	z.Dir = dir
	z.Stdout = os.Stdout
	z.Stderr = os.Stderr
	z.Stdin = os.Stdin

	return z.Run()
}
