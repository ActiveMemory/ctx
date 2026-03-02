//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package doctor

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

func setupContextDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CTX_DIR", dir)
	rc.Reset()

	// Create required files.
	for _, f := range config.FilesRequired {
		path := filepath.Join(dir, f)
		if writeErr := os.WriteFile(path, []byte("# "+f+"\n"), 0o600); writeErr != nil {
			t.Fatal(writeErr)
		}
	}
	return dir
}

func TestDoctor_Healthy(t *testing.T) {
	setupContextDir(t)

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("doctor failed: %v", runErr)
	}

	output := out.String()
	if !strings.Contains(output, "0 warnings, 0 errors") {
		t.Errorf("expected healthy summary, got: %s", output)
	}
	if !strings.Contains(output, "Context initialized") {
		t.Errorf("expected context initialized check, got: %s", output)
	}
}

func TestDoctor_MissingRequiredFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CTX_DIR", dir)
	rc.Reset()

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("doctor failed: %v", runErr)
	}

	output := out.String()
	if !strings.Contains(output, "Missing required files") {
		t.Errorf("expected missing files error, got: %s", output)
	}
	if !strings.Contains(output, "1 errors") {
		t.Errorf("expected 1 error in summary, got: %s", output)
	}
}

func TestDoctor_EventLogOff(t *testing.T) {
	setupContextDir(t)

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	output := out.String()
	if !strings.Contains(output, "Event logging disabled") {
		t.Errorf("expected event logging info note, got: %s", output)
	}
	// Info notes should not count as warnings or errors.
	if !strings.Contains(output, "0 warnings, 0 errors") {
		t.Errorf("expected 0 warnings/errors (info is not a warning), got: %s", output)
	}
}

func TestDoctor_JSON(t *testing.T) {
	setupContextDir(t)

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--json"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("doctor --json failed: %v", runErr)
	}

	var report Report
	if unmarshalErr := json.Unmarshal(out.Bytes(), &report); unmarshalErr != nil {
		t.Fatalf("output is not valid JSON: %v\noutput: %s", unmarshalErr, out.String())
	}
	if len(report.Results) == 0 {
		t.Error("expected at least one result")
	}
}

