//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config"
)

func TestMarkdownSessionParser_Tool(t *testing.T) {
	p := NewMarkdownSessionParser()
	if got := p.Tool(); got != config.ToolMarkdown {
		t.Errorf("Tool() = %q, want %q", got, config.ToolMarkdown)
	}
}

func TestMarkdownSessionParser_Matches(t *testing.T) {
	p := NewMarkdownSessionParser()
	dir := t.TempDir()

	tests := []struct {
		name    string
		file    string
		content string
		want    bool
	}{
		{
			name:    "valid session header",
			file:    "2026-01-15-fix-api.md",
			content: "# Session: 2026-01-15 — Fix API\n\n## What Was Done\n- Fixed endpoint\n",
			want:    true,
		},
		{
			name:    "valid session header with hyphen separator",
			file:    "2026-01-15-fix-api.md",
			content: "# Session: 2026-01-15 - Fix API\n\n## What Was Done\n- Fixed endpoint\n",
			want:    true,
		},
		{
			name:    "valid Turkish header",
			file:    "2026-01-15-duzeltme.md",
			content: "# Oturum: 2026-01-15 — API Düzeltme\n\n## What Was Done\n- Fixed\n",
			want:    true,
		},
		{
			name:    "valid date-only header",
			file:    "2026-01-15-work.md",
			content: "# 2026-01-15 — Morning Work\n\n## What Was Done\n",
			want:    true,
		},
		{
			name:    "non-markdown file",
			file:    "session.txt",
			content: "# Session: 2026-01-15 — Fix API\n",
			want:    false,
		},
		{
			name:    "README.md",
			file:    "README.md",
			content: "# Session: 2026-01-15 — Readme\n",
			want:    false,
		},
		{
			name:    "markdown without session header",
			file:    "notes.md",
			content: "# Random Notes\n\nSome notes here.\n",
			want:    false,
		},
		{
			name:    "empty markdown",
			file:    "empty.md",
			content: "",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.file)
			if err := os.WriteFile(path, []byte(tt.content), 0600); err != nil {
				t.Fatal(err)
			}
			if got := p.Matches(path); got != tt.want {
				t.Errorf("Matches(%q) = %v, want %v", tt.file, got, tt.want)
			}
		})
	}
}

func TestMarkdownSessionParser_ParseFile(t *testing.T) {
	p := NewMarkdownSessionParser()

	// Set up directory structure: project/.context/sessions/
	dir := t.TempDir()
	sessionsDir := filepath.Join(dir, ".context", "sessions")
	if err := os.MkdirAll(sessionsDir, 0750); err != nil {
		t.Fatal(err)
	}

	content := `# Session: 2026-01-15 — Fix API Rate Limiting

## What Was Done
- Added rate limiter middleware
- Updated tests for rate limiting

## Decisions
- Chose token bucket algorithm over sliding window

## Learnings
- Go's rate package is simpler than custom implementations

## Next Steps
- Add monitoring dashboard for rate limits
`

	path := filepath.Join(sessionsDir, "2026-01-15-fix-api.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	sessions, err := p.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}

	s := sessions[0]

	if s.ID != "2026-01-15-fix-api" {
		t.Errorf("ID = %q, want %q", s.ID, "2026-01-15-fix-api")
	}
	if s.Tool != config.ToolMarkdown {
		t.Errorf("Tool = %q, want %q", s.Tool, config.ToolMarkdown)
	}
	if s.FirstUserMsg != "Fix API Rate Limiting" {
		t.Errorf("FirstUserMsg = %q, want %q", s.FirstUserMsg, "Fix API Rate Limiting")
	}
	if s.StartTime.Year() != 2026 || s.StartTime.Month() != 1 || s.StartTime.Day() != 15 {
		t.Errorf("StartTime = %v, want 2026-01-15", s.StartTime)
	}
	if s.Project == "" {
		// Project should be inferred from parent dir of .context/sessions/
		t.Log("Project not inferred (expected in this temp dir structure)")
	}
	if s.TurnCount != 1 {
		t.Errorf("TurnCount = %d, want 1", s.TurnCount)
	}
	if len(s.Messages) < 1 {
		t.Fatal("expected at least 1 message")
	}
}

func TestMarkdownSessionParser_ParseFile_NoHeader(t *testing.T) {
	p := NewMarkdownSessionParser()
	dir := t.TempDir()

	path := filepath.Join(dir, "no-header.md")
	content := "# Random Document\n\nNo session header here.\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	sessions, err := p.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(sessions) != 0 {
		t.Errorf("expected 0 sessions for non-session file, got %d", len(sessions))
	}
}

func TestMarkdownSessionParser_ParseLine(t *testing.T) {
	p := NewMarkdownSessionParser()
	msg, sessID, err := p.ParseLine([]byte("anything"))
	if msg != nil || sessID != "" || err != nil {
		t.Error("ParseLine should return nil, empty, nil for markdown parser")
	}
}

