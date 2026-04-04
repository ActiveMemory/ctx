//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	cfgSysinfo "github.com/ActiveMemory/ctx/internal/config/sysinfo"
)

// Evaluate checks a snapshot against resource thresholds and returns any
// alerts. Unsupported or zero-total resources are silently skipped.
//
// Thresholds:
//   - Memory: WARNING >= 80%, DANGER >= 90%
//   - Swap:   WARNING >= 50%, DANGER >= 75%
//   - Disk:   WARNING >= 85%, DANGER >= 95%
//   - Load:   WARNING >= 0.8x CPUs, DANGER >= 1.5x CPUs
//
// Parameters:
//   - snap: System resource snapshot to evaluate
//
// Returns:
//   - []ResourceAlert: Alerts for any resources exceeding thresholds
func Evaluate(snap Snapshot) []ResourceAlert {
	var alerts []ResourceAlert

	// Memory
	if snap.Memory.Supported && snap.Memory.TotalBytes > 0 {
		pct := percent(snap.Memory.UsedBytes, snap.Memory.TotalBytes)
		msg := fmt.Sprintf(desc.Text(text.DescKeyResourcesAlertMemory),
			pct, FormatGiB(snap.Memory.UsedBytes), FormatGiB(snap.Memory.TotalBytes))
		if pct >= stats.ThresholdMemoryDangerPct {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityDanger, Resource: cfgSysinfo.ResourceMemory, Message: msg,
			})
		} else if pct >= stats.ThresholdMemoryWarnPct {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityWarning,
				Resource: cfgSysinfo.ResourceMemory,
				Message:  msg,
			})
		}
	}

	// Swap
	if snap.Memory.Supported && snap.Memory.SwapTotalBytes > 0 {
		pct := percent(snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes)
		msg := fmt.Sprintf(
			desc.Text(text.DescKeyResourcesAlertSwap),
			pct,
			FormatGiB(snap.Memory.SwapUsedBytes),
			FormatGiB(snap.Memory.SwapTotalBytes),
		)
		if pct >= stats.ThresholdSwapDangerPct {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityDanger, Resource: cfgSysinfo.ResourceSwap, Message: msg,
			})
		} else if pct >= stats.ThresholdSwapWarnPct {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityWarning, Resource: cfgSysinfo.ResourceSwap, Message: msg,
			})
		}
	}

	// Disk
	if snap.Disk.Supported && snap.Disk.TotalBytes > 0 {
		pct := percent(snap.Disk.UsedBytes, snap.Disk.TotalBytes)
		msg := fmt.Sprintf(desc.Text(text.DescKeyResourcesAlertDisk),
			pct, FormatGiB(snap.Disk.UsedBytes), FormatGiB(snap.Disk.TotalBytes))
		if pct >= stats.ThresholdDiskDangerPct {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityDanger, Resource: cfgSysinfo.ResourceDisk, Message: msg,
			})
		} else if pct >= stats.ThresholdDiskWarnPct {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityWarning, Resource: cfgSysinfo.ResourceDisk, Message: msg,
			})
		}
	}

	// Load (1m)
	if snap.Load.Supported && snap.Load.NumCPU > 0 {
		ratio := snap.Load.Load1 / float64(snap.Load.NumCPU)
		msg := fmt.Sprintf(desc.Text(text.DescKeyResourcesAlertLoad), ratio)
		if ratio >= stats.ThresholdLoadDangerRatio {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityDanger, Resource: cfgSysinfo.ResourceLoad, Message: msg,
			})
		} else if ratio >= stats.ThresholdLoadWarnRatio {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityWarning, Resource: cfgSysinfo.ResourceLoad, Message: msg,
			})
		}
	}

	return alerts
}

// FormatGiB formats bytes as a GiB value with one decimal place (e.g. "14.7").
//
// Parameters:
//   - bytes: Value in bytes to format
//
// Returns:
//   - string: Formatted GiB string (e.g. "14.7")
func FormatGiB(bytes uint64) string {
	gib := float64(bytes) / stats.ThresholdBytesPerGiB
	return fmt.Sprintf(stats.FormatGiB, gib)
}
