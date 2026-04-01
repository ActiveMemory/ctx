//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

const (
	// ProcLoadavg is the Linux procfs path for load averages.
	ProcLoadavg = "/proc/loadavg"
	// ProcMeminfo is the Linux procfs path for memory information.
	ProcMeminfo = "/proc/meminfo"
	// LoadavgFmt is the scanf format for parsing /proc/loadavg fields.
	LoadavgFmt = "%f %f %f"
	// MemInfoSuffix is the unit suffix in /proc/meminfo values.
	MemInfoSuffix = " kB"
	// BytesPerKB converts kilobytes to bytes.
	BytesPerKB = 1024
)

// Meminfo field keys from /proc/meminfo.
const (
	FieldMemTotal     = "MemTotal"
	FieldMemAvailable = "MemAvailable"
	FieldMemFree      = "MemFree"
	FieldBuffers      = "Buffers"
	FieldCached       = "Cached"
	FieldSwapTotal    = "SwapTotal"
	FieldSwapFree     = "SwapFree"
)
