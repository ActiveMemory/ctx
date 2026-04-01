//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package vscode

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	cfgVscode "github.com/ActiveMemory/ctx/internal/config/vscode"
)

// testCmd returns a cobra.Command that captures output.
func testCmd(buf *bytes.Buffer) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd
}

func TestWriteMCPJSON_CreatesFile(t *testing.T) {
	tmp := t.TempDir()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.MkdirAll(cfgVscode.Dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	if err := writeMCPJSON(cmd); err != nil {
		t.Fatalf("writeMCPJSON() error = %v", err)
	}

	target := filepath.Join(cfgVscode.Dir, "mcp.json")
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

	output := buf.String()
	if len(output) == 0 {
		t.Error("expected output message for created file")
	}
}

func TestWriteMCPJSON_SkipsExisting(t *testing.T) {
	tmp := t.TempDir()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	target := filepath.Join(cfgVscode.Dir, "mcp.json")
	if err := os.MkdirAll(cfgVscode.Dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	existing := []byte(`{"servers":{"other":{}}}`)
	if err := os.WriteFile(target, existing, 0o644); err != nil {
		t.Fatalf("write existing: %v", err)
	}

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	if err := writeMCPJSON(cmd); err != nil {
		t.Fatalf("writeMCPJSON() error = %v", err)
	}

	// File should not be overwritten
	data, _ := os.ReadFile(target)
	if string(data) != string(existing) {
		t.Error("writeMCPJSON overwrote existing file")
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("exists")) {
		t.Errorf("expected 'exists' in output, got %q", output)
	}
}

func TestCreateVSCodeArtifacts_CreatesMCPJSON(t *testing.T) {
	tmp := t.TempDir()

	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	var buf bytes.Buffer
	cmd := testCmd(&buf)

	if err := CreateVSCodeArtifacts(cmd); err != nil {
		t.Fatalf("CreateVSCodeArtifacts() error = %v", err)
	}

	// Verify mcp.json was created as part of the artifacts
	target := filepath.Join(cfgVscode.Dir, "mcp.json")
	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Error("CreateVSCodeArtifacts did not create mcp.json")
	}

	// Verify extensions.json was also created
	extTarget := filepath.Join(cfgVscode.Dir, "extensions.json")
	if _, err := os.Stat(extTarget); os.IsNotExist(err) {
		t.Error("CreateVSCodeArtifacts did not create extensions.json")
	}

	// Verify tasks.json was also created
	taskTarget := filepath.Join(cfgVscode.Dir, "tasks.json")
	if _, err := os.Stat(taskTarget); os.IsNotExist(err) {
		t.Error("CreateVSCodeArtifacts did not create tasks.json")
	}
}