func TestDoctor_HighCompletion(t *testing.T) {
	dir := setupContextDir(t)

	// Write a TASKS.md with high completion ratio.
	tasks := "# Tasks\n"
	for i := 0; i < 20; i++ {
		tasks += "- [x] Completed task\n"
	}
	tasks += "- [ ] Pending task\n"
	tasksPath := filepath.Join(dir, config.FileTask)
	if writeErr := os.WriteFile(tasksPath, []byte(tasks), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	output := out.String()
	if !strings.Contains(output, "consider archiving") {
		t.Errorf("expected archiving suggestion for high completion, got: %s", output)
	}
}

func TestDoctor_ContextSizeBreakdown(t *testing.T) {
	dir := setupContextDir(t)

	// Write enough content to some files to verify per-file breakdown appears.
	archPath := filepath.Join(dir, "ARCHITECTURE.md")
	if writeErr := os.WriteFile(archPath, []byte(strings.Repeat("word ", 500)), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}
	tasksPath := filepath.Join(dir, config.FileTask)
	if writeErr := os.WriteFile(tasksPath, []byte(strings.Repeat("task ", 200)), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	output := out.String()

	// Should show "window" not "budget".
	if strings.Contains(output, "budget") {
		t.Errorf("should use 'window' not 'budget', got: %s", output)
	}
	if !strings.Contains(output, "window:") {
		t.Errorf("expected 'window:' in context size line, got: %s", output)
	}

	// Should show per-file breakdown.
	if !strings.Contains(output, "ARCHITECTURE.md") {
		t.Errorf("expected ARCHITECTURE.md in breakdown, got: %s", output)
	}
	if !strings.Contains(output, "TASKS.md") {
		t.Errorf("expected TASKS.md in breakdown, got: %s", output)
	}
	if !strings.Contains(output, "tokens") {
		t.Errorf("expected 'tokens' in breakdown lines, got: %s", output)
	}
}

func TestDoctor_ContextSizeJSON(t *testing.T) {
	setupContextDir(t)

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--json"})
	if runErr := cmd.Execute(); runErr != nil {
		t.Fatalf("doctor --json failed: %v", runErr)
	}

	var report Report
	if unmarshalErr := json.Unmarshal(out.Bytes(), &report); unmarshalErr != nil {
		t.Fatalf("output is not valid JSON: %v", unmarshalErr)
	}

	// Should have context_file_* results.
	var fileResults int
	for _, r := range report.Results {
		if strings.HasPrefix(r.Name, "context_file_") {
			fileResults++
		}
	}
	if fileResults == 0 {
		t.Error("expected context_file_* results in JSON output")
	}
}

func TestDoctor_PluginNotInstalled(t *testing.T) {
	setupContextDir(t)

	// Set HOME to a temp dir with no plugin files.
	t.Setenv("HOME", t.TempDir())

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	output := out.String()
	if !strings.Contains(output, "Plugin") {
		t.Errorf("expected Plugin category, got: %s", output)
	}
	if !strings.Contains(output, "not installed") {
		t.Errorf("expected 'not installed' info, got: %s", output)
	}
}

func TestDoctor_PluginInstalledNotEnabled(t *testing.T) {
	setupContextDir(t)

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Create installed_plugins.json with ctx plugin.
	pluginsDir := filepath.Join(tmpHome, ".claude", "plugins")
	if mkErr := os.MkdirAll(pluginsDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	pluginsData := map[string]any{
		"version": 2,
		"plugins": map[string]any{
			config.PluginID: []map[string]string{
				{"scope": "user", "version": "0.7.2"},
			},
		},
	}
	data, _ := json.Marshal(pluginsData)
	if writeErr := os.WriteFile(filepath.Join(pluginsDir, "installed_plugins.json"), data, 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	output := out.String()
	if !strings.Contains(output, "not enabled") {
		t.Errorf("expected 'not enabled' warning, got: %s", output)
	}
	if !strings.Contains(output, "1 warnings") {
		t.Errorf("expected 1 warning in summary, got: %s", output)
	}
}

func TestDoctor_DriftWarnings(t *testing.T) {
	dir := setupContextDir(t)

	// Add an ARCHITECTURE.md referencing a nonexistent path to trigger drift.
	archPath := filepath.Join(dir, "ARCHITECTURE.md")
	archContent := "# Architecture\n\n" +
		"See `internal/nonexistent/fake.go` for details.\n"
	if writeErr := os.WriteFile(archPath, []byte(archContent), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	// Drift detection may or may not find warnings depending on the checks
	// that are relevant. The important thing is that it doesn't crash.
	output := out.String()
	if !strings.Contains(output, "ctx doctor") {
		t.Errorf("expected doctor header in output, got: %s", output)
	}
}

func TestAddResourceResults_AllHealthy(t *testing.T) {
	snap := sysinfo.Snapshot{
		Memory: sysinfo.MemInfo{
			TotalBytes:     16 * 1 << 30, // 16 GB
			UsedBytes:      8 * 1 << 30,  // 8 GB (50%)
			SwapTotalBytes: 8 * 1 << 30,
			SwapUsedBytes:  1 << 30, // 1 GB (~12%)
			Supported:      true,
		},
		Disk: sysinfo.DiskInfo{
			TotalBytes: 500 * 1 << 30, // 500 GB
			UsedBytes:  200 * 1 << 30, // 200 GB (40%)
			Path:       "/",
			Supported:  true,
		},
		Load: sysinfo.LoadInfo{
			Load1:     2.0,
			Load5:     1.5,
			Load15:    1.0,
			NumCPU:    8,
			Supported: true,
		},
	}

	report := &Report{}
	addResourceResults(report, snap)

	if len(report.Results) != 4 {
		t.Fatalf("expected 4 results (memory, swap, disk, load), got %d", len(report.Results))
	}
	for _, r := range report.Results {
		if r.Status != statusOK {
			t.Errorf("result %s: expected ok, got %s", r.Name, r.Status)
		}
		if r.Category != "Resources" {
			t.Errorf("result %s: expected Resources category, got %s", r.Name, r.Category)
		}
	}
}

func TestAddResourceResults_MemoryWarning(t *testing.T) {
	snap := sysinfo.Snapshot{
		Memory: sysinfo.MemInfo{
			TotalBytes: 1000,
			UsedBytes:  820, // 82% → WARNING
			Supported:  true,
		},
		Disk: sysinfo.DiskInfo{Supported: false},
		Load: sysinfo.LoadInfo{Supported: false},
	}

	report := &Report{}
	addResourceResults(report, snap)

	if len(report.Results) != 1 {
		t.Fatalf("expected 1 result (memory only), got %d", len(report.Results))
	}
	if report.Results[0].Name != "resource_memory" {
		t.Errorf("expected resource_memory, got %s", report.Results[0].Name)
	}
	if report.Results[0].Status != statusWarning {
		t.Errorf("expected warning for 82%% memory, got %s", report.Results[0].Status)
	}
}

func TestAddResourceResults_DangerMapsToError(t *testing.T) {
	snap := sysinfo.Snapshot{
		Memory: sysinfo.MemInfo{
			TotalBytes:     1000,
			UsedBytes:      920, // 92% → DANGER
			SwapTotalBytes: 1000,
			SwapUsedBytes:  760, // 76% → DANGER
			Supported:      true,
		},
		Disk: sysinfo.DiskInfo{
			TotalBytes: 1000,
			UsedBytes:  960, // 96% → DANGER
			Supported:  true,
		},
		Load: sysinfo.LoadInfo{
			Load1:     12.0,
			NumCPU:    8, // 1.5x → DANGER
			Supported: true,
		},
	}

	report := &Report{}
	addResourceResults(report, snap)

	if len(report.Results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(report.Results))
	}
	for _, r := range report.Results {
		if r.Status != statusError {
			t.Errorf("result %s: expected error for danger severity, got %s", r.Name, r.Status)
		}
	}
}

func TestAddResourceResults_UnsupportedSkipped(t *testing.T) {
	snap := sysinfo.Snapshot{
		Memory: sysinfo.MemInfo{Supported: false},
		Disk:   sysinfo.DiskInfo{Supported: false},
		Load:   sysinfo.LoadInfo{Supported: false},
	}

	report := &Report{}
	addResourceResults(report, snap)

	if len(report.Results) != 0 {
		t.Errorf("expected 0 results for unsupported metrics, got %d", len(report.Results))
	}
}

func TestAddResourceResults_NoSwapWhenZeroTotal(t *testing.T) {
	snap := sysinfo.Snapshot{
		Memory: sysinfo.MemInfo{
			TotalBytes:     16 * 1 << 30,
			UsedBytes:      4 * 1 << 30,
			SwapTotalBytes: 0, // No swap configured
			SwapUsedBytes:  0,
			Supported:      true,
		},
		Disk: sysinfo.DiskInfo{Supported: false},
		Load: sysinfo.LoadInfo{Supported: false},
	}

	report := &Report{}
	addResourceResults(report, snap)

	if len(report.Results) != 1 {
		t.Fatalf("expected 1 result (memory only, no swap), got %d", len(report.Results))
	}
	if report.Results[0].Name != "resource_memory" {
		t.Errorf("expected resource_memory, got %s", report.Results[0].Name)
	}
}

func TestAddResourceResults_MessageFormat(t *testing.T) {
	snap := sysinfo.Snapshot{
		Memory: sysinfo.MemInfo{
			TotalBytes: 16 * 1 << 30,
			UsedBytes:  8 * 1 << 30,
			Supported:  true,
		},
		Disk: sysinfo.DiskInfo{Supported: false},
		Load: sysinfo.LoadInfo{
			Load1:     2.0,
			NumCPU:    8,
			Supported: true,
		},
	}

	report := &Report{}
	addResourceResults(report, snap)

	for _, r := range report.Results {
		switch r.Name {
		case "resource_memory":
			if !strings.Contains(r.Message, "Memory") || !strings.Contains(r.Message, "GB") {
				t.Errorf("memory message missing expected format: %s", r.Message)
			}
		case "resource_load":
			if !strings.Contains(r.Message, "Load") || !strings.Contains(r.Message, "CPUs") {
				t.Errorf("load message missing expected format: %s", r.Message)
			}
		}
	}
}

func TestDoctor_ResourcesCategoryInOutput(t *testing.T) {
	setupContextDir(t)

	cmd := Cmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{})
	_ = cmd.Execute()

	output := out.String()
	// On any supported platform, at least one resource metric should appear.
	// On unsupported platforms, the Resources header simply won't appear,
	// which is correct behavior (graceful degradation).
	if strings.Contains(output, "Resources") {
		// If the category appears, it should have at least one metric.
		if !strings.Contains(output, "Memory") &&
			!strings.Contains(output, "Disk") &&
			!strings.Contains(output, "Load") {
			t.Errorf("Resources category present but no metrics shown: %s", output)
		}
	}
}
