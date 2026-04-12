//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sysinfo provides constants for system information
// collection commands and keys.
package sysinfo

// macOS system command names.
const (
	// CmdSysctl is the sysctl command name.
	CmdSysctl = "sysctl"
	// CmdVMStat is the vm_stat command name.
	CmdVMStat = "vm_stat"
)

// macOS sysctl keys and output patterns.
const (
	// KeyLoadAvg is the sysctl key for load averages.
	KeyLoadAvg = "vm.loadavg"
	// KeyHWMemsize is the sysctl key for total physical memory.
	KeyHWMemsize = "hw.memsize"
	// KeyVMSwapUsage is the sysctl key for swap usage.
	KeyVMSwapUsage = "vm.swapusage"
	// TrimBraces is the brace wrapper trimmed from sysctl
	// vm.loadavg output (e.g. "{ 0.52 0.41 0.38 }").
	TrimBraces = "{ }"
	// FlagNoNewline suppresses the key name in sysctl output.
	FlagNoNewline = "-n"
	// FmtLoadAvg is the Sscanf format for parsing three
	// load average floats.
	FmtLoadAvg = "%f %f %f"
)

// vm_stat output parsing constants.
const (
	// MarkerPageSize is the sentinel substring in vm_stat
	// output that precedes the page size value.
	MarkerPageSize = "page size of"
	// LabelPagesFree is the vm_stat line label for free pages.
	LabelPagesFree = "Pages free"
	// LabelPagesInactive is the vm_stat line label for
	// inactive pages.
	LabelPagesInactive = "Pages inactive"
)

// Swap usage parsing constants.
const (
	// SuffixMB is the megabyte suffix in sysctl swap output.
	SuffixMB = "M"
	// LabelTotal is the swap usage field name for total swap.
	LabelTotal = "total"
	// LabelUsed is the swap usage field name for used swap.
	LabelUsed = "used"
)
