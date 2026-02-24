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
type MemInfo struct {
	TotalBytes     uint64
	UsedBytes      uint64
	SwapTotalBytes uint64
	SwapUsedBytes  uint64
	Supported      bool
}

// DiskInfo holds filesystem usage for a given path.
type DiskInfo struct {
	TotalBytes uint64
	UsedBytes  uint64
	Path       string
	Supported  bool
}

// LoadInfo holds system load averages and CPU count.
type LoadInfo struct {
	Load1     float64
	Load5     float64
	Load15    float64
	NumCPU    int
	Supported bool
}

// Snapshot captures a point-in-time view of system resources.
type Snapshot struct {
	Memory MemInfo
	Disk   DiskInfo
	Load   LoadInfo
}

// ResourceAlert describes a single threshold breach.
type ResourceAlert struct {
	Severity Severity
	Resource string // "memory", "swap", "disk", "load"
	Message  string
}

// Collect gathers a resource snapshot. The path argument determines which
// filesystem is checked for disk usage (typically the working directory).
func Collect(path string) Snapshot {
	return Snapshot{
		Memory: collectMemory(),
		Disk:   collectDisk(path),
		Load:   collectLoad(),
	}
}

// MaxSeverity returns the highest severity among the given alerts.
// Returns SeverityOK when the slice is empty.
func MaxSeverity(alerts []ResourceAlert) Severity {
	max := SeverityOK
	for _, a := range alerts {
		if a.Severity > max {
			max = a.Severity
		}
	}
	return max
}
