//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"

	cfgAgy "github.com/ActiveMemory/ctx/internal/config/agy"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgParser "github.com/ActiveMemory/ctx/internal/config/parser"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// Ensure Agy implements Session.
var _ Session = (*Agy)(nil)

// NewAgy creates a new Antigravity session parser.
//
// Returns:
//   - *Agy: a parser instance for Antigravity JSONL transcript files
func NewAgy() *Agy {
	return &Agy{}
}

// Tool returns the tool identifier for this parser.
//
// Returns:
//   - string: the identifier "agy"
func (p *Agy) Tool() string {
	return cfgSession.ToolAgy
}

// Matches returns true if the file appears to be an Antigravity transcript.
//
// Checks if the file has a .jsonl extension, is not transcript_full.jsonl,
// and its first non-empty line contains an Antigravity step.
//
// Parameters:
//   - path: file path to check
//
// Returns:
//   - bool: true if this parser can handle the file
func (p *Agy) Matches(path string) bool {
	if !strings.HasSuffix(path, file.ExtJSONL) {
		return false
	}

	if filepath.Base(path) == cfgAgy.FileTranscriptFullJSONL {
		return false
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

	linesChecked := 0
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		linesChecked++
		if linesChecked > cfgParser.LinesToPeek {
			break
		}

		var raw agyRawStep
		if unmarshalErr := json.Unmarshal(line, &raw); unmarshalErr != nil {
			return false
		}

		if validSource(raw.Source) && validStepType(raw.Type) {
			return true
		}
	}

	return false
}

// ParseFile reads an Antigravity transcript and returns the session it holds.
//
// Parameters:
//   - path: path to the transcript JSONL file
//
// Returns:
//   - []*entity.Session: the parsed session (at most one), or nil
//   - error: non-nil if the file cannot be opened or read
func (p *Agy) ParseFile(path string) ([]*entity.Session, error) {
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
		Tool:       cfgSession.ToolAgy,
		SourceFile: path,
	}

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var raw agyRawStep
		if unmarshalErr := json.Unmarshal(line, &raw); unmarshalErr != nil {
			// Skip malformed or partial trailing lines for live transcripts.
			continue
		}

		p.applyStep(s, raw)
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, errParser.ScanFile(scanErr)
	}

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

	s.ID = resolveSessionID(path)
	s.CWD = resolveWorkspace(path, s.ID)

	if !s.StartTime.IsZero() && !s.EndTime.IsZero() {
		s.Duration = s.EndTime.Sub(s.StartTime)
	}
	s.TotalTokens = s.TotalTokensIn + s.TotalTokensOut

	return []*entity.Session{s}, nil
}

// ParseLine parses a single Antigravity transcript line into a Message.
//
// Parameters:
//   - line: raw JSONL line bytes to parse
//
// Returns:
//   - *entity.Message: the parsed message, or nil if none
//   - string: session identifier (empty for Antigravity lines)
//   - error: non-nil if JSON unmarshaling fails
func (p *Agy) ParseLine(line []byte) (*entity.Message, string, error) {
	trimmed := bytes.TrimSpace(line)
	if len(trimmed) == 0 {
		return nil, "", nil
	}

	var raw agyRawStep
	if unmarshalErr := json.Unmarshal(trimmed, &raw); unmarshalErr != nil {
		return nil, "", errParser.Unmarshal(unmarshalErr)
	}

	switch raw.Type {
	case cfgAgy.StepTypeUserInput:
		text, ok := cleanUserContent(raw.Content)
		if !ok || text == "" {
			return nil, "", nil
		}
		msg := &entity.Message{
			ID:        strconv.Itoa(raw.StepIndex),
			Timestamp: raw.CreatedAt,
			Role:      cfgSession.RoleUser,
			Text:      text,
		}
		return msg, "", nil

	case cfgAgy.StepTypePlannerResponse:
		toolUses := convertToolCalls(raw.ToolCalls)
		text := strings.TrimSpace(raw.Content)
		thinking := strings.TrimSpace(raw.Thinking)
		if text == "" && thinking == "" && len(toolUses) == 0 {
			return nil, "", nil
		}
		inTokens := raw.InputTokens + raw.CacheReadTokens
		outTokens := raw.OutputTokens
		msg := &entity.Message{
			ID:        strconv.Itoa(raw.StepIndex),
			Timestamp: raw.CreatedAt,
			Role:      cfgSession.RoleAssistant,
			Text:      text,
			Thinking:  thinking,
			ToolUses:  toolUses,
			TokensIn:  inTokens,
			TokensOut: outTokens,
		}
		return msg, "", nil

	default:
		return nil, "", nil
	}
}
