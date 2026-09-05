//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/entity"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

const (
	// codexFixtureProbe is the sanitized real rollout captured from
	// codex 0.148.0 (`codex exec`, one tool call, one reply).
	codexFixtureProbe = "rollout-2026-08-23T12-02-47-" +
		"01a03001-0374-7603-8c61-787b22366b4f.jsonl"
	// codexFixtureInjectedOnly carries developer and Codex-injected
	// user items but no user prose.
	codexFixtureInjectedOnly = "rollout-2026-08-23T13-00-00-" +
		"11111111-0000-4000-8000-000000000001.jsonl"
	// codexFixtureMalformed mixes function_call / local_shell_call
	// items, a compacted line, and two malformed lines.
	codexFixtureMalformed = "rollout-2026-08-23T14-00-00-" +
		"22222222-0000-4000-8000-000000000002.jsonl"

	// claudeFixtureValid is the Claude Code transcript used by the
	// schema validator tests; a negative case for Codex matching.
	claudeFixtureValid = "valid.jsonl"
)

// codexFixture returns the path of a Codex fixture under
// testdata/codex.
func codexFixture(name string) string {
	return filepath.Join("testdata", "codex", name)
}

// claudeFixture returns the path of a Claude Code fixture from the
// sibling schema package's testdata.
func claudeFixture(name string) string {
	return filepath.Join("..", "schema", "testdata", name)
}

// copyCodexFixture copies a Codex fixture to dst so tests can exercise
// matching and scanning under a different name or directory.
func copyCodexFixture(t *testing.T, name, dst string) {
	t.Helper()
	data, readErr := ctxIo.SafeReadUserFile(codexFixture(name))
	if readErr != nil {
		t.Fatal(readErr)
	}
	if writeErr := ctxIo.SafeWriteFile(dst, data, 0600); writeErr != nil {
		t.Fatal(writeErr)
	}
}

// parseCodexFixture parses a fixture and fails the test unless exactly
// one session comes back.
func parseCodexFixture(t *testing.T, name string) *entity.Session {
	t.Helper()
	sessions, parseErr := NewCodex().ParseFile(codexFixture(name))
	if parseErr != nil {
		t.Fatalf("ParseFile(%s): %v", name, parseErr)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session from %s, got %d", name, len(sessions))
	}
	return sessions[0]
}

// textMessages returns the messages of the given role that carry
// prose (tool-use and tool-result messages carry none).
func textMessages(s *entity.Session, role string) []entity.Message {
	var out []entity.Message
	for _, m := range s.Messages {
		if m.Role == role && m.Text != "" {
			out = append(out, m)
		}
	}
	return out
}

