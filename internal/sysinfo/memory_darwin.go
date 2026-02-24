//go:build darwin

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"os/exec"
	"strconv"
	"strings"
)

func collectMemory() MemInfo {
	// Total physical memory
	out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		return MemInfo{Supported: false}
	}
	totalBytes, err := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		return MemInfo{Supported: false}
	}

	// Memory page stats via vm_stat
	var usedBytes uint64
	out, err = exec.Command("vm_stat").Output()
	if err == nil {
		usedBytes = parseVMStat(string(out), totalBytes)
	}

	// Swap via sysctl
	var swapTotal, swapUsed uint64
	out, err = exec.Command("sysctl", "-n", "vm.swapusage").Output()
	if err == nil {
		swapTotal, swapUsed = parseSwapUsage(string(out))
	}

	return MemInfo{
		TotalBytes:     totalBytes,
		UsedBytes:      usedBytes,
		SwapTotalBytes: swapTotal,
		SwapUsedBytes:  swapUsed,
		Supported:      true,
	}
}

// parseVMStat extracts used memory from vm_stat output.
// Used = Total - (free + inactive) * pageSize.
func parseVMStat(output string, totalBytes uint64) uint64 {
	var pageSize uint64 = 16384 // default on Apple Silicon
	pages := make(map[string]uint64)

	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "page size of") {
			for _, word := range strings.Fields(line) {
				if n, err := strconv.ParseUint(word, 10, 64); err == nil && n > 0 {
					pageSize = n
					break
				}
			}
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(parts[1]), "."))
		if n, err := strconv.ParseUint(val, 10, 64); err == nil {
			pages[key] = n
		}
	}

	freeBytes := (pages["Pages free"] + pages["Pages inactive"]) * pageSize
	if freeBytes >= totalBytes {
		return 0
	}
	return totalBytes - freeBytes
}

// parseSwapUsage parses sysctl vm.swapusage output.
// Format: "total = 2048.00M  used = 123.45M  free = 1924.55M  (encrypted)"
func parseSwapUsage(output string) (total, used uint64) {
	parseMB := func(s string) uint64 {
		s = strings.TrimSuffix(strings.TrimSpace(s), "M")
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0
		}
		return uint64(f * 1024 * 1024)
	}

	fields := strings.Fields(output)
	for i, f := range fields {
		if f == "=" && i > 0 && i+1 < len(fields) {
			switch fields[i-1] {
			case "total":
				total = parseMB(fields[i+1])
			case "used":
				used = parseMB(fields[i+1])
			}
		}
	}
	return total, used
}
