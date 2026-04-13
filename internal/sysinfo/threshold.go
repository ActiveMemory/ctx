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

	type byteCheck struct {
		supported bool
		used      uint64
		total     uint64
		descKey   string
		resource  string
		dangerPct float64
		warnPct   float64
	}

	checks := []byteCheck{
		{
			snap.Memory.Supported,
			snap.Memory.UsedBytes,
			snap.Memory.TotalBytes,
			text.DescKeyResourcesAlertMemory,
			cfgSysinfo.ResourceMemory,
			stats.ThresholdMemoryDangerPct,
			stats.ThresholdMemoryWarnPct,
		},
		{
			snap.Memory.Supported,
			snap.Memory.SwapUsedBytes,
			snap.Memory.SwapTotalBytes,
			text.DescKeyResourcesAlertSwap,
			cfgSysinfo.ResourceSwap,
			stats.ThresholdSwapDangerPct,
			stats.ThresholdSwapWarnPct,
		},
		{
			snap.Disk.Supported,
			snap.Disk.UsedBytes,
			snap.Disk.TotalBytes,
			text.DescKeyResourcesAlertDisk,
			cfgSysinfo.ResourceDisk,
			stats.ThresholdDiskDangerPct,
			stats.ThresholdDiskWarnPct,
		},
	}

	for _, c := range checks {
		if !c.supported || c.total == 0 {
			continue
		}
		pct := percent(c.used, c.total)
		msg := fmt.Sprintf(
			desc.Text(c.descKey), pct,
			FormatGiB(c.used), FormatGiB(c.total),
		)
		if pct >= c.dangerPct {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityDanger,
				Resource: c.resource,
				Message:  msg,
			})
		} else if pct >= c.warnPct {
			alerts = append(alerts, ResourceAlert{
				Severity: SeverityWarning,
				Resource: c.resource,
				Message:  msg,
			})
		}
	}

	// Load (5m) — 5-minute average smooths transient build/test spikes.
	if snap.Load.Supported && snap.Load.NumCPU > 0 {
		ratio := snap.Load.Load5 / float64(snap.Load.NumCPU)
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