func TestIsSessionHeader(t *testing.T) {
	tests := []struct {
		name string
		line string
		want bool
	}{
		{"session with em dash", "# Session: 2026-01-15 — Fix API", true},
		{"session with hyphen", "# Session: 2026-01-15 - Fix API", true},
		{"turkish header", "# Oturum: 2026-01-15 — Düzeltme", true},
		{"date only with em dash", "# 2026-01-15 — Morning Work", true},
		{"date only with hyphen", "# 2026-01-15 - Morning Work", true},
		{"date only no topic", "# 2026-01-15", true},
		{"not a header", "## Session: 2026-01-15 — Fix API", false},
		{"no hash", "Session: 2026-01-15 — Fix API", false},
		{"random h1", "# Random Title", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSessionHeader(tt.line); got != tt.want {
				t.Errorf("isSessionHeader(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestParseSessionHeader(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantDate  string
		wantTopic string
	}{
		{
			name:      "session with em dash",
			line:      "# Session: 2026-01-15 — Fix API Rate Limiting",
			wantDate:  "2026-01-15",
			wantTopic: "Fix API Rate Limiting",
		},
		{
			name:      "session with hyphen",
			line:      "# Session: 2026-01-15 - Fix API",
			wantDate:  "2026-01-15",
			wantTopic: "Fix API",
		},
		{
			name:      "turkish header",
			line:      "# Oturum: 2026-01-15 — Düzeltme",
			wantDate:  "2026-01-15",
			wantTopic: "Düzeltme",
		},
		{
			name:      "date only with topic",
			line:      "# 2026-01-15 — Morning Work",
			wantDate:  "2026-01-15",
			wantTopic: "Morning Work",
		},
		{
			name:      "date only no topic",
			line:      "# 2026-01-15",
			wantDate:  "",
			wantTopic: "2026-01-15",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, topic := parseSessionHeader(tt.line)
			if date != tt.wantDate {
				t.Errorf("date = %q, want %q", date, tt.wantDate)
			}
			if topic != tt.wantTopic {
				t.Errorf("topic = %q, want %q", topic, tt.wantTopic)
			}
		})
	}
}

func TestParseSessionDate(t *testing.T) {
	tests := []struct {
		name    string
		dateStr string
		wantOK  bool
	}{
		{"valid date", "2026-01-15", true},
		{"invalid format", "01-15-2026", false},
		{"empty", "", false},
		{"partial", "2026-01", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSessionDate(tt.dateStr)
			gotOK := !result.IsZero()
			if gotOK != tt.wantOK {
				t.Errorf("parseSessionDate(%q) isZero = %v, want non-zero = %v",
					tt.dateStr, result.IsZero(), tt.wantOK)
			}
		})
	}
}

func TestExtractSections(t *testing.T) {
	lines := []string{
		"# Session: 2026-01-15 — Test",
		"",
		"## What Was Done",
		"- Item 1",
		"- Item 2",
		"",
		"## Decisions",
		"- Decision A",
		"",
		"## Learnings",
		"- Learning X",
	}

	sections := extractSections(lines)

	if len(sections) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(sections))
	}

	if _, ok := sections["What Was Done"]; !ok {
		t.Error("missing 'What Was Done' section")
	}
	if _, ok := sections["Decisions"]; !ok {
		t.Error("missing 'Decisions' section")
	}
	if _, ok := sections["Learnings"]; !ok {
		t.Error("missing 'Learnings' section")
	}
}

func TestScanDirectory_WithMarkdown(t *testing.T) {
	dir := t.TempDir()

	// Create a valid markdown session file
	mdContent := `# Session: 2026-02-10 — Debug Parser

## What Was Done
- Fixed parser edge case

## Learnings
- Edge cases in date parsing
`
	mdFile := filepath.Join(dir, "2026-02-10-debug-parser.md")
	if err := os.WriteFile(mdFile, []byte(mdContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Create a valid Claude JSONL session file
	jsonlContent := `{"uuid":"m1","sessionId":"s1","slug":"test","type":"user","timestamp":"2026-01-20T10:00:00Z","cwd":"/test","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"hello"}]}}`
	jsonlFile := filepath.Join(dir, "session.jsonl")
	if err := os.WriteFile(jsonlFile, []byte(jsonlContent), 0600); err != nil {
		t.Fatal(err)
	}

	sessions, err := ScanDirectory(dir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions (1 markdown + 1 jsonl), got %d", len(sessions))
	}

	// Check that both tools are represented
	tools := make(map[string]bool)
	for _, s := range sessions {
		tools[s.Tool] = true
	}
	if !tools[config.ToolMarkdown] {
		t.Error("expected markdown session in results")
	}
	if !tools[config.ToolClaudeCode] {
		t.Error("expected claude-code session in results")
	}
}

func TestRegisteredTools_IncludesMarkdown(t *testing.T) {
	tools := RegisteredTools()
	found := false
	for _, tool := range tools {
		if tool == config.ToolMarkdown {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected %q in registered tools", config.ToolMarkdown)
	}
}

func TestGetParser_Markdown(t *testing.T) {
	p := Parser(config.ToolMarkdown)
	if p == nil {
		t.Fatalf("expected parser for %q", config.ToolMarkdown)
	}
	if p.Tool() != config.ToolMarkdown {
		t.Errorf("Tool() = %q, want %q", p.Tool(), config.ToolMarkdown)
	}
}
