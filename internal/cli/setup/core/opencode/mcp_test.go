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
	if ctxServer["command"] != "ctx" {
		t.Errorf("command = %q, want ctx", ctxServer["command"])
	}
	if ctxServer["type"] != "local" {
		t.Errorf("type = %q, want local", ctxServer["type"])
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
