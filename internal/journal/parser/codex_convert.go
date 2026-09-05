//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// applyLine folds one rollout line into the session under
// construction.
//
// Parameters:
//   - s: session being built
//   - raw: decoded line envelope
func (p *Codex) applyLine(s *entity.Session, raw codexRawLine) {
	switch raw.Type {
	case cfgCodex.LineTypeSessionMeta:
		p.applyMeta(s, raw.Payload)
	case cfgCodex.LineTypeTurnContext:
		var tc codexRawTurnContext
		if unmarshalErr := json.Unmarshal(
			raw.Payload, &tc,
		); unmarshalErr == nil && tc.Model != "" {
			s.Model = tc.Model
		}
	case cfgCodex.LineTypeEventMsg:
		p.applyEvent(s, raw.Payload)
	case cfgCodex.LineTypeResponseItem:
		p.appendMessage(s, p.convertItem(raw))
	}

	if raw.Timestamp.IsZero() {
		return
	}
	if s.StartTime.IsZero() {
		s.StartTime = raw.Timestamp
	}
	s.EndTime = raw.Timestamp
}

// applyMeta seeds session identity and context from a session_meta
// payload.
//
// Parameters:
//   - s: session being built
//   - payload: raw session_meta payload
func (p *Codex) applyMeta(s *entity.Session, payload json.RawMessage) {
	var meta codexRawSessionMeta
	if unmarshalErr := json.Unmarshal(payload, &meta); unmarshalErr != nil {
		return
	}

	s.ID = codexSessionID(meta)
	s.CWD = meta.CWD
	if meta.CWD != "" {
		s.Project = filepath.Base(meta.CWD)
	}
	s.Entrypoint = meta.Originator
	if meta.Git != nil {
		s.GitBranch = meta.Git.Branch
	}
	if !meta.Timestamp.IsZero() {
		s.StartTime = meta.Timestamp
	}
}

// applyEvent folds an event_msg payload into the session. Only
// token_count events carry session-level data; the usage block is
// cumulative, so the latest event overwrites the totals.
//
// Parameters:
//   - s: session being built
//   - payload: raw event_msg payload
func (p *Codex) applyEvent(s *entity.Session, payload json.RawMessage) {
	var ev codexRawEvent
	if unmarshalErr := json.Unmarshal(payload, &ev); unmarshalErr != nil {
		return
	}
	switch ev.Type {
	case cfgCodex.EventTokenCount:
		if ev.Info == nil {
			return
		}
		// input_tokens is cumulative AND includes cached prompt
		// tokens; subtract the cache so totals are comparable
		// with the Claude parser's TotalTokensIn.
		in := ev.Info.TotalTokenUsage.InputTokens -
			ev.Info.TotalTokenUsage.CachedInputTokens
		if in < 0 {
			in = 0
		}
		s.TotalTokensIn = in
		s.TotalTokensOut = ev.Info.TotalTokenUsage.OutputTokens
	case cfgCodex.EventItemCompleted:
		if ev.Item == nil ||
			ev.Item.Type != cfgCodex.ItemCommandExecution {
			return
		}
		if ev.Item.ExitCode != nil && *ev.Item.ExitCode != 0 {
			s.HasErrors = true
		}
	}
}

// appendMessage adds a converted message to the session and, for
// the first text-bearing user message, sets the preview. TurnCount
// and the prose gate are computed after the scan (see ParseFile)
// so TurnCount matches the Claude parser's semantics.
//
// Parameters:
//   - s: session being built
//   - msg: converted message, or nil to skip
func (p *Codex) appendMessage(s *entity.Session, msg *entity.Message) {
	if msg == nil {
		return
	}
	s.Messages = append(s.Messages, *msg)

	if !msg.BelongsToUser() || msg.Text == "" {
		return
	}
	if s.FirstUserMsg == "" {
		preview := msg.Text
		if len(preview) > session.PreviewMaxLen {
			preview = truncateRunes(
				preview, session.PreviewMaxLen,
			) + token.Ellipsis
		}
		s.FirstUserMsg = preview
	}
}

// truncateRunes cuts a string at the last rune boundary at or
// before max bytes, so previews never end in a split rune.
//
// Parameters:
//   - s: string to truncate
//   - max: maximum byte length
//
// Returns:
//   - string: prefix ending on a rune boundary
func truncateRunes(s string, max int) string {
	if len(s) <= max {
		return s
	}
	cut := 0
	for i := range s {
		if i > max {
			break
		}
		cut = i
	}
	return s[:cut]
}

