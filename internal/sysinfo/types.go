//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

// Severity represents the urgency level of a resource alert.
type Severity int

const (
	SeverityOK      Severity = iota // No concern
	SeverityWarning                 // Approaching limits
	SeverityDanger                  // Critically low resources
)

// String returns the lowercase label for the severity level.
func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "warning"
	case SeverityDanger:
		return "danger"
	default:
		return "ok"
	}
}

// MemInfo holds memory and swap usage metrics.
//
// Fields:
//   - TotalBytes: Total physical memory
//   - UsedBytes: Used physical memory
//   - SwapTotalBytes: Total swap space
//   - SwapUsedBytes: Used swap space
//   - Supported: Whether memory info is available on this platform
type MemInfo struct {
	TotalBytes     uint64
	UsedBytes      uint64
	SwapTotalBytes uint64
	SwapUsedBytes  uint64
	Supported      bool
}

// DiskInfo holds filesystem usage for a given path.
//
// Fields:
//   - TotalBytes: Total filesystem capacity
//   - UsedBytes: Used filesystem space
//   - Path: Filesystem mount path
//   - Supported: Whether disk info is available on this platform
//   - Err: Collection error (nil on success)
type DiskInfo struct {
	TotalBytes uint64
	UsedBytes  uint64
	Path       string
	Supported  bool
	Err        error
}

// LoadInfo holds system load averages and CPU count.
//
// Fields:
//   - Load1: 1-minute load average
//   - Load5: 5-minute load average
//   - Load15: 15-minute load average
//   - NumCPU: Number of logical CPUs
//   - Supported: Whether load info is available on this platform
type LoadInfo struct {
	Load1     float64
	Load5     float64
	Load15    float64
	NumCPU    int
	Supported bool
}

// Snapshot captures a point-in-time view of system resources.
//
// Fields:
//   - Memory: Memory and swap metrics
//   - Disk: Filesystem usage for the project root
//   - Load: System load averages
type Snapshot struct {
	Memory MemInfo
	Disk   DiskInfo
	Load   LoadInfo
}

// Resource name constants for threshold evaluation.
const (
	ResourceMemory = "memory"
	ResourceSwap   = "swap"
	ResourceDisk   = "disk"
	ResourceLoad   = "load"
)

// ResourceAlert describes a single threshold breach.
//
// Fields:
//   - Severity: Alert urgency (OK, Warning, Danger)
//   - Resource: Which resource breached (memory, swap, disk, load)
//   - Message: Human-readable description
type ResourceAlert struct {
	Severity Severity
	Resource string
	Message  string
}
