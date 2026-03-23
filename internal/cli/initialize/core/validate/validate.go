//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/err/initialize"
	writeInit "github.com/ActiveMemory/ctx/internal/write/initialize"
)

// CheckCtxInPath verifies that the ctx binary is in PATH.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if ctx is not found in PATH
func CheckCtxInPath(cmd *cobra.Command) error {
	if os.Getenv(env.SkipPathCheck) == env.True {
		return nil
	}
	_, err := exec.LookPath(cli.Binary)
	if err != nil {
		writeInit.ErrCtxNotInPath(cmd)
		return initialize.CtxNotInPath()
	}
	return nil
}