func TestCodexParser_Matches(t *testing.T) {
	p := NewCodex()

	if !p.Matches(codexFixture(codexFixtureProbe)) {
		t.Error("should match a rollout-*.jsonl fixture")
	}

	// The same content under a non-rollout name is matched by sniffing
	// the leading session_meta line.
	dir := t.TempDir()
	renamed := filepath.Join(dir, "session.jsonl")
	copyCodexFixture(t, codexFixtureProbe, renamed)
	if !p.Matches(renamed) {
		t.Error("should match a renamed rollout via session_meta sniffing")
	}

	// Negative: Claude Code transcript.
	if p.Matches(claudeFixture(claudeFixtureValid)) {
		t.Error("should not match a Claude Code transcript")
	}
	if !NewClaudeCode().Matches(claudeFixture(claudeFixtureValid)) {
		t.Fatal("sanity: Claude parser should match its own fixture")
	}

	// Negative: Markdown file.
	mdFile := filepath.Join(dir, "session.md")
	if writeErr := os.WriteFile(
		mdFile, []byte("# Session: 2026-08-23 - Topic"), 0600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	if p.Matches(mdFile) {
		t.Error("should not match a Markdown file")
	}

	// Negative: JSONL that is neither a rollout nor session_meta-led.
	other := filepath.Join(dir, "other.jsonl")
	if writeErr := os.WriteFile(
		other, []byte(`{"foo":"bar"}`), 0600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	if p.Matches(other) {
		t.Error("should not match arbitrary JSONL")
	}
}

func TestCodexParser_ParseFile_Probe(t *testing.T) {
	s := parseCodexFixture(t, codexFixtureProbe)

	if s.Tool != session.ToolCodex {
		t.Errorf("Tool = %q, want %q", s.Tool, session.ToolCodex)
	}
	if s.ID != "01a03001-0374-7603-8c61-787b22366b4f" {
		t.Errorf("ID = %q", s.ID)
	}
	if s.CWD != "/home/user/projects/probe" {
		t.Errorf("CWD = %q", s.CWD)
	}
	if s.Project != "probe" {
		t.Errorf("Project = %q", s.Project)
	}
	if s.GitBranch != "main" {
		t.Errorf("GitBranch = %q", s.GitBranch)
	}
	if s.Model != "gpt-5.6-sol" {
		t.Errorf("Model = %q", s.Model)
	}
	if s.Entrypoint != "codex_exec" {
		t.Errorf("Entrypoint = %q", s.Entrypoint)
	}
	if s.SourceFile != codexFixture(codexFixtureProbe) {
		t.Errorf("SourceFile = %q", s.SourceFile)
	}

	wantStart := time.Date(2026, 8, 23, 19, 2, 47, 954_000_000, time.UTC)
	if !s.StartTime.Equal(wantStart) {
		t.Errorf("StartTime = %v, want %v", s.StartTime, wantStart)
	}
	wantEnd := time.Date(2026, 8, 23, 19, 2, 55, 692_000_000, time.UTC)
	if !s.EndTime.Equal(wantEnd) {
		t.Errorf("EndTime = %v, want %v", s.EndTime, wantEnd)
	}
	if s.Duration != wantEnd.Sub(wantStart) {
		t.Errorf("Duration = %v", s.Duration)
	}

	// Exactly one user prose message: the injected
	// <recommended_plugins>/<environment_context> item is dropped.
	wantPrompt := "Read README.md and reply with exactly the single word OK"
	users := textMessages(s, cfgCodex.RoleUser)
	if len(users) != 1 {
		t.Fatalf("expected 1 user message, got %d", len(users))
	}
	if users[0].Text != wantPrompt {
		t.Errorf("user text = %q", users[0].Text)
	}
	// 1 prose turn + 1 tool-result user message (Claude parity).
	if s.TurnCount != 2 {
		t.Errorf("TurnCount = %d, want 2", s.TurnCount)
	}
	if s.FirstUserMsg != wantPrompt {
		t.Errorf("FirstUserMsg = %q", s.FirstUserMsg)
	}

	// Exactly one assistant prose message.
	assistants := textMessages(s, cfgCodex.RoleAssistant)
	if len(assistants) != 1 {
		t.Fatalf("expected 1 assistant message, got %d", len(assistants))
	}
	if assistants[0].Text != "OK" {
		t.Errorf("assistant text = %q", assistants[0].Text)
	}

	// Tool use (custom_tool_call "exec") and its result.
	tools := s.AllToolUses()
	if len(tools) != 1 {
		t.Fatalf("expected 1 tool use, got %d", len(tools))
	}
	if tools[0].Name != "exec" {
		t.Errorf("tool name = %q, want exec", tools[0].Name)
	}
	if tools[0].ID != "call_MM3BLBRuv3UU8upNk9DNU7oA" {
		t.Errorf("tool call id = %q", tools[0].ID)
	}
	var results []entity.ToolResult
	for _, m := range s.Messages {
		results = append(results, m.ToolResults...)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 tool result, got %d", len(results))
	}
	if results[0].ToolUseID != tools[0].ID {
		t.Errorf("tool result id = %q, want %q", results[0].ToolUseID, tools[0].ID)
	}
	if results[0].Content != "Script completed\nWall time 1.0 seconds\nOutput:\nprobe\n" {
		t.Errorf("tool result content = %q", results[0].Content)
	}

	// Message order: user, tool-use, tool-result, assistant (developer,
	// injected, reasoning, and event lines contribute nothing).
	if len(s.Messages) != 4 {
		t.Fatalf("expected 4 messages, got %d", len(s.Messages))
	}
	if !s.Messages[0].BelongsToUser() || !s.Messages[1].UsesTools() ||
		len(s.Messages[2].ToolResults) != 1 ||
		!s.Messages[3].BelongsToAssistant() {
		t.Errorf("unexpected message order: %+v", s.Messages)
	}

	// Tokens come from the LAST cumulative token_count event.
	if s.TotalTokensIn != 32328-26112 {
		t.Errorf("TotalTokensIn = %d, want cache-adjusted 6216",
			s.TotalTokensIn)
	}
	if s.TotalTokensOut != 116 {
		t.Errorf("TotalTokensOut = %d, want 116", s.TotalTokensOut)
	}
	if s.TotalTokens != 6216+116 {
		t.Errorf("TotalTokens = %d", s.TotalTokens)
	}
	if s.HasErrors {
		t.Error("HasErrors should be false")
	}
}

func TestCodexParser_ParseFile_InjectedOnly(t *testing.T) {
	sessions, parseErr := NewCodex().ParseFile(
		codexFixture(codexFixtureInjectedOnly),
	)
	if parseErr != nil {
		t.Fatalf("ParseFile: %v", parseErr)
	}
	if sessions != nil {
		t.Errorf("expected nil session for injected-only rollout, got %+v", sessions)
	}
}

func TestCodexParser_ParseFile_Malformed(t *testing.T) {
	s := parseCodexFixture(t, codexFixtureMalformed)

	// session_id fallback when payload.id is absent.
	if s.ID != "22222222-0000-4000-8000-000000000002" {
		t.Errorf("ID = %q", s.ID)
	}
	if s.GitBranch != "feat/codex" {
		t.Errorf("GitBranch = %q", s.GitBranch)
	}
	if s.Model != "gpt-5.5" {
		t.Errorf("Model = %q", s.Model)
	}

	// The injected part is dropped, the prose part survives.
	users := textMessages(s, cfgCodex.RoleUser)
	if len(users) != 1 || users[0].Text != "List the files in this directory" {
		t.Fatalf("user messages = %+v", users)
	}

	// function_call and local_shell_call both become tool uses.
	tools := s.AllToolUses()
	if len(tools) != 2 {
		t.Fatalf("expected 2 tool uses, got %d", len(tools))
	}
	if tools[0].Name != "shell" || tools[0].ID != "call_fc_1" {
		t.Errorf("function_call tool = %+v", tools[0])
	}
	if tools[0].Input != `{"command":["ls","-la"],"workdir":"/home/user/projects/probe"}` {
		t.Errorf("function_call input = %q", tools[0].Input)
	}
	if tools[1].Name != cfgCodex.ItemTypeLocalShellCall || tools[1].ID != "call_ls_1" {
		t.Errorf("local_shell_call tool = %+v", tools[1])
	}
	if tools[1].Input != `{"type":"exec","command":["git","status","--short"],"timeout_ms":10000}` {
		t.Errorf("local_shell_call input = %q", tools[1].Input)
	}

	// String-shaped function_call_output.
	var results []entity.ToolResult
	for _, m := range s.Messages {
		results = append(results, m.ToolResults...)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 tool results, got %d", len(results))
	}
	if results[0].ToolUseID != "call_fc_1" ||
		results[0].Content != "total 0\ndrwxr-xr-x  2 user user  64 Aug 23 21:00 .\n" {
		t.Errorf("function_call_output = %+v", results[0])
	}
	if results[1].ToolUseID != "call_ls_1" {
		t.Errorf("second tool result = %+v", results[1])
	}

	assistants := textMessages(s, cfgCodex.RoleAssistant)
	if len(assistants) != 1 || assistants[0].Text != "The directory is empty." {
		t.Fatalf("assistant messages = %+v", assistants)
	}

	// user, tool-use, tool-result, tool-use, tool-result, assistant;
	// the two malformed lines, compacted, and reasoning are skipped.
	if len(s.Messages) != 6 {
		t.Errorf("expected 6 messages, got %d", len(s.Messages))
	}
	// 1 prose turn + 2 tool-result user messages (Claude parity).
	if s.TurnCount != 3 {
		t.Errorf("TurnCount = %d, want 3", s.TurnCount)
	}
	if s.TotalTokensIn != 800 || s.TotalTokensOut != 20 {
		t.Errorf("tokens = %d/%d, want cache-adjusted 800/20",
			s.TotalTokensIn, s.TotalTokensOut)
	}

	wantEnd := time.Date(2026, 8, 23, 21, 0, 10, 500_000_000, time.UTC)
	if !s.EndTime.Equal(wantEnd) {
		t.Errorf("EndTime = %v, want %v", s.EndTime, wantEnd)
	}
}

func TestCodexParser_ParseLine(t *testing.T) {
	p := NewCodex()

	tests := []struct {
		name     string
		line     string
		wantMsg  bool
		wantSess string
		wantErr  bool
		wantRole string
		wantText string
	}{
		{
			name: "empty line",
			line: "",
		},
		{
			name:    "invalid JSON",
			line:    "not json at all",
			wantErr: true,
		},
		{
			name:     "session_meta reports the session ID",
			line:     `{"timestamp":"2026-08-23T19:02:48.054Z","type":"session_meta","payload":{"id":"sess-1","cwd":"/tmp/x"}}`,
			wantSess: "sess-1",
		},
		{
			name:     "session_meta falls back to session_id",
			line:     `{"timestamp":"2026-08-23T19:02:48.054Z","type":"session_meta","payload":{"session_id":"sess-2","cwd":"/tmp/x"}}`,
			wantSess: "sess-2",
		},
		{
			name: "event_msg is skipped",
			line: `{"timestamp":"2026-08-23T19:02:48.054Z","type":"event_msg","payload":{"type":"task_started"}}`,
		},
		{
			name: "developer message is skipped",
			line: `{"timestamp":"2026-08-23T19:02:50.014Z","type":"response_item","payload":{"type":"message","id":"m1","role":"developer","content":[{"type":"input_text","text":"preamble"}]}}`,
		},
		{
			name: "injected user item is skipped",
			line: `{"timestamp":"2026-08-23T19:02:50.014Z","type":"response_item","payload":{"type":"message","id":"m2","role":"user","content":[{"type":"input_text","text":"  <environment_context>\nx\n</environment_context>"}]}}`,
		},
		{
			name:     "user message",
			line:     `{"timestamp":"2026-08-23T19:02:50.030Z","type":"response_item","payload":{"type":"message","id":"m3","role":"user","content":[{"type":"input_text","text":"hello"},{"type":"input_text","text":"world"}]}}`,
			wantMsg:  true,
			wantRole: cfgCodex.RoleUser,
			wantText: "hello\nworld",
		},
		{
			name:     "assistant message",
			line:     `{"timestamp":"2026-08-23T19:02:55.268Z","type":"response_item","payload":{"type":"message","id":"m4","role":"assistant","content":[{"type":"output_text","text":"OK"}]}}`,
			wantMsg:  true,
			wantRole: cfgCodex.RoleAssistant,
			wantText: "OK",
		},
		{
			name:     "function_call becomes a tool use",
			line:     `{"timestamp":"2026-08-23T19:02:53.343Z","type":"response_item","payload":{"type":"function_call","id":"fc","call_id":"c1","name":"shell","arguments":"{}"}}`,
			wantMsg:  true,
			wantRole: cfgCodex.RoleAssistant,
		},
		{
			name:     "function_call_output becomes a tool result",
			line:     `{"timestamp":"2026-08-23T19:02:54.338Z","type":"response_item","payload":{"type":"function_call_output","id":"fco","call_id":"c1","output":"done"}}`,
			wantMsg:  true,
			wantRole: cfgCodex.RoleUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, sessID, parseErr := p.ParseLine([]byte(tt.line))
			if (parseErr != nil) != tt.wantErr {
				t.Fatalf("ParseLine() error = %v, wantErr %v", parseErr, tt.wantErr)
			}
			if sessID != tt.wantSess {
				t.Errorf("sessionID = %q, want %q", sessID, tt.wantSess)
			}
			if !tt.wantMsg {
				if msg != nil {
					t.Errorf("ParseLine() returned %+v, want nil", msg)
				}
				return
			}
			if msg == nil {
				t.Fatal("ParseLine() returned nil message, want non-nil")
				return
			}
			if msg.Role != tt.wantRole {
				t.Errorf("msg.Role = %q, want %q", msg.Role, tt.wantRole)
			}
			if msg.Text != tt.wantText {
				t.Errorf("msg.Text = %q, want %q", msg.Text, tt.wantText)
			}
		})
	}
}

func TestCodexSessionDirs(t *testing.T) {
	// Absent: CODEX_HOME points at a directory with no sessions/.
	empty := t.TempDir()
	t.Setenv(cfgCodex.EnvHome, empty)
	if got := CodexSessionDirs(); got != nil {
		t.Errorf("expected nil for missing sessions dir, got %v", got)
	}

	// Present: CODEX_HOME/sessions exists.
	home := t.TempDir()
	sessions := filepath.Join(home, cfgCodex.DirSessions)
	if mkErr := os.MkdirAll(sessions, 0750); mkErr != nil {
		t.Fatal(mkErr)
	}
	t.Setenv(cfgCodex.EnvHome, home)
	got := CodexSessionDirs()
	if len(got) != 1 || got[0] != sessions {
		t.Errorf("CodexSessionDirs() = %v, want [%s]", got, sessions)
	}

	// A file at the sessions path is not a directory.
	fileHome := t.TempDir()
	if writeErr := os.WriteFile(
		filepath.Join(fileHome, cfgCodex.DirSessions), []byte("x"), 0600,
	); writeErr != nil {
		t.Fatal(writeErr)
	}
	t.Setenv(cfgCodex.EnvHome, fileHome)
	if got := CodexSessionDirs(); got != nil {
		t.Errorf("expected nil when sessions is a file, got %v", got)
	}
}

func TestCodexParser_RegistryDispatch(t *testing.T) {
	sessions, parseErr := ParseFile(codexFixture(codexFixtureProbe))
	if parseErr != nil {
		t.Fatalf("ParseFile (auto-detect): %v", parseErr)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Tool != session.ToolCodex {
		t.Errorf("Tool = %q, want %q", sessions[0].Tool, session.ToolCodex)
	}

	if find(session.ToolCodex) == nil {
		t.Error("expected Codex parser in registry")
	}

	// The Claude fixture still dispatches to the Claude parser.
	claudeSessions, claudeErr := ParseFile(claudeFixture(claudeFixtureValid))
	if claudeErr != nil {
		t.Fatalf("ParseFile (Claude): %v", claudeErr)
	}
	if len(claudeSessions) == 0 || claudeSessions[0].Tool != session.ToolClaudeCode {
		t.Errorf("Claude fixture dispatched to %+v", claudeSessions)
	}
}

func TestCodexParser_ScanNestedSessionsDir(t *testing.T) {
	// Rollouts nest as sessions/YYYY/MM/DD; ScanDirectory must find
	// them through CodexSessionDirs' single root.
	home := t.TempDir()
	day := filepath.Join(home, cfgCodex.DirSessions, "2026", "08", "23")
	if mkErr := os.MkdirAll(day, 0750); mkErr != nil {
		t.Fatal(mkErr)
	}
	copyCodexFixture(t, codexFixtureProbe, filepath.Join(day, codexFixtureProbe))
	t.Setenv(cfgCodex.EnvHome, home)

	dirs := CodexSessionDirs()
	if len(dirs) != 1 {
		t.Fatalf("CodexSessionDirs() = %v", dirs)
	}
	sessions, scanErr := ScanDirectory(dirs[0])
	if scanErr != nil {
		t.Fatalf("ScanDirectory: %v", scanErr)
	}
	if len(sessions) != 1 || sessions[0].Tool != session.ToolCodex {
		t.Errorf("ScanDirectory found %+v", sessions)
	}
}
