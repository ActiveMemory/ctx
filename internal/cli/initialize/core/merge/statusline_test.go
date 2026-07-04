//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// setupProject creates a temp project root with .context/, chdirs
// into it, and resets the rc cache. Cleanup restores everything.
func setupProject(t *testing.T) string {
	t.Helper()
	tmpDir, mkErr := os.MkdirTemp("", "ctx-merge-statusline-*")
	if mkErr != nil {
		t.Fatalf("failed to create temp dir: %v", mkErr)
	}
	// Resolve symlinks (macOS /var vs /private/var) so path
	// comparisons hold.
	resolved, resolveErr := filepath.EvalSymlinks(tmpDir)
	if resolveErr == nil {
		tmpDir = resolved
	}
	if mkdirErr := os.Mkdir(filepath.Join(tmpDir, ".context"), 0o750); mkdirErr != nil {
		t.Fatalf("failed to create .context: %v", mkdirErr)
	}
	origDir, _ := os.Getwd()
	if chdirErr := os.Chdir(tmpDir); chdirErr != nil {
		t.Fatalf("failed to chdir: %v", chdirErr)
	}
	rc.Reset()
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
		_ = os.RemoveAll(tmpDir)
	})
	return tmpDir
}

// runStatusLine invokes SettingsStatusLine with a throwaway command.
func runStatusLine(t *testing.T) {
	t.Helper()
	cmd := &cobra.Command{}
	cmd.SetOut(os.NewFile(0, os.DevNull))
	if mergeErr := SettingsStatusLine(cmd); mergeErr != nil {
		t.Fatalf("SettingsStatusLine failed: %v", mergeErr)
	}
}

// readSettings parses settings.local.json into a raw key map.
func readSettings(t *testing.T) map[string]json.RawMessage {
	t.Helper()
	content, readErr := os.ReadFile(cfgClaude.Settings)
	if readErr != nil {
		t.Fatalf("failed to read settings: %v", readErr)
	}
	raw := map[string]json.RawMessage{}
	if unmarshalErr := json.Unmarshal(content, &raw); unmarshalErr != nil {
		t.Fatalf("settings not valid JSON: %v", unmarshalErr)
	}
	return raw
}

func TestStatusLineDeployedIntoFreshProject(t *testing.T) {
	setupProject(t)
	runStatusLine(t)

	raw := readSettings(t)
	entry, exists := raw[cfgClaude.FieldStatusLine]
	if !exists {
		t.Fatal("statusLine key missing after deploy")
	}
	if !strings.Contains(string(entry), cfgClaude.StatusLineCommand) {
		t.Errorf("statusLine %s does not run %q",
			entry, cfgClaude.StatusLineCommand)
	}
}

func TestStatusLineBacksUpForeignEntry(t *testing.T) {
	tmpDir := setupProject(t)
	foreign := `{"type":"command","command":"my-custom-line.sh"}`
	if mkdirErr := os.MkdirAll(".claude", 0o750); mkdirErr != nil {
		t.Fatalf("failed to create .claude: %v", mkdirErr)
	}
	settings := `{"statusLine": ` + foreign + `, "env": {"FOO": "bar"}}`
	if writeErr := os.WriteFile(cfgClaude.Settings, []byte(settings), 0o600); writeErr != nil {
		t.Fatalf("failed to seed settings: %v", writeErr)
	}

	runStatusLine(t)

	raw := readSettings(t)
	if !strings.Contains(string(raw[cfgClaude.FieldStatusLine]),
		cfgClaude.StatusLineCommand) {
		t.Errorf("statusLine not replaced: %s", raw[cfgClaude.FieldStatusLine])
	}
	if _, exists := raw["env"]; !exists {
		t.Error("unmodeled env key dropped by statusline merge")
	}

	backupPath := filepath.Join(
		tmpDir, ".context", dir.State, cfgClaude.PreviousStatusLine,
	)
	backup, readErr := os.ReadFile(filepath.Clean(backupPath))
	if readErr != nil {
		t.Fatalf("backup not written: %v", readErr)
	}
	if !strings.Contains(string(backup), "my-custom-line.sh") {
		t.Errorf("backup %s lost the foreign command", backup)
	}
}

func TestStatusLineIdempotent(t *testing.T) {
	setupProject(t)
	runStatusLine(t)
	before := readSettings(t)
	runStatusLine(t)
	after := readSettings(t)
	if string(before[cfgClaude.FieldStatusLine]) !=
		string(after[cfgClaude.FieldStatusLine]) {
		t.Error("second deploy changed the statusLine entry")
	}
}

