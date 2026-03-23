//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

// SeverityFor returns the severity level for a given resource name
// from an alert list. Returns SeverityOK if no alert matches.
//
// Parameters:
//   - alerts: list of resource alerts to search
//   - resource: resource name to match (e.g., "memory", "disk")
//
// Returns:
//   - Severity: the severity level for the resource
func SeverityFor(alerts []ResourceAlert, resource string) Severity {
	for _, a := range alerts {
		if a.Resource == resource {
			return a.Severity
		}
	}
	return SeverityOK
}

// Collect gathers a resource snapshot.
//
// The path argument determines which filesystem is checked for disk usage
// (typically the working directory).
//
// Parameters:
//   - path: Filesystem path for disk usage check
//
// Returns:
//   - Snapshot: Memory, disk, and load metrics
func Collect(path string) Snapshot {
	return Snapshot{
		Memory: collectMemory(),
		Disk:   collectDisk(path),
		Load:   collectLoad(),
	}
}

// MaxSeverity returns the highest severity among the given alerts.
//
// Returns SeverityOK when the slice is empty.
//
// Parameters:
//   - alerts: Resource alerts to evaluate
//
// Returns:
//   - Severity: Highest severity found, or SeverityOK if empty
func MaxSeverity(alerts []ResourceAlert) Severity {
	highest := SeverityOK
	for _, a := range alerts {
		if a.Severity > highest {
			highest = a.Severity
		}
	}
	return highest
}
