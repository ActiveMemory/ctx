//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

// Linux procfs path constants. These constants are consumed
// by Linux-specific source files (memory_linux.go,
// load_linux.go) and are not visible on non-Linux builds.
const (
	// ProcLoadavg is the Linux procfs path for load averages.
	ProcLoadavg = "/proc/loadavg"
	// ProcMeminfo is the Linux procfs path for memory information.
	ProcMeminfo = "/proc/meminfo"
	// ProcPressureMemory is the Linux procfs path for the
	// memory Pressure Stall Information (PSI) signal. Absent
	// when PSI is disabled or the kernel predates 4.20.
	ProcPressureMemory = "/proc/pressure/memory"
	// LoadavgFmt is the scanf format for parsing /proc/loadavg fields.
	LoadavgFmt = "%f %f %f"
	// MemInfoSuffix is the unit suffix in /proc/meminfo values.
	MemInfoSuffix = " kB"
	// BytesPerKB converts kilobytes to bytes.
	BytesPerKB = 1024
)

// Meminfo field keys from /proc/meminfo. These constants are
// consumed by Linux-specific source files and are not visible
// on non-Linux builds.
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

// Pressure Stall Information (PSI) parsing tokens from
// /proc/pressure/memory. Each line reads, for example:
//
//	some avg10=0.00 avg60=0.00 avg300=0.00 total=0
//	full avg10=0.00 avg60=0.00 avg300=0.00 total=0
//
// These constants are consumed by Linux-specific source
// files and are not visible on non-Linux builds.
const (
	// PSILineSome is the line prefix for the "some" pressure
	// row: the share of time at least one task stalled.
	PSILineSome = "some"
	// PSILineFull is the line prefix for the "full" pressure
	// row: the share of time every runnable task stalled.
	PSILineFull = "full"
	// PSIFieldAvg10 is the field name (with its "=" delimiter)
	// for the 10-second pressure average within a PSI line.
	PSIFieldAvg10 = "avg10="
)
