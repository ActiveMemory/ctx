//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestBootstrap_TextOutput(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	if err := runBootstrap(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)

	if !strings.Contains(out, "context_dir:") {
		t.Errorf("expected output to contain 'context_dir:', got: %s", out)
	}
	if !strings.Contains(out, "ctx bootstrap") {
		t.Errorf("expected output to contain 'ctx bootstrap', got: %s", out)
	}
	if !strings.Contains(out, "CONSTITUTION.md") {
		t.Errorf("expected output to contain 'CONSTITUTION.md', got: %s", out)
	}
	if !strings.Contains(out, "TASKS.md") {
		t.Errorf("expected output to contain 'TASKS.md', got: %s", out)
	}
	if !strings.Contains(out, "DECISIONS.md") {
		t.Errorf("expected output to contain 'DECISIONS.md', got: %s", out)
	}
	if !strings.Contains(out, "Rules:") {
		t.Errorf("expected output to contain 'Rules:', got: %s", out)
	}
}

func TestBootstrap_JSONOutput(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()
	setupContextDir(t)

	cmd := newTestCmd()
	cmd.Flags().Bool("json", true, "")
	_ = cmd.Flags().Set("json", "true")

	if err := runBootstrap(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)

	var result struct {
		ContextDir string   `json:"context_dir"`
		Files      []string `json:"files"`
		Rules      []string `json:"rules"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("invalid JSON output: %v\noutput: %s", err, out)
	}

	if result.ContextDir == "" {
		t.Error("expected context_dir to be non-empty")
	}
	if len(result.Files) == 0 {
		t.Error("expected files to be non-empty")
	}
	if len(result.Rules) == 0 {
		t.Error("expected rules to be non-empty")
	}

	// Verify known files are present
	fileSet := make(map[string]bool)
	for _, f := range result.Files {
		fileSet[f] = true
	}
	for _, expected := range []string{"CONSTITUTION.md", "TASKS.md", "DECISIONS.md"} {
		if !fileSet[expected] {
			t.Errorf("expected files to contain %q, got: %v", expected, result.Files)
		}
	}
}

func TestBootstrap_CustomDir(t *testing.T) {
	origDir, _ := os.Getwd()
	workDir := t.TempDir()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()

	customDir := "my-custom-context"
	rc.OverrideContextDir(customDir)
	defer rc.Reset()

	// Create the custom context dir with a file
	if err := os.MkdirAll(customDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(customDir+"/CONSTITUTION.md", []byte("# test\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cmd := newTestCmd()
	if err := runBootstrap(cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := cmdOutput(cmd)
	if !strings.Contains(out, customDir) {
		t.Errorf("expected output to contain custom dir %q, got: %s", customDir, out)
	}
}

func TestBootstrap_MissingDir(t *testing.T) {
	origDir, _ := os.Getwd()
	_ = os.Chdir(t.TempDir())
	defer func() { _ = os.Chdir(origDir) }()

	rc.OverrideContextDir("nonexistent-dir")
	defer rc.Reset()

	cmd := newTestCmd()
	err := runBootstrap(cmd)
	if err == nil {
		t.Fatal("expected error for missing directory, got nil")
	}
	if !strings.Contains(err.Error(), "context directory not found") {
		t.Errorf("expected error to contain 'context directory not found', got: %v", err)
	}
	if !strings.Contains(err.Error(), "ctx init") {
		t.Errorf("expected error to mention 'ctx init', got: %v", err)
	}
}
