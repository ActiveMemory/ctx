//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeSteering creates a steering file in dir with the given content.
func writeSteering(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const steeringAlways = `---
name: api-rules
description: REST API conventions
inclusion: always
priority: 10
---
Use RESTful conventions.
`

const steeringCursorOnly = `---
name: cursor-only
description: Cursor-specific rules
inclusion: auto
tools: [cursor]
priority: 50
---
Cursor body content.
`

const steeringManual = `---
name: manual-rule
description: Manual rule
inclusion: manual
priority: 50
---
Manual body.
`

func TestSyncTool_CursorFormat(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")
	writeSteering(t, steeringDir, "api-rules", steeringAlways)

	report, err := SyncTool(steeringDir, root, "cursor")
	if err != nil {
		t.Fatalf("SyncTool: %v", err)
	}
	if len(report.Written) != 1 || report.Written[0] != "api-rules" {
		t.Errorf("expected 1 written file, got %v", report.Written)
	}

	out := filepath.Join(root, ".cursor", "rules", "api-rules.mdc")
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	content := string(data)

	// Verify Cursor frontmatter.
	if !strings.Contains(content, "alwaysApply: true") {
		t.Error("cursor output should have alwaysApply: true for always inclusion")
	}
	if !strings.Contains(content, "description: REST API conventions") {
		t.Error("cursor output should contain description")
	}
	if !strings.Contains(content, "Use RESTful conventions.") {
		t.Error("cursor output should contain body")
	}
}

func TestSyncTool_ClineFormat(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")
	writeSteering(t, steeringDir, "api-rules", steeringAlways)

	report, err := SyncTool(steeringDir, root, "cline")
	if err != nil {
		t.Fatalf("SyncTool: %v", err)
	}
	if len(report.Written) != 1 {
		t.Errorf("expected 1 written, got %v", report.Written)
	}

	out := filepath.Join(root, ".clinerules", "api-rules.md")
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	content := string(data)

	// Cline: plain markdown, no frontmatter.
	if strings.Contains(content, "---") {
		t.Error("cline output should not contain frontmatter delimiters")
	}
	if !strings.HasPrefix(content, "# api-rules") {
		t.Error("cline output should start with # <name>")
	}
	if !strings.Contains(content, "Use RESTful conventions.") {
		t.Error("cline output should contain body")
	}
}

func TestSyncTool_KiroFormat(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")
	writeSteering(t, steeringDir, "api-rules", steeringAlways)

	report, err := SyncTool(steeringDir, root, "kiro")
	if err != nil {
		t.Fatalf("SyncTool: %v", err)
	}
	if len(report.Written) != 1 {
		t.Errorf("expected 1 written, got %v", report.Written)
	}

	out := filepath.Join(root, ".kiro", "steering", "api-rules.md")
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "name: api-rules") {
		t.Error("kiro output should contain name field")
	}
	if !strings.Contains(content, "mode: always") {
		t.Error("kiro output should map inclusion to mode")
	}
	if !strings.Contains(content, "Use RESTful conventions.") {
		t.Error("kiro output should contain body")
	}
}

func TestSyncTool_SkipsExcludedTool(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")
	writeSteering(t, steeringDir, "cursor-only", steeringCursorOnly)

	report, err := SyncTool(steeringDir, root, "kiro")
	if err != nil {
		t.Fatalf("SyncTool: %v", err)
	}
	if len(report.Written) != 0 {
		t.Errorf("expected 0 written for excluded tool, got %v", report.Written)
	}
	if len(report.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %v", report.Skipped)
	}
}

func TestSyncTool_IdempotentSkipsUnchanged(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")
	writeSteering(t, steeringDir, "api-rules", steeringAlways)

	// First sync writes the file.
	r1, err := SyncTool(steeringDir, root, "cursor")
	if err != nil {
		t.Fatalf("first sync: %v", err)
	}
	if len(r1.Written) != 1 {
		t.Fatalf("first sync should write 1 file, got %v", r1.Written)
	}

	// Second sync should skip (unchanged).
	r2, err := SyncTool(steeringDir, root, "cursor")
	if err != nil {
		t.Fatalf("second sync: %v", err)
	}
	if len(r2.Written) != 0 {
		t.Errorf("second sync should write 0 files (idempotent), got %v", r2.Written)
	}
	if len(r2.Skipped) != 1 {
		t.Errorf("second sync should skip 1 file, got %v", r2.Skipped)
	}
}

