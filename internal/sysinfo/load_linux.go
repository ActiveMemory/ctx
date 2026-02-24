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

func collectLoad() LoadInfo {
	f, err := os.Open("/proc/loadavg")
	if err != nil {
		return LoadInfo{Supported: false}
	}
	defer func() { _ = f.Close() }()
	return parseLoadavg(f)
}

// parseLoadavg parses /proc/loadavg content into a LoadInfo struct.
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
