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
	// FieldMemTotal is the total physical memory field.
	FieldMemTotal = "MemTotal"
	// FieldMemAvailable is the available memory field (kernel 3.14+).
	FieldMemAvailable = "MemAvailable"
	// FieldMemFree is the free memory field (fallback for older kernels).
	FieldMemFree = "MemFree"
	// FieldBuffers is the kernel buffer memory field.
	FieldBuffers = "Buffers"
	// FieldCached is the page cache memory field.
	FieldCached = "Cached"
	// FieldSwapTotal is the total swap space field.
	FieldSwapTotal = "SwapTotal"
	// FieldSwapFree is the free swap space field.
	FieldSwapFree = "SwapFree"
)