func TestSyncTool_UnsupportedToolReturnsError(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")
	writeSteering(t, steeringDir, "api-rules", steeringAlways)

	_, err := SyncTool(steeringDir, root, "claude")
	if err == nil {
		t.Fatal("expected error for unsupported sync tool")
	}
	if !strings.Contains(err.Error(), "unsupported sync tool") {
		t.Errorf("error should mention unsupported tool, got: %v", err)
	}
}

func TestSyncTool_PathTraversalRejected(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")

	// Create a steering file whose name field contains deep path
	// traversal that escapes the project root. filepath.Join cleans
	// the path, so we need enough ".." segments to escape past
	// .cursor/rules/ and the project root itself.
	traversalContent := `---
name: ../../../../etc/evil
description: path traversal attempt
inclusion: always
priority: 50
---
evil content
`
	writeSteering(t, steeringDir, "evil", traversalContent)

	report, err := SyncTool(steeringDir, root, "cursor")
	if err != nil {
		t.Fatalf("SyncTool: %v", err)
	}
	if len(report.Errors) == 0 {
		t.Error("path traversal should produce a boundary validation error")
	}
	if len(report.Written) != 0 {
		t.Errorf("path traversal file should not be written, got %v", report.Written)
	}
}

func TestSyncAll_SyncsToAllTools(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")
	writeSteering(t, steeringDir, "api-rules", steeringAlways)

	report, err := SyncAll(steeringDir, root)
	if err != nil {
		t.Fatalf("SyncAll: %v", err)
	}

	// Should write to cursor, cline, and kiro.
	if len(report.Written) != 3 {
		t.Errorf("expected 3 written (one per tool), got %d: %v", len(report.Written), report.Written)
	}

	// Verify all output files exist.
	paths := []string{
		filepath.Join(root, ".cursor", "rules", "api-rules.mdc"),
		filepath.Join(root, ".clinerules", "api-rules.md"),
		filepath.Join(root, ".kiro", "steering", "api-rules.md"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("expected output file to exist: %s", p)
		}
	}
}

func TestSyncAll_SkipsToolExcludedFiles(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")
	writeSteering(t, steeringDir, "cursor-only", steeringCursorOnly)

	report, err := SyncAll(steeringDir, root)
	if err != nil {
		t.Fatalf("SyncAll: %v", err)
	}

	// Only cursor should get the file; cline and kiro should skip.
	if len(report.Written) != 1 {
		t.Errorf("expected 1 written (cursor only), got %d: %v", len(report.Written), report.Written)
	}
	if len(report.Skipped) != 2 {
		t.Errorf("expected 2 skipped (cline + kiro), got %d: %v", len(report.Skipped), report.Skipped)
	}
}

func TestSyncTool_CursorAlwaysApplyFalseForNonAlways(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")
	writeSteering(t, steeringDir, "manual-rule", steeringManual)

	_, err := SyncTool(steeringDir, root, "cursor")
	if err != nil {
		t.Fatalf("SyncTool: %v", err)
	}

	out := filepath.Join(root, ".cursor", "rules", "manual-rule.mdc")
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(data), "alwaysApply: false") {
		t.Error("cursor output should have alwaysApply: false for manual inclusion")
	}
}

func TestSyncTool_KiroModeMapping(t *testing.T) {
	root := t.TempDir()
	steeringDir := filepath.Join(root, ".context", "steering")

	autoContent := `---
name: auto-rule
description: Auto rule
inclusion: auto
priority: 50
---
Auto body.
`
	writeSteering(t, steeringDir, "auto-rule", autoContent)

	_, err := SyncTool(steeringDir, root, "kiro")
	if err != nil {
		t.Fatalf("SyncTool: %v", err)
	}

	out := filepath.Join(root, ".kiro", "steering", "auto-rule.md")
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(data), "mode: auto") {
		t.Error("kiro output should map auto inclusion to mode: auto")
	}
}
