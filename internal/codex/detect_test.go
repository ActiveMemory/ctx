//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"os"
	"path/filepath"
	"testing"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
)

// fakeHome returns a temp dir wired as $CODEX_HOME.
func fakeHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv(cfgCodex.EnvHome, home)
	return home
}

// fakeCodexBinary puts an executable `codex` on PATH (or, when
// present is false, an empty PATH).
func fakeCodexBinary(t *testing.T, present bool) {
	t.Helper()
	binDir := t.TempDir()
	if present {
		bin := filepath.Join(binDir, cfgCodex.Binary)
		if err := os.WriteFile(bin, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatalf("seed fake codex: %v", err)
		}
	}
	t.Setenv("PATH", binDir)
}

func installPlugin(t *testing.T, home string) {
	t.Helper()
	dir := filepath.Join(
		home, cfgCodex.DirPlugins, cfgCodex.DirPluginCache,
		cfgCodex.MarketplaceID, cfgCodex.PluginName,
		cfgCodex.PluginVersionLocal,
	)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir plugin cache: %v", err)
	}
}

func writeConfig(t *testing.T, home, content string) {
	t.Helper()
	if err := os.WriteFile(
		filepath.Join(home, cfgCodex.FileConfigTOML), []byte(content), 0o644,
	); err != nil {
		t.Fatalf("seed config.toml: %v", err)
	}
}

func TestHome_HonorsEnv(t *testing.T) {
	home := fakeHome(t)
	if got := Home(); got != home {
		t.Fatalf("Home() = %q, want %q", got, home)
	}
}

func TestHome_DefaultsToDotCodex(t *testing.T) {
	t.Setenv(cfgCodex.EnvHome, "")
	userHome, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no user home")
	}
	want := filepath.Join(userHome, cfgCodex.DirHome)
	if got := Home(); got != want {
		t.Fatalf("Home() = %q, want %q", got, want)
	}
}

func TestPluginInstalled(t *testing.T) {
	home := fakeHome(t)
	if PluginInstalled(home) {
		t.Fatal("empty home reported installed")
	}
	if PluginInstalled("") {
		t.Fatal("empty path reported installed")
	}
	// Plugin dir without a version subdir is not an install.
	if err := os.MkdirAll(filepath.Join(
		home, cfgCodex.DirPlugins, cfgCodex.DirPluginCache,
		cfgCodex.MarketplaceID, cfgCodex.PluginName,
	), 0o755); err != nil {
		t.Fatal(err)
	}
	if PluginInstalled(home) {
		t.Fatal("plugin dir without version reported installed")
	}
	installPlugin(t, home)
	if !PluginInstalled(home) {
		t.Fatal("cached plugin not detected")
	}
}

func TestPluginEnabled(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    bool
	}{
		{"no file", "", false},
		{"no table", "model = \"gpt-5\"\n", false},
		{"header only", cfgCodex.TOMLHeaderPluginCtx + "\n", true},
		{"enabled true", cfgCodex.TOMLHeaderPluginCtx + "\nenabled = true\n", true},
		{"enabled true with comment", cfgCodex.TOMLHeaderPluginCtx + "\nenabled = true # keep\n", true},
		{"enabled false", cfgCodex.TOMLHeaderPluginCtx + "\nenabled = false\n", false},
		{"indented header", "  " + cfgCodex.TOMLHeaderPluginCtx + "  \n  enabled=true\n", true},
		{
			"enabled false belongs to another table",
			cfgCodex.TOMLHeaderPluginCtx + "\n[plugins.\"other@x\"]\nenabled = false\n",
			true,
		},
		{
			"table later in file",
			"[mcp_servers.foo]\ncommand = \"foo\"\n\n" + cfgCodex.TOMLHeaderPluginCtx + "\n",
			true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			home := fakeHome(t)
			if tc.content != "" {
				writeConfig(t, home, tc.content)
			}
			if got := PluginEnabled(home); got != tc.want {
				t.Fatalf("PluginEnabled = %v, want %v", got, tc.want)
			}
		})
	}
	if PluginEnabled("") {
		t.Fatal("empty home reported enabled")
	}
}

