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
	"strconv"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/stats"
	cfgSysinfo "github.com/ActiveMemory/ctx/internal/config/sysinfo"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// collectMemory reads physical and swap memory usage
// from /proc/meminfo on Linux, then overlays the OS-native
// memory pressure signal from /proc/pressure/memory.
//
// Returns a MemInfo with Supported=false if /proc/meminfo cannot be opened.
//
// Returns:
//   - MemInfo: Physical and swap memory statistics
func collectMemory() MemInfo {
	f, openErr := ctxIo.SafeOpenUserFile(cfgSysinfo.ProcMeminfo)
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
	info := parseMeminfo(f)
	info.Pressure, info.PressureSupported = collectPressure()
	return info
}

// collectPressure reads the memory Pressure Stall Information
// (PSI) signal from /proc/pressure/memory on Linux.
//
// Returns SeverityOK with supported=false when the PSI file is
// absent or unreadable (PSI disabled or kernel before 4.20), so
// callers raise no alert.
//
// Returns:
//   - Severity: Mapped pressure severity (SeverityOK when unsupported)
//   - bool: Whether the PSI signal is available
func collectPressure() (Severity, bool) {
	f, openErr := ctxIo.SafeOpenUserFile(cfgSysinfo.ProcPressureMemory)
	if openErr != nil {
		return SeverityOK, false
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			ctxLog.Warn(
				warn.Close, cfgSysinfo.ProcPressureMemory, closeErr,
			)
		}
	}()
	return parsePressure(f)
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
		val := strings.TrimSuffix(
			strings.TrimSpace(parts[1]),
			cfgSysinfo.MemInfoSuffix,
		)
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

// parsePressure parses /proc/pressure/memory content into a
// memory pressure severity.
//
// Reads the 10-second pressure averages from the "some" and
// "full" lines. "some" measures the share of the window in
// which at least one task stalled on memory; "full" measures
// the share in which every runnable task stalled. Danger is
// raised when full.avg10 meets ThresholdMemPressureFullDangerPct
// and takes precedence; otherwise Warning is raised when
// some.avg10 meets ThresholdMemPressureSomeWarnPct.
//
// Returns supported=false when no PSI line could be parsed (the
// file was empty or malformed), so callers raise no alert.
//
// Parameters:
//   - r: Reader providing /proc/pressure/memory content
//
// Returns:
//   - Severity: Mapped pressure severity (SeverityOK when below thresholds)
//   - bool: Whether a PSI line was successfully parsed
func parsePressure(r io.Reader) (Severity, bool) {
	var someAvg10, fullAvg10 float64
	parsed := false

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		avg10, ok := pressureAvg10(fields)
		if !ok {
			continue
		}
		switch fields[0] {
		case cfgSysinfo.PSILineSome:
			someAvg10 = avg10
			parsed = true
		case cfgSysinfo.PSILineFull:
			fullAvg10 = avg10
			parsed = true
		}
	}

	if !parsed {
		return SeverityOK, false
	}
	if fullAvg10 >= stats.ThresholdMemPressureFullDangerPct {
		return SeverityDanger, true
	}
	if someAvg10 >= stats.ThresholdMemPressureSomeWarnPct {
		return SeverityWarning, true
	}
	return SeverityOK, true
}

// pressureAvg10 extracts the avg10 value from the fields of a
// single PSI line (e.g. "some avg10=0.00 avg60=0.00 ...").
//
// Parameters:
//   - fields: Whitespace-split tokens of one PSI line
//
// Returns:
//   - float64: Parsed avg10 value (0 when absent)
//   - bool: Whether an avg10 field was found and parsed
func pressureAvg10(fields []string) (float64, bool) {
	for _, field := range fields {
		raw, found := strings.CutPrefix(field, cfgSysinfo.PSIFieldAvg10)
		if !found {
			continue
		}
		avg10, parseErr := strconv.ParseFloat(raw, 64)
		if parseErr != nil {
			return 0, false
		}
		return avg10, true
	}
	return 0, false
}
