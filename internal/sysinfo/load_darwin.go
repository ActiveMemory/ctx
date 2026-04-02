//go:build darwin

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"fmt"
	"runtime"
	"strings"

	execSysinfo "github.com/ActiveMemory/ctx/internal/exec/sysinfo"
)

// collectLoad queries system load averages on macOS via sysctl.
//
// Parses the output of `sysctl -n vm.loadavg` (format: "{ 0.52 0.41 0.38 }")
// into 1-, 5-, and 15-minute load averages. Returns a LoadInfo with
// Supported=false if the command fails or output cannot be parsed.
//
// Returns:
//   - LoadInfo: System load averages and CPU count
func collectLoad() LoadInfo {
	out, cmdErr := execSysinfo.Sysctl("-n", "vm.loadavg")
	if cmdErr != nil {
		return LoadInfo{Supported: false}
	}
	// Output: "{ 0.52 0.41 0.38 }"
	s := strings.Trim(strings.TrimSpace(string(out)), "{ }")
	var load1, load5, load15 float64
	_, scanErr := fmt.Sscanf(
		s, "%f %f %f", &load1, &load5, &load15,
	)
	if scanErr != nil {
		return LoadInfo{Supported: false}
	}
	return LoadInfo{
		Load1:     load1,
		Load5:     load5,
		Load15:    load15,
		NumCPU:    runtime.NumCPU(),
		Supported: true,
	}
}
