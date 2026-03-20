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

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/spf13/cobra"

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
		return desc.TextDesc(text.DescKeyResourcesStatusWarn)
	case sysinfo.SeverityDanger:
		return desc.TextDesc(text.DescKeyResourcesStatusDanger)
	default:
		return desc.TextDesc(text.DescKeyResourcesStatusOk)
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
	left := fmt.Sprintf("%-7s  %s", label, values)
	pad := stats.ResourcesStatusCol - len(left)
	if pad < 1 {
		pad = 1
	}
	return left + strings.Repeat(" ", pad) + status
}

// OutputResourcesText prints system resource information in human-readable
// table format with status indicators and alert summaries.
//
// Parameters:
//   - cmd: Cobra command for output
//   - snap: collected system resource snapshot
//   - alerts: evaluated resource alerts
func OutputResourcesText(cmd *cobra.Command, snap sysinfo.Snapshot, alerts []sysinfo.ResourceAlert) {
	cmd.Println(desc.TextDesc(text.DescKeyResourcesHeader))
	cmd.Println(desc.TextDesc(text.DescKeyResourcesSeparator))
	cmd.Println()

	// Memory line
	if snap.Memory.Supported {
		pct := PctOf(snap.Memory.UsedBytes, snap.Memory.TotalBytes)
		values := fmt.Sprintf("%5s / %5s GB (%d%%)",
			sysinfo.FormatGiB(snap.Memory.UsedBytes),
			sysinfo.FormatGiB(snap.Memory.TotalBytes),
			pct)
		sev := SeverityFor(alerts, "memory")
		cmd.Println(FormatResourceLine("Memory:", values, StatusText(sev)))
	}

	// Swap line
	if snap.Memory.Supported {
		pct := PctOf(snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes)
		values := fmt.Sprintf("%5s / %5s GB (%d%%)",
			sysinfo.FormatGiB(snap.Memory.SwapUsedBytes),
			sysinfo.FormatGiB(snap.Memory.SwapTotalBytes),
			pct)
		sev := SeverityFor(alerts, "swap")
		cmd.Println(FormatResourceLine("Swap:", values, StatusText(sev)))
	}

	// Disk line
	if snap.Disk.Supported {
		pct := PctOf(snap.Disk.UsedBytes, snap.Disk.TotalBytes)
		values := fmt.Sprintf("%5s / %5s GB (%d%%)",
			sysinfo.FormatGiB(snap.Disk.UsedBytes),
			sysinfo.FormatGiB(snap.Disk.TotalBytes),
			pct)
		sev := SeverityFor(alerts, "disk")
		cmd.Println(FormatResourceLine("Disk:", values, StatusText(sev)))
	}

	// Load line
	if snap.Load.Supported {
		ratio := 0.0
		if snap.Load.NumCPU > 0 {
			ratio = snap.Load.Load1 / float64(snap.Load.NumCPU)
		}
		values := fmt.Sprintf("%5.2f / %5.2f / %5.2f  (%d CPUs, ratio %.2f)",
			snap.Load.Load1, snap.Load.Load5, snap.Load.Load15,
			snap.Load.NumCPU, ratio)
		sev := SeverityFor(alerts, "load")
		cmd.Println(FormatResourceLine("Load:", values, StatusText(sev)))
	}

	// Summary
	cmd.Println()
	if len(alerts) == 0 {
		cmd.Println(desc.TextDesc(text.DescKeyResourcesAllClear))
	} else {
		cmd.Println(desc.TextDesc(text.DescKeyResourcesAlerts))
		for _, a := range alerts {
			if a.Severity == sysinfo.SeverityDanger {
				cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyResourcesAlertDanger), a.Message))
			} else {
				cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyResourcesAlertWarning), a.Message))
			}
		}
	}
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
	type jsonAlert struct {
		Severity string `json:"severity"`
		Resource string `json:"resource"`
		Message  string `json:"message"`
	}
	type jsonOutput struct {
		Memory struct {
			TotalBytes uint64 `json:"total_bytes"`
			UsedBytes  uint64 `json:"used_bytes"`
			Percent    int    `json:"percent"`
			Supported  bool   `json:"supported"`
		} `json:"memory"`
		Swap struct {
			TotalBytes uint64 `json:"total_bytes"`
			UsedBytes  uint64 `json:"used_bytes"`
			Percent    int    `json:"percent"`
			Supported  bool   `json:"supported"`
		} `json:"swap"`
		Disk struct {
			TotalBytes uint64 `json:"total_bytes"`
			UsedBytes  uint64 `json:"used_bytes"`
			Percent    int    `json:"percent"`
			Path       string `json:"path"`
			Supported  bool   `json:"supported"`
		} `json:"disk"`
		Load struct {
			Load1     float64 `json:"load1"`
			Load5     float64 `json:"load5"`
			Load15    float64 `json:"load15"`
			NumCPU    int     `json:"num_cpu"`
			Ratio     float64 `json:"ratio"`
			Supported bool    `json:"supported"`
		} `json:"load"`
		Alerts      []jsonAlert `json:"alerts"`
		MaxSeverity string      `json:"max_severity"`
	}

	out := jsonOutput{}

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

	out.Alerts = make([]jsonAlert, 0, len(alerts))
	for _, a := range alerts {
		out.Alerts = append(out.Alerts, jsonAlert{
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
