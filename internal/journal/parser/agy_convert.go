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
	"time"

	"github.com/ActiveMemory/ctx/internal/agy"
	cfgAgy "github.com/ActiveMemory/ctx/internal/config/agy"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgHttp "github.com/ActiveMemory/ctx/internal/config/http"
	cfgParser "github.com/ActiveMemory/ctx/internal/config/parser"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// applyStep applies a parsed step to the session accumulating data.
//
// Parameters:
//   - s: target session
//   - raw: raw transcript step
func (p *Agy) applyStep(s *entity.Session, raw agyRawStep) {
	switch raw.Type {
	case cfgAgy.StepTypeUserInput:
		text, ok := cleanUserContent(raw.Content)
		if !ok {
			return
		}
		msg := entity.Message{
			ID:        strconv.Itoa(raw.StepIndex),
			Timestamp: raw.CreatedAt,
			Role:      cfgSession.RoleUser,
			Text:      text,
		}
		s.Messages = append(s.Messages, msg)
		if s.FirstUserMsg == "" && text != "" {
			s.FirstUserMsg = msg.Preview(cfgSession.PreviewMaxLen)
		}
		updateSessionTimes(s, raw.CreatedAt)

	case cfgAgy.StepTypePlannerResponse:
		toolUses := convertToolCalls(raw.ToolCalls)
		text := strings.TrimSpace(raw.Content)
		thinking := strings.TrimSpace(raw.Thinking)
		if text == "" && thinking == "" && len(toolUses) == 0 {
			return
		}
		inTokens := raw.InputTokens + raw.CacheReadTokens
		outTokens := raw.OutputTokens
		s.TotalTokensIn += inTokens
		s.TotalTokensOut += outTokens
		msg := entity.Message{
			ID:        strconv.Itoa(raw.StepIndex),
			Timestamp: raw.CreatedAt,
			Role:      cfgSession.RoleAssistant,
			Text:      text,
			Thinking:  thinking,
			ToolUses:  toolUses,
			TokensIn:  inTokens,
			TokensOut: outTokens,
		}
		s.Messages = append(s.Messages, msg)
		updateSessionTimes(s, raw.CreatedAt)
	}
}

// updateSessionTimes tracks earliest and latest timestamps on a session.
//
// Parameters:
//   - s: target session
//   - ts: timestamp from the step
func updateSessionTimes(s *entity.Session, ts time.Time) {
	if ts.IsZero() {
		return
	}
	if s.StartTime.IsZero() || ts.Before(s.StartTime) {
		s.StartTime = ts
	}
	if ts.After(s.EndTime) {
		s.EndTime = ts
	}
}

// convertToolCalls transforms Antigravity raw tool calls to entity.ToolUse.
//
// Parameters:
//   - calls: raw tool calls slice
//
// Returns:
//   - []entity.ToolUse: normalized tool use slice
func convertToolCalls(calls []agyRawToolCall) []entity.ToolUse {
	var toolUses []entity.ToolUse
	for _, tc := range calls {
		toolUses = append(toolUses, entity.ToolUse{
			Name:  tc.Name,
			Input: string(tc.Args),
		})
	}
	return toolUses
}

// cleanUserContent extracts user text from a USER_INPUT step.
//
// Parameters:
//   - content: raw content string from the transcript step
//
// Returns:
//   - string: cleaned user prose
//   - bool: true if the step represents user content, false to skip
func cleanUserContent(content string) (string, bool) {
	if strings.Contains(content, cfgAgy.TagUserRequestOpen) &&
		strings.Contains(content, cfgAgy.TagUserRequestClose) {
		start := strings.Index(content, cfgAgy.TagUserRequestOpen) +
			len(cfgAgy.TagUserRequestOpen)
		end := strings.Index(content, cfgAgy.TagUserRequestClose)
		if end > start {
			text := strings.TrimSpace(content[start:end])
			text = stripAdditionalMetadata(text)
			return text, true
		}
	}

	if strings.Contains(content, cfgAgy.TagSystemMessageOpen) {
		return "", false
	}

	text := stripAdditionalMetadata(content)
	return strings.TrimSpace(text), true
}

// stripAdditionalMetadata strips <ADDITIONAL_METADATA> tags and content.
//
// Parameters:
//   - content: content string potentially containing metadata tags
//
// Returns:
//   - string: content with metadata tags removed
func stripAdditionalMetadata(content string) string {
	for strings.Contains(content, cfgAgy.TagAdditionalMetadataOpen) &&
		strings.Contains(content, cfgAgy.TagAdditionalMetadataClose) {
		start := strings.Index(content, cfgAgy.TagAdditionalMetadataOpen)
		end := strings.Index(content, cfgAgy.TagAdditionalMetadataClose) +
			len(cfgAgy.TagAdditionalMetadataClose)
		if end > start {
			content = content[:start] + content[end:]
		} else {
			break
		}
	}
	return strings.TrimSpace(content)
}

