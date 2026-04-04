//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	cfgDrift "github.com/ActiveMemory/ctx/internal/config/drift"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/steering"
)

// **Validates: Requirements 19.7**

func TestCheckSteeringTools(t *testing.T) {
	tests := []struct {
		name         string
		files        map[string]string // steering file name → content
		wantWarnings int
		wantPassed   bool
	}{
		{
			name:         "no steering directory",
			files:        nil,
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name: "valid tool identifiers",
			files: map[string]string{
				"api.md": "---\nname: api\ntools: [claude, cursor]\n---\nBody\n",
			},
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name: "empty tools list (all tools)",
			files: map[string]string{
				"api.md": "---\nname: api\n---\nBody\n",
			},
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name: "invalid tool identifier",
			files: map[string]string{
				"api.md": "---\nname: api\ntools: [claude, vscode]\n---\nBody\n",
			},
			wantWarnings: 1,
			wantPassed:   false,
		},
		{
			name: "multiple invalid tools in one file",
			files: map[string]string{
				"api.md": "---\nname: api\ntools: [vscode, neovim]\n---\nBody\n",
			},
			wantWarnings: 2,
			wantPassed:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			steeringDir := filepath.Join(tmpDir, ".context", "steering")
			if tt.files != nil {
				if err := os.MkdirAll(steeringDir, 0o755); err != nil {
					t.Fatal(err)
				}
				for name, content := range tt.files {
					if err := os.WriteFile(
						filepath.Join(steeringDir, name),
						[]byte(content), 0o644,
					); err != nil {
						t.Fatal(err)
					}
				}
			}

			// chdir so rc reads .ctxrc from tmpDir
			origDir := chdir(t, tmpDir)
			defer func() { _ = os.Chdir(origDir) }()

			writeCtxRC(t, tmpDir, fmt.Sprintf("steering:\n  dir: %s\n", steeringDir))
			rc.Reset()
			defer rc.Reset()

			report := &Report{
				Warnings:   []Issue{},
				Violations: []Issue{},
				Passed:     []cfgDrift.CheckName{},
			}

			checkSteeringTools(report)

			if len(report.Warnings) != tt.wantWarnings {
				t.Errorf("expected %d warnings, got %d", tt.wantWarnings, len(report.Warnings))
				for _, w := range report.Warnings {
					t.Logf("  warning: %s", w.Message)
				}
			}

			for _, w := range report.Warnings {
				if w.Type != cfgDrift.IssueInvalidTool {
					t.Errorf("expected issue type %q, got %q", cfgDrift.IssueInvalidTool, w.Type)
				}
			}

			passed := checkPassed(report, cfgDrift.CheckSteeringTools)
			if passed != tt.wantPassed {
				t.Errorf("expected passed=%v, got passed=%v", tt.wantPassed, passed)
			}
		})
	}
}

func TestCheckHookPerms(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T, hooksDir string)
		wantWarnings int
		wantPassed   bool
	}{
		{
			name:         "no hooks directory",
			setup:        func(_ *testing.T, _ string) {},
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name: "all hooks executable",
			setup: func(t *testing.T, hooksDir string) {
				t.Helper()
				dir := filepath.Join(hooksDir, "pre-tool-use")
				mustMkdir(t, dir)
				mustWriteFile(t, filepath.Join(dir, "check.sh"),
					"#!/bin/bash\necho ok", 0o755)
			},
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name: "hook missing executable bit",
			setup: func(t *testing.T, hooksDir string) {
				t.Helper()
				dir := filepath.Join(hooksDir, "session-start")
				mustMkdir(t, dir)
				mustWriteFile(t, filepath.Join(dir, "init.sh"),
					"#!/bin/bash\necho ok", 0o644)
			},
			wantWarnings: 1,
			wantPassed:   false,
		},
		{
			name: "mixed executable and non-executable",
			setup: func(t *testing.T, hooksDir string) {
				t.Helper()
				dir := filepath.Join(hooksDir, "post-tool-use")
				mustMkdir(t, dir)
				mustWriteFile(t, filepath.Join(dir, "lint.sh"),
					"#!/bin/bash\necho ok", 0o755)
				mustWriteFile(t, filepath.Join(dir, "broken.sh"),
					"#!/bin/bash\necho ok", 0o644)
			},
			wantWarnings: 1,
			wantPassed:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			hooksDir := filepath.Join(tmpDir, ".context", "hooks")
			tt.setup(t, hooksDir)

			origDir := chdir(t, tmpDir)
			defer func() { _ = os.Chdir(origDir) }()

			writeCtxRC(t, tmpDir, fmt.Sprintf("hooks:\n  dir: %s\n", hooksDir))
			rc.Reset()
			defer rc.Reset()

			report := &Report{
				Warnings:   []Issue{},
				Violations: []Issue{},
				Passed:     []cfgDrift.CheckName{},
			}

			checkHookPerms(report)

			if len(report.Warnings) != tt.wantWarnings {
				t.Errorf("expected %d warnings, got %d", tt.wantWarnings, len(report.Warnings))
				for _, w := range report.Warnings {
					t.Logf("  warning: %s (file=%s)", w.Message, w.File)
				}
			}

			for _, w := range report.Warnings {
				if w.Type != cfgDrift.IssueHookNoExec {
					t.Errorf("expected issue type %q, got %q", cfgDrift.IssueHookNoExec, w.Type)
				}
			}

			passed := checkPassed(report, cfgDrift.CheckHookPerms)
			if passed != tt.wantPassed {
				t.Errorf("expected passed=%v, got passed=%v", tt.wantPassed, passed)
			}
		})
	}
}

