//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

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

// Text formats and prints system resource information as a human-readable
// table with status indicators and alert summaries.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - snap: collected system resource snapshot
//   - alerts: evaluated resource alerts
func Text(cmd *cobra.Command, snap sysinfo.Snapshot, alerts []sysinfo.ResourceAlert) {
	if cmd == nil {
		return
	}
	for _, line := range formatText(snap, alerts) {
		cmd.Println(line)
	}
}

// JSON writes system resource information as formatted JSON to the
// command's output writer.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - snap: collected system resource snapshot
//   - alerts: evaluated resource alerts
//
// Returns:
//   - error: Non-nil on JSON encoding failure
func JSON(cmd *cobra.Command, snap sysinfo.Snapshot, alerts []sysinfo.ResourceAlert) error {
	if cmd == nil {
		return nil
	}
	out := jsonOutput{}

	out.Memory.TotalBytes = snap.Memory.TotalBytes
	out.Memory.UsedBytes = snap.Memory.UsedBytes
	out.Memory.Percent = pctOf(snap.Memory.UsedBytes, snap.Memory.TotalBytes)
	out.Memory.Supported = snap.Memory.Supported

	out.Swap.TotalBytes = snap.Memory.SwapTotalBytes
	out.Swap.UsedBytes = snap.Memory.SwapUsedBytes
	out.Swap.Percent = pctOf(snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes)
	out.Swap.Supported = snap.Memory.Supported

	out.Disk.TotalBytes = snap.Disk.TotalBytes
	out.Disk.UsedBytes = snap.Disk.UsedBytes
	out.Disk.Percent = pctOf(snap.Disk.UsedBytes, snap.Disk.TotalBytes)
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

func pctOf(used, total uint64) int {
	if total == 0 {
		return 0
	}
	return int(float64(used) / float64(total) * 100)
}

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

func formatLine(label, values, status string) string {
	left := fmt.Sprintf(
		fmt.Sprintf("%%-%ds  %%s", stats.ResourcesLabelWidth), label, values)
	pad := stats.ResourcesStatusCol - len(left)
	if pad < 1 {
		pad = 1
	}
	return left + strings.Repeat(" ", pad) + status
}

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