// validSource reports whether source is a known Antigravity step source.
//
// Parameters:
//   - source: step source string
//
// Returns:
//   - bool: true if recognized
func validSource(source string) bool {
	return source == cfgAgy.SourceUserExplicit ||
		source == cfgAgy.SourceModel ||
		source == cfgAgy.SourceSystem
}

// validStepType reports whether stepType is a known Antigravity step type.
//
// Parameters:
//   - stepType: step type string
//
// Returns:
//   - bool: true if recognized
func validStepType(stepType string) bool {
	return stepType == cfgAgy.StepTypeUserInput ||
		stepType == cfgAgy.StepTypePlannerResponse ||
		stepType == cfgAgy.StepTypeGeneric
}

// resolveSessionID resolves the session ID from path.
//
// Parameters:
//   - path: transcript file path
//
// Returns:
//   - string: resolved session identifier
func resolveSessionID(path string) string {
	logsDir := filepath.Dir(path)
	if filepath.Base(logsDir) == cfgAgy.DirLogs {
		sysGenDir := filepath.Dir(logsDir)
		if filepath.Base(sysGenDir) == cfgAgy.DirSystemGenerated {
			idDir := filepath.Dir(sysGenDir)
			base := filepath.Base(idDir)
			if base != "" && base != token.Dot {
				return base
			}
		}
	}
	return strings.TrimSuffix(filepath.Base(path), file.ExtJSONL)
}

// findHistoryFile locates history.jsonl in ancestor directories or home.
//
// Parameters:
//   - sessionPath: path to the transcript file
//
// Returns:
//   - string: path to history.jsonl, or empty if not found
func findHistoryFile(sessionPath string) string {
	d4 := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(sessionPath))))
	c4 := filepath.Join(d4, cfgAgy.FileHistoryJSONL)
	if info, statErr := io.SafeStat(c4); statErr == nil && !info.IsDir() {
		return c4
	}

	d5 := filepath.Dir(d4)
	c5 := filepath.Join(d5, cfgAgy.FileHistoryJSONL)
	if info, statErr := io.SafeStat(c5); statErr == nil && !info.IsDir() {
		return c5
	}

	if home := agy.Home(); home != "" {
		cHome := filepath.Join(
			home, cfgAgy.DirAntigravityCLI, cfgAgy.FileHistoryJSONL,
		)
		if info, statErr := io.SafeStat(cHome); statErr == nil && !info.IsDir() {
			return cHome
		}
	}

	return ""
}

// resolveWorkspace extracts the workspace for sessionID from history.jsonl.
//
// Parameters:
//   - sessionPath: path to the transcript file
//   - sessionID: session identifier to match
//
// Returns:
//   - string: resolved workspace directory or empty
func resolveWorkspace(sessionPath, sessionID string) string {
	histPath := findHistoryFile(sessionPath)
	if histPath == "" {
		return ""
	}

	f, scanner, openErr := openScanner(histPath, cfgParser.BufMaxSizeSchema)
	if openErr != nil {
		return ""
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			logWarn.Warn(cfgWarn.Close, histPath, closeErr)
		}
	}()

	var lastWorkspace string
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}

		var entry agyRawHistory
		if unmarshalErr := json.Unmarshal(line, &entry); unmarshalErr != nil {
			continue
		}

		if entry.ConversationID != sessionID {
			continue
		}

		ws := extractWorkspace(entry)
		if ws != "" {
			lastWorkspace = ws
		}
	}

	return lastWorkspace
}

// extractWorkspace reads the workspace path from a history entry.
//
// Parameters:
//   - entry: history entry
//
// Returns:
//   - string: extracted workspace path or empty
func extractWorkspace(entry agyRawHistory) string {
	if entry.Workspace != "" {
		return parseWorkspacePath(entry.Workspace)
	}

	if len(entry.Workspaces) > 0 {
		var wsList []string
		if json.Unmarshal(entry.Workspaces, &wsList) == nil &&
			len(wsList) > 0 && wsList[0] != "" {
			return parseWorkspacePath(wsList[0])
		}
		var single string
		if json.Unmarshal(entry.Workspaces, &single) == nil && single != "" {
			return parseWorkspacePath(single)
		}
	}

	return ""
}

// parseWorkspacePath normalizes a workspace path or file URI.
//
// Parameters:
//   - raw: raw workspace string
//
// Returns:
//   - string: normalized local filesystem path
func parseWorkspacePath(raw string) string {
	if strings.HasPrefix(raw, cfgHttp.PrefixFile) {
		return fileURIToPath(raw)
	}
	return filepath.Clean(raw)
}
