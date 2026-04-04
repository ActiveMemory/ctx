//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

// Severity label strings for display output.
const (
	// LabelOK is the severity label for no concern.
	LabelOK = "ok"
	// LabelWarning is the severity label for
	// approaching limits.
	LabelWarning = "warning"
	// LabelDanger is the severity label for critically
	// low resources.
	LabelDanger = "danger"
)

// Resource name constants for threshold evaluation.
const (
	// ResourceMemory is the resource name for physical
	// memory.
	ResourceMemory = "memory"
	// ResourceSwap is the resource name for swap space.
	ResourceSwap = "swap"
	// ResourceDisk is the resource name for filesystem
	// usage.
	ResourceDisk = "disk"
	// ResourceLoad is the resource name for system load.
	ResourceLoad = "load"
)
