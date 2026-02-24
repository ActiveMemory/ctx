//go:build linux

//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"strings"
	"testing"
)

func TestParseMeminfo_Typical(t *testing.T) {
	input := `MemTotal:       16384000 kB
MemFree:         1234567 kB
MemAvailable:    2345678 kB
Buffers:          123456 kB
Cached:          1234567 kB
SwapTotal:       8192000 kB
SwapFree:        2048000 kB
`
	info := parseMeminfo(strings.NewReader(input))
	if !info.Supported {
		t.Fatal("expected Supported = true")
	}
	wantTotal := uint64(16384000) * 1024
	if info.TotalBytes != wantTotal {
		t.Errorf("TotalBytes = %d, want %d", info.TotalBytes, wantTotal)
	}
	wantUsed := wantTotal - uint64(2345678)*1024
	if info.UsedBytes != wantUsed {
		t.Errorf("UsedBytes = %d, want %d", info.UsedBytes, wantUsed)
	}
	wantSwapTotal := uint64(8192000) * 1024
	if info.SwapTotalBytes != wantSwapTotal {
		t.Errorf("SwapTotalBytes = %d, want %d", info.SwapTotalBytes, wantSwapTotal)
	}
	wantSwapUsed := wantSwapTotal - uint64(2048000)*1024
	if info.SwapUsedBytes != wantSwapUsed {
		t.Errorf("SwapUsedBytes = %d, want %d", info.SwapUsedBytes, wantSwapUsed)
	}
}

func TestParseMeminfo_NoMemAvailable(t *testing.T) {
	// Kernels < 3.14 don't have MemAvailable
	input := `MemTotal:       16384000 kB
MemFree:         1000000 kB
Buffers:          200000 kB
Cached:           300000 kB
SwapTotal:              0 kB
SwapFree:               0 kB
`
	info := parseMeminfo(strings.NewReader(input))
	if !info.Supported {
		t.Fatal("expected Supported = true")
	}
	// Should fall back to MemFree + Buffers + Cached
	wantAvailable := uint64(1000000+200000+300000) * 1024
	wantUsed := uint64(16384000)*1024 - wantAvailable
	if info.UsedBytes != wantUsed {
		t.Errorf("UsedBytes = %d, want %d", info.UsedBytes, wantUsed)
	}
}

func TestParseMeminfo_Empty(t *testing.T) {
	info := parseMeminfo(strings.NewReader(""))
	// Zero values, but still marked supported (parser doesn't fail on empty)
	if info.TotalBytes != 0 {
		t.Errorf("TotalBytes = %d, want 0", info.TotalBytes)
	}
}

func TestParseLoadavg_Typical(t *testing.T) {
	input := "0.52 0.41 0.38 1/234 5678\n"
	info := parseLoadavg(strings.NewReader(input))
	if !info.Supported {
		t.Fatal("expected Supported = true")
	}
	if info.Load1 != 0.52 {
		t.Errorf("Load1 = %f, want 0.52", info.Load1)
	}
	if info.Load5 != 0.41 {
		t.Errorf("Load5 = %f, want 0.41", info.Load5)
	}
	if info.Load15 != 0.38 {
		t.Errorf("Load15 = %f, want 0.38", info.Load15)
	}
	if info.NumCPU < 1 {
		t.Errorf("NumCPU = %d, want >= 1", info.NumCPU)
	}
}

func TestParseLoadavg_HighLoad(t *testing.T) {
	input := "12.50 9.30 6.10 5/890 12345\n"
	info := parseLoadavg(strings.NewReader(input))
	if info.Load1 != 12.50 {
		t.Errorf("Load1 = %f, want 12.50", info.Load1)
	}
}

func TestParseLoadavg_Empty(t *testing.T) {
	info := parseLoadavg(strings.NewReader(""))
	if info.Supported {
		t.Error("expected Supported = false for empty input")
	}
}
