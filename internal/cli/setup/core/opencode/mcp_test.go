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
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func testCmd(buf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

func chdirTemp(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	orig, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}

func readMCP(t *testing.T) map[string]interface{} {
	t.Helper()
	raw, err := os.ReadFile("opencode.json")
	if err != nil {
		t.Fatalf("read opencode.json: %v", err)
	}
	parsed := map[string]interface{}{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("opencode.json not valid JSON: %v", err)
	}
	return parsed
}

func TestEnsureMCPConfig_CreatesFile(t *testing.T) {
	chdirTemp(t)

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig: %v", err)
	}

	parsed := readMCP(t)
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
	// We wrap the launch in `sh -c` so $PWD can be substituted into
	// CTX_DIR at MCP spawn time. ctx rejects relative CTX_DIR values
	// (see internal/rc.ContextDir), and OpenCode has no path templating
	// in opencode.json — the shell wrapper is how we get an absolute
	// path that follows the user's checkout.
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
		`exec env`,                // replace shell with ctx
		`CTX_DIR="$PWD/.context"`, // absolute path resolved at spawn
		`ctx mcp serve`,           // the actual MCP server invocation
	}
	for _, s := range wantSubs {
		if !strings.Contains(script, s) {
			t.Errorf("launch script missing %q\nfull script: %s", s, script)
		}
	}
	if _, hasArgs := ctxServer["args"]; hasArgs {
		t.Error("args field must not be set; OpenCode schema folds args into command array")
	}
	if _, hasEnv := ctxServer["environment"]; hasEnv {
		t.Error("environment field must not be set; CTX_DIR is computed from $PWD inside the sh wrapper. A literal CTX_DIR='.context' here would be rejected by ctx as non-absolute.")
	}
	enabled, ok := ctxServer["enabled"].(bool)
	if !ok || !enabled {
		t.Errorf("enabled = %v, want true", ctxServer["enabled"])
	}
}

func TestEnsureMCPConfig_TreatsEmptyFileAsAbsent(t *testing.T) {
	chdirTemp(t)

	if err := os.WriteFile(
		"opencode.json", []byte("   \n\t  "), 0o644,
	); err != nil {
		t.Fatalf("seed empty file: %v", err)
	}

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig on empty file: %v", err)
	}

	parsed := readMCP(t)
	if _, ok := parsed["mcp"].(map[string]interface{}); !ok {
		t.Fatal("mcp key not registered after empty-file path")
	}
}

func TestEnsureMCPConfig_PreservesExistingKeys(t *testing.T) {
	chdirTemp(t)

	seed := []byte(`{"theme":"dark","mcp":{"other":{"type":"local"}}}`)
	if err := os.WriteFile("opencode.json", seed, 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig: %v", err)
	}

	parsed := readMCP(t)
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

func TestEnsureMCPConfig_SkipsWhenCtxAlreadyRegistered(t *testing.T) {
	chdirTemp(t)

	seed := []byte(`{"mcp":{"ctx":{"command":"custom"}}}`)
	if err := os.WriteFile("opencode.json", seed, 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err != nil {
		t.Fatalf("ensureMCPConfig: %v", err)
	}

	got, _ := os.ReadFile("opencode.json")
	if string(got) != string(seed) {
		t.Errorf(
			"file rewritten when ctx already registered: %s", got,
		)
	}
	if !bytes.Contains(buf.Bytes(), []byte("skipped")) {
		t.Errorf(
			"expected 'skipped' in output, got %q", buf.String(),
		)
	}
}

func TestEnsureMCPConfig_RejectsMalformedJSON(t *testing.T) {
	chdirTemp(t)

	if err := os.WriteFile(
		"opencode.json", []byte("{not json"), 0o644,
	); err != nil {
		t.Fatalf("seed: %v", err)
	}

	var buf bytes.Buffer
	if err := ensureMCPConfig(testCmd(&buf)); err == nil {
		t.Fatal("expected error on malformed JSON, got nil")
	}

	// Verify we did not clobber the user's broken-but-extant file.
	got, _ := os.ReadFile("opencode.json")
	if !bytes.Contains(got, []byte("{not json")) {
		t.Errorf("original malformed file overwritten: %s", got)
	}
}
