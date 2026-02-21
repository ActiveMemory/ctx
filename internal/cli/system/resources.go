//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/sysinfo"
	"github.com/spf13/cobra"
)

// resourcesCmd returns the "ctx system resources" subcommand.
func resourcesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resources",
		Short: "Show system resource usage (memory, swap, disk, load)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResources(cmd)
		},
	}
	cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}

// statusCol is the column where the status indicator starts.
const statusCol = 52

func runResources(cmd *cobra.Command) error {
	snap := sysinfo.Collect(".")
	alerts := sysinfo.Evaluate(snap)

	jsonFlag, _ := cmd.Flags().GetBool("json")
	if jsonFlag {
		return outputResourcesJSON(cmd, snap, alerts)
	}
	outputResourcesText(cmd, snap, alerts)
	return nil
}

func outputResourcesText(cmd *cobra.Command, snap sysinfo.Snapshot, alerts []sysinfo.ResourceAlert) {
	cmd.Println("System Resources")
	cmd.Println("====================")
	cmd.Println()

	// Memory line
	if snap.Memory.Supported {
		pct := pctOf(snap.Memory.UsedBytes, snap.Memory.TotalBytes)
		values := fmt.Sprintf("%5s / %5s GB (%d%%)",
			sysinfo.FormatGiB(snap.Memory.UsedBytes),
			sysinfo.FormatGiB(snap.Memory.TotalBytes),
			pct)
		sev := severityFor(alerts, "memory")
		cmd.Println(formatLine("Memory:", values, statusText(sev)))
	}

	// Swap line
	if snap.Memory.Supported {
		pct := pctOf(snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes)
		values := fmt.Sprintf("%5s / %5s GB (%d%%)",
			sysinfo.FormatGiB(snap.Memory.SwapUsedBytes),
			sysinfo.FormatGiB(snap.Memory.SwapTotalBytes),
			pct)
		sev := severityFor(alerts, "swap")
		cmd.Println(formatLine("Swap:", values, statusText(sev)))
	}

	// Disk line
	if snap.Disk.Supported {
		pct := pctOf(snap.Disk.UsedBytes, snap.Disk.TotalBytes)
		values := fmt.Sprintf("%5s / %5s GB (%d%%)",
			sysinfo.FormatGiB(snap.Disk.UsedBytes),
			sysinfo.FormatGiB(snap.Disk.TotalBytes),
			pct)
		sev := severityFor(alerts, "disk")
		cmd.Println(formatLine("Disk:", values, statusText(sev)))
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
		sev := severityFor(alerts, "load")
		cmd.Println(formatLine("Load:", values, statusText(sev)))
	}

	// Summary
	cmd.Println()
	if len(alerts) == 0 {
		cmd.Println("All clear \u2014 no resource warnings.")
	} else {
		cmd.Println("Alerts:")
		for _, a := range alerts {
			icon := "\u26a0"
			if a.Severity == sysinfo.SeverityDanger {
				icon = "\u2716"
			}
			cmd.Println(fmt.Sprintf("  %s %s", icon, a.Message))
		}
	}
}

func outputResourcesJSON(cmd *cobra.Command, snap sysinfo.Snapshot, alerts []sysinfo.ResourceAlert) error {
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

func formatLine(label, values, status string) string {
	left := fmt.Sprintf("%-7s  %s", label, values)
	pad := statusCol - len(left)
	if pad < 1 {
		pad = 1
	}
	return left + strings.Repeat(" ", pad) + status
}

func statusText(sev sysinfo.Severity) string {
	switch sev {
	case sysinfo.SeverityWarning:
		return "\u26a0 WARNING"
	case sysinfo.SeverityDanger:
		return "\u2716 DANGER"
	default:
		return "\u2713 ok"
	}
}

func severityFor(alerts []sysinfo.ResourceAlert, resource string) sysinfo.Severity {
	for _, a := range alerts {
		if a.Resource == resource {
			return a.Severity
		}
	}
	return sysinfo.SeverityOK
}

func pctOf(used, total uint64) int {
	if total == 0 {
		return 0
	}
	return int(float64(used) / float64(total) * 100)
}
