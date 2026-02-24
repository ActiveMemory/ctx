//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import "testing"

const giB = 1 << 30

func TestMaxSeverity_NoAlerts(t *testing.T) {
	if got := MaxSeverity(nil); got != SeverityOK {
		t.Errorf("MaxSeverity(nil) = %v, want ok", got)
	}
}

func TestMaxSeverity_SingleWarning(t *testing.T) {
	alerts := []ResourceAlert{{Severity: SeverityWarning}}
	if got := MaxSeverity(alerts); got != SeverityWarning {
		t.Errorf("MaxSeverity = %v, want warning", got)
	}
}

func TestMaxSeverity_MixedSeverities(t *testing.T) {
	alerts := []ResourceAlert{
		{Severity: SeverityOK},
		{Severity: SeverityDanger},
		{Severity: SeverityWarning},
	}
	if got := MaxSeverity(alerts); got != SeverityDanger {
		t.Errorf("MaxSeverity = %v, want danger", got)
	}
}

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		sev  Severity
		want string
	}{
		{SeverityOK, "ok"},
		{SeverityWarning, "warning"},
		{SeverityDanger, "danger"},
	}
	for _, tt := range tests {
		if got := tt.sev.String(); got != tt.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tt.sev, got, tt.want)
		}
	}
}

func TestCollect_DoesNotPanic(t *testing.T) {
	snap := Collect(".")
	// Should return a valid snapshot on any platform
	_ = snap
}

func TestFormatGiB(t *testing.T) {
	tests := []struct {
		bytes uint64
		want  string
	}{
		{0, "0.0"},
		{1 * giB, "1.0"},
		{16 * giB, "16.0"},
	}
	for _, tt := range tests {
		if got := FormatGiB(tt.bytes); got != tt.want {
			t.Errorf("FormatGiB(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestEvaluate_AllClear(t *testing.T) {
	snap := Snapshot{
		Memory: MemInfo{
			TotalBytes:     16 * giB,
			UsedBytes:      4 * giB,
			SwapTotalBytes: 8 * giB,
			SwapUsedBytes:  0,
			Supported:      true,
		},
		Disk: DiskInfo{
			TotalBytes: 500 * giB,
			UsedBytes:  180 * giB,
			Path:       "/",
			Supported:  true,
		},
		Load: LoadInfo{
			Load1:     0.52,
			Load5:     0.41,
			Load15:    0.38,
			NumCPU:    8,
			Supported: true,
		},
	}
	alerts := Evaluate(snap)
	if len(alerts) != 0 {
		t.Errorf("expected no alerts, got %d: %v", len(alerts), alerts)
	}
}

func TestEvaluate_UnsupportedResourcesSkipped(t *testing.T) {
	snap := Snapshot{
		Memory: MemInfo{Supported: false},
		Disk:   DiskInfo{Supported: false},
		Load:   LoadInfo{Supported: false},
	}
	alerts := Evaluate(snap)
	if len(alerts) != 0 {
		t.Errorf("expected no alerts for unsupported resources, got %d", len(alerts))
	}
}

func TestEvaluate_ZeroTotalSkipped(t *testing.T) {
	snap := Snapshot{
		Memory: MemInfo{TotalBytes: 0, Supported: true},
		Disk:   DiskInfo{TotalBytes: 0, Supported: true},
		Load:   LoadInfo{NumCPU: 0, Supported: true},
	}
	alerts := Evaluate(snap)
	if len(alerts) != 0 {
		t.Errorf("expected no alerts for zero totals, got %d", len(alerts))
	}
}

func TestEvaluate_NoSwapNoAlert(t *testing.T) {
	snap := Snapshot{
		Memory: MemInfo{
			TotalBytes:     16 * giB,
			UsedBytes:      4 * giB,
			SwapTotalBytes: 0,
			Supported:      true,
		},
	}
	alerts := Evaluate(snap)
	for _, a := range alerts {
		if a.Resource == "swap" {
			t.Errorf("unexpected swap alert when swap total is 0")
		}
	}
}
