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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/agy"
	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	cfgAgy "github.com/ActiveMemory/ctx/internal/config/agy"
	cfgAsset "github.com/ActiveMemory/ctx/internal/config/asset"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
)

func testCmd(buf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	return cmd
}

func withTempProjectDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	t.Setenv(cfgAgy.EnvHome, t.TempDir())
	t.Setenv(cfgAgy.EnvConfigDir, "")
	return tmp
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func seedFile(t *testing.T, path string, content []byte) {
	t.Helper()
	clean := filepath.Clean(path)
	if err := os.MkdirAll(filepath.Dir(clean), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(clean, content, 0o644); err != nil {
		t.Fatalf("seed %s: %v", path, err)
	}
}

func skillPath(name string) string {
	return filepath.Join(cfgSetup.SkillsPathAgy, name, cfgHook.FileSKILLMd)
}

func mcpConfigPath() string {
	return filepath.Join(agy.ConfigDir(), cfgAsset.FileMCPConfigJSON)
}

func TestDeploy_FreshProjectCreatesAllArtifacts(t *testing.T) {
	withTempProjectDir(t)

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("Deploy: %v", err)
	}

	embedded, err := agent.AgyHooksJSON()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(readFile(t, cfgSetup.HooksPathAgy), embedded) {
		t.Fatal("hooks.json is not the embedded manifest")
	}

	mcpData := readFile(t, mcpConfigPath())
	var mcpCfg struct {
		MCPServers map[string]struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		} `json:"mcpServers"`
	}
	if unmarshalErr := json.Unmarshal(mcpData, &mcpCfg); unmarshalErr != nil {
		t.Fatalf("unmarshal mcp_config.json: %v", unmarshalErr)
	}
	server, ok := mcpCfg.MCPServers["ctx"]
	if !ok || server.Command != "ctx" {
		t.Fatalf("mcp_config.json missing ctx server: %s", string(mcpData))
	}

	if _, statErr := os.Stat(cfgHook.FileAgentsMd); statErr != nil {
		t.Fatal("AGENTS.md not created")
	}

	skills, err := agent.AgySkills()
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) == 0 {
		t.Fatal("no embedded Antigravity skills")
	}
	for name, content := range skills {
		if !bytes.Equal(readFile(t, skillPath(name)), content) {
			t.Fatalf("skill %s not deployed verbatim", name)
		}
	}

	out := buf.String()
	if !strings.Contains(out, "Antigravity will now:") {
		t.Fatalf("summary missing in output:\n%s", out)
	}
	if strings.Contains(out, "skipped") || strings.Contains(out, "warning") {
		t.Fatalf("fresh deploy reported skips or warnings:\n%s", out)
	}
}

func TestDeploy_SecondRunIsIdempotent(t *testing.T) {
	withTempProjectDir(t)

	if err := Deploy(testCmd(&bytes.Buffer{})); err != nil {
		t.Fatalf("first Deploy: %v", err)
	}
	hooks := readFile(t, cfgSetup.HooksPathAgy)
	mcp := readFile(t, mcpConfigPath())
	agents := readFile(t, cfgHook.FileAgentsMd)
	skill := readFile(t, skillPath("ctx-agent"))

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("second Deploy: %v", err)
	}
	if !bytes.Equal(hooks, readFile(t, cfgSetup.HooksPathAgy)) {
		t.Fatal("hooks.json rewritten on second run")
	}
	if !bytes.Equal(mcp, readFile(t, mcpConfigPath())) {
		t.Fatal("mcp_config.json rewritten on second run")
	}
	if !bytes.Equal(agents, readFile(t, cfgHook.FileAgentsMd)) {
		t.Fatal("AGENTS.md rewritten on second run")
	}
	if !bytes.Equal(skill, readFile(t, skillPath("ctx-agent"))) {
		t.Fatal("skill rewritten on second run")
	}

	out := buf.String()
	for _, path := range []string{
		cfgSetup.HooksPathAgy, mcpConfigPath(), skillPath("ctx-agent"),
	} {
		if !strings.Contains(out, path+" (up to date, skipped)") {
			t.Fatalf("expected skip line for %s:\n%s", path, out)
		}
	}
}

