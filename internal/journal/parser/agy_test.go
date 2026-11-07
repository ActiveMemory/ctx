//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	cfgAgy "github.com/ActiveMemory/ctx/internal/config/agy"
	"github.com/ActiveMemory/ctx/internal/config/session"
)

func TestAgy_Tool(t *testing.T) {
	p := NewAgy()
	if got := p.Tool(); got != session.ToolAgy {
		t.Errorf("Tool() = %q, want %q", got, session.ToolAgy)
	}
}

func TestAgy_Matches(t *testing.T) {
	p := NewAgy()
	tempDir := t.TempDir()

	// 1. Valid transcript.jsonl with USER_INPUT
	validPath := filepath.Join(tempDir, cfgAgy.FileTranscriptJSONL)
	userStep := `{"step_index":0,"source":"USER_EXPLICIT","type":"USER_INPUT",` +
		`"status":"DONE","created_at":"2026-11-06T23:30:04Z",` +
		`"content":"hello"}` + "\n"
	if err := os.WriteFile(validPath, []byte(userStep), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if !p.Matches(validPath) {
		t.Errorf("Matches(%q) = false, want true", validPath)
	}

	// 2. Reject transcript_full.jsonl even with valid content
	fullPath := filepath.Join(tempDir, cfgAgy.FileTranscriptFullJSONL)
	if err := os.WriteFile(fullPath, []byte(userStep), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if p.Matches(fullPath) {
		t.Errorf("Matches(%q) = true, want false (full transcript rejected)", fullPath)
	}

	// 3. Reject non-.jsonl extension
	txtPath := filepath.Join(tempDir, "transcript.txt")
	if err := os.WriteFile(txtPath, []byte(userStep), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if p.Matches(txtPath) {
		t.Errorf("Matches(%q) = true, want false for non-jsonl", txtPath)
	}

	// 4. Reject Claude Code JSONL file
	claudePath := filepath.Join(tempDir, "claude.jsonl")
	claudeLine := `{"parentUuid":"1","sessionId":"s1","type":"user",` +
		`"message":{"role":"user","content":"hi"}}` + "\n"
	if err := os.WriteFile(claudePath, []byte(claudeLine), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if p.Matches(claudePath) {
		t.Errorf("Matches(%q) = true, want false for Claude file", claudePath)
	}

	// 5. Reject Codex rollout JSONL file
	codexPath := filepath.Join(tempDir, "rollout-test.jsonl")
	codexLine := `{"timestamp":"2026-08-23T19:02:47Z","type":"session_meta",` +
		`"payload":{"id":"cx-1","cwd":"/dir"}}` + "\n"
	if err := os.WriteFile(codexPath, []byte(codexLine), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if p.Matches(codexPath) {
		t.Errorf("Matches(%q) = true, want false for Codex file", codexPath)
	}

	// 6. Reject non-existent and empty files
	if p.Matches(filepath.Join(tempDir, "nonexistent.jsonl")) {
		t.Error("Matches(nonexistent) = true, want false")
	}
	emptyPath := filepath.Join(tempDir, "empty.jsonl")
	if err := os.WriteFile(emptyPath, []byte(""), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if p.Matches(emptyPath) {
		t.Error("Matches(empty) = true, want false")
	}
}

func TestAgy_ParseFile(t *testing.T) {
	p := NewAgy()
	tempDir := t.TempDir()

	// Build directory tree:
	// <tempDir>/antigravity-cli/history.jsonl
	// <tempDir>/antigravity-cli/brain/<id>/.system_generated/logs/transcript.jsonl
	convID := "2526ca85-83ee-4a91-b4ea-6be6c421d65b"
	agyCliDir := filepath.Join(tempDir, "antigravity-cli")
	logsDir := filepath.Join(
		agyCliDir, "brain", convID, ".system_generated", "logs",
	)
	if err := os.MkdirAll(logsDir, 0o750); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	wsDir := filepath.Join(tempDir, "my-workspace")
	if err := os.MkdirAll(wsDir, 0o750); err != nil {
		t.Fatalf("MkdirAll(wsDir): %v", err)
	}

	// Write history.jsonl
	histPath := filepath.Join(agyCliDir, cfgAgy.FileHistoryJSONL)
	histContent := `{"conversationId":"other-id","workspace":"/other"}` + "\n" +
		`{"conversationId":"` + convID + `","workspace":"` +
		strings.ReplaceAll(wsDir, "\\", "\\\\") + `"}` + "\n"
	if err := os.WriteFile(histPath, []byte(histContent), 0o600); err != nil {
		t.Fatalf("WriteFile(history): %v", err)
	}

	// Write transcript.jsonl with:
	// - USER_INPUT with <USER_REQUEST> and <ADDITIONAL_METADATA>
	// - System notification with <SYSTEM_MESSAGE> (should be skipped)
	// - PLANNER_RESPONSE with content, thinking, tool_calls, tokens
	// - Partial trailing line (tolerated)
	transcriptPath := filepath.Join(logsDir, cfgAgy.FileTranscriptJSONL)
	lines := []string{
		`{"step_index":0,"source":"USER_EXPLICIT","type":"USER_INPUT",` +
			`"status":"DONE","created_at":"2026-11-06T23:30:00Z",` +
			`"content":"<USER_REQUEST>\nFix the broken parser\n</USER_REQUEST>\n` +
			`<ADDITIONAL_METADATA>\nLocal time: 2026-11-06T23:30:00Z\n</ADDITIONAL_METADATA>"}`,
		`{"step_index":1,"source":"SYSTEM","type":"USER_INPUT",` +
			`"status":"DONE","created_at":"2026-11-06T23:30:02Z",` +
			`"content":"The following is a <SYSTEM_MESSAGE> notification\n` +
			`<SYSTEM_MESSAGE>internal</SYSTEM_MESSAGE>"}`,
		`{"step_index":2,"source":"MODEL","type":"PLANNER_RESPONSE",` +
			`"status":"DONE","created_at":"2026-11-06T23:31:00Z",` +
			`"input_tokens":500,"cache_read_tokens":100,"output_tokens":80,` +
			`"content":"I fixed the parser.","thinking":"Reviewing code...",` +
			`"tool_calls":[{"name":"edit_file","args":{"Path":"parser.go"}}]}`,
		`{"step_index":3,"source":"USER_EXPLICIT","type":"USER_IN`, // Partial trailing line
	}
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(transcriptPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(transcript): %v", err)
	}

	sessions, err := p.ParseFile(transcriptPath)
	if err != nil {
		t.Fatalf("ParseFile(%q): %v", transcriptPath, err)
	}
	if len(sessions) != 1 {
		t.Fatalf("len(sessions) = %d, want 1", len(sessions))
	}

	s := sessions[0]
	if s.ID != convID {
		t.Errorf("ID = %q, want %q", s.ID, convID)
	}
	if s.CWD != filepath.Clean(wsDir) {
		t.Errorf("CWD = %q, want %q", s.CWD, filepath.Clean(wsDir))
	}
	if s.Tool != session.ToolAgy {
		t.Errorf("Tool = %q, want %q", s.Tool, session.ToolAgy)
	}
	if s.TurnCount != 1 {
		t.Errorf("TurnCount = %d, want 1 (system msg was skipped)", s.TurnCount)
	}
	if s.FirstUserMsg != "Fix the broken parser" {
		t.Errorf("FirstUserMsg = %q, want %q", s.FirstUserMsg, "Fix the broken parser")
	}
	if s.TotalTokensIn != 600 {
		t.Errorf("TotalTokensIn = %d, want 600 (500+100)", s.TotalTokensIn)
	}
	if s.TotalTokensOut != 80 {
		t.Errorf("TotalTokensOut = %d, want 80", s.TotalTokensOut)
	}
	if s.TotalTokens != 680 {
		t.Errorf("TotalTokens = %d, want 680", s.TotalTokens)
	}

	wantStart, _ := time.Parse(time.RFC3339, "2026-11-06T23:30:00Z")
	wantEnd, _ := time.Parse(time.RFC3339, "2026-11-06T23:31:00Z")
	if !s.StartTime.Equal(wantStart) {
		t.Errorf("StartTime = %v, want %v", s.StartTime, wantStart)
	}
	if !s.EndTime.Equal(wantEnd) {
		t.Errorf("EndTime = %v, want %v", s.EndTime, wantEnd)
	}
	if s.Duration != time.Minute {
		t.Errorf("Duration = %v, want 1m", s.Duration)
	}

	if len(s.Messages) != 2 {
		t.Fatalf("len(Messages) = %d, want 2", len(s.Messages))
	}

	// Verify User Message
	uMsg := s.Messages[0]
	if uMsg.ID != "0" {
		t.Errorf("user msg ID = %q, want 0", uMsg.ID)
	}
	if uMsg.Role != session.RoleUser {
		t.Errorf("user msg Role = %q, want user", uMsg.Role)
	}
	if uMsg.Text != "Fix the broken parser" {
		t.Errorf("user msg Text = %q, want 'Fix the broken parser'", uMsg.Text)
	}

	// Verify Assistant Message
	aMsg := s.Messages[1]
	if aMsg.ID != "2" {
		t.Errorf("assistant msg ID = %q, want 2", aMsg.ID)
	}
	if aMsg.Role != session.RoleAssistant {
		t.Errorf("assistant msg Role = %q, want assistant", aMsg.Role)
	}
	if aMsg.Text != "I fixed the parser." {
		t.Errorf("assistant msg Text = %q", aMsg.Text)
	}
	if aMsg.Thinking != "Reviewing code..." {
		t.Errorf("assistant msg Thinking = %q", aMsg.Thinking)
	}
	if aMsg.TokensIn != 600 {
		t.Errorf("assistant msg TokensIn = %d, want 600", aMsg.TokensIn)
	}
	if aMsg.TokensOut != 80 {
		t.Errorf("assistant msg TokensOut = %d, want 80", aMsg.TokensOut)
	}
	if len(aMsg.ToolUses) != 1 {
		t.Fatalf("len(ToolUses) = %d, want 1", len(aMsg.ToolUses))
	}
	if aMsg.ToolUses[0].Name != "edit_file" {
		t.Errorf("ToolUse Name = %q, want edit_file", aMsg.ToolUses[0].Name)
	}
}

func TestAgy_ParseFile_ZeroProse(t *testing.T) {
	p := NewAgy()
	tempDir := t.TempDir()
	transcriptPath := filepath.Join(tempDir, cfgAgy.FileTranscriptJSONL)

	// Transcript containing only a system message (no user prose)
	content := `{"step_index":0,"source":"SYSTEM","type":"USER_INPUT",` +
		`"status":"DONE","created_at":"2026-11-06T23:30:00Z",` +
		`"content":"<SYSTEM_MESSAGE>notification only</SYSTEM_MESSAGE>"}` + "\n"
	if err := os.WriteFile(transcriptPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	sessions, err := p.ParseFile(transcriptPath)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if sessions != nil {
		t.Errorf("ParseFile with zero prose = %v, want nil", sessions)
	}
}

func TestAgy_ParseFile_URIWorkspace(t *testing.T) {
	p := NewAgy()
	tempDir := t.TempDir()

	convID := "test-uri-conv"
	brainDir := filepath.Join(
		tempDir, "antigravity-cli", "brain", convID, ".system_generated", "logs",
	)
	if err := os.MkdirAll(brainDir, 0o750); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	histPath := filepath.Join(tempDir, "antigravity-cli", cfgAgy.FileHistoryJSONL)
	histContent := `{"conversationId":"` + convID +
		`","workspace":"file:///projects/my-test"}` + "\n"
	if err := os.WriteFile(histPath, []byte(histContent), 0o600); err != nil {
		t.Fatalf("WriteFile(history): %v", err)
	}

	transcriptPath := filepath.Join(brainDir, cfgAgy.FileTranscriptJSONL)
	content := `{"step_index":0,"source":"USER_EXPLICIT","type":"USER_INPUT",` +
		`"status":"DONE","created_at":"2026-11-06T23:30:00Z","content":"hello"}` + "\n" +
		`{"step_index":1,"source":"MODEL","type":"PLANNER_RESPONSE",` +
		`"status":"DONE","created_at":"2026-11-06T23:30:01Z","content":"world"}` + "\n"
	if err := os.WriteFile(transcriptPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(transcript): %v", err)
	}

	sessions, err := p.ParseFile(transcriptPath)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("len(sessions) = %d, want 1", len(sessions))
	}

	if sessions[0].CWD == "" {
		t.Errorf("CWD was not resolved from file:// URI")
	}
}

func TestAgy_ParseLine(t *testing.T) {
	p := NewAgy()

	// 1. Empty line
	msg, id, err := p.ParseLine(nil)
	if err != nil || msg != nil || id != "" {
		t.Errorf("ParseLine(nil) = (%v, %q, %v), want (nil, '', nil)", msg, id, err)
	}

	// 2. USER_INPUT line
	userLine := []byte(`{"step_index":5,"source":"USER_EXPLICIT","type":"USER_INPUT",` +
		`"status":"DONE","created_at":"2026-11-06T23:30:00Z",` +
		`"content":"<USER_REQUEST>run tests</USER_REQUEST>"}`)
	msg, id, err = p.ParseLine(userLine)
	if err != nil {
		t.Fatalf("ParseLine(user): %v", err)
	}
	if msg == nil || msg.Text != "run tests" || msg.Role != session.RoleUser || id != "" {
		t.Errorf("ParseLine(user) = %+v, id=%q", msg, id)
	}

	// 3. PLANNER_RESPONSE line
	plannerLine := []byte(`{"step_index":6,"source":"MODEL","type":"PLANNER_RESPONSE",` +
		`"status":"DONE","created_at":"2026-11-06T23:31:00Z",` +
		`"input_tokens":100,"output_tokens":50,"content":"done","thinking":"thought",` +
		`"tool_calls":[{"name":"run_command","args":{"CommandLine":"go test"}}]}`)
	msg, id, err = p.ParseLine(plannerLine)
	if err != nil {
		t.Fatalf("ParseLine(planner): %v", err)
	}
	if msg == nil || msg.Text != "done" || msg.Role != session.RoleAssistant {
		t.Errorf("ParseLine(planner) = %+v", msg)
	}
	if msg.TokensIn != 100 || msg.TokensOut != 50 || len(msg.ToolUses) != 1 {
		t.Errorf("ParseLine(planner) tokens/tools unexpected: %+v", msg)
	}

	// 4. GENERIC line
	genericLine := []byte(`{"step_index":7,"source":"MODEL","type":"GENERIC",` +
		`"status":"DONE","created_at":"2026-11-06T23:32:00Z","content":"artifact"}`)
	msg, _, err = p.ParseLine(genericLine)
	if err != nil || msg != nil {
		t.Errorf("ParseLine(generic) = (%v, %v), want (nil, nil)", msg, err)
	}
}

func TestAgy_ScanDirectory_And_FindSessionsForCWD(t *testing.T) {
	tempDir := t.TempDir()
	convID := "agy-scan-test-conv"

	wsDir := filepath.Join(tempDir, "workspace-target")
	if err := os.MkdirAll(wsDir, 0o750); err != nil {
		t.Fatalf("MkdirAll(wsDir): %v", err)
	}

	agyHome := filepath.Join(tempDir, "gemini-home")
	logsDir := filepath.Join(
		agyHome, cfgAgy.DirAntigravityCLI, cfgAgy.DirBrain,
		convID, cfgAgy.DirSystemGenerated, cfgAgy.DirLogs,
	)
	if err := os.MkdirAll(logsDir, 0o750); err != nil {
		t.Fatalf("MkdirAll(logsDir): %v", err)
	}

	histPath := filepath.Join(
		agyHome, cfgAgy.DirAntigravityCLI, cfgAgy.FileHistoryJSONL,
	)
	histLine := `{"conversationId":"` + convID + `","workspace":"` +
		strings.ReplaceAll(wsDir, "\\", "\\\\") + `"}` + "\n"
	if err := os.WriteFile(histPath, []byte(histLine), 0o600); err != nil {
		t.Fatalf("WriteFile(history): %v", err)
	}

	transcriptPath := filepath.Join(logsDir, cfgAgy.FileTranscriptJSONL)
	transcriptContent := `{"step_index":0,"source":"USER_EXPLICIT","type":"USER_INPUT",` +
		`"status":"DONE","created_at":"2026-11-06T12:00:00Z","content":"Do the work"}` + "\n" +
		`{"step_index":1,"source":"MODEL","type":"PLANNER_RESPONSE",` +
		`"status":"DONE","created_at":"2026-11-06T12:01:00Z","input_tokens":40,` +
		`"output_tokens":20,"content":"Work complete"}` + "\n"
	if err := os.WriteFile(transcriptPath, []byte(transcriptContent), 0o600); err != nil {
		t.Fatalf("WriteFile(transcript): %v", err)
	}

	// 1. ScanDirectory
	brainDir := filepath.Join(agyHome, cfgAgy.DirAntigravityCLI, cfgAgy.DirBrain)
	sessions, err := ScanDirectory(brainDir)
	if err != nil {
		t.Fatalf("ScanDirectory(%q): %v", brainDir, err)
	}
	if len(sessions) != 1 {
		t.Fatalf("len(sessions) = %d, want 1", len(sessions))
	}
	if sessions[0].ID != convID {
		t.Errorf("session.ID = %q, want %q", sessions[0].ID, convID)
	}
	if sessions[0].CWD != filepath.Clean(wsDir) {
		t.Errorf("session.CWD = %q, want %q", sessions[0].CWD, filepath.Clean(wsDir))
	}

	// 2. FindSessionsForCWD via AGY_HOME
	t.Setenv("USERPROFILE", tempDir)
	t.Setenv("HOME", tempDir)
	t.Setenv("APPDATA", tempDir)
	t.Setenv(cfgAgy.EnvHome, agyHome)

	cwdSessions, cwdErr := FindSessionsForCWD(filepath.Clean(wsDir))
	if cwdErr != nil {
		t.Fatalf("FindSessionsForCWD: %v", cwdErr)
	}
	if len(cwdSessions) != 1 {
		t.Fatalf("len(cwdSessions) = %d, want 1", len(cwdSessions))
	}
	if cwdSessions[0].ID != convID {
		t.Errorf("cwdSessions[0].ID = %q, want %q", cwdSessions[0].ID, convID)
	}
}

func TestAgySessionDirs(t *testing.T) {
	tempDir := t.TempDir()

	// Case 1: uninitialized/missing brain dir
	t.Setenv(cfgAgy.EnvHome, tempDir)
	if got := AgySessionDirs(); got != nil {
		t.Errorf("AgySessionDirs() = %v, want nil when dir missing", got)
	}

	// Case 2: existing brain dir
	brainDir := filepath.Join(tempDir, cfgAgy.DirAntigravityCLI, cfgAgy.DirBrain)
	if err := os.MkdirAll(brainDir, 0o750); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	got := AgySessionDirs()
	if len(got) != 1 || got[0] != brainDir {
		t.Errorf("AgySessionDirs() = %v, want [%s]", got, brainDir)
	}
}
