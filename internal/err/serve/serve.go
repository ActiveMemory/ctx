//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package serve

import "fmt"

// NoRunningHub wraps a PID file read failure.
//
// Parameters:
//   - cause: the underlying read error
//
// Returns:
//   - error: "no running hub: <cause>"
func NoRunningHub(cause error) error {
	return fmt.Errorf("no running hub: %w", cause)
}

// InvalidPID wraps a PID file parse failure.
//
// Parameters:
//   - cause: the underlying parse error
//
// Returns:
//   - error: "invalid PID file: <cause>"
func InvalidPID(cause error) error {
	return fmt.Errorf("invalid PID file: %w", cause)
}

// Kill wraps a process kill failure.
//
// Parameters:
//   - pid: process ID that failed to stop
//   - cause: the underlying kill error
//
// Returns:
//   - error: "kill <pid>: <cause>"
func Kill(pid int, cause error) error {
	return fmt.Errorf("kill %d: %w", pid, cause)
}