func TestDeploy_PreservesForeignHookGroup(t *testing.T) {
	withTempProjectDir(t)

	foreign := []byte(`{
  "other-tool": {
    "SessionStart": []
  }
}`)
	seedFile(t, cfgSetup.HooksPathAgy, foreign)

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("Deploy: %v", err)
	}

	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(readFile(t, cfgSetup.HooksPathAgy), &parsed); err != nil {
		t.Fatal(err)
	}
	if _, ok := parsed["other-tool"]; !ok {
		t.Fatal("other-tool hook group was deleted")
	}
	if _, ok := parsed["ctx"]; !ok {
		t.Fatal("ctx hook group was not added")
	}
	if !strings.Contains(buf.String(), "(merged)") {
		t.Fatalf("expected merged notice:\n%s", buf.String())
	}
}

func TestDeploySkills_RefreshesStaleManagedSkill(t *testing.T) {
	withTempProjectDir(t)
	target := skillPath("ctx-agent")
	seedFile(t, target, []byte("---\nname: ctx-agent\n---\nstale body\n"))

	var buf bytes.Buffer
	if err := deploySkills(testCmd(&buf)); err != nil {
		t.Fatalf("deploySkills: %v", err)
	}
	skills, _ := agent.AgySkills()
	if !bytes.Equal(readFile(t, target), skills["ctx-agent"]) {
		t.Fatal("stale managed skill not refreshed")
	}
	if strings.Contains(buf.String(), target+" (up to date") {
		t.Fatalf("expected refresh, got skip:\n%s", buf.String())
	}
}

func TestDeploySkills_RejectsForeignSkill(t *testing.T) {
	withTempProjectDir(t)
	target := skillPath("ctx-agent")
	foreign := []byte("---\nname: foreign-skill\n---\nbody\n")
	seedFile(t, target, foreign)

	var buf bytes.Buffer
	if err := deploySkills(testCmd(&buf)); err != nil {
		t.Fatalf("deploySkills: %v", err)
	}
	if !bytes.Equal(readFile(t, target), foreign) {
		t.Fatal("foreign skill overwritten")
	}
	if !strings.Contains(buf.String(), target+" (not ctx-managed, skipped)") {
		t.Fatalf("expected rejection notice:\n%s", buf.String())
	}
	if _, statErr := os.Stat(skillPath("ctx-remember")); statErr != nil {
		t.Fatal("sibling skill not deployed after rejection")
	}
}

func TestDeploySkills_RejectsSymlinkTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink behavior varies on Windows in this environment")
	}
	withTempProjectDir(t)
	target := skillPath("ctx-agent")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	realFile := filepath.Join(t.TempDir(), "outside.md")
	if err := os.WriteFile(realFile, []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(realFile, target); err != nil {
		t.Fatal(err)
	}
	if err := deploySkills(testCmd(&bytes.Buffer{})); err == nil {
		t.Fatal("expected symlink rejection, got nil")
	}
}

func TestDeploy_PluginEnabledShortCircuitsToAgentsMd(t *testing.T) {
	withTempProjectDir(t)
	pluginFile := filepath.Join(
		agy.ConfigDir(), cfgAgy.DirPlugins, cfgAgy.PluginName, cfgAsset.FilePluginJSON,
	)
	seedFile(t, pluginFile, []byte(`{"name":"ctx"}`))

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("Deploy: %v", err)
	}
	if _, statErr := os.Stat(cfgHook.FileAgentsMd); statErr != nil {
		t.Fatal("AGENTS.md not deployed in plugin mode")
	}
	for _, path := range []string{
		cfgSetup.HooksPathAgy, mcpConfigPath(), cfgSetup.SkillsPathAgy,
	} {
		if _, statErr := os.Stat(path); statErr == nil {
			t.Fatalf("%s deployed although plugin is enabled", path)
		}
	}
	if !strings.Contains(buf.String(), "plugin is enabled") {
		t.Fatalf("expected plugin-active notice:\n%s", buf.String())
	}
}

func TestDeploy_PluginDisabledDeploysEverything(t *testing.T) {
	withTempProjectDir(t)
	pluginFile := filepath.Join(
		agy.ConfigDir(), cfgAgy.DirPlugins, cfgAgy.PluginName, cfgAsset.FilePluginJSON,
	)
	seedFile(t, pluginFile, []byte(`{"name":"ctx","disabled":true}`))

	if err := Deploy(testCmd(&bytes.Buffer{})); err != nil {
		t.Fatalf("Deploy: %v", err)
	}
	if _, statErr := os.Stat(cfgSetup.HooksPathAgy); statErr != nil {
		t.Fatal("hooks.json not deployed although plugin is disabled")
	}
}
