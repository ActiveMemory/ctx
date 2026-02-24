//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package permissions

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

// setupDir creates a temp dir with .claude/, chdirs into it, and returns cleanup.
func setupDir(t *testing.T) {
	t.Helper()
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	if err := os.MkdirAll(config.DirClaude, config.PermExec); err != nil {
		t.Fatal(err)
	}
}

// writeSettings writes JSON content to settings.local.json.
func writeSettings(t *testing.T, content string) {
	t.Helper()
	if err := os.WriteFile(config.FileSettings, []byte(content), config.PermFile); err != nil {
		t.Fatal(err)
	}
}

// writeGolden writes JSON content to settings.golden.json.
func writeGolden(t *testing.T, content string) {
	t.Helper()
	if err := os.WriteFile(config.FileSettingsGolden, []byte(content), config.PermFile); err != nil {
		t.Fatal(err)
	}
}

// runCmd executes a permissions subcommand and captures output.
func runCmd(args ...string) (string, error) {
	cmd := Cmd()
	cmd.SetArgs(args)
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	return buf.String(), err
}

func TestSnapshotCreatesGoldenFile(t *testing.T) {
	setupDir(t)

	settings := `{"permissions":{"allow":["Bash(ctx:*)"]}}`
	writeSettings(t, settings)

	out, err := runCmd("snapshot")
	if err != nil {
		t.Fatalf("snapshot error: %v", err)
	}
	if !strings.Contains(out, "Saved") {
		t.Errorf("output = %q, want 'Saved'", out)
	}

	// Verify golden file is a byte-for-byte copy.
	golden, err := os.ReadFile(config.FileSettingsGolden)
	if err != nil {
		t.Fatal(err)
	}
	if string(golden) != settings {
		t.Errorf("golden = %q, want %q", string(golden), settings)
	}
}

func TestSnapshotOverwritesExisting(t *testing.T) {
	setupDir(t)

	writeSettings(t, `{"permissions":{"allow":["A"]}}`)
	writeGolden(t, `{"permissions":{"allow":["OLD"]}}`)

	out, err := runCmd("snapshot")
	if err != nil {
		t.Fatalf("snapshot error: %v", err)
	}
	if !strings.Contains(out, "Updated") {
		t.Errorf("output = %q, want 'Updated'", out)
	}

	golden, err := os.ReadFile(config.FileSettingsGolden)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(golden), `"A"`) {
		t.Error("golden should contain new content")
	}
}

func TestSnapshotNoSettingsFile(t *testing.T) {
	setupDir(t)

	_, err := runCmd("snapshot")
	if err == nil {
		t.Fatal("expected error when settings.local.json is missing")
	}
	if !strings.Contains(err.Error(), "settings.local.json") {
		t.Errorf("error = %q, want mention of settings.local.json", err.Error())
	}
}

func TestRestoreFromGolden(t *testing.T) {
	setupDir(t)

	golden := `{"permissions":{"allow":["Bash(ctx:*)","Bash(go test:*)"]}}`
	local := `{"permissions":{"allow":["Bash(ctx:*)","Bash(go test:*)","Bash(rm -rf:*)"]}}`
	writeGolden(t, golden)
	writeSettings(t, local)

	out, err := runCmd("restore")
	if err != nil {
		t.Fatalf("restore error: %v", err)
	}
	if !strings.Contains(out, "Dropped") {
		t.Errorf("output should mention dropped permissions: %q", out)
	}
	if !strings.Contains(out, "Restored from golden") {
		t.Errorf("output should confirm restore: %q", out)
	}

	// Verify settings file now matches golden.
	data, err := os.ReadFile(config.FileSettings)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != golden {
		t.Errorf("settings = %q, want %q", string(data), golden)
	}
}

