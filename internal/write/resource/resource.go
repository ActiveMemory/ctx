//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

import (
	"encoding/json"

	"github.com/spf13/cobra"

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