func TestDetect(t *testing.T) {
	t.Run("absent", func(t *testing.T) {
		fakeHome(t)
		fakeCodexBinary(t, false)
		if got := Detect(); got != StateAbsent {
			t.Fatalf("Detect = %v, want StateAbsent", got)
		}
	})
	t.Run("not installed", func(t *testing.T) {
		fakeHome(t)
		fakeCodexBinary(t, true)
		if got := Detect(); got != StatePluginNotInstalled {
			t.Fatalf("Detect = %v, want StatePluginNotInstalled", got)
		}
	})
	t.Run("installed not enabled", func(t *testing.T) {
		home := fakeHome(t)
		fakeCodexBinary(t, true)
		installPlugin(t, home)
		if got := Detect(); got != StatePluginInstalledNotEnabled {
			t.Fatalf("Detect = %v, want StatePluginInstalledNotEnabled", got)
		}
	})
	t.Run("ready", func(t *testing.T) {
		home := fakeHome(t)
		fakeCodexBinary(t, true)
		installPlugin(t, home)
		writeConfig(t, home, cfgCodex.TOMLHeaderPluginCtx+"\nenabled = true\n")
		if got := Detect(); got != StatePluginReady {
			t.Fatalf("Detect = %v, want StatePluginReady", got)
		}
	})
}

func TestState_String(t *testing.T) {
	for _, s := range []State{
		StateAbsent, StatePluginNotInstalled,
		StatePluginInstalledNotEnabled, StatePluginReady,
	} {
		if s.String() == "" {
			t.Fatalf("State(%d).String() is empty: text key missing", s)
		}
	}
}

func TestUnwired(t *testing.T) {
	origDir, _ := os.Getwd()
	project := t.TempDir()
	if err := os.Chdir(project); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	home := fakeHome(t)
	fakeCodexBinary(t, false)
	if Unwired() {
		t.Fatal("no codex binary: want wired (no hint)")
	}

	fakeCodexBinary(t, true)
	if !Unwired() {
		t.Fatal("codex present, nothing wired: want unwired")
	}

	// Enabled in config but no Codex-variant cache (stale flag or
	// legacy Claude variant): still unwired, the hint must show.
	writeConfig(t, home, cfgCodex.TOMLHeaderPluginCtx+"\n")
	if !Unwired() {
		t.Fatal("plugin enabled without codex-variant cache: want unwired")
	}

	seedVariantCache(t, home)
	if Unwired() {
		t.Fatal("plugin enabled: want wired")
	}

	writeConfig(t, home, "")
	if err := os.MkdirAll(cfgCodex.Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(cfgCodex.Dir, cfgCodex.FileHooksJSON), []byte("{}"), 0o644,
	); err != nil {
		t.Fatal(err)
	}
	if Unwired() {
		t.Fatal("hooks.json present: want wired")
	}
}

func TestSkillManaged(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    bool
	}{
		{"managed", "---\nname: ctx-agent\ndescription: x\n---\nbody\n", true},
		{"managed quoted", "---\nname: \"ctx-agent\"\n---\n", true},
		{"other name", "---\nname: ctx-other\n---\n", false},
		{"no frontmatter", "# ctx-agent\nname: ctx-agent\n", false},
		{"name after frontmatter", "---\ndescription: x\n---\nname: ctx-agent\n", false},
		{"empty", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := SkillManaged([]byte(tc.content), "ctx-agent"); got != tc.want {
				t.Fatalf("SkillManaged = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestPluginNativeVariant distinguishes the Codex plugin cache
// layout (.codex-plugin/) from the legacy Claude Code variant
// (.claude-plugin/) that a stale marketplace source delivers.
func TestPluginNativeVariant(t *testing.T) {
	home := fakeHome(t)
	version := filepath.Join(
		home, cfgCodex.DirPlugins, cfgCodex.DirPluginCache,
		cfgCodex.MarketplaceID, cfgCodex.PluginName, "0.8.1",
	)

	if PluginNativeVariant(home) {
		t.Fatal("no cache: want false")
	}

	// Claude variant: .claude-plugin/ manifest dir.
	if mkErr := os.MkdirAll(
		filepath.Join(version, ".claude-plugin"), 0o755,
	); mkErr != nil {
		t.Fatal(mkErr)
	}
	if PluginNativeVariant(home) {
		t.Fatal("claude variant: want false")
	}

	// Codex variant appears (e.g. after reinstall).
	if mkErr := os.MkdirAll(
		filepath.Join(version, cfgCodex.DirPluginManifest), 0o755,
	); mkErr != nil {
		t.Fatal(mkErr)
	}
	if !PluginNativeVariant(home) {
		t.Fatal("codex variant present: want true")
	}
}

// seedVariantCache creates a cached plugin copy with the Codex
// manifest dir so PluginNativeVariant reports true.
func seedVariantCache(t *testing.T, home string) {
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
