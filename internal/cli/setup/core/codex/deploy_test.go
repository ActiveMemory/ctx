//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
)

func testCmd(buf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	return cmd
}

// withTempProjectDir chdirs into a fresh project dir and points
// $CODEX_HOME at an empty dir so the user's real Codex install
// never leaks into the test.
func withTempProjectDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })
	t.Setenv(cfgCodex.EnvHome, t.TempDir())
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
	return filepath.Join(cfgSetup.SkillsPathCodex, name, cfgHook.FileSKILLMd)
}

func TestDeploy_FreshProjectCreatesAllArtifacts(t *testing.T) {
	withTempProjectDir(t)

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("Deploy: %v", err)
	}

	embedded, err := agent.CodexHooksJSON()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(readFile(t, cfgSetup.HooksPathCodex), embedded) {
		t.Fatal("hooks.json is not the embedded manifest")
	}
	toml := string(readFile(t, cfgSetup.MCPConfigPathCodex))
	if !strings.Contains(toml, cfgCodex.TOMLHeaderMCPCtx) ||
		!strings.Contains(toml, `args = ["mcp", "serve"]`) {
		t.Fatalf("config.toml missing MCP table:\n%s", toml)
	}
	if _, statErr := os.Stat(cfgHook.FileAgentsMd); statErr != nil {
		t.Fatal("AGENTS.md not created")
	}
	skills, err := agent.CodexSkills()
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) == 0 {
		t.Fatal("no embedded Codex skills")
	}
	for name, content := range skills {
		if !bytes.Equal(readFile(t, skillPath(name)), content) {
			t.Fatalf("skill %s not deployed verbatim", name)
		}
	}
	out := buf.String()
	if !strings.Contains(out, "/hooks") {
		t.Fatalf("summary lacks the /hooks trust reminder:\n%s", out)
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
	hooks := readFile(t, cfgSetup.HooksPathCodex)
	toml := readFile(t, cfgSetup.MCPConfigPathCodex)
	agents := readFile(t, cfgHook.FileAgentsMd)
	skill := readFile(t, skillPath("ctx-agent"))

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("second Deploy: %v", err)
	}
	if !bytes.Equal(hooks, readFile(t, cfgSetup.HooksPathCodex)) {
		t.Fatal("hooks.json rewritten on second run")
	}
	if !bytes.Equal(toml, readFile(t, cfgSetup.MCPConfigPathCodex)) {
		t.Fatal("config.toml rewritten on second run")
	}
	if !bytes.Equal(agents, readFile(t, cfgHook.FileAgentsMd)) {
		t.Fatal("AGENTS.md rewritten on second run")
	}
	if !bytes.Equal(skill, readFile(t, skillPath("ctx-agent"))) {
		t.Fatal("skill rewritten on second run")
	}
	out := buf.String()
	for _, path := range []string{cfgSetup.HooksPathCodex, cfgSetup.MCPConfigPathCodex, skillPath("ctx-agent")} {
		if !strings.Contains(out, path+" (up to date, skipped)") {
			t.Fatalf("expected skip line for %s:\n%s", path, out)
		}
	}
	if strings.Contains(out, "✓ "+cfgSetup.HooksPathCodex) {
		t.Fatalf("second run reported a write:\n%s", out)
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
	skills, _ := agent.CodexSkills()
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
	foreign := []byte("---\nname: my-own-skill\n---\nmine\n")
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
	// Other skills still deploy.
	if _, statErr := os.Stat(skillPath("ctx-remember")); statErr != nil {
		t.Fatal("sibling skill not deployed after a rejection")
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

// seedCodexVariantCache creates an installed-plugin cache copy
// carrying the Codex manifest dir, so PluginNativeVariant is true.
func seedCodexVariantCache(t *testing.T, home string) {
	t.Helper()
	manifest := filepath.Join(
		home, cfgCodex.DirPlugins, cfgCodex.DirPluginCache,
		cfgCodex.MarketplaceID, cfgCodex.PluginName, "0.8.1",
		cfgCodex.DirPluginManifest,
	)
	if mkErr := os.MkdirAll(filepath.Clean(manifest), 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}
}

func TestDeploy_PluginEnabledShortCircuitsToAgentsMd(t *testing.T) {
	withTempProjectDir(t)
	home := os.Getenv(cfgCodex.EnvHome)
	seedFile(t, filepath.Join(home, cfgCodex.FileConfigTOML),
		[]byte(cfgCodex.TOMLHeaderPluginCtx+"\nenabled = true\n"))
	seedCodexVariantCache(t, home)

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("Deploy: %v", err)
	}
	if _, statErr := os.Stat(cfgHook.FileAgentsMd); statErr != nil {
		t.Fatal("AGENTS.md not deployed in plugin mode")
	}
	for _, path := range []string{
		cfgSetup.HooksPathCodex, cfgSetup.MCPConfigPathCodex, cfgSetup.SkillsPathCodex,
	} {
		if _, statErr := os.Stat(path); statErr == nil {
			t.Fatalf("%s deployed although the plugin is enabled", path)
		}
	}
	if !strings.Contains(buf.String(), "plugin is enabled") {
		t.Fatalf("expected plugin-active notice:\n%s", buf.String())
	}
}

// TestDeploy_PluginWrongVariantDeploysEverything covers the
// customer without Claude Code whose plugin install silently
// delivered the legacy Claude variant: Deploy must not
// short-circuit; it warns and deploys the project-local route.
func TestDeploy_PluginWrongVariantDeploysEverything(t *testing.T) {
	withTempProjectDir(t)
	home := os.Getenv(cfgCodex.EnvHome)
	seedFile(t, filepath.Join(home, cfgCodex.FileConfigTOML),
		[]byte(cfgCodex.TOMLHeaderPluginCtx+"\nenabled = true\n"))
	// Claude-variant cache: manifest dir is .claude-plugin/.
	claudeManifest := filepath.Join(
		home, cfgCodex.DirPlugins, cfgCodex.DirPluginCache,
		cfgCodex.MarketplaceID, cfgCodex.PluginName, "0.8.1",
		".claude-plugin",
	)
	if mkErr := os.MkdirAll(filepath.Clean(claudeManifest), 0o755); mkErr != nil {
		t.Fatal(mkErr)
	}

	var buf bytes.Buffer
	if err := Deploy(testCmd(&buf)); err != nil {
		t.Fatalf("Deploy: %v", err)
	}
	if _, statErr := os.Stat(cfgSetup.HooksPathCodex); statErr != nil {
		t.Fatal("hooks.json not deployed despite wrong-variant plugin")
	}
	if _, statErr := os.Stat(cfgSetup.MCPConfigPathCodex); statErr != nil {
		t.Fatal("config.toml not deployed despite wrong-variant plugin")
	}
	if !strings.Contains(buf.String(), "not the Codex variant") {
		t.Fatalf("expected wrong-variant warning:\n%s", buf.String())
	}
}

func TestDeploy_PluginDisabledDeploysEverything(t *testing.T) {
	withTempProjectDir(t)
	home := os.Getenv(cfgCodex.EnvHome)
	seedFile(t, filepath.Join(home, cfgCodex.FileConfigTOML),
		[]byte(cfgCodex.TOMLHeaderPluginCtx+"\nenabled = false\n"))

	if err := Deploy(testCmd(&bytes.Buffer{})); err != nil {
		t.Fatalf("Deploy: %v", err)
	}
	if _, statErr := os.Stat(cfgSetup.HooksPathCodex); statErr != nil {
		t.Fatal("hooks.json not deployed although the plugin is disabled")
	}
}
