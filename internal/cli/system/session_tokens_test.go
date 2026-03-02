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

func TestFindJSONLPath_Found(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)
	t.Setenv("HOME", tmpDir)

	// Create a fake JSONL file in the expected location
	sessionID := "test-session-abc123"
	projectDir := filepath.Join(tmpDir, ".claude", "projects", "testproj")
	if mkErr := os.MkdirAll(projectDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	jsonlPath := filepath.Join(projectDir, sessionID+".jsonl")
	if writeErr := os.WriteFile(jsonlPath, []byte(`{}`+"\n"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	got, findErr := findJSONLPath(sessionID)
	if findErr != nil {
		t.Fatalf("unexpected error: %v", findErr)
	}
	if got != jsonlPath {
		t.Errorf("findJSONLPath() = %q, want %q", got, jsonlPath)
	}
}

func TestFindJSONLPath_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)
	t.Setenv("HOME", tmpDir)

	got, findErr := findJSONLPath("nonexistent-session")
	if findErr != nil {
		t.Fatalf("unexpected error: %v", findErr)
	}
	if got != "" {
		t.Errorf("expected empty path, got %q", got)
	}
}

func TestFindJSONLPath_Cached(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_RUNTIME_DIR", tmpDir)
	t.Setenv("HOME", tmpDir)

	sessionID := "test-cached-session"
	projectDir := filepath.Join(tmpDir, ".claude", "projects", "testproj")
	if mkErr := os.MkdirAll(projectDir, 0o750); mkErr != nil {
		t.Fatal(mkErr)
	}
	jsonlPath := filepath.Join(projectDir, sessionID+".jsonl")
	if writeErr := os.WriteFile(jsonlPath, []byte(`{}`+"\n"), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	// First call populates cache
	first, findErr := findJSONLPath(sessionID)
	if findErr != nil {
		t.Fatalf("first call error: %v", findErr)
	}

	// Verify cache file exists
	cacheFile := filepath.Join(tmpDir, "ctx", "jsonl-path-"+sessionID)
	if _, statErr := os.Stat(cacheFile); statErr != nil {
		t.Fatalf("cache file not created: %v", statErr)
	}

	// Second call should return same result (from cache)
	second, findErr := findJSONLPath(sessionID)
	if findErr != nil {
		t.Fatalf("second call error: %v", findErr)
	}
	if first != second {
		t.Errorf("cached result mismatch: first=%q second=%q", first, second)
	}
}

func TestParseLastUsageAndModel_ValidData(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")

	content := `{"type":"human","message":{"role":"user","content":"hello"}}
{"type":"assistant","message":{"model":"claude-opus-4-6","role":"assistant","content":"hi","usage":{"input_tokens":50000,"output_tokens":500,"cache_creation_input_tokens":8000,"cache_read_input_tokens":100000}}}
`
	if writeErr := os.WriteFile(jsonlPath, []byte(content), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	info, parseErr := parseLastUsageAndModel(jsonlPath)
	if parseErr != nil {
		t.Fatalf("unexpected error: %v", parseErr)
	}

	// 50000 + 8000 + 100000 = 158000
	wantTokens := 158000
	if info.Tokens != wantTokens {
		t.Errorf("Tokens = %d, want %d", info.Tokens, wantTokens)
	}
	if info.Model != "claude-opus-4-6" {
		t.Errorf("Model = %q, want %q", info.Model, "claude-opus-4-6")
	}
}

func TestParseLastUsageAndModel_NoUsage(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")

	content := `{"type":"human","message":{"role":"user","content":"hello"}}
`
	if writeErr := os.WriteFile(jsonlPath, []byte(content), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	info, parseErr := parseLastUsageAndModel(jsonlPath)
	if parseErr != nil {
		t.Fatalf("unexpected error: %v", parseErr)
	}
	if info.Tokens != 0 {
		t.Errorf("Tokens = %d, want 0", info.Tokens)
	}
	if info.Model != "" {
		t.Errorf("Model = %q, want empty", info.Model)
	}
}

func TestParseLastUsageAndModel_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")

	if writeErr := os.WriteFile(jsonlPath, []byte(""), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	info, parseErr := parseLastUsageAndModel(jsonlPath)
	if parseErr != nil {
		t.Fatalf("unexpected error: %v", parseErr)
	}
	if info.Tokens != 0 {
		t.Errorf("Tokens = %d, want 0", info.Tokens)
	}
}

func TestParseLastUsageAndModel_MultipleMessages(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")

	// Two assistant messages — should pick the last one
	content := `{"type":"assistant","message":{"model":"claude-sonnet-4-5","role":"assistant","content":"first","usage":{"input_tokens":10000,"output_tokens":100,"cache_creation_input_tokens":0,"cache_read_input_tokens":5000}}}
{"type":"human","message":{"role":"user","content":"more"}}
{"type":"assistant","message":{"model":"claude-opus-4-6-20260205","role":"assistant","content":"second","usage":{"input_tokens":80000,"output_tokens":200,"cache_creation_input_tokens":4000,"cache_read_input_tokens":60000}}}
`
	if writeErr := os.WriteFile(jsonlPath, []byte(content), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	info, parseErr := parseLastUsageAndModel(jsonlPath)
	if parseErr != nil {
		t.Fatalf("unexpected error: %v", parseErr)
	}

	// Last assistant: 80000 + 4000 + 60000 = 144000
	wantTokens := 144000
	if info.Tokens != wantTokens {
		t.Errorf("Tokens = %d, want %d", info.Tokens, wantTokens)
	}
	if info.Model != "claude-opus-4-6-20260205" {
		t.Errorf("Model = %q, want %q", info.Model, "claude-opus-4-6-20260205")
	}
}

func TestParseLastUsageAndModel_NoModelField(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")

	// JSONL without model field (older format)
	content := `{"type":"assistant","message":{"role":"assistant","content":"hi","usage":{"input_tokens":50000,"output_tokens":500,"cache_creation_input_tokens":8000,"cache_read_input_tokens":100000}}}
`
	if writeErr := os.WriteFile(jsonlPath, []byte(content), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	info, parseErr := parseLastUsageAndModel(jsonlPath)
	if parseErr != nil {
		t.Fatalf("unexpected error: %v", parseErr)
	}

	if info.Tokens != 158000 {
		t.Errorf("Tokens = %d, want 158000", info.Tokens)
	}
	if info.Model != "" {
		t.Errorf("Model = %q, want empty", info.Model)
	}
}

func TestReadSessionTokenInfo_EmptySessionID(t *testing.T) {
	info, readErr := readSessionTokenInfo("")
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if info.Tokens != 0 {
		t.Errorf("Tokens = %d, want 0", info.Tokens)
	}
}

func TestReadSessionTokenInfo_UnknownSessionID(t *testing.T) {
	info, readErr := readSessionTokenInfo("unknown")
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if info.Tokens != 0 {
		t.Errorf("Tokens = %d, want 0", info.Tokens)
	}
}

func TestModelContextWindow(t *testing.T) {
	tests := []struct {
		model string
		want  int
	}{
		// 1M-capable models
		{"claude-opus-4-6", contextWindow1M},
		{"claude-opus-4-6-20260205", contextWindow1M},
		{"claude-sonnet-4-6", contextWindow1M},
		{"claude-sonnet-4-6-20260217", contextWindow1M},
		{"claude-sonnet-4-5", contextWindow1M},
		{"claude-sonnet-4-5-20250929", contextWindow1M},
		{"claude-sonnet-4", contextWindow1M},
		{"claude-sonnet-4-0", contextWindow1M},
		{"claude-sonnet-4-20250514", contextWindow1M},

		// 200k models
		{"claude-opus-4-5", rc.DefaultContextWindow},
		{"claude-opus-4-5-20251101", rc.DefaultContextWindow},
		{"claude-opus-4-1", rc.DefaultContextWindow},
		{"claude-opus-4-1-20250805", rc.DefaultContextWindow},
		{"claude-haiku-4-5-20251001", rc.DefaultContextWindow},
		{"claude-3-haiku-20240307", rc.DefaultContextWindow},
		{"claude-3-5-sonnet-20241022", rc.DefaultContextWindow},

		// Unknown / empty
		{"", 0},
		{"gpt-4o", 0},
		{"some-custom-model", 0},
	}

	for _, tt := range tests {
		name := tt.model
		if name == "" {
			name = "empty"
		}
		t.Run(name, func(t *testing.T) {
			got := modelContextWindow(tt.model)
			if got != tt.want {
				t.Errorf("modelContextWindow(%q) = %d, want %d", tt.model, got, tt.want)
			}
		})
	}
}

func TestEffectiveContextWindow(t *testing.T) {
	rc.Reset()

	// Tier 1: known 1M model
	got := effectiveContextWindow("claude-opus-4-6")
	if got != contextWindow1M {
		t.Errorf("with 1M model: got %d, want %d", got, contextWindow1M)
	}

	// Tier 1: known 200k model
	got = effectiveContextWindow("claude-opus-4-5")
	if got != rc.DefaultContextWindow {
		t.Errorf("with 200k model: got %d, want %d", got, rc.DefaultContextWindow)
	}

	// Tier 2/3: empty model falls to rc.ContextWindow()
	got = effectiveContextWindow("")
	if got != rc.DefaultContextWindow {
		t.Errorf("with empty model: got %d, want %d", got, rc.DefaultContextWindow)
	}

	// Tier 2/3: non-Claude model falls to rc.ContextWindow()
	got = effectiveContextWindow("gpt-4o")
	if got != rc.DefaultContextWindow {
		t.Errorf("with unknown model: got %d, want %d", got, rc.DefaultContextWindow)
	}
}

func TestFormatTokenCount(t *testing.T) {
	tests := []struct {
		tokens int
		want   string
	}{
		{500, "500"},
		{1200, "1.2k"},
		{9900, "9.9k"},
		{10000, "10k"},
		{52000, "52k"},
		{164000, "164k"},
		{200000, "200k"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatTokenCount(tt.tokens)
			if got != tt.want {
				t.Errorf("formatTokenCount(%d) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}

func TestFormatWindowSize(t *testing.T) {
	tests := []struct {
		size int
		want string
	}{
		{500, "500"},
		{128000, "128k"},
		{200000, "200k"},
		{1000000, "1000k"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatWindowSize(tt.size)
			if got != tt.want {
				t.Errorf("formatWindowSize(%d) = %q, want %q", tt.size, got, tt.want)
			}
		})
	}
}
