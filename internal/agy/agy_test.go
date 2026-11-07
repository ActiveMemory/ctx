//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	cfgAgy "github.com/ActiveMemory/ctx/internal/config/agy"
	cfgAsset "github.com/ActiveMemory/ctx/internal/config/asset"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
)

func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

func TestHome_EnvOverride(t *testing.T) {
	temp := t.TempDir()
	t.Setenv(cfgAgy.EnvHome, temp)
	if got := Home(); got != temp {
		t.Fatalf("Home() = %q, want %q", got, temp)
	}
}

func TestHome_Default(t *testing.T) {
	t.Setenv(cfgAgy.EnvHome, "")
	home := Home()
	if home == "" {
		t.Fatal("Home() returned empty string")
	}
	if !strings.HasSuffix(home, cfgAgy.DirGemini) {
		t.Fatalf("Home() = %q, want suffix %q", home, cfgAgy.DirGemini)
	}
}

func TestConfigDir_EnvOverride(t *testing.T) {
	temp := t.TempDir()
	t.Setenv(cfgAgy.EnvConfigDir, temp)
	if got := ConfigDir(); got != temp {
		t.Fatalf("ConfigDir() = %q, want %q", got, temp)
	}
}

func TestConfigDir_FromHome(t *testing.T) {
	temp := t.TempDir()
	t.Setenv(cfgAgy.EnvConfigDir, "")
	t.Setenv(cfgAgy.EnvHome, temp)
	want := filepath.Join(temp, cfgAgy.DirConfig)
	if got := ConfigDir(); got != want {
		t.Fatalf("ConfigDir() = %q, want %q", got, want)
	}
}

