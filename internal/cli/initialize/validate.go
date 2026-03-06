//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
)

// checkCtxInPath verifies that ctx is available in PATH.
//
// The hooks use "ctx" expecting it to be in PATH, so init should fail
// if the user hasn't installed ctx globally yet.
// Set CTX_SKIP_PATH_CHECK=1 to skip this check (used in tests).
//
// Parameters:
//   - cmd: Cobra command for error output stream
//
// Returns:
//   - error: non-nil if ctx is not found in PATH
func checkCtxInPath(cmd *cobra.Command) error {
	if os.Getenv(config.EnvSkipPathCheck) == config.EnvTrue {
		return nil
	}

	_, err := exec.LookPath("ctx")
	if err != nil {
		write.ErrCtxNotInPath(cmd)
		return ctxerr.CtxNotInPath()
	}
	return nil
}
