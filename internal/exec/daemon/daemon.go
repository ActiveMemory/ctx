//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package daemon

import (
	"os/exec"
	"syscall"
)

// Start launches a detached background process with the
// given binary and arguments.
//
// Parameters:
//   - binPath: absolute path to the executable
//   - args: command-line arguments
//
// Returns:
//   - int: PID of the started process
//   - error: non-nil if the process fails to start
func Start(binPath string, args []string) (int, error) {
	proc := exec.Command( //nolint:gosec // caller-controlled
		binPath, args...,
	)
	proc.Stdout = nil
	proc.Stderr = nil
	proc.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if startErr := proc.Start(); startErr != nil {
		return 0, startErr
	}

	return proc.Process.Pid, nil
}
