//go:build darwin

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"strconv"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
	execSysinfo "github.com/ActiveMemory/ctx/internal/exec/sysinfo"
)

// sysctl key and flag constants for macOS memory queries.
const (
	// flagNoNewline suppresses the key name in sysctl output.
	flagNoNewline = "-n"
	// keyHWMemsize is the sysctl key for total physical memory.
	keyHWMemsize = "hw.memsize"
	// keyVMSwapUsage is the sysctl key for swap usage.
	keyVMSwapUsage = "vm.swapusage"
)

// vm_stat output parsing constants.
const (
	// markerPageSize is the sentinel substring in vm_stat
	// output that precedes the page size value.
	markerPageSize = "page size of"
	// labelPagesFree is the vm_stat line label for free pages.
	labelPagesFree = "Pages free"
	// labelPagesInactive is the vm_stat line label for
	// inactive pages.
	labelPagesInactive = "Pages inactive"
)

// Swap usage parsing constants.
const (
	// suffixMB is the megabyte suffix in sysctl swap output.
	suffixMB = "M"
	// labelTotal is the swap usage field name for total swap.
	labelTotal = "total"
	// labelUsed is the swap usage field name for used swap.
	labelUsed = "used"
)

// defaultPageSize is the default memory page size on Apple
// Silicon (bytes).
const defaultPageSize = 16384

// bytesPerKB is the number of bytes in a kilobyte.
const bytesPerKB = 1024

// collectMemory queries physical and swap memory usage on macOS.
//
// Uses `sysctl -n hw.memsize` for total RAM, `vm_stat` for page-level
// usage, and `sysctl -n vm.swapusage` for swap statistics. Returns a
// MemInfo with Supported=false if the total memory cannot be determined.
//
// Returns:
//   - MemInfo: Physical and swap memory statistics
func collectMemory() MemInfo {
	// Total physical memory
	out, memErr := execSysinfo.Sysctl(flagNoNewline, keyHWMemsize)
	if memErr != nil {
		return MemInfo{Supported: false}
	}
	totalBytes, parseErr := strconv.ParseUint(
		strings.TrimSpace(string(out)), 10, 64,
	)
	if parseErr != nil {
		return MemInfo{Supported: false}
	}

	// Memory page stats via vm_stat
	var usedBytes uint64
	out, vmStatErr := execSysinfo.VMStat()
	if vmStatErr == nil {
		usedBytes = parseVMStat(string(out), totalBytes)
	}

	// Swap via sysctl
	var swapTotal, swapUsed uint64
	out, swapErr := execSysinfo.Sysctl(flagNoNewline, keyVMSwapUsage)
	if swapErr == nil {
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
//
// Computes used bytes as Total - (free + inactive) * pageSize.
// Defaults to 16384-byte pages (Apple Silicon) if page size is not
// found in the output.
//
// Parameters:
//   - output: Raw output from the vm_stat command
//   - totalBytes: Total physical memory in bytes
//
// Returns:
//   - uint64: Estimated used memory in bytes
func parseVMStat(output string, totalBytes uint64) uint64 {
	var pageSize uint64 = defaultPageSize
	pages := make(map[string]uint64)

	for _, line := range strings.Split(output, token.NewlineLF) {
		if strings.Contains(line, markerPageSize) {
			for _, word := range strings.Fields(line) {
				n, parseErr := strconv.ParseUint(word, 10, 64)
				if parseErr == nil && n > 0 {
					pageSize = n
					break
				}
			}
			continue
		}
		parts := strings.SplitN(line, token.Colon, 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		raw := strings.TrimSpace(parts[1])
		val := strings.TrimSpace(
			strings.TrimSuffix(raw, token.Dot),
		)
		if n, parseErr := strconv.ParseUint(val, 10, 64); parseErr == nil {
			pages[key] = n
		}
	}

	freeBytes := (pages[labelPagesFree] + pages[labelPagesInactive]) * pageSize
	if freeBytes >= totalBytes {
		return 0
	}
	return totalBytes - freeBytes
}

// parseSwapUsage parses sysctl vm.swapusage output.
//
// Expected format:
//
//	"total = 2048.00M  used = 123.45M  free = 1924.55M"
//
// Values are parsed as megabytes and converted to bytes.
//
// Parameters:
//   - output: Raw output from `sysctl -n vm.swapusage`
//
// Returns:
//   - total: Total swap space in bytes
//   - used: Used swap space in bytes
func parseSwapUsage(output string) (total, used uint64) {
	parseMB := func(s string) uint64 {
		s = strings.TrimSuffix(strings.TrimSpace(s), suffixMB)
		f, parseErr := strconv.ParseFloat(s, 64)
		if parseErr != nil {
			return 0
		}
		return uint64(f * bytesPerKB * bytesPerKB)
	}

	fields := strings.Fields(output)
	for i, f := range fields {
		if f == token.KeyValueSep && i > 0 && i+1 < len(fields) {
			switch fields[i-1] {
			case labelTotal:
				total = parseMB(fields[i+1])
			case labelUsed:
				used = parseMB(fields[i+1])
			}
		}
	}
	return total, used
}
