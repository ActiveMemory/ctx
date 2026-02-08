//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestClaudeCodeParser_Matches(t *testing.T) {
	parser := NewClaudeCodeParser()

	// Create temp directory
	dir := t.TempDir()

	// Test: non-JSONL file
	txtFile := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(txtFile, []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}
	if parser.Matches(txtFile) {
		t.Error("should not parse .txt files")
	}

	// Test: JSONL file without Claude Code format
	badJSONL := filepath.Join(dir, "bad.jsonl")
	if err := os.WriteFile(badJSONL, []byte(`{"foo": "bar"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if parser.Matches(badJSONL) {
		t.Error("should not parse non-Claude Code JSONL")
	}

	// Test: Valid Claude Code JSONL
	goodJSONL := filepath.Join(dir, "good.jsonl")
	content := `{"uuid":"abc","sessionId":"session-1","slug":"test-session","type":"user","timestamp":"2026-01-20T10:00:00Z","cwd":"/home/test","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"hello"}]}}`
	if err := os.WriteFile(goodJSONL, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if !parser.Matches(goodJSONL) {
		t.Error("should parse valid Claude Code JSONL")
	}
}

func TestClaudeCodeParser_ParseFile(t *testing.T) {
	parser := NewClaudeCodeParser()
	dir := t.TempDir()

	// Create test JSONL with multiple messages
	jsonlFile := filepath.Join(dir, "session.jsonl")
	content := `{"uuid":"msg1","sessionId":"sess-1","slug":"test-session","type":"user","timestamp":"2026-01-20T10:00:00Z","cwd":"/home/test/project","gitBranch":"main","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"What is 2+2?"}]}}
{"uuid":"msg2","parentUuid":"msg1","sessionId":"sess-1","slug":"test-session","type":"assistant","timestamp":"2026-01-20T10:00:30Z","cwd":"/home/test/project","gitBranch":"main","version":"2.1.0","message":{"model":"claude-opus-4-5","role":"assistant","content":[{"type":"thinking","thinking":"Let me calculate..."},{"type":"text","text":"2+2 equals 4."}],"usage":{"input_tokens":100,"output_tokens":50}}}
{"uuid":"msg3","parentUuid":"msg2","sessionId":"sess-1","slug":"test-session","type":"user","timestamp":"2026-01-20T10:01:00Z","cwd":"/home/test/project","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"Thanks!"}]}}`

	if err := os.WriteFile(jsonlFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	sessions, err := parser.ParseFile(jsonlFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}

	session := sessions[0]

	// Check session fields
	if session.ID != "sess-1" {
		t.Errorf("expected session ID 'sess-1', got '%s'", session.ID)
	}
	if session.Slug != "test-session" {
		t.Errorf("expected slug 'test-session', got '%s'", session.Slug)
	}
	if session.Tool != "claude-code" {
		t.Errorf("expected tool 'claude-code', got '%s'", session.Tool)
	}
	if session.Project != "project" {
		t.Errorf("expected project 'project', got '%s'", session.Project)
	}
	if session.GitBranch != "main" {
		t.Errorf("expected git branch 'main', got '%s'", session.GitBranch)
	}
	if session.Model != "claude-opus-4-5" {
		t.Errorf("expected model 'claude-opus-4-5', got '%s'", session.Model)
	}

	// Check message count
	if len(session.Messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(session.Messages))
	}

	// Check turn count (user messages)
	if session.TurnCount != 2 {
		t.Errorf("expected 2 turns, got %d", session.TurnCount)
	}

	// Check token stats
	if session.TotalTokensIn != 100 {
		t.Errorf("expected 100 input tokens, got %d", session.TotalTokensIn)
	}
	if session.TotalTokensOut != 50 {
		t.Errorf("expected 50 output tokens, got %d", session.TotalTokensOut)
	}

	// Check first user message preview
	if session.FirstUserMsg != "What is 2+2?" {
		t.Errorf("expected first user msg 'What is 2+2?', got '%s'", session.FirstUserMsg)
	}

	// Check message content
	msg1 := session.Messages[0]
	if !msg1.BelongsToUser() {
		t.Error("first message should be user")
	}
	if msg1.Text != "What is 2+2?" {
		t.Errorf("expected 'What is 2+2?', got '%s'", msg1.Text)
	}

	msg2 := session.Messages[1]
	if !msg2.BelongsToAssistant() {
		t.Error("second message should be assistant")
	}
	if msg2.Thinking != "Let me calculate..." {
		t.Errorf("expected thinking content, got '%s'", msg2.Thinking)
	}
	if msg2.Text != "2+2 equals 4." {
		t.Errorf("expected '2+2 equals 4.', got '%s'", msg2.Text)
	}

	// Check duration
	expectedDuration := time.Minute
	if session.Duration != expectedDuration {
		t.Errorf("expected duration %v, got %v", expectedDuration, session.Duration)
	}
}

func TestClaudeCodeParser_ParseFile_WithToolUse(t *testing.T) {
	parser := NewClaudeCodeParser()
	dir := t.TempDir()

	jsonlFile := filepath.Join(dir, "tools.jsonl")
	content := `{"uuid":"msg1","sessionId":"sess-2","slug":"tool-session","type":"user","timestamp":"2026-01-20T10:00:00Z","cwd":"/home/test","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"List files"}]}}
{"uuid":"msg2","sessionId":"sess-2","slug":"tool-session","type":"assistant","timestamp":"2026-01-20T10:00:10Z","cwd":"/home/test","version":"2.1.0","message":{"role":"assistant","content":[{"type":"tool_use","id":"tool1","name":"bash","input":{"command":"ls -la"}}]}}
{"uuid":"msg3","sessionId":"sess-2","slug":"tool-session","type":"user","timestamp":"2026-01-20T10:00:11Z","cwd":"/home/test","version":"2.1.0","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"tool1","content":"file1.txt\nfile2.txt"}]}}`

	if err := os.WriteFile(jsonlFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	sessions, err := parser.ParseFile(jsonlFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}

	session := sessions[0]

	// Check tool uses
	allTools := session.AllToolUses()
	if len(allTools) != 1 {
		t.Fatalf("expected 1 tool use, got %d", len(allTools))
	}
	if allTools[0].Name != "bash" {
		t.Errorf("expected tool name 'bash', got '%s'", allTools[0].Name)
	}

	// Check tool result in message
	msg3 := session.Messages[2]
	if len(msg3.ToolResults) != 1 {
		t.Fatalf("expected 1 tool result, got %d", len(msg3.ToolResults))
	}
	if msg3.ToolResults[0].ToolUseID != "tool1" {
		t.Errorf("expected tool_use_id 'tool1', got '%s'", msg3.ToolResults[0].ToolUseID)
	}
}

func TestClaudeCodeParser_ParseFile_MultipleSessions(t *testing.T) {
	parser := NewClaudeCodeParser()
	dir := t.TempDir()

	// JSONL with two different sessions
	jsonlFile := filepath.Join(dir, "multi.jsonl")
	content := `{"uuid":"a1","sessionId":"sess-A","slug":"session-a","type":"user","timestamp":"2026-01-20T09:00:00Z","cwd":"/home/test","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"Hello A"}]}}
{"uuid":"b1","sessionId":"sess-B","slug":"session-b","type":"user","timestamp":"2026-01-20T10:00:00Z","cwd":"/home/test","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"Hello B"}]}}
{"uuid":"a2","sessionId":"sess-A","slug":"session-a","type":"assistant","timestamp":"2026-01-20T09:00:30Z","cwd":"/home/test","version":"2.1.0","message":{"role":"assistant","content":[{"type":"text","text":"Hi A!"}]}}`

	if err := os.WriteFile(jsonlFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	sessions, err := parser.ParseFile(jsonlFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}

	// Sessions should be sorted by start time
	if sessions[0].ID != "sess-A" {
		t.Errorf("expected first session 'sess-A', got '%s'", sessions[0].ID)
	}
	if sessions[1].ID != "sess-B" {
		t.Errorf("expected second session 'sess-B', got '%s'", sessions[1].ID)
	}
}

func TestClaudeCodeParser_ParseFile_SkipsMalformed(t *testing.T) {
	parser := NewClaudeCodeParser()
	dir := t.TempDir()

	// Mix of valid and invalid lines
	jsonlFile := filepath.Join(dir, "mixed.jsonl")
	content := `not json at all
{"uuid":"msg1","sessionId":"sess-1","slug":"test","type":"user","timestamp":"2026-01-20T10:00:00Z","cwd":"/test","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"valid"}]}}
{"incomplete": true
{"uuid":"msg2","sessionId":"sess-1","slug":"test","type":"assistant","timestamp":"2026-01-20T10:00:30Z","cwd":"/test","version":"2.1.0","message":{"role":"assistant","content":[{"type":"text","text":"also valid"}]}}`

	if err := os.WriteFile(jsonlFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	sessions, err := parser.ParseFile(jsonlFile)
	if err != nil {
		t.Fatalf("ParseFile should not fail on malformed lines: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if len(sessions[0].Messages) != 2 {
		t.Errorf("expected 2 valid messages, got %d", len(sessions[0].Messages))
	}
}

func TestScanDirectory(t *testing.T) {
	dir := t.TempDir()

	// Create subdirectory structure
	subdir := filepath.Join(dir, "subdir")
	if err := os.MkdirAll(subdir, 0750); err != nil {
		t.Fatal(err)
	}

	// Create session files
	file1 := filepath.Join(dir, "session1.jsonl")
	content1 := `{"uuid":"m1","sessionId":"s1","slug":"first","type":"user","timestamp":"2026-01-20T08:00:00Z","cwd":"/test","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"first"}]}}`
	if err := os.WriteFile(file1, []byte(content1), 0600); err != nil {
		t.Fatal(err)
	}

	file2 := filepath.Join(subdir, "session2.jsonl")
	content2 := `{"uuid":"m2","sessionId":"s2","slug":"second","type":"user","timestamp":"2026-01-20T10:00:00Z","cwd":"/test","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"second"}]}}`
	if err := os.WriteFile(file2, []byte(content2), 0600); err != nil {
		t.Fatal(err)
	}

	// Non-session file (should be ignored)
	nonSession := filepath.Join(dir, "readme.txt")
	if err := os.WriteFile(nonSession, []byte("ignore me"), 0600); err != nil {
		t.Fatal(err)
	}

	sessions, err := ScanDirectory(dir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}

	// Should be sorted newest first
	if sessions[0].ID != "s2" {
		t.Errorf("expected newest session first (s2), got %s", sessions[0].ID)
	}
	if sessions[1].ID != "s1" {
		t.Errorf("expected oldest session last (s1), got %s", sessions[1].ID)
	}
}

func TestParseFile_AutoDetect(t *testing.T) {
	dir := t.TempDir()

	jsonlFile := filepath.Join(dir, "auto.jsonl")
	content := `{"uuid":"m1","sessionId":"s1","slug":"auto","type":"user","timestamp":"2026-01-20T10:00:00Z","cwd":"/test","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"auto detect"}]}}`
	if err := os.WriteFile(jsonlFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	sessions, err := ParseFile(jsonlFile)
	if err != nil {
		t.Fatalf("ParseFile (auto-detect) failed: %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Tool != "claude-code" {
		t.Errorf("expected tool 'claude-code', got '%s'", sessions[0].Tool)
	}
}

func TestRegisteredTools(t *testing.T) {
	tools := RegisteredTools()
	if len(tools) == 0 {
		t.Error("expected at least one registered tool")
	}

	found := false
	for _, tool := range tools {
		if tool == "claude-code" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'claude-code' in registered tools")
	}
}

func TestGetParser(t *testing.T) {
	parser := Parser("claude-code")
	if parser == nil {
		t.Error("expected parser for 'claude-code'")
	}
	if parser.Tool() != "claude-code" {
		t.Errorf("expected tool 'claude-code', got '%s'", parser.Tool())
	}

	unknown := Parser("unknown-tool")
	if unknown != nil {
		t.Error("expected nil for unknown tool")
	}
}

func TestFindSessions_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	sessions, err := FindSessions()
	if err != nil {
		t.Fatalf("FindSessions failed: %v", err)
	}

	t.Logf("Found %d sessions", len(sessions))

	for i, s := range sessions {
		if i >= 3 {
			t.Logf("... and %d more", len(sessions)-3)
			break
		}
		preview := s.FirstUserMsg
		if len(preview) > 50 {
			preview = preview[:50] + "..."
		}
		t.Logf("%d. %s | %s | %d turns | %s", i+1, s.Slug, s.StartTime.Format("2006-01-02"), s.TurnCount, preview)
	}
}

func TestDebugSession(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	sessions, _ := FindSessions()
	for _, s := range sessions {
		if strings.HasPrefix(s.ID, "4c050ec4") {
			t.Logf("Session: %s", s.ID)
			t.Logf("Messages: %d", len(s.Messages))
			for i, m := range s.Messages {
				if i > 5 {
					break
				}
				t.Logf("  %d. %s: text=%d chars, tools=%d", i, m.Role, len(m.Text), len(m.ToolUses))
				if len(m.ToolUses) > 0 {
					t.Logf("      tool: %s, input: %.100s", m.ToolUses[0].Name, m.ToolUses[0].Input)
				}
			}
			break
		}
	}
}
