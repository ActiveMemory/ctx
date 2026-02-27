//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/rc"
)

func TestRenderTemplate_Static(t *testing.T) {
	got := renderTemplate("hello world", nil, "fallback")
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestRenderTemplate_WithVars(t *testing.T) {
	got := renderTemplate("count is {{.Count}}", map[string]any{"Count": 42}, "fallback")
	if got != "count is 42" {
		t.Errorf("expected 'count is 42', got %q", got)
	}
}

func TestRenderTemplate_Empty(t *testing.T) {
	got := renderTemplate("", nil, "fallback")
	if got != "" {
		t.Errorf("expected empty string for empty template, got %q", got)
	}
}

func TestRenderTemplate_WhitespaceOnly(t *testing.T) {
	got := renderTemplate("   \n\t  ", nil, "fallback")
	if got != "" {
		t.Errorf("expected empty string for whitespace-only template, got %q", got)
	}
}

func TestRenderTemplate_MalformedTemplate(t *testing.T) {
	got := renderTemplate("{{.Bad", nil, "fallback")
	if got != "fallback" {
		t.Errorf("expected fallback for malformed template, got %q", got)
	}
}

func TestRenderTemplate_UnknownVariable(t *testing.T) {
	got := renderTemplate("value is {{.Missing}}", map[string]any{}, "fallback")
	if got != "value is <no value>" {
		t.Errorf("expected '<no value>' for unknown variable, got %q", got)
	}
}

func TestRenderTemplate_NilVars(t *testing.T) {
	got := renderTemplate("static text", nil, "fallback")
	if got != "static text" {
		t.Errorf("expected 'static text', got %q", got)
	}
}

func TestLoadMessage_EmbeddedDefault(t *testing.T) {
	// qa-reminder/gate.txt exists as an embedded template
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()

	got := loadMessage("qa-reminder", "gate", nil, "should-not-see-fallback")
	if got == "should-not-see-fallback" {
		t.Error("expected embedded template to be used instead of fallback")
	}
	if got == "" {
		t.Error("expected non-empty result from embedded template")
	}
}

func TestLoadMessage_UserOverride(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()

	// Create user override
	overrideDir := filepath.Join(rc.ContextDir(), "hooks", "messages", "qa-reminder")
	if err := os.MkdirAll(overrideDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(overrideDir, "gate.txt"), []byte("custom message"), 0o600); err != nil {
		t.Fatal(err)
	}

	got := loadMessage("qa-reminder", "gate", nil, "fallback")
	if got != "custom message" {
		t.Errorf("expected user override 'custom message', got %q", got)
	}
}

func TestLoadMessage_UserOverrideEmpty(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()

	// Create empty user override (intentional silence)
	overrideDir := filepath.Join(rc.ContextDir(), "hooks", "messages", "qa-reminder")
	if err := os.MkdirAll(overrideDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(overrideDir, "gate.txt"), []byte(""), 0o600); err != nil {
		t.Fatal(err)
	}

	got := loadMessage("qa-reminder", "gate", nil, "fallback")
	if got != "" {
		t.Errorf("expected empty string for empty override (silence), got %q", got)
	}
}

func TestLoadMessage_FallbackWhenNoTemplate(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()

	got := loadMessage("nonexistent-hook", "nonexistent-variant", nil, "my fallback")
	if got != "my fallback" {
		t.Errorf("expected fallback 'my fallback', got %q", got)
	}
}

func TestLoadMessage_TemplateWithVars(t *testing.T) {
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()

	// Create user override with template variable
	overrideDir := filepath.Join(rc.ContextDir(), "hooks", "messages", "test-hook")
	if err := os.MkdirAll(overrideDir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(overrideDir, "test.txt"),
		[]byte("Hello {{.Name}}, count={{.Count}}"), 0o600); err != nil {
		t.Fatal(err)
	}

	got := loadMessage("test-hook", "test",
		map[string]any{"Name": "World", "Count": 5}, "fallback")
	if got != "Hello World, count=5" {
		t.Errorf("expected 'Hello World, count=5', got %q", got)
	}
}

func TestLoadMessage_EmbeddedWithVars(t *testing.T) {
	// check-persistence/nudge.txt has {{.PromptsSinceNudge}}
	workDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(workDir)
	defer func() { _ = os.Chdir(origDir) }()
	rc.Reset()

	got := loadMessage("check-persistence", "nudge",
		map[string]any{"PromptsSinceNudge": 15}, "fallback")
	if got == "fallback" {
		t.Error("expected embedded template, got fallback")
	}
	if got == "" {
		t.Error("expected non-empty result")
	}
	// Template should have rendered the variable
	if !contains(got, "15+ prompts") {
		t.Errorf("expected rendered variable '15+ prompts' in output, got %q", got)
	}
}

func TestBoxLines(t *testing.T) {
	got := boxLines("line one\nline two\nline three")
	expected := "│ line one\n│ line two\n│ line three\n"
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestBoxLines_TrailingNewline(t *testing.T) {
	// Trailing newlines should be trimmed before splitting
	got := boxLines("single line\n")
	expected := "│ single line\n"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestBoxLines_EmptyString(t *testing.T) {
	got := boxLines("")
	expected := "│ \n"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
