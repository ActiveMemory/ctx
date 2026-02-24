//go:build darwin

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func collectLoad() LoadInfo {
	out, err := exec.Command("sysctl", "-n", "vm.loadavg").Output()
	if err != nil {
		return LoadInfo{Supported: false}
	}
	// Output: "{ 0.52 0.41 0.38 }"
	s := strings.Trim(strings.TrimSpace(string(out)), "{ }")
	var load1, load5, load15 float64
	if _, err := fmt.Sscanf(s, "%f %f %f", &load1, &load5, &load15); err != nil {
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
