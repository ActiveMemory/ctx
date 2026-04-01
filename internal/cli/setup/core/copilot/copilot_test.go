//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package copilot

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// testCmd returns a cobra.Command that captures output.
func testCmd(buf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

func TestEnsureVSCodeMCP_CreatesFile(t *testing.T) {
	tmp := t.TempDir()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	if err := ensureVSCodeMCP(cmd); err != nil {
		t.Fatalf("ensureVSCodeMCP() error = %v", err)
	}

	target := filepath.Join(".vscode", "mcp.json")
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read mcp.json: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("mcp.json is not valid JSON: %v", err)
	}

	servers, ok := parsed["servers"].(map[string]interface{})
	if !ok {
		t.Fatal("mcp.json missing 'servers' key")
	}

	ctxServer, ok := servers["ctx"].(map[string]interface{})
	if !ok {
		t.Fatal("mcp.json missing 'servers.ctx' key")
	}

	if ctxServer["command"] != "ctx" {
		t.Errorf("expected command 'ctx', got %q", ctxServer["command"])
	}

	args, ok := ctxServer["args"].([]interface{})
	if !ok || len(args) != 2 || args[0] != "mcp" || args[1] != "serve" {
		t.Errorf("expected args [mcp, serve], got %v", ctxServer["args"])
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("\u2713")) {
		t.Errorf("expected success marker in output, got %q", output)
	}
}

func TestEnsureVSCodeMCP_SkipsExisting(t *testing.T) {
	tmp := t.TempDir()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	vsDir := ".vscode"
	target := filepath.Join(vsDir, "mcp.json")
	if err := os.MkdirAll(vsDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	existing := []byte(`{"servers":{"custom":{"command":"other"}}}`)
	if err := os.WriteFile(target, existing, 0o644); err != nil {
		t.Fatalf("write existing: %v", err)
	}

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	if err := ensureVSCodeMCP(cmd); err != nil {
		t.Fatalf("ensureVSCodeMCP() error = %v", err)
	}

	// File should not be overwritten
	data, _ := os.ReadFile(target)
	if string(data) != string(existing) {
		t.Error("ensureVSCodeMCP overwrote existing file")
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("skipped")) {
		t.Errorf("expected 'skipped' in output, got %q", output)
	}
}

func TestEnsureVSCodeMCP_CreatesVSCodeDir(t *testing.T) {
	tmp := t.TempDir()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// Ensure .vscode/ does NOT exist beforehand
	vsDir := filepath.Join(tmp, ".vscode")
	if _, err := os.Stat(vsDir); err == nil {
		t.Fatal(".vscode should not exist yet")
	}

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	if err := ensureVSCodeMCP(cmd); err != nil {
		t.Fatalf("ensureVSCodeMCP() error = %v", err)
	}

	// .vscode/ should now exist
	info, err := os.Stat(".vscode")
	if err != nil {
		t.Fatalf(".vscode dir was not created: %v", err)
	}
	if !info.IsDir() {
		t.Error(".vscode should be a directory")
	}
}
