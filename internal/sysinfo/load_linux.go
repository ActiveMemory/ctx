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

	cfgSysinfo "github.com/ActiveMemory/ctx/internal/config/sysinfo"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// collectLoad reads system load averages from /proc/loadavg on Linux.
//
// Returns a LoadInfo with Supported=false if /proc/loadavg cannot be
// opened or its content cannot be parsed.
//
// Returns:
//   - LoadInfo: System load averages and CPU count
func collectLoad() LoadInfo {
	f, openErr := os.Open(cfgSysinfo.ProcLoadavg)
	if openErr != nil {
		return LoadInfo{Supported: false}
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			ctxLog.Warn(
				warn.Close, cfgSysinfo.ProcLoadavg, closeErr,
			)
		}
	}()
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
	_, scanErr := fmt.Fscanf(r, cfgSysinfo.LoadavgFmt, &load1, &load5, &load15)
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
