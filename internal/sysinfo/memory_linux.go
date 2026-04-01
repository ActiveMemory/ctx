//go:build linux

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	cfgSysinfo "github.com/ActiveMemory/ctx/internal/config/sysinfo"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// collectMemory reads physical and swap memory usage
// from /proc/meminfo on Linux.
//
// Returns a MemInfo with Supported=false if /proc/meminfo cannot be opened.
//
// Returns:
//   - MemInfo: Physical and swap memory statistics
func collectMemory() MemInfo {
	f, openErr := os.Open(cfgSysinfo.ProcMeminfo)
	if openErr != nil {
		return MemInfo{Supported: false}
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			ctxLog.Warn(
				warn.Close, cfgSysinfo.ProcMeminfo, closeErr,
			)
		}
	}()
	return parseMeminfo(f)
}

// parseMeminfo parses /proc/meminfo content into a MemInfo struct.
//
// Reads key-value pairs in "Key: value kB" format. Used memory is
// computed as Total - Available (with a fallback to Free + Buffers +
// Cached for kernels before 3.14 that lack MemAvailable).
//
// Parameters:
//   - r: Reader providing /proc/meminfo content
//
// Returns:
//   - MemInfo: Parsed memory and swap statistics
func parseMeminfo(r io.Reader) MemInfo {
	vals := make(map[string]uint64)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), token.Colon, 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSuffix(strings.TrimSpace(parts[1]), cfgSysinfo.MemInfoSuffix)
		n, parseErr := strconv.ParseUint(
			strings.TrimSpace(val), 10, 64,
		)
		if parseErr == nil {
			vals[key] = n * cfgSysinfo.BytesPerKB
		}
	}

	total := vals[cfgSysinfo.FieldMemTotal]
	available := vals[cfgSysinfo.FieldMemAvailable]
	if available == 0 {
		// Fallback for kernels without MemAvailable (< 3.14)
		available = vals[cfgSysinfo.FieldMemFree] +
			vals[cfgSysinfo.FieldBuffers] + vals[cfgSysinfo.FieldCached]
	}

	var used uint64
	if total > available {
		used = total - available
	}

	swapTotal := vals[cfgSysinfo.FieldSwapTotal]
	swapFree := vals[cfgSysinfo.FieldSwapFree]
	var swapUsed uint64
	if swapTotal > swapFree {
		swapUsed = swapTotal - swapFree
	}

	return MemInfo{
		TotalBytes:     total,
		UsedBytes:      used,
		SwapTotalBytes: swapTotal,
		SwapUsedBytes:  swapUsed,
		Supported:      true,
	}
}
