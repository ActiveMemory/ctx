//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/stats"

	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

// PctOf calculates the percentage of used relative to total.
// Returns 0 when total is zero to avoid division by zero.
//
// Parameters:
//   - used: the consumed amount
//   - total: the total capacity
//
// Returns:
//   - int: percentage (0-100)
func PctOf(used, total uint64) int {
	if total == 0 {
		return 0
	}
	return int(float64(used) / float64(total) * 100)
}

// SeverityFor returns the severity level for a given resource name
// from an alert list. Returns SeverityOK if no alert matches.
//
// Parameters:
//   - alerts: list of resource alerts to search
//   - resource: resource name to match (e.g., "memory", "disk")
//
// Returns:
//   - sysinfo.Severity: the severity level for the resource
func SeverityFor(alerts []sysinfo.ResourceAlert, resource string) sysinfo.Severity {
	for _, a := range alerts {
		if a.Resource == resource {
			return a.Severity
		}
	}
	return sysinfo.SeverityOK
}

// StatusText returns a human-readable status indicator string for a
// severity level, using the embedded text assets.
//
// Parameters:
//   - sev: severity level
//
// Returns:
//   - string: formatted status text (e.g., "✓ ok", "⚠ WARNING", "✖ DANGER")
func StatusText(sev sysinfo.Severity) string {
	switch sev {
	case sysinfo.SeverityWarning:
		return desc.Text(text.DescKeyResourcesStatusWarn)
	case sysinfo.SeverityDanger:
		return desc.Text(text.DescKeyResourcesStatusDanger)
	default:
		return desc.Text(text.DescKeyResourcesStatusOk)
	}
}

// FormatResourceLine builds a single resource output line with
// left-aligned label+values and a right-aligned status indicator
// at the configured column position.
//
// Parameters:
//   - label: resource label (e.g., "Memory:")
//   - values: formatted resource values
//   - status: status indicator text
//
// Returns:
//   - string: formatted line with aligned status
func FormatResourceLine(label, values, status string) string {
	left := fmt.Sprintf(fmt.Sprintf("%%-%ds  %%s", stats.ResourcesLabelWidth), label, values)
	pad := stats.ResourcesStatusCol - len(left)
	if pad < 1 {
		pad = 1
	}
	return left + strings.Repeat(" ", pad) + status
}

func resourceValueFormat() string {
	return desc.Text(text.DescKeyResourcesValueFormat)
}

