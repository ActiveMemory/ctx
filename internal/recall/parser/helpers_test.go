//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"testing"
)

func TestMessage_UsesTools(t *testing.T) {
	tests := []struct {
		name string
		msg  Message
		want bool
	}{
		{"empty message", Message{}, false},
		{"text only", Message{Text: "hello"}, false},
		{"with tool uses", Message{ToolUses: []ToolUse{{Name: "Bash"}}}, true},
		{"multiple tools", Message{ToolUses: []ToolUse{{Name: "Read"}, {Name: "Write"}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msg.UsesTools()
			if got != tt.want {
				t.Errorf("UsesTools() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_Preview(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		maxLen int
		want   string
	}{
		{"short text", "hello", 10, "hello"},
		{"exact length", "hello", 5, "hello"},
		{"truncated", "hello world foo bar", 10, "hello worl..."},
		{"empty", "", 10, ""},
		{"zero maxLen", "hello", 0, "..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{Text: tt.text}
			got := msg.Preview(tt.maxLen)
			if got != tt.want {
				t.Errorf("Preview(%d) = %q, want %q", tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestSession_UserMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []Message
		want     int
	}{
		{"empty session", nil, 0},
		{"all user", []Message{{Role: "user"}, {Role: "user"}}, 2},
		{"mixed roles", []Message{
			{Role: "user"},
			{Role: "assistant"},
			{Role: "user"},
			{Role: "assistant"},
		}, 2},
		{"all assistant", []Message{{Role: "assistant"}, {Role: "assistant"}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{Messages: tt.messages}
			got := s.UserMessages()
			if len(got) != tt.want {
				t.Errorf("UserMessages() returned %d, want %d", len(got), tt.want)
			}
			for _, m := range got {
				if m.Role != "user" {
					t.Errorf("UserMessages() returned message with role %q", m.Role)
				}
			}
		})
	}
}

func TestSession_AssistantMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []Message
		want     int
	}{
		{"empty session", nil, 0},
		{"all assistant", []Message{{Role: "assistant"}, {Role: "assistant"}}, 2},
		{"mixed roles", []Message{
			{Role: "user"},
			{Role: "assistant"},
			{Role: "user"},
			{Role: "assistant"},
		}, 2},
		{"all user", []Message{{Role: "user"}, {Role: "user"}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{Messages: tt.messages}
			got := s.AssistantMessages()
			if len(got) != tt.want {
				t.Errorf("AssistantMessages() returned %d, want %d", len(got), tt.want)
			}
			for _, m := range got {
				if m.Role != "assistant" {
					t.Errorf("AssistantMessages() returned message with role %q", m.Role)
				}
			}
		})
	}
}

func TestGetPathRelativeToHome(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"Linux path", "/home/jose/projects/ctx", "projects/ctx"},
		{"macOS path", "/Users/jose/projects/ctx", "projects/ctx"},
		{"not under home", "/var/log/syslog", ""},
		{"empty path", "", ""},
		{"home only Linux", "/home/jose", ""},
		{"home only macOS", "/Users/jose", ""},
		{"deep path", "/home/user/a/b/c/d", "a/b/c/d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPathRelativeToHome(tt.path)
			if got != tt.want {
				t.Errorf("getPathRelativeToHome(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	p := NewClaudeCodeParser()

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
			name: "non-message type",
			line: `{"uuid":"x","sessionId":"s1","type":"file-history-snapshot","timestamp":"2026-01-20T10:00:00Z","cwd":"/test","message":{"role":"user","content":[]}}`,
		},
		{
			name:     "valid user message",
			line:     `{"uuid":"m1","sessionId":"sess-1","slug":"test","type":"user","timestamp":"2026-01-20T10:00:00Z","cwd":"/test","version":"2.1.0","message":{"role":"user","content":[{"type":"text","text":"hello world"}]}}`,
			wantMsg:  true,
			wantSess: "sess-1",
			wantRole: "user",
			wantText: "hello world",
		},
		{
			name:     "valid assistant message",
			line:     `{"uuid":"m2","sessionId":"sess-1","slug":"test","type":"assistant","timestamp":"2026-01-20T10:00:30Z","cwd":"/test","version":"2.1.0","message":{"role":"assistant","content":[{"type":"text","text":"I can help"}]}}`,
			wantMsg:  true,
			wantSess: "sess-1",
			wantRole: "assistant",
			wantText: "I can help",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, sessID, err := p.ParseLine([]byte(tt.line))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantMsg {
				if msg == nil {
					t.Fatal("ParseLine() returned nil message, want non-nil")
				}
				if sessID != tt.wantSess {
					t.Errorf("sessionID = %q, want %q", sessID, tt.wantSess)
				}
				if msg.Role != tt.wantRole {
					t.Errorf("msg.Role = %q, want %q", msg.Role, tt.wantRole)
				}
				if msg.Text != tt.wantText {
					t.Errorf("msg.Text = %q, want %q", msg.Text, tt.wantText)
				}
			} else if !tt.wantErr {
				if msg != nil {
					t.Errorf("ParseLine() returned non-nil message, want nil")
				}
			}
		})
	}
}
