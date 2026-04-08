//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"encoding/json"
	"testing"
)

func TestParseOriginKind(t *testing.T) {
	tests := []struct {
		name string
		raw  json.RawMessage
		want string
	}{
		{
			"nil",
			nil,
			"",
		},
		{
			"task notification",
			json.RawMessage(`{"kind":"task-notification"}`),
			"task-notification",
		},
		{
			"empty object",
			json.RawMessage(`{}`),
			"",
		},
		{
			"invalid json",
			json.RawMessage(`not json`),
			"",
		},
		{
			"string value",
			json.RawMessage(`"just a string"`),
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseOriginKind(tt.raw)
			if got != tt.want {
				t.Errorf(
					"parseOriginKind() = %q, want %q",
					got, tt.want)
			}
		})
	}
}

func TestConvertMessageEnvelopeFields(t *testing.T) {
	p := NewClaudeCode()

	raw := claudeRawMessage{
		UUID:                    "test-uuid",
		Type:                    "user",
		PlanContent:             "# My Plan\n\nDo the thing.",
		IsApiErrorMessage:       false,
		SourceToolAssistantUUID: "parent-uuid-123",
		ToolUseResult:           "Error: EISDIR",
		Entrypoint:              "cli",
		Origin: json.RawMessage(
			`{"kind":"task-notification"}`),
		Message: claudeRawContent{
			Role:    "user",
			Content: json.RawMessage(`"hello"`),
		},
	}

	msg := p.convertMessage(raw)

	if msg.PlanContent != raw.PlanContent {
		t.Errorf("PlanContent = %q, want %q",
			msg.PlanContent, raw.PlanContent)
	}
	if msg.IsApiError {
		t.Error("IsApiError should be false")
	}
	if msg.SourceToolAssistantUUID != raw.SourceToolAssistantUUID {
		t.Errorf("SourceToolAssistantUUID = %q, want %q",
			msg.SourceToolAssistantUUID,
			raw.SourceToolAssistantUUID)
	}
	if msg.ToolUseResult != raw.ToolUseResult {
		t.Errorf("ToolUseResult = %q, want %q",
			msg.ToolUseResult, raw.ToolUseResult)
	}
	if msg.Origin != "task-notification" {
		t.Errorf("Origin = %q, want %q",
			msg.Origin, "task-notification")
	}
}

func TestConvertMessageApiError(t *testing.T) {
	p := NewClaudeCode()

	raw := claudeRawMessage{
		UUID:              "err-uuid",
		Type:              "assistant",
		IsApiErrorMessage: true,
		Message: claudeRawContent{
			Role:    "assistant",
			Content: json.RawMessage(`"rate limited"`),
		},
	}

	msg := p.convertMessage(raw)

	if !msg.IsApiError {
		t.Error("IsApiError should be true")
	}
}

func TestBuildSessionEntrypoint(t *testing.T) {
	p := NewClaudeCode()

	rawMsgs := []claudeRawMessage{
		{
			UUID:       "m1",
			SessionID:  "s1",
			Type:       "user",
			CWD:        "/test",
			Version:    "2.1.90",
			Entrypoint: "ide",
			Message: claudeRawContent{
				Role:    "user",
				Content: json.RawMessage(`"hello"`),
			},
		},
		{
			UUID:       "m2",
			SessionID:  "s1",
			Type:       "assistant",
			CWD:        "/test",
			Version:    "2.1.90",
			Entrypoint: "ide",
			Message: claudeRawContent{
				Role:    "assistant",
				Content: json.RawMessage(`"hi"`),
				Model:   "claude-sonnet-4-6",
			},
		},
	}

	sess := p.buildSession("s1", rawMsgs, "/test/s1.jsonl")
	if sess == nil {
		t.Fatal("expected session, got nil")
	}
	if sess.Entrypoint != "ide" {
		t.Errorf(
			"Entrypoint = %q, want %q",
			sess.Entrypoint, "ide")
	}
}