func TestStatusLineDisableRestoresBackup(t *testing.T) {
	setupProject(t)
	foreign := `{"type":"command","command":"my-custom-line.sh"}`
	if mkdirErr := os.MkdirAll(".claude", 0o750); mkdirErr != nil {
		t.Fatalf("failed to create .claude: %v", mkdirErr)
	}
	settings := `{"statusLine": ` + foreign + `}`
	if writeErr := os.WriteFile(cfgClaude.Settings, []byte(settings), 0o600); writeErr != nil {
		t.Fatalf("failed to seed settings: %v", writeErr)
	}
	runStatusLine(t) // deploys ours, backs up foreign

	rcBody := "statusline:\n  enabled: false\n"
	if writeErr := os.WriteFile(".ctxrc", []byte(rcBody), 0o600); writeErr != nil {
		t.Fatalf("failed to write .ctxrc: %v", writeErr)
	}
	rc.Reset()
	runStatusLine(t) // disable path

	raw := readSettings(t)
	if !strings.Contains(string(raw[cfgClaude.FieldStatusLine]),
		"my-custom-line.sh") {
		t.Errorf("foreign statusLine not restored: %s",
			raw[cfgClaude.FieldStatusLine])
	}
}

func TestStatusLineDisableRemovesWithoutBackup(t *testing.T) {
	setupProject(t)
	runStatusLine(t) // deploys ours; nothing to back up

	rcBody := "statusline:\n  enabled: false\n"
	if writeErr := os.WriteFile(".ctxrc", []byte(rcBody), 0o600); writeErr != nil {
		t.Fatalf("failed to write .ctxrc: %v", writeErr)
	}
	rc.Reset()
	runStatusLine(t)

	raw := readSettings(t)
	if _, exists := raw[cfgClaude.FieldStatusLine]; exists {
		t.Errorf("statusLine survived disable with no backup: %s",
			raw[cfgClaude.FieldStatusLine])
	}
}

func TestStatusLineDisableLeavesForeignAlone(t *testing.T) {
	setupProject(t)
	foreign := `{"type":"command","command":"my-custom-line.sh"}`
	if mkdirErr := os.MkdirAll(".claude", 0o750); mkdirErr != nil {
		t.Fatalf("failed to create .claude: %v", mkdirErr)
	}
	settings := `{"statusLine": ` + foreign + `}`
	if writeErr := os.WriteFile(cfgClaude.Settings, []byte(settings), 0o600); writeErr != nil {
		t.Fatalf("failed to seed settings: %v", writeErr)
	}

	rcBody := "statusline:\n  enabled: false\n"
	if writeErr := os.WriteFile(".ctxrc", []byte(rcBody), 0o600); writeErr != nil {
		t.Fatalf("failed to write .ctxrc: %v", writeErr)
	}
	rc.Reset()
	runStatusLine(t)

	raw := readSettings(t)
	if !strings.Contains(string(raw[cfgClaude.FieldStatusLine]),
		"my-custom-line.sh") {
		t.Errorf("foreign statusLine modified on disable: %s",
			raw[cfgClaude.FieldStatusLine])
	}
}

func TestPermissionsMergePreservesUnknownKeys(t *testing.T) {
	setupProject(t)
	if mkdirErr := os.MkdirAll(".claude", 0o750); mkdirErr != nil {
		t.Fatalf("failed to create .claude: %v", mkdirErr)
	}
	settings := `{
		"env": {"FOO": "bar"},
		"statusLine": {"type": "command", "command": "my-custom-line.sh"},
		"permissions": {"allow": ["Bash(ls:*)"]}
	}`
	if writeErr := os.WriteFile(cfgClaude.Settings, []byte(settings), 0o600); writeErr != nil {
		t.Fatalf("failed to seed settings: %v", writeErr)
	}

	cmd := &cobra.Command{}
	cmd.SetOut(os.NewFile(0, os.DevNull))
	if mergeErr := SettingsPermissions(cmd); mergeErr != nil {
		t.Fatalf("SettingsPermissions failed: %v", mergeErr)
	}

	raw := readSettings(t)
	if _, exists := raw["env"]; !exists {
		t.Error("permissions merge dropped the env key")
	}
	if !strings.Contains(string(raw[cfgClaude.FieldStatusLine]),
		"my-custom-line.sh") {
		t.Error("permissions merge dropped or rewrote statusLine")
	}
	if !strings.Contains(string(raw[cfgClaude.FieldPermissions]),
		"Bash(ls:*)") {
		t.Error("permissions merge lost the user's existing allow entry")
	}
}
