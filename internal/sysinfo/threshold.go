//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import "fmt"

// Evaluate checks a snapshot against resource thresholds and returns any
// alerts. Unsupported or zero-total resources are silently skipped.
func Evaluate(snap Snapshot) []ResourceAlert {
	var alerts []ResourceAlert

	// Memory: WARNING >= 80%, DANGER >= 90%
	if snap.Memory.Supported && snap.Memory.TotalBytes > 0 {
		pct := percent(snap.Memory.UsedBytes, snap.Memory.TotalBytes)
		if pct >= 90 {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityDanger,
				Resource: "memory",
				Message:  fmt.Sprintf("Memory %.0f%% used (%s / %s GB)", pct, FormatGiB(snap.Memory.UsedBytes), FormatGiB(snap.Memory.TotalBytes)),
			})
		} else if pct >= 80 {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityWarning,
				Resource: "memory",
				Message:  fmt.Sprintf("Memory %.0f%% used (%s / %s GB)", pct, FormatGiB(snap.Memory.UsedBytes), FormatGiB(snap.Memory.TotalBytes)),
			})
		}
	}

	// Swap: WARNING >= 50%, DANGER >= 75%
	if snap.Memory.Supported && snap.Memory.SwapTotalBytes > 0 {
		pct := percent(snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes)
		if pct >= 75 {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityDanger,
				Resource: "swap",
				Message:  fmt.Sprintf("Swap %.0f%% used (%s / %s GB)", pct, FormatGiB(snap.Memory.SwapUsedBytes), FormatGiB(snap.Memory.SwapTotalBytes)),
			})
		} else if pct >= 50 {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityWarning,
				Resource: "swap",
				Message:  fmt.Sprintf("Swap %.0f%% used (%s / %s GB)", pct, FormatGiB(snap.Memory.SwapUsedBytes), FormatGiB(snap.Memory.SwapTotalBytes)),
			})
		}
	}

	// Disk: WARNING >= 85%, DANGER >= 95%
	if snap.Disk.Supported && snap.Disk.TotalBytes > 0 {
		pct := percent(snap.Disk.UsedBytes, snap.Disk.TotalBytes)
		if pct >= 95 {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityDanger,
				Resource: "disk",
				Message:  fmt.Sprintf("Disk %.0f%% used (%s / %s GB)", pct, FormatGiB(snap.Disk.UsedBytes), FormatGiB(snap.Disk.TotalBytes)),
			})
		} else if pct >= 85 {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityWarning,
				Resource: "disk",
				Message:  fmt.Sprintf("Disk %.0f%% used (%s / %s GB)", pct, FormatGiB(snap.Disk.UsedBytes), FormatGiB(snap.Disk.TotalBytes)),
			})
		}
	}

	// Load (1m): WARNING >= 0.8x CPUs, DANGER >= 1.5x CPUs
	if snap.Load.Supported && snap.Load.NumCPU > 0 {
		ratio := snap.Load.Load1 / float64(snap.Load.NumCPU)
		if ratio >= 1.5 {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityDanger,
				Resource: "load",
				Message:  fmt.Sprintf("Load %.2fx CPU count", ratio),
			})
		} else if ratio >= 0.8 {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityWarning,
				Resource: "load",
				Message:  fmt.Sprintf("Load %.2fx CPU count", ratio),
			})
		}
	}

	return alerts
}

// FormatGiB formats bytes as a GiB value with one decimal place (e.g. "14.7").
func FormatGiB(bytes uint64) string {
	gib := float64(bytes) / (1 << 30)
	return fmt.Sprintf("%.1f", gib)
}

func percent(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(used) / float64(total) * 100
}
