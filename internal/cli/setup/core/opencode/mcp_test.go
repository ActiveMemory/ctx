//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func testCmd(buf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

func setOpenCodeHome(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("OPENCODE_HOME", tmp)
	return tmp
}

func configPath(home string) string {
	return filepath.Join(home, "opencode.json")
}

func readMCP(t *testing.T, path string) map[string]interface{} {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	parsed := map[string]interface{}{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("config not valid JSON: %v", err)
	}
	return parsed
}

func TestEnsureMCPConfig_CreatesFile(t *testing.T) {
	home := setOpenCodeHome(t)

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig: %v", err)
	}

	parsed := readMCP(t, configPath(home))
	servers, ok := parsed["mcp"].(map[string]interface{})
	if !ok {
		t.Fatal("missing mcp key")
	}
	ctxServer, ok := servers["ctx"].(map[string]interface{})
	if !ok {
		t.Fatal("missing mcp.ctx key")
	}
	if ctxServer["type"] != "local" {
		t.Errorf("type = %q, want local", ctxServer["type"])
	}
	cmdArr, ok := ctxServer["command"].([]interface{})
	if !ok {
		t.Fatalf("command must be an array per OpenCode schema, got %T", ctxServer["command"])
	}
	if got := len(cmdArr); got != 3 {
		t.Fatalf("command length = %d, want 3 (sh -c <script>)", got)
	}
	if cmdArr[0] != "sh" || cmdArr[1] != "-c" {
		t.Errorf("command prefix = [%q %q], want [sh -c]", cmdArr[0], cmdArr[1])
	}
	script, ok := cmdArr[2].(string)
	if !ok {
		t.Fatalf("command[2] must be a script string, got %T", cmdArr[2])
	}
	wantSubs := []string{
		`exec env`,
		`CTX_DIR="$PWD/.context"`,
		`'mcp' 'serve'`,
	}
	for _, s := range wantSubs {
		if !strings.Contains(script, s) {
			t.Errorf("launch script missing %q\nfull script: %s", s, script)
		}
	}
	// Binary path should be absolute (resolved at setup time).
	if !strings.Contains(script, "/") {
		t.Errorf("binary path should be absolute, got script: %s", script)
	}
	if _, hasArgs := ctxServer["args"]; hasArgs {
		t.Error("args field must not be set; OpenCode schema folds args into command array")
	}
	if _, hasEnv := ctxServer["environment"]; hasEnv {
		t.Error("environment field must not be set; CTX_DIR is computed from $PWD inside the sh wrapper")
	}
	enabled, ok := ctxServer["enabled"].(bool)
	if !ok || !enabled {
		t.Errorf("enabled = %v, want true", ctxServer["enabled"])
	}
}

func TestEnsureMCPConfig_TreatsEmptyFileAsAbsent(t *testing.T) {
	home := setOpenCodeHome(t)

	if err := os.WriteFile(
		configPath(home), []byte("   \n\t  "), 0o644,
	); err != nil {
		t.Fatalf("seed empty file: %v", err)
	}

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig on empty file: %v", err)
	}

	parsed := readMCP(t, configPath(home))
	if _, ok := parsed["mcp"].(map[string]interface{}); !ok {
		t.Fatal("mcp key not registered after empty-file path")
	}
}

func TestEnsureMCPConfig_PreservesExistingKeys(t *testing.T) {
	home := setOpenCodeHome(t)

	seed := []byte(`{"theme":"dark","mcp":{"other":{"type":"local"}}}`)
	if err := os.WriteFile(configPath(home), seed, 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig: %v", err)
	}

	parsed := readMCP(t, configPath(home))
	if parsed["theme"] != "dark" {
		t.Errorf("theme not preserved: %v", parsed["theme"])
	}
	servers, _ := parsed["mcp"].(map[string]interface{})
	if _, ok := servers["other"]; !ok {
		t.Error("existing mcp.other entry was lost")
	}
	if _, ok := servers["ctx"]; !ok {
		t.Error("ctx server not added alongside existing entries")
	}
}

