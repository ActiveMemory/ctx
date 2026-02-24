//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

func TestOutputResourcesText_AllClear(t *testing.T) {
	cmd := newTestCmd()
	snap := sysinfo.Snapshot{
		Memory: sysinfo.MemInfo{
			TotalBytes:     16 * giB,
			UsedBytes:      4 * giB,
			SwapTotalBytes: 8 * giB,
			SwapUsedBytes:  0,
			Supported:      true,
		},
		Disk: sysinfo.DiskInfo{
			TotalBytes: 500 * giB,
			UsedBytes:  180 * giB,
			Path:       "/",
			Supported:  true,
		},
		Load: sysinfo.LoadInfo{
			Load1:     0.52,
			Load5:     0.41,
			Load15:    0.38,
			NumCPU:    8,
			Supported: true,
		},
	}
	alerts := sysinfo.Evaluate(snap)
	outputResourcesText(cmd, snap, alerts)
	out := cmdOutput(cmd)

	if !strings.Contains(out, "System Resources") {
		t.Error("missing header")
	}
	if !strings.Contains(out, "Memory:") {
		t.Error("missing Memory line")
	}
	if !strings.Contains(out, "Swap:") {
		t.Error("missing Swap line")
	}
	if !strings.Contains(out, "Disk:") {
		t.Error("missing Disk line")
	}
	if !strings.Contains(out, "Load:") {
		t.Error("missing Load line")
	}
	if !strings.Contains(out, "All clear") {
		t.Errorf("expected 'All clear' message, got:\n%s", out)
	}
	// Verify no alerts shown
	if strings.Contains(out, "Alerts:") {
		t.Error("unexpected Alerts section for all-clear snapshot")
	}
}

func TestOutputResourcesText_WithAlerts(t *testing.T) {
	cmd := newTestCmd()
	snap := sysinfo.Snapshot{
		Memory: sysinfo.MemInfo{
			TotalBytes:     16 * giB,
			UsedBytes:      15 * giB, // ~94%
			SwapTotalBytes: 8 * giB,
			SwapUsedBytes:  7 * giB, // ~88%
			Supported:      true,
		},
		Disk: sysinfo.DiskInfo{
			TotalBytes: 500 * giB,
			UsedBytes:  180 * giB,
			Path:       "/",
			Supported:  true,
		},
		Load: sysinfo.LoadInfo{
			Load1:     12.50,
			Load5:     9.30,
			Load15:    6.10,
			NumCPU:    8,
			Supported: true,
		},
	}
	alerts := sysinfo.Evaluate(snap)
	outputResourcesText(cmd, snap, alerts)
	out := cmdOutput(cmd)

	if !strings.Contains(out, "DANGER") {
		t.Error("expected DANGER indicator")
	}
	if !strings.Contains(out, "Alerts:") {
		t.Error("expected Alerts section")
	}
	if strings.Contains(out, "All clear") {
		t.Error("unexpected 'All clear' with active alerts")
	}
}

func TestOutputResourcesJSON_ValidJSON(t *testing.T) {
	cmd := newTestCmd()
	snap := sysinfo.Snapshot{
		Memory: sysinfo.MemInfo{
			TotalBytes:     16 * giB,
			UsedBytes:      4 * giB,
			SwapTotalBytes: 8 * giB,
			Supported:      true,
		},
		Disk: sysinfo.DiskInfo{
			TotalBytes: 500 * giB,
			UsedBytes:  180 * giB,
			Supported:  true,
		},
		Load: sysinfo.LoadInfo{
			Load1:     0.5,
			NumCPU:    8,
			Supported: true,
		},
	}
	alerts := sysinfo.Evaluate(snap)
	if err := outputResourcesJSON(cmd, snap, alerts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := cmdOutput(cmd)
	if !strings.Contains(out, `"max_severity"`) {
		t.Error("missing max_severity field")
	}
	if !strings.Contains(out, `"memory"`) {
		t.Error("missing memory field")
	}
}

func TestOutputResourcesText_UnsupportedResourcesHidden(t *testing.T) {
	cmd := newTestCmd()
	snap := sysinfo.Snapshot{
		Memory: sysinfo.MemInfo{Supported: false},
		Disk:   sysinfo.DiskInfo{Supported: false},
		Load:   sysinfo.LoadInfo{Supported: false},
	}
	outputResourcesText(cmd, snap, nil)
	out := cmdOutput(cmd)

	if strings.Contains(out, "Memory:") {
		t.Error("unsupported Memory should be hidden")
	}
	if strings.Contains(out, "Disk:") {
		t.Error("unsupported Disk should be hidden")
	}
	if strings.Contains(out, "Load:") {
		t.Error("unsupported Load should be hidden")
	}
}

const giB = 1 << 30
