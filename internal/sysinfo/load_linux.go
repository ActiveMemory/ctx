//go:build linux

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

// collectLoad reads system load averages from /proc/loadavg on Linux.
//
// Returns a LoadInfo with Supported=false if /proc/loadavg cannot be
// opened or its content cannot be parsed.
//
// Returns:
//   - LoadInfo: System load averages and CPU count
func collectLoad() LoadInfo {
	f, err := os.Open("/proc/loadavg")
	if err != nil {
		return LoadInfo{Supported: false}
	}
	defer func() { _ = f.Close() }()
	return parseLoadavg(f)
}

// parseLoadavg parses /proc/loadavg content into a LoadInfo struct.
//
// Expects space-separated 1-, 5-, and 15-minute load averages as the
// first three fields. Returns Supported=false if parsing fails.
//
// Parameters:
//   - r: Reader providing /proc/loadavg content
//
// Returns:
//   - LoadInfo: Parsed load averages and CPU count
func parseLoadavg(r io.Reader) LoadInfo {
	var load1, load5, load15 float64
	_, err := fmt.Fscanf(r, "%f %f %f", &load1, &load5, &load15)
	if err != nil {
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
