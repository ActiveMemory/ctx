//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"os/exec"

	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/err/initialize"
	initialize2 "github.com/ActiveMemory/ctx/internal/write/initialize"
	"github.com/spf13/cobra"
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
	_, err := exec.LookPath("ctx")
	if err != nil {
		initialize2.ErrCtxNotInPath(cmd)
		return initialize.CtxNotInPath()
	}
	return nil
}
