//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestEnsureProviderConfig_WritesProviderBaseURL(t *testing.T) {
	home := t.TempDir()
	t.Setenv("OPENCODE_HOME", home)
	if err := EnsureProviderConfig(testProviderCmd(), "openai", "https://example.com/v1"); err != nil {
		t.Fatalf("EnsureProviderConfig() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(home, "opencode.json"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), `"baseURL": "https://example.com/v1"`) {
		t.Fatalf("opencode.json = %s", string(data))
	}
}

func TestEnsureProviderConfig_PreservesExistingKeys(t *testing.T) {
	home := t.TempDir()
	t.Setenv("OPENCODE_HOME", home)
	path := filepath.Join(home, "opencode.json")
	if err := os.WriteFile(path, []byte(`{"theme":"dark"}`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := EnsureProviderConfig(testProviderCmd(), "openai-compatible", "https://proxy.example"); err != nil {
		t.Fatalf("EnsureProviderConfig() error = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(data), `"theme": "dark"`) || !strings.Contains(string(data), `"provider"`) {
		t.Fatalf("opencode.json = %s", string(data))
	}
}

func testProviderCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.SetOut(&bytes.Buffer{})
	return cmd
}