func TestEnsureMCPConfig_SkipsWhenCtxAlreadyMatches(t *testing.T) {
	home := setOpenCodeHome(t)

	var seedBuf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&seedBuf)); err != nil {
		t.Fatalf("seed ensureMCPConfig: %v", err)
	}
	seed, err := os.ReadFile(configPath(home))
	if err != nil {
		t.Fatalf("read seeded config: %v", err)
	}

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig: %v", err)
	}

	got, _ := os.ReadFile(configPath(home))
	if string(got) != string(seed) {
		t.Errorf(
			"file rewritten when ctx config already matched: %s", got,
		)
	}
	if !bytes.Contains(buf.Bytes(), []byte("skipped")) {
		t.Errorf(
			"expected 'skipped' in output, got %q", buf.String(),
		)
	}
}

func TestEnsureMCPConfig_RefreshesStaleCtxServer(t *testing.T) {
	home := setOpenCodeHome(t)

	seed := []byte(`{"mcp":{"ctx":{"type":"local","command":["ctx","mcp","serve"],"enabled":false}}}`)
	if err := os.WriteFile(configPath(home), seed, 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig: %v", err)
	}

	parsed := readMCP(t, configPath(home))
	servers, _ := parsed["mcp"].(map[string]interface{})
	ctxServer, _ := servers["ctx"].(map[string]interface{})
	cmdArr, ok := ctxServer["command"].([]interface{})
	if !ok || len(cmdArr) != 3 {
		t.Fatalf("command = %T %v, want refreshed [sh -c script]", ctxServer["command"], ctxServer["command"])
	}
	if enabled, _ := ctxServer["enabled"].(bool); !enabled {
		t.Fatalf("enabled = %v, want true after refresh", ctxServer["enabled"])
	}
	if bytes.Contains(buf.Bytes(), []byte("skipped")) {
		t.Fatalf("expected refresh to rewrite file, got skipped output %q", buf.String())
	}
}

func TestEnsureMCPConfig_QuotesBinaryPathInLaunchScript(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	cmdArr := launchCommand()
	if got := len(cmdArr); got != 3 {
		t.Fatalf("command length = %d, want 3", got)
	}
	script := cmdArr[2]
	if !strings.Contains(script, `'ctx' 'mcp' 'serve'`) {
		t.Fatalf("launch script not safely quoted: %s", script)
	}
}

// TestEnsureMCPConfig_ResolvesBinaryToAbsolutePath covers the
// LookPath-success branch that the QuotesBinaryPath test deliberately
// skips. With a fake `ctx` binary on PATH, launchCommand should embed
// the absolute path so OpenCode can spawn the MCP child even from
// non-interactive shells whose PATH may not contain ctx.
func TestEnsureMCPConfig_ResolvesBinaryToAbsolutePath(t *testing.T) {
	binDir := t.TempDir()
	fake := filepath.Join(binDir, "ctx")
	if err := os.WriteFile(fake, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("seed fake ctx: %v", err)
	}
	t.Setenv("PATH", binDir)

	cmdArr := launchCommand()
	if got := len(cmdArr); got != 3 {
		t.Fatalf("command length = %d, want 3", got)
	}
	script := cmdArr[2]
	wantQuoted := "'" + fake + "' 'mcp' 'serve'"
	if !strings.Contains(script, wantQuoted) {
		t.Fatalf("launch script does not embed absolute binary path: got %q, want substring %q", script, wantQuoted)
	}
}

func TestEnsureMCPConfig_ReturnsOnNonNotExistReadError(t *testing.T) {
	home := setOpenCodeHome(t)
	configDir := configPath(home)
	if err := os.Mkdir(configDir, 0o755); err != nil {
		t.Fatalf("mkdir config path: %v", err)
	}

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err == nil {
		t.Fatal("expected read error for directory target, got nil")
	}
}

func TestEnsureMCPConfig_RejectsMalformedJSON(t *testing.T) {
	home := setOpenCodeHome(t)

	if err := os.WriteFile(
		configPath(home), []byte("{not json"), 0o644,
	); err != nil {
		t.Fatalf("seed: %v", err)
	}

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err == nil {
		t.Fatal("expected error on malformed JSON, got nil")
	}

	got, _ := os.ReadFile(configPath(home))
	if !bytes.Contains(got, []byte("{not json")) {
		t.Errorf("original malformed file overwritten: %s", got)
	}
}