// convertItem converts a response_item line to a Message.
//
// Parameters:
//   - raw: decoded line envelope whose payload is a response item
//
// Returns:
//   - *entity.Message: the message, or nil if the item carries none
func (p *Codex) convertItem(raw codexRawLine) *entity.Message {
	var item codexRawItem
	if unmarshalErr := json.Unmarshal(raw.Payload, &item); unmarshalErr != nil {
		return nil
	}

	switch item.Type {
	case cfgCodex.ItemTypeMessage:
		return p.convertMessage(raw.Timestamp, item)
	case cfgCodex.ItemTypeFunctionCall:
		return codexToolUse(raw.Timestamp, item, item.Name, item.Arguments)
	case cfgCodex.ItemTypeCustomToolCall:
		return codexToolUse(raw.Timestamp, item, item.Name, item.Input)
	case cfgCodex.ItemTypeLocalShellCall:
		return codexToolUse(raw.Timestamp, item, item.Type, string(item.Action))
	case cfgCodex.ItemTypeFunctionCallOutput,
		cfgCodex.ItemTypeCustomToolCallOutput:
		return &entity.Message{
			ID:        item.ID,
			Timestamp: raw.Timestamp,
			Role:      cfgCodex.RoleUser,
			ToolResults: []entity.ToolResult{{
				ToolUseID: item.CallID,
				Content:   codexOutputText(item.Output),
			}},
		}
	default:
		return nil
	}
}

// convertMessage converts a message item to a user or assistant
// Message. Developer messages and user items made only of
// Codex-injected parts yield nil.
//
// Parameters:
//   - ts: line timestamp
//   - item: decoded message item
//
// Returns:
//   - *entity.Message: the message, or nil if nothing remains
func (p *Codex) convertMessage(
	ts time.Time, item codexRawItem,
) *entity.Message {
	var text string
	switch item.Role {
	case cfgCodex.RoleUser:
		text = codexJoinParts(item.Content, cfgCodex.ContentInputText, true)
	case cfgCodex.RoleAssistant:
		text = codexJoinParts(item.Content, cfgCodex.ContentOutputText, false)
	default:
		return nil
	}
	if text == "" {
		return nil
	}
	return &entity.Message{
		ID:        item.ID,
		Timestamp: ts,
		Role:      item.Role,
		Text:      text,
	}
}

// codexToolUse builds the assistant message that represents a tool
// invocation, mirroring how Claude Code tool_use blocks are shaped.
//
// Parameters:
//   - ts: line timestamp
//   - item: decoded call item (for ID and call_id)
//   - name: tool name
//   - input: tool input as written in the rollout
//
// Returns:
//   - *entity.Message: assistant message with a single ToolUse
func codexToolUse(
	ts time.Time, item codexRawItem, name, input string,
) *entity.Message {
	return &entity.Message{
		ID:        item.ID,
		Timestamp: ts,
		Role:      cfgCodex.RoleAssistant,
		ToolUses: []entity.ToolUse{{
			ID:    item.CallID,
			Name:  name,
			Input: input,
		}},
	}
}

// codexJoinParts concatenates the text of content parts of the
// given type. When dropInjected is set, parts that open with one of
// the Codex-injected tags are omitted.
//
// Parameters:
//   - parts: content parts of a message
//   - partType: content type to keep (input_text or output_text)
//   - dropInjected: filter Codex-injected parts
//
// Returns:
//   - string: the joined text, empty if no part survives
func codexJoinParts(
	parts []codexRawPart, partType string, dropInjected bool,
) string {
	var kept []string
	for _, part := range parts {
		if part.Type != partType {
			continue
		}
		if dropInjected && codexInjected(part.Text) {
			continue
		}
		kept = append(kept, part.Text)
	}
	return strings.Join(kept, token.NewlineLF)
}

// codexInjected reports whether a user content part is one Codex
// injects itself (environment, instructions, permissions).
//
// Parameters:
//   - text: the part text
//
// Returns:
//   - bool: true if the trimmed text opens with an injected tag
func codexInjected(text string) bool {
	trimmed := strings.TrimSpace(text)
	for _, prefix := range cfgCodex.InjectedUserPrefixes {
		if strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}

// codexOutputText decodes a tool output, which Codex writes either
// as a plain string or as an array of content parts.
//
// Parameters:
//   - output: raw output JSON
//
// Returns:
//   - string: the output text (raw bytes when neither shape decodes)
func codexOutputText(output json.RawMessage) string {
	if len(output) == 0 {
		return ""
	}

	var plain string
	if strErr := json.Unmarshal(output, &plain); strErr == nil {
		return plain
	}

	var parts []codexRawPart
	if partsErr := json.Unmarshal(output, &parts); partsErr == nil {
		var joined strings.Builder
		for _, part := range parts {
			joined.WriteString(part.Text)
		}
		return joined.String()
	}

	return string(output)
}

// codexSessionID picks the session identifier from a session_meta
// payload, preferring the newer `id` field over `session_id`.
//
// Parameters:
//   - meta: decoded session_meta payload
//
// Returns:
//   - string: the session ID, empty if neither field is set
func codexSessionID(meta codexRawSessionMeta) string {
	if meta.ID != "" {
		return meta.ID
	}
	return meta.SessionID
}