// FormatResourcesText formats system resource information as human-readable
// lines with status indicators and alert summaries.
//
// Parameters:
//   - snap: collected system resource snapshot
//   - alerts: evaluated resource alerts
//
// Returns:
//   - []string: formatted output lines
func FormatResourcesText(snap sysinfo.Snapshot, alerts []sysinfo.ResourceAlert) []string {
	var lines []string
	lines = append(lines, desc.Text(text.DescKeyResourcesHeader))
	lines = append(lines, desc.Text(text.DescKeyResourcesSeparator))
	lines = append(lines, "")

	// Memory line
	if snap.Memory.Supported {
		pct := PctOf(snap.Memory.UsedBytes, snap.Memory.TotalBytes)
		values := fmt.Sprintf(resourceValueFormat(),
			sysinfo.FormatGiB(snap.Memory.UsedBytes),
			sysinfo.FormatGiB(snap.Memory.TotalBytes),
			pct)
		sev := SeverityFor(alerts, sysinfo.ResourceMemory)
		lines = append(lines, FormatResourceLine(desc.Text(text.DescKeyResourcesLabelMemory), values, StatusText(sev)))
	}

	// Swap line
	if snap.Memory.Supported {
		pct := PctOf(snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes)
		values := fmt.Sprintf(resourceValueFormat(),
			sysinfo.FormatGiB(snap.Memory.SwapUsedBytes),
			sysinfo.FormatGiB(snap.Memory.SwapTotalBytes),
			pct)
		sev := SeverityFor(alerts, sysinfo.ResourceSwap)
		lines = append(lines, FormatResourceLine(desc.Text(text.DescKeyResourcesLabelSwap), values, StatusText(sev)))
	}

	// Disk line
	if snap.Disk.Supported {
		pct := PctOf(snap.Disk.UsedBytes, snap.Disk.TotalBytes)
		values := fmt.Sprintf(resourceValueFormat(),
			sysinfo.FormatGiB(snap.Disk.UsedBytes),
			sysinfo.FormatGiB(snap.Disk.TotalBytes),
			pct)
		sev := SeverityFor(alerts, sysinfo.ResourceDisk)
		lines = append(lines, FormatResourceLine(desc.Text(text.DescKeyResourcesLabelDisk), values, StatusText(sev)))
	}

	// Load line
	if snap.Load.Supported {
		ratio := 0.0
		if snap.Load.NumCPU > 0 {
			ratio = snap.Load.Load1 / float64(snap.Load.NumCPU)
		}
		values := fmt.Sprintf(desc.Text(text.DescKeyResourcesLoadFormat),
			snap.Load.Load1, snap.Load.Load5, snap.Load.Load15,
			snap.Load.NumCPU, ratio)
		sev := SeverityFor(alerts, sysinfo.ResourceLoad)
		lines = append(lines, FormatResourceLine(desc.Text(text.DescKeyResourcesLabelLoad), values, StatusText(sev)))
	}

	// Summary
	lines = append(lines, "")
	if len(alerts) == 0 {
		lines = append(lines, desc.Text(text.DescKeyResourcesAllClear))
	} else {
		lines = append(lines, desc.Text(text.DescKeyResourcesAlerts))
		for _, a := range alerts {
			if a.Severity == sysinfo.SeverityDanger {
				lines = append(lines, fmt.Sprintf(desc.Text(text.DescKeyResourcesAlertDanger), a.Message))
			} else {
				lines = append(lines, fmt.Sprintf(desc.Text(text.DescKeyResourcesAlertWarning), a.Message))
			}
		}
	}
	return lines
}

// OutputResourcesJSON writes system resource information as formatted
// JSON to the command's output writer.
//
// Parameters:
//   - cmd: Cobra command for output
//   - snap: collected system resource snapshot
//   - alerts: evaluated resource alerts
//
// Returns:
//   - error: Non-nil on JSON encoding failure
func OutputResourcesJSON(cmd *cobra.Command, snap sysinfo.Snapshot, alerts []sysinfo.ResourceAlert) error {
	out := ResourceJSONOutput{}

	out.Memory.TotalBytes = snap.Memory.TotalBytes
	out.Memory.UsedBytes = snap.Memory.UsedBytes
	out.Memory.Percent = PctOf(snap.Memory.UsedBytes, snap.Memory.TotalBytes)
	out.Memory.Supported = snap.Memory.Supported

	out.Swap.TotalBytes = snap.Memory.SwapTotalBytes
	out.Swap.UsedBytes = snap.Memory.SwapUsedBytes
	out.Swap.Percent = PctOf(snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes)
	out.Swap.Supported = snap.Memory.Supported

	out.Disk.TotalBytes = snap.Disk.TotalBytes
	out.Disk.UsedBytes = snap.Disk.UsedBytes
	out.Disk.Percent = PctOf(snap.Disk.UsedBytes, snap.Disk.TotalBytes)
	out.Disk.Path = snap.Disk.Path
	out.Disk.Supported = snap.Disk.Supported

	out.Load.Load1 = snap.Load.Load1
	out.Load.Load5 = snap.Load.Load5
	out.Load.Load15 = snap.Load.Load15
	out.Load.NumCPU = snap.Load.NumCPU
	if snap.Load.NumCPU > 0 {
		out.Load.Ratio = snap.Load.Load1 / float64(snap.Load.NumCPU)
	}
	out.Load.Supported = snap.Load.Supported

	out.Alerts = make([]ResourceJSONAlert, 0, len(alerts))
	for _, a := range alerts {
		out.Alerts = append(out.Alerts, ResourceJSONAlert{
			Severity: a.Severity.String(),
			Resource: a.Resource,
			Message:  a.Message,
		})
	}
	out.MaxSeverity = sysinfo.MaxSeverity(alerts).String()

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