func TestPluginInstalledAndEnabled(t *testing.T) {
	cfgDir := t.TempDir()
	if PluginInstalled(cfgDir) {
		t.Fatal("PluginInstalled returned true for empty dir")
	}
	if PluginEnabled(cfgDir) {
		t.Fatal("PluginEnabled returned true for empty dir")
	}

	pluginFile := filepath.Join(
		cfgDir, cfgAgy.DirPlugins, cfgAgy.PluginName, cfgAsset.FilePluginJSON,
	)
	if err := os.MkdirAll(filepath.Dir(pluginFile), 0o755); err != nil {
		t.Fatal(err)
	}

	// Enabled plugin
	if err := os.WriteFile(pluginFile, []byte(`{"name":"ctx"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if !PluginInstalled(cfgDir) {
		t.Fatal("PluginInstalled returned false after writing plugin.json")
	}
	if !PluginEnabled(cfgDir) {
		t.Fatal("PluginEnabled returned false for enabled plugin.json")
	}

	// Explicitly disabled plugin
	if err := os.WriteFile(
		pluginFile, []byte(`{"name":"ctx","disabled":true}`), 0o644,
	); err != nil {
		t.Fatal(err)
	}
	if !PluginInstalled(cfgDir) {
		t.Fatal("PluginInstalled returned false for disabled plugin.json")
	}
	if PluginEnabled(cfgDir) {
		t.Fatal("PluginEnabled returned true for disabled plugin.json")
	}
}

func mockAgyOnPath(t *testing.T) {
	t.Helper()
	binDir := t.TempDir()
	binName := cfgAgy.Binary
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	target := filepath.Join(binDir, binName)
	if err := os.WriteFile(target, []byte("mock"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func TestDetect_States(t *testing.T) {
	// Absent
	t.Setenv("PATH", t.TempDir())
	if got := Detect(); got != StateAbsent {
		t.Fatalf("Detect() = %v, want StateAbsent", got)
	}

	// Mock binary
	mockAgyOnPath(t)
	cfgDir := t.TempDir()
	t.Setenv(cfgAgy.EnvConfigDir, cfgDir)

	// PluginNotInstalled
	if got := Detect(); got != StatePluginNotInstalled {
		t.Fatalf("Detect() = %v, want StatePluginNotInstalled", got)
	}

	pluginFile := filepath.Join(
		cfgDir, cfgAgy.DirPlugins, cfgAgy.PluginName, cfgAsset.FilePluginJSON,
	)
	if err := os.MkdirAll(filepath.Dir(pluginFile), 0o755); err != nil {
		t.Fatal(err)
	}

	// PluginInstalledNotEnabled
	if err := os.WriteFile(
		pluginFile, []byte(`{"name":"ctx","disabled":true}`), 0o644,
	); err != nil {
		t.Fatal(err)
	}
	if got := Detect(); got != StatePluginInstalledNotEnabled {
		t.Fatalf("Detect() = %v, want StatePluginInstalledNotEnabled", got)
	}

	// PluginReady
	if err := os.WriteFile(
		pluginFile, []byte(`{"name":"ctx"}`), 0o644,
	); err != nil {
		t.Fatal(err)
	}
	if got := Detect(); got != StatePluginReady {
		t.Fatalf("Detect() = %v, want StatePluginReady", got)
	}
}

func TestState_String(t *testing.T) {
	states := []State{
		StateAbsent,
		StatePluginNotInstalled,
		StatePluginInstalledNotEnabled,
		StatePluginReady,
	}
	for _, s := range states {
		str := s.String()
		if str == "" {
			t.Errorf("State(%d).String() is empty", s)
		}
	}
}

func TestProjectConfigured(t *testing.T) {
	origDir, _ := os.Getwd()
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	if ProjectConfigured() {
		t.Fatal("ProjectConfigured() = true in empty dir")
	}

	hooksPath := cfgSetup.HooksPathAgy
	if err := os.MkdirAll(filepath.Dir(hooksPath), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(hooksPath, []byte(`{"other":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if ProjectConfigured() {
		t.Fatal("ProjectConfigured() = true without ctx key")
	}

	if err := os.WriteFile(hooksPath, []byte(`{"ctx":{"SessionStart":[]}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if !ProjectConfigured() {
		t.Fatal("ProjectConfigured() = false with ctx key")
	}
}

func TestUnwired(t *testing.T) {
	origDir, _ := os.Getwd()
	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	// Not on PATH -> not unwired
	t.Setenv("PATH", t.TempDir())
	if Unwired() {
		t.Fatal("Unwired() = true when binary not on PATH")
	}

	// On PATH, no project hooks, no plugin -> unwired
	mockAgyOnPath(t)
	cfgDir := t.TempDir()
	t.Setenv(cfgAgy.EnvConfigDir, cfgDir)
	if !Unwired() {
		t.Fatal("Unwired() = false when binary on PATH and neither route configured")
	}

	// Plugin enabled -> not unwired
	pluginFile := filepath.Join(
		cfgDir, cfgAgy.DirPlugins, cfgAgy.PluginName, cfgAsset.FilePluginJSON,
	)
	if err := os.MkdirAll(filepath.Dir(pluginFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(pluginFile, []byte(`{"name":"ctx"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if Unwired() {
		t.Fatal("Unwired() = true when plugin enabled")
	}

	// Plugin disabled, but project configured -> not unwired
	if err := os.WriteFile(
		pluginFile, []byte(`{"name":"ctx","disabled":true}`), 0o644,
	); err != nil {
		t.Fatal(err)
	}
	hooksPath := cfgSetup.HooksPathAgy
	if err := os.MkdirAll(filepath.Dir(hooksPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hooksPath, []byte(`{"ctx":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if Unwired() {
		t.Fatal("Unwired() = true when project configured")
	}
}

func TestMergeHooks(t *testing.T) {
	embedded := []byte(`{
  "ctx": {
    "SessionStart": [
      {
        "type": "command",
        "command": "ctx system context-load-gate"
      }
    ]
  }
}`)

	// Empty -> OutcomeCreated
	out, outcome, err := MergeHooks(nil, embedded)
	if err != nil {
		t.Fatalf("MergeHooks(nil): %v", err)
	}
	if outcome != OutcomeCreated {
		t.Fatalf("outcome = %v, want OutcomeCreated", outcome)
	}
	if !bytes.Equal(out, embedded) {
		t.Fatal("output does not equal embedded")
	}

	// Idempotent -> OutcomeSkipped
	out2, outcome2, err2 := MergeHooks(embedded, embedded)
	if err2 != nil {
		t.Fatalf("MergeHooks(embedded): %v", err2)
	}
	if outcome2 != OutcomeSkipped {
		t.Fatalf("outcome = %v, want OutcomeSkipped", outcome2)
	}
	if !bytes.Equal(out2, embedded) {
		t.Fatal("output does not equal input")
	}

	// Preserves foreign group & replaces stale ctx -> OutcomeMerged
	existing := []byte(`{
  "other-tool": {
    "SessionStart": []
  },
  "ctx": {
    "SessionStart": [
      {
        "type": "command",
        "command": "old-command"
      }
    ]
  }
}`)
	out3, outcome3, err3 := MergeHooks(existing, embedded)
	if err3 != nil {
		t.Fatalf("MergeHooks(existing): %v", err3)
	}
	if outcome3 != OutcomeMerged {
		t.Fatalf("outcome = %v, want OutcomeMerged", outcome3)
	}
	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(out3, &parsed); err != nil {
		t.Fatal(err)
	}
	if _, hasOther := parsed["other-tool"]; !hasOther {
		t.Fatal("other-tool group was dropped")
	}
	if _, hasCtx := parsed["ctx"]; !hasCtx {
		t.Fatal("ctx group missing")
	}

	// Second run with out3 should skip
	out4, outcome4, err4 := MergeHooks(out3, embedded)
	if err4 != nil {
		t.Fatal(err4)
	}
	if outcome4 != OutcomeSkipped {
		t.Fatalf("outcome = %v, want OutcomeSkipped", outcome4)
	}
	if !bytes.Equal(out3, out4) {
		t.Fatal("re-merged output was rewritten")
	}

	// Empty object -> OutcomeCreated
	_, outcome5, err5 := MergeHooks([]byte(`{}`), embedded)
	if err5 != nil {
		t.Fatal(err5)
	}
	if outcome5 != OutcomeCreated {
		t.Fatalf("outcome = %v, want OutcomeCreated", outcome5)
	}

	// Invalid JSON -> error
	if _, _, err := MergeHooks([]byte(`{invalid`), embedded); err == nil {
		t.Fatal("expected error on invalid JSON")
	}
}

func TestEnsureMCPConfig(t *testing.T) {
	// Empty -> OutcomeCreated
	out, outcome, err := EnsureMCPConfig(nil)
	if err != nil {
		t.Fatalf("EnsureMCPConfig(nil): %v", err)
	}
	if outcome != OutcomeCreated {
		t.Fatalf("outcome = %v, want OutcomeCreated", outcome)
	}
	var parsed map[string]map[string]map[string]any
	if err := json.Unmarshal(out, &parsed); err != nil {
		t.Fatal(err)
	}
	ctxServer, ok := parsed[cfgHook.KeyMCPServers]["ctx"]
	if !ok {
		t.Fatal("mcpServers.ctx missing in created config")
	}
	if ctxServer[cfgHook.KeyCommand] != "ctx" {
		t.Errorf("command = %v, want ctx", ctxServer[cfgHook.KeyCommand])
	}

	// Second run with out -> OutcomeSkipped
	out2, outcome2, err2 := EnsureMCPConfig(out)
	if err2 != nil {
		t.Fatal(err2)
	}
	if outcome2 != OutcomeSkipped {
		t.Fatalf("outcome = %v, want OutcomeSkipped", outcome2)
	}
	if !bytes.Equal(out, out2) {
		t.Fatal("second run rewritten")
	}

	// Merging into existing foreign server -> OutcomeMerged
	existing := []byte(`{
  "mcpServers": {
    "foreign": {
      "command": "foreign-cmd"
    }
  }
}`)
	out3, outcome3, err3 := EnsureMCPConfig(existing)
	if err3 != nil {
		t.Fatal(err3)
	}
	if outcome3 != OutcomeMerged {
		t.Fatalf("outcome = %v, want OutcomeMerged", outcome3)
	}
	var parsed3 map[string]map[string]map[string]any
	if err := json.Unmarshal(out3, &parsed3); err != nil {
		t.Fatal(err)
	}
	if _, hasForeign := parsed3[cfgHook.KeyMCPServers]["foreign"]; !hasForeign {
		t.Fatal("foreign server dropped")
	}
	if _, hasCtx := parsed3[cfgHook.KeyMCPServers]["ctx"]; !hasCtx {
		t.Fatal("ctx server not added")
	}
}

func TestSkillManaged(t *testing.T) {
	managed := []byte("---\nname: ctx-remember\n---\nbody\n")
	if !SkillManaged(managed, "ctx-remember") {
		t.Fatal("SkillManaged returned false for managed skill")
	}
	if SkillManaged(managed, "ctx-agent") {
		t.Fatal("SkillManaged returned true for mismatched skill name")
	}
	foreign := []byte("---\nname: other-skill\n---\nbody\n")
	if SkillManaged(foreign, "ctx-remember") {
		t.Fatal("SkillManaged returned true for foreign skill")
	}
}