func TestRestoreShowsDiff(t *testing.T) {
	setupDir(t)

	goldenJSON := settingsJSON(t, []string{"A", "B", "C"})
	localJSON := settingsJSON(t, []string{"A", "C", "D", "E"})
	writeGolden(t, goldenJSON)
	writeSettings(t, localJSON)

	out, err := runCmd("restore")
	if err != nil {
		t.Fatalf("restore error: %v", err)
	}

	// Should drop D and E (in local but not golden).
	if !strings.Contains(out, "Dropped 2") {
		t.Errorf("expected 'Dropped 2', got: %q", out)
	}
	if !strings.Contains(out, "- D") {
		t.Errorf("expected '- D' in output: %q", out)
	}
	if !strings.Contains(out, "- E") {
		t.Errorf("expected '- E' in output: %q", out)
	}

	// Should restore B (in golden but not local).
	if !strings.Contains(out, "Restored 1") {
		t.Errorf("expected 'Restored 1', got: %q", out)
	}
	if !strings.Contains(out, "+ B") {
		t.Errorf("expected '+ B' in output: %q", out)
	}
}

func TestRestoreNoGoldenFile(t *testing.T) {
	setupDir(t)

	writeSettings(t, `{"permissions":{"allow":["A"]}}`)

	_, err := runCmd("restore")
	if err == nil {
		t.Fatal("expected error when golden file is missing")
	}
	if !strings.Contains(err.Error(), "settings.golden.json") {
		t.Errorf("error = %q, want mention of golden file", err.Error())
	}
}

func TestRestoreAlreadyClean(t *testing.T) {
	setupDir(t)

	content := `{"permissions":{"allow":["A","B"]}}`
	writeGolden(t, content)
	writeSettings(t, content)

	out, err := runCmd("restore")
	if err != nil {
		t.Fatalf("restore error: %v", err)
	}
	if !strings.Contains(out, "already match") {
		t.Errorf("output = %q, want 'already match'", out)
	}
}

func TestRestoreNoLocalFile(t *testing.T) {
	setupDir(t)

	golden := `{"permissions":{"allow":["Bash(ctx:*)"]}}`
	writeGolden(t, golden)

	out, err := runCmd("restore")
	if err != nil {
		t.Fatalf("restore error: %v", err)
	}
	if !strings.Contains(out, "no local settings") {
		t.Errorf("output = %q, want mention of no local settings", out)
	}

	// Verify settings file was created from golden.
	data, err := os.ReadFile(config.FileSettings)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != golden {
		t.Errorf("settings = %q, want %q", string(data), golden)
	}
}

func TestDiffStringSlices(t *testing.T) {
	tests := []struct {
		name         string
		golden       []string
		local        []string
		wantRestored []string
		wantDropped  []string
	}{
		{
			name:         "identical",
			golden:       []string{"A", "B", "C"},
			local:        []string{"A", "B", "C"},
			wantRestored: nil,
			wantDropped:  nil,
		},
		{
			name:         "local has extras",
			golden:       []string{"A", "B"},
			local:        []string{"A", "B", "C", "D"},
			wantRestored: nil,
			wantDropped:  []string{"C", "D"},
		},
		{
			name:         "golden has extras",
			golden:       []string{"A", "B", "C"},
			local:        []string{"A"},
			wantRestored: []string{"B", "C"},
			wantDropped:  nil,
		},
		{
			name:         "mixed",
			golden:       []string{"A", "B", "C"},
			local:        []string{"A", "C", "D", "E"},
			wantRestored: []string{"B"},
			wantDropped:  []string{"D", "E"},
		},
		{
			name:         "both empty",
			golden:       nil,
			local:        nil,
			wantRestored: nil,
			wantDropped:  nil,
		},
		{
			name:         "golden empty",
			golden:       nil,
			local:        []string{"A"},
			wantRestored: nil,
			wantDropped:  []string{"A"},
		},
		{
			name:         "local empty",
			golden:       []string{"A"},
			local:        nil,
			wantRestored: []string{"A"},
			wantDropped:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restored, dropped := diffStringSlices(tt.golden, tt.local)
			if !reflect.DeepEqual(restored, tt.wantRestored) {
				t.Errorf("restored = %v, want %v", restored, tt.wantRestored)
			}
			if !reflect.DeepEqual(dropped, tt.wantDropped) {
				t.Errorf("dropped = %v, want %v", dropped, tt.wantDropped)
			}
		})
	}
}

