//go:build darwin

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import "os/exec"

// Sysctl runs sysctl with the given arguments and returns stdout.
//
// Parameters:
//   - args: sysctl flags and key names (e.g. "-n", "hw.memsize")
//
// Returns:
//   - []byte: raw stdout output
//   - error: non-nil if the command fails
func Sysctl(args ...string) ([]byte, error) {
	//nolint:gosec // fixed command, no user input
	return exec.Command("sysctl", args...).Output()
}

// VMStat runs vm_stat and returns stdout.
//
// Returns:
//   - []byte: raw stdout output
//   - error: non-nil if the command fails
func VMStat() ([]byte, error) {
	return exec.Command("vm_stat").Output()
}
