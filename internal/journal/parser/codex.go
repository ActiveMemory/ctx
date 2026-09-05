//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"encoding/json"
	"path/filepath"
	"strings"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgParser "github.com/ActiveMemory/ctx/internal/config/parser"
	"github.com/ActiveMemory/ctx/internal/config/session"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// Ensure Codex implements Session.
var _ Session = (*Codex)(nil)

// NewCodex creates a new Codex rollout parser.
//
// Returns:
//   - *Codex: a parser instance for Codex rollout JSONL files
func NewCodex() *Codex {
	return &Codex{}
}

// Tool returns the tool identifier for this parser.
//
// Returns:
//   - string: the identifier "codex"
func (p *Codex) Tool() string {
	return session.ToolCodex
}

// Matches returns true if the file appears to be a Codex rollout.
//
// A rollout is a .jsonl file whose basename starts with `rollout-`
// or whose first non-empty line is a `session_meta` envelope. Claude
// Code transcripts carry a top-level `sessionId` and `type` of
// user/assistant on every line and never satisfy either check.
//
// Parameters:
//   - path: file path to check
//
// Returns:
//   - bool: true if this parser can handle the file
func (p *Codex) Matches(path string) bool {
	if !strings.HasSuffix(path, file.ExtJSONL) {
		return false
	}

	if strings.HasPrefix(filepath.Base(path), cfgCodex.RolloutPrefix) {
		return true
	}

	f, scanner, openErr := openScanner(path, cfgParser.BufMaxSizeSchema)
	if openErr != nil {
		return false
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			logWarn.Warn(cfgWarn.Close, path, closeErr)
		}
	}()

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var raw codexRawLine
		if unmarshalErr := json.Unmarshal(line, &raw); unmarshalErr != nil {
			return false
		}
		return raw.Type == cfgCodex.LineTypeSessionMeta
	}

	return false
}

// ParseFile reads a Codex rollout and returns the session it holds.
//
// Each rollout is exactly one session. Lines are applied in file
// order: `session_meta` seeds identity and context, `turn_context`
// tracks the model, `token_count` events carry cumulative usage
// (the last one wins), and `response_item` lines become messages.
// Malformed lines are skipped, as in the Claude Code parser. A
// rollout with no user prose (only Codex-injected items) yields no
// session.
//
// Parameters:
//   - path: path to the rollout JSONL file
//
// Returns:
//   - []*entity.Session: the parsed session (at most one), or nil
//   - error: non-nil if the file cannot be opened or read
func (p *Codex) ParseFile(path string) ([]*entity.Session, error) {
	f, scanner, openErr := openScanner(path, cfgParser.BufMaxSizeSchema)
	if openErr != nil {
		return nil, errParser.OpenFile(openErr)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			logWarn.Warn(cfgWarn.Close, path, closeErr)
		}
	}()

	s := &entity.Session{
		Tool:       session.ToolCodex,
		SourceFile: path,
	}

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var raw codexRawLine
		if unmarshalErr := json.Unmarshal(line, &raw); unmarshalErr != nil {
			// Skip malformed lines, don't fail the entire file
			continue
		}

		p.applyLine(s, raw)
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, errParser.ScanFile(scanErr)
	}

	// TurnCount matches the Claude parser: every user-role
	// message counts, including tool results. Sessions with no
	// user PROSE (only injected/tool traffic) are still skipped.
	prose := 0
	for _, msg := range s.Messages {
		if !msg.BelongsToUser() {
			continue
		}
		s.TurnCount++
		if msg.Text != "" {
			prose++
		}
	}
	if prose == 0 {
		return nil, nil
	}

	if s.ID == "" {
		s.ID = strings.TrimSuffix(filepath.Base(path), file.ExtJSONL)
	}
	if !s.StartTime.IsZero() && !s.EndTime.IsZero() {
		s.Duration = s.EndTime.Sub(s.StartTime)
	}
	s.TotalTokens = s.TotalTokensIn + s.TotalTokensOut

	return []*entity.Session{s}, nil
}

// ParseLine parses a single rollout line into a Message.
//
// Only `response_item` lines carry conversation content: messages,
// tool calls, and tool outputs convert to a Message; developer
// messages, Codex-injected user items, and reasoning items return
// nil. A `session_meta` line returns no message but reports the
// session ID. Every other line type (event_msg, turn_context,
// world_state, compacted) returns nil with an empty session ID.
//
// Parameters:
//   - line: raw JSONL line bytes to parse
//
// Returns:
//   - *entity.Message: the parsed message, or nil if the line carries none
//   - string: the session ID (session_meta lines only)
//   - error: non-nil if JSON unmarshaling fails
func (p *Codex) ParseLine(line []byte) (*entity.Message, string, error) {
	if len(line) == 0 {
		return nil, "", nil
	}

	var raw codexRawLine
	if unmarshalErr := json.Unmarshal(line, &raw); unmarshalErr != nil {
		return nil, "", errParser.Unmarshal(unmarshalErr)
	}

	switch raw.Type {
	case cfgCodex.LineTypeSessionMeta:
		var meta codexRawSessionMeta
		if metaErr := json.Unmarshal(raw.Payload, &meta); metaErr != nil {
			return nil, "", errParser.Unmarshal(metaErr)
		}
		return nil, codexSessionID(meta), nil
	case cfgCodex.LineTypeResponseItem:
		return p.convertItem(raw), "", nil
	default:
		return nil, "", nil
	}
}