func TestCheckSyncStaleness(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T, tmpDir, steeringDir string)
		wantWarnings int
		wantPassed   bool
	}{
		{
			name:         "no steering files",
			setup:        func(_ *testing.T, _, _ string) {},
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name: "synced files up to date",
			setup: func(t *testing.T, tmpDir, steeringDir string) {
				t.Helper()
				// Create a steering file.
				mustMkdir(t, steeringDir)
				mustWriteFile(t, filepath.Join(steeringDir, "api.md"),
					"---\nname: api\ndescription: API rules\ninclusion: always\npriority: 50\n---\nAPI body\n", 0o644)

				// Sync to all tools so all native files are up to date.
				_, err := steering.SyncAll(steeringDir, tmpDir)
				if err != nil {
					t.Fatal(err)
				}
			},
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name: "synced file is stale",
			setup: func(t *testing.T, tmpDir, steeringDir string) {
				t.Helper()
				mustMkdir(t, steeringDir)
				mustWriteFile(t, filepath.Join(steeringDir, "api.md"),
					"---\nname: api\ndescription: API rules\ninclusion: always\npriority: 50\n---\nAPI body\n", 0o644)

				// Sync all tools first.
				_, err := steering.SyncAll(steeringDir, tmpDir)
				if err != nil {
					t.Fatal(err)
				}

				// Now modify the source steering file — all synced files become stale.
				mustWriteFile(t, filepath.Join(steeringDir, "api.md"),
					"---\nname: api\ndescription: Updated API rules\ninclusion: always\npriority: 50\n---\nUpdated body\n", 0o644)
			},
			// All 3 syncable tools (cursor, cline, kiro) will report stale.
			wantWarnings: 3,
			wantPassed:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			steeringDir := filepath.Join(tmpDir, ".context", "steering")
			tt.setup(t, tmpDir, steeringDir)

			origDir := chdir(t, tmpDir)
			defer func() { _ = os.Chdir(origDir) }()

			writeCtxRC(t, tmpDir, fmt.Sprintf("steering:\n  dir: %s\n", steeringDir))
			rc.Reset()
			defer rc.Reset()

			report := &Report{
				Warnings:   []Issue{},
				Violations: []Issue{},
				Passed:     []cfgDrift.CheckName{},
			}

			checkSyncStaleness(report)

			if len(report.Warnings) != tt.wantWarnings {
				t.Errorf("expected %d warnings, got %d", tt.wantWarnings, len(report.Warnings))
				for _, w := range report.Warnings {
					t.Logf("  warning: %s (file=%s path=%s)", w.Message, w.File, w.Path)
				}
			}

			for _, w := range report.Warnings {
				if w.Type != cfgDrift.IssueStaleSyncFile {
					t.Errorf("expected issue type %q, got %q", cfgDrift.IssueStaleSyncFile, w.Type)
				}
			}

			passed := checkPassed(report, cfgDrift.CheckSyncStaleness)
			if passed != tt.wantPassed {
				t.Errorf("expected passed=%v, got passed=%v", tt.wantPassed, passed)
			}
		})
	}
}

func TestCheckRCTool(t *testing.T) {
	tests := []struct {
		name         string
		rcContent    string
		wantWarnings int
		wantPassed   bool
	}{
		{
			name:         "no tool configured",
			rcContent:    "",
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name:         "valid tool",
			rcContent:    "tool: kiro\n",
			wantWarnings: 0,
			wantPassed:   true,
		},
		{
			name:         "invalid tool",
			rcContent:    "tool: vscode\n",
			wantWarnings: 1,
			wantPassed:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			origDir := chdir(t, tmpDir)
			defer func() { _ = os.Chdir(origDir) }()

			writeCtxRC(t, tmpDir, tt.rcContent)
			rc.Reset()
			defer rc.Reset()

			report := &Report{
				Warnings:   []Issue{},
				Violations: []Issue{},
				Passed:     []cfgDrift.CheckName{},
			}

			checkRCTool(report)

			if len(report.Warnings) != tt.wantWarnings {
				t.Errorf("expected %d warnings, got %d", tt.wantWarnings, len(report.Warnings))
				for _, w := range report.Warnings {
					t.Logf("  warning: %s", w.Message)
				}
			}

			for _, w := range report.Warnings {
				if w.Type != cfgDrift.IssueInvalidTool {
					t.Errorf("expected issue type %q, got %q", cfgDrift.IssueInvalidTool, w.Type)
				}
			}

			passed := checkPassed(report, cfgDrift.CheckRCTool)
			if passed != tt.wantPassed {
				t.Errorf("expected passed=%v, got passed=%v", tt.wantPassed, passed)
			}
		})
	}
}

// --- helpers ---

func checkPassed(report *Report, check cfgDrift.CheckName) bool {
	for _, p := range report.Passed {
		if p == check {
			return true
		}
	}
	return false
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func mustWriteFile(t *testing.T, path, content string, perm os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		t.Fatal(err)
	}
}

func writeCtxRC(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, ".ctxrc"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func chdir(t *testing.T, dir string) string {
	t.Helper()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	return origDir
}