// settingsJSON builds a settings JSON string from a list of allow entries.
func settingsJSON(t *testing.T, allow []string) string {
	t.Helper()
	type perms struct {
		Allow []string `json:"allow"`
	}
	type settings struct {
		Permissions perms `json:"permissions"`
	}
	b, err := json.Marshal(settings{Permissions: perms{Allow: allow}})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// settingsJSONWithDeny builds a settings JSON string with both allow and deny.
func settingsJSONWithDeny(t *testing.T, allow, deny []string) string {
	t.Helper()
	type perms struct {
		Allow []string `json:"allow"`
		Deny  []string `json:"deny"`
	}
	type settings struct {
		Permissions perms `json:"permissions"`
	}
	b, err := json.Marshal(settings{Permissions: perms{Allow: allow, Deny: deny}})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestRestoreShowsDenyDiff(t *testing.T) {
	setupDir(t)

	goldenJSON := settingsJSONWithDeny(t, []string{"A"}, []string{"Bash(sudo *)", "Bash(git push *)"})
	localJSON := settingsJSONWithDeny(t, []string{"A"}, []string{"Bash(sudo *)", "Bash(rm -rf *)"})
	writeGolden(t, goldenJSON)
	writeSettings(t, localJSON)

	out, err := runCmd("restore")
	if err != nil {
		t.Fatalf("restore error: %v", err)
	}

	// Should drop "Bash(rm -rf *)" (in local deny but not golden deny)
	if !strings.Contains(out, "Dropped 1 session deny") {
		t.Errorf("expected 'Dropped 1 session deny', got: %q", out)
	}
	if !strings.Contains(out, "Bash(rm -rf *)") {
		t.Errorf("expected dropped deny rule in output: %q", out)
	}

	// Should restore "Bash(git push *)" (in golden deny but not local deny)
	if !strings.Contains(out, "Restored 1 deny") {
		t.Errorf("expected 'Restored 1 deny', got: %q", out)
	}
	if !strings.Contains(out, "Bash(git push *)") {
		t.Errorf("expected restored deny rule in output: %q", out)
	}
}

func TestSnapshotPreservesExactBytes(t *testing.T) {
	setupDir(t)

	// Settings with intentional formatting (trailing newline, specific indent).
	content := "{\n  \"permissions\": {\n    \"allow\": [\n      \"Bash(ctx:*)\"\n    ]\n  }\n}\n"
	writeSettings(t, content)

	if _, err := runCmd("snapshot"); err != nil {
		t.Fatal(err)
	}

	golden, err := os.ReadFile(config.FileSettingsGolden)
	if err != nil {
		t.Fatal(err)
	}
	if string(golden) != content {
		t.Error("golden should be byte-for-byte identical to source")
	}
}

func TestRestorePreservesExactBytes(t *testing.T) {
	setupDir(t)

	// Golden with specific formatting.
	goldenContent := "{\n  \"permissions\": {\n    \"allow\": [\n      \"A\"\n    ]\n  }\n}\n"
	writeGolden(t, goldenContent)
	writeSettings(t, `{"permissions":{"allow":["A","B"]}}`)

	if _, err := runCmd("restore"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(config.FileSettings)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != goldenContent {
		t.Error("restored settings should be byte-for-byte identical to golden")
	}
}

func TestCmdHasSubcommands(t *testing.T) {
	cmd := Cmd()
	if cmd.Use != "permissions" {
		t.Errorf("cmd.Use = %q, want 'permissions'", cmd.Use)
	}

	names := make(map[string]bool)
	for _, sub := range cmd.Commands() {
		names[sub.Name()] = true
	}
	if !names["snapshot"] {
		t.Error("missing snapshot subcommand")
	}
	if !names["restore"] {
		t.Error("missing restore subcommand")
	}
}

// Verify golden file path is under .claude/ (not .context/).
func TestGoldenFilePath(t *testing.T) {
	if !strings.HasPrefix(config.FileSettingsGolden, config.DirClaude+"/") {
		t.Errorf("FileSettingsGolden = %q, want prefix %q", config.FileSettingsGolden, config.DirClaude+"/")
	}
}
