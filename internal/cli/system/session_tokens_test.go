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

func TestParseLastUsage_ValidData(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")

	// Write a JSONL file with an assistant message containing usage data
	content := `{"type":"human","message":{"role":"user","content":"hello"}}
{"type":"assistant","message":{"role":"assistant","content":"hi","usage":{"input_tokens":50000,"output_tokens":500,"cache_creation_input_tokens":8000,"cache_read_input_tokens":100000}}}
`
	if writeErr := os.WriteFile(jsonlPath, []byte(content), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	tokens, parseErr := parseLastUsage(jsonlPath)
	if parseErr != nil {
		t.Fatalf("unexpected error: %v", parseErr)
	}

	// 50000 + 8000 + 100000 = 158000
	want := 158000
	if tokens != want {
		t.Errorf("parseLastUsage() = %d, want %d", tokens, want)
	}
}

func TestParseLastUsage_NoUsage(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")

	content := `{"type":"human","message":{"role":"user","content":"hello"}}
`
	if writeErr := os.WriteFile(jsonlPath, []byte(content), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	tokens, parseErr := parseLastUsage(jsonlPath)
	if parseErr != nil {
		t.Fatalf("unexpected error: %v", parseErr)
	}
	if tokens != 0 {
		t.Errorf("parseLastUsage() = %d, want 0", tokens)
	}
}

func TestParseLastUsage_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")

	if writeErr := os.WriteFile(jsonlPath, []byte(""), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	tokens, parseErr := parseLastUsage(jsonlPath)
	if parseErr != nil {
		t.Fatalf("unexpected error: %v", parseErr)
	}
	if tokens != 0 {
		t.Errorf("parseLastUsage() = %d, want 0", tokens)
	}
}

func TestParseLastUsage_MultipleMessages(t *testing.T) {
	tmpDir := t.TempDir()
	jsonlPath := filepath.Join(tmpDir, "test.jsonl")

	// Two assistant messages — should pick the last one
	content := `{"type":"assistant","message":{"role":"assistant","content":"first","usage":{"input_tokens":10000,"output_tokens":100,"cache_creation_input_tokens":0,"cache_read_input_tokens":5000}}}
{"type":"human","message":{"role":"user","content":"more"}}
{"type":"assistant","message":{"role":"assistant","content":"second","usage":{"input_tokens":80000,"output_tokens":200,"cache_creation_input_tokens":4000,"cache_read_input_tokens":60000}}}
`
	if writeErr := os.WriteFile(jsonlPath, []byte(content), 0o600); writeErr != nil {
		t.Fatal(writeErr)
	}

	tokens, parseErr := parseLastUsage(jsonlPath)
	if parseErr != nil {
		t.Fatalf("unexpected error: %v", parseErr)
	}

	// Last assistant: 80000 + 4000 + 60000 = 144000
	want := 144000
	if tokens != want {
		t.Errorf("parseLastUsage() = %d, want %d", tokens, want)
	}
}

func TestReadSessionTokenUsage_EmptySessionID(t *testing.T) {
	tokens, readErr := readSessionTokenUsage("")
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if tokens != 0 {
		t.Errorf("readSessionTokenUsage(\"\") = %d, want 0", tokens)
	}
}

func TestReadSessionTokenUsage_UnknownSessionID(t *testing.T) {
	tokens, readErr := readSessionTokenUsage("unknown")
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if tokens != 0 {
		t.Errorf("readSessionTokenUsage(\"unknown\") = %d, want 0", tokens)
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
