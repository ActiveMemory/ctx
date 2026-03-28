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
)

// collectMemory reads physical and swap memory usage from /proc/meminfo on Linux.
//
// Returns a MemInfo with Supported=false if /proc/meminfo cannot be opened.
//
// Returns:
//   - MemInfo: Physical and swap memory statistics
func collectMemory() MemInfo {
	f, openErr := os.Open("/proc/meminfo")
	if openErr != nil {
		return MemInfo{Supported: false}
	}
	defer func() { _ = f.Close() }()
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
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSuffix(strings.TrimSpace(parts[1]), " kB")
		if n, parseErr := strconv.ParseUint(strings.TrimSpace(val), 10, 64); parseErr == nil {
			vals[key] = n * 1024 // kB → bytes
		}
	}

	total := vals["MemTotal"]
	available := vals["MemAvailable"]
	if available == 0 {
		// Fallback for kernels without MemAvailable (< 3.14)
		available = vals["MemFree"] + vals["Buffers"] + vals["Cached"]
	}

	var used uint64
	if total > available {
		used = total - available
	}

	swapTotal := vals["SwapTotal"]
	swapFree := vals["SwapFree"]
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
