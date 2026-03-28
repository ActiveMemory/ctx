//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

// pctOf calculates the integer percentage of used relative to total.
//
// Parameters:
//   - used: Numerator value
//   - total: Denominator value; returns 0 when zero to avoid division by zero
//
// Returns:
//   - int: Percentage as an integer (0-100)
func pctOf(used, total uint64) int {
	if total == 0 {
		return 0
	}
	return int(float64(used) / float64(total) * 100)
}

// statusText returns a human-readable label for a severity level.
//
// Parameters:
//   - sev: Severity to convert
//
// Returns:
//   - string: Localized status label (ok, warning, or danger)
func statusText(sev sysinfo.Severity) string {
	switch sev {
	case sysinfo.SeverityWarning:
		return desc.Text(text.DescKeyResourcesStatusWarn)
	case sysinfo.SeverityDanger:
		return desc.Text(text.DescKeyResourcesStatusDanger)
	default:
		return desc.Text(text.DescKeyResourcesStatusOk)
	}
}

// formatLine builds a fixed-width table row with label, values, and status columns.
//
// Parameters:
//   - label: Left-aligned resource name
//   - values: Usage figures placed after the label
//   - status: Right-aligned status indicator
//
// Returns:
//   - string: Padded single-line row
func formatLine(label, values, status string) string {
	left := fmt.Sprintf(
		fmt.Sprintf(desc.Text(text.DescKeyResourcesRowFormat), stats.ResourcesLabelWidth), label, values)
	pad := stats.ResourcesStatusCol - len(left)
	if pad < 1 {
		pad = 1
	}
	return left + strings.Repeat(" ", pad) + status
}

// formatText renders the resource snapshot and alerts as a human-readable text table.
//
// Parameters:
//   - snap: System resource snapshot with memory, disk, and inode data
//   - alerts: Threshold alerts to annotate each row
//
// Returns:
//   - []string: Lines of formatted text, including header and separator
func formatText(snap sysinfo.Snapshot, alerts []sysinfo.ResourceAlert) []string {
	var lines []string
	lines = append(lines, desc.Text(text.DescKeyResourcesHeader))
	lines = append(lines, desc.Text(text.DescKeyResourcesSeparator))
	lines = append(lines, "")

	type gibEntry struct {
		supported   bool
		used, total uint64
		resource    string
		labelKey    string
	}
	gibEntries := []gibEntry{
		{snap.Memory.Supported, snap.Memory.UsedBytes, snap.Memory.TotalBytes,
			sysinfo.ResourceMemory, text.DescKeyResourcesLabelMemory},
		{snap.Memory.Supported, snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes,
			sysinfo.ResourceSwap, text.DescKeyResourcesLabelSwap},
		{snap.Disk.Supported, snap.Disk.UsedBytes, snap.Disk.TotalBytes,
			sysinfo.ResourceDisk, text.DescKeyResourcesLabelDisk},
	}
	valueFmt := desc.Text(text.DescKeyResourcesValueFormat)
	for _, e := range gibEntries {
		if !e.supported {
			continue
		}
		pct := pctOf(e.used, e.total)
		values := fmt.Sprintf(valueFmt,
			sysinfo.FormatGiB(e.used), sysinfo.FormatGiB(e.total), pct)
		sev := sysinfo.SeverityFor(alerts, e.resource)
		lines = append(lines, formatLine(desc.Text(e.labelKey), values, statusText(sev)))
	}

	if snap.Load.Supported {
		ratio := 0.0
		if snap.Load.NumCPU > 0 {
			ratio = snap.Load.Load1 / float64(snap.Load.NumCPU)
		}
		values := fmt.Sprintf(desc.Text(text.DescKeyResourcesLoadFormat),
			snap.Load.Load1, snap.Load.Load5, snap.Load.Load15,
			snap.Load.NumCPU, ratio)
		sev := sysinfo.SeverityFor(alerts, sysinfo.ResourceLoad)
		lines = append(lines, formatLine(
			desc.Text(text.DescKeyResourcesLabelLoad), values, statusText(sev)))
	}

	lines = append(lines, "")
	if len(alerts) == 0 {
		lines = append(lines, desc.Text(text.DescKeyResourcesAllClear))
	} else {
		lines = append(lines, desc.Text(text.DescKeyResourcesAlerts))
		for _, a := range alerts {
			if a.Severity == sysinfo.SeverityDanger {
				lines = append(lines,
					fmt.Sprintf(desc.Text(text.DescKeyResourcesAlertDanger), a.Message))
			} else {
				lines = append(lines,
					fmt.Sprintf(desc.Text(text.DescKeyResourcesAlertWarning), a.Message))
			}
		}
	}
	return lines
}
