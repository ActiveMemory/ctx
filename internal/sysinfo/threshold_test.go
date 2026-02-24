//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import "testing"

func TestEvaluate_MemoryBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		used     uint64
		total    uint64
		wantSev  Severity
		wantN    int
		resource string
	}{
		{"79.9% no alert", 799, 1000, SeverityOK, 0, "memory"},
		{"80% warning", 800, 1000, SeverityWarning, 1, "memory"},
		{"89% warning", 890, 1000, SeverityWarning, 1, "memory"},
		{"90% danger", 900, 1000, SeverityDanger, 1, "memory"},
		{"100% danger", 1000, 1000, SeverityDanger, 1, "memory"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snap := Snapshot{
				Memory: MemInfo{
					TotalBytes: tt.total,
					UsedBytes:  tt.used,
					Supported:  true,
				},
			}
			alerts := Evaluate(snap)
			memAlerts := filterByResource(alerts, tt.resource)
			if len(memAlerts) != tt.wantN {
				t.Fatalf("expected %d alerts, got %d: %v", tt.wantN, len(memAlerts), memAlerts)
			}
			if tt.wantN > 0 && memAlerts[0].Severity != tt.wantSev {
				t.Errorf("severity = %v, want %v", memAlerts[0].Severity, tt.wantSev)
			}
		})
	}
}

func TestEvaluate_SwapBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		used    uint64
		total   uint64
		wantSev Severity
		wantN   int
	}{
		{"49% no alert", 490, 1000, SeverityOK, 0},
		{"50% warning", 500, 1000, SeverityWarning, 1},
		{"74% warning", 740, 1000, SeverityWarning, 1},
		{"75% danger", 750, 1000, SeverityDanger, 1},
		{"100% danger", 1000, 1000, SeverityDanger, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snap := Snapshot{
				Memory: MemInfo{
					TotalBytes:     1000,
					UsedBytes:      100, // below memory threshold
					SwapTotalBytes: tt.total,
					SwapUsedBytes:  tt.used,
					Supported:      true,
				},
			}
			alerts := Evaluate(snap)
			swapAlerts := filterByResource(alerts, "swap")
			if len(swapAlerts) != tt.wantN {
				t.Fatalf("expected %d alerts, got %d: %v", tt.wantN, len(swapAlerts), swapAlerts)
			}
			if tt.wantN > 0 && swapAlerts[0].Severity != tt.wantSev {
				t.Errorf("severity = %v, want %v", swapAlerts[0].Severity, tt.wantSev)
			}
		})
	}
}

func TestEvaluate_DiskBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		used    uint64
		total   uint64
		wantSev Severity
		wantN   int
	}{
		{"84% no alert", 840, 1000, SeverityOK, 0},
		{"85% warning", 850, 1000, SeverityWarning, 1},
		{"94% warning", 940, 1000, SeverityWarning, 1},
		{"95% danger", 950, 1000, SeverityDanger, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snap := Snapshot{
				Disk: DiskInfo{
					TotalBytes: tt.total,
					UsedBytes:  tt.used,
					Path:       "/",
					Supported:  true,
				},
			}
			alerts := Evaluate(snap)
			diskAlerts := filterByResource(alerts, "disk")
			if len(diskAlerts) != tt.wantN {
				t.Fatalf("expected %d alerts, got %d: %v", tt.wantN, len(diskAlerts), diskAlerts)
			}
			if tt.wantN > 0 && diskAlerts[0].Severity != tt.wantSev {
				t.Errorf("severity = %v, want %v", diskAlerts[0].Severity, tt.wantSev)
			}
		})
	}
}

func TestEvaluate_LoadBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		load1   float64
		ncpu    int
		wantSev Severity
		wantN   int
	}{
		{"ratio 0.79 no alert", 6.32, 8, SeverityOK, 0},
		{"ratio 0.80 warning", 6.40, 8, SeverityWarning, 1},
		{"ratio 1.49 warning", 11.92, 8, SeverityWarning, 1},
		{"ratio 1.50 danger", 12.00, 8, SeverityDanger, 1},
		{"ratio 2.00 danger", 16.00, 8, SeverityDanger, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snap := Snapshot{
				Load: LoadInfo{
					Load1:     tt.load1,
					NumCPU:    tt.ncpu,
					Supported: true,
				},
			}
			alerts := Evaluate(snap)
			loadAlerts := filterByResource(alerts, "load")
			if len(loadAlerts) != tt.wantN {
				t.Fatalf("expected %d alerts, got %d: %v", tt.wantN, len(loadAlerts), loadAlerts)
			}
			if tt.wantN > 0 && loadAlerts[0].Severity != tt.wantSev {
				t.Errorf("severity = %v, want %v", loadAlerts[0].Severity, tt.wantSev)
			}
		})
	}
}

func TestEvaluate_AllDanger(t *testing.T) {
	snap := Snapshot{
		Memory: MemInfo{
			TotalBytes:     16 * giB,
			UsedBytes:      15 * giB,
			SwapTotalBytes: 8 * giB,
			SwapUsedBytes:  7 * giB,
			Supported:      true,
		},
		Disk: DiskInfo{
			TotalBytes: 500 * giB,
			UsedBytes:  490 * giB,
			Path:       "/",
			Supported:  true,
		},
		Load: LoadInfo{
			Load1:     12.0,
			NumCPU:    8,
			Supported: true,
		},
	}
	alerts := Evaluate(snap)
	if len(alerts) != 4 {
		t.Fatalf("expected 4 alerts, got %d: %v", len(alerts), alerts)
	}
	if MaxSeverity(alerts) != SeverityDanger {
		t.Errorf("max severity = %v, want danger", MaxSeverity(alerts))
	}
}

func filterByResource(alerts []ResourceAlert, resource string) []ResourceAlert {
	var out []ResourceAlert
	for _, a := range alerts {
		if a.Resource == resource {
			out = append(out, a)
		}
	}
	return out
}
