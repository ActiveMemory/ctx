//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/parser"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/entity"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
)

// ClaudeCode parses Claude Code JSONL session files.
//
// Claude Code stores sessions as JSONL files where each line is a
// self-contained JSON object representing a message. Messages are
// linked via parentUuid and grouped by sessionId.
type ClaudeCode struct{}

// NewClaudeCode creates a new Claude Code session parser.
//
// Returns:
//   - *ClaudeCode: A parser instance for Claude Code JSONL files
func NewClaudeCode() *ClaudeCode {
	return &ClaudeCode{}
}

// Tool returns the tool identifier for this parser.
//
// Returns:
//   - string: The identifier "claude-code"
func (p *ClaudeCode) Tool() string {
	return session.ToolClaudeCode
}

// Matches returns true if the file appears to be a Claude Code session file.
//
// Checks if the file has a .jsonl extension and contains Claude Code message
// format (sessionId and slug fields in the first few lines).
//
// Parameters:
//   - path: File path to check
//
// Returns:
//   - bool: True if this parser can handle the file
func (p *ClaudeCode) Matches(path string) bool {
	// Check extension
	if !strings.HasSuffix(path, file.ExtJSONL) {
		return false
	}

	// Peek at the first few lines to detect the Claude Code format
	f, openErr := os.Open(filepath.Clean(path))
	if openErr != nil {
		return false
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	// Check the first N lines for Claude Code message structure
	// (early lines can be file-history-snapshot which should be skipped)
	for i := 0; i < parser.LinesToPeek && scanner.Scan(); i++ {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var raw claudeRawMessage
		if unmarshalErr := json.Unmarshal(line, &raw); unmarshalErr != nil {
			continue
		}

		// Claude Code messages have sessionId and type (user/assistant)
		// Note: slug field was removed in newer Claude Code versions
		if raw.SessionID != "" && (raw.Type == claude.RoleUser ||
			raw.Type == claude.RoleAssistant) {
			return true
		}
	}

	return false
}

// ParseFile reads a Claude Code JSONL file and returns all sessions.
//
// Parses each line as a JSON message, groups messages by session ID, and
// constructs Session objects with timing, token counts, and message content.
// Sessions are sorted by start time.
//
// Parameters:
//   - path: Path to the JSONL file to parse
//
// Returns:
//   - []*entity.Session: All sessions found in the file, sorted by start time
//   - error: Non-nil if the file cannot be opened or read
func (p *ClaudeCode) ParseFile(path string) ([]*entity.Session, error) {
	f, openErr := os.Open(filepath.Clean(path))
	if openErr != nil {
		return nil, errParser.OpenFile(openErr)
	}
	defer func() { _ = f.Close() }()

	// Group messages by session ID
	sessionMsgs := make(map[string][]claudeRawMessage)

	scanner := bufio.NewScanner(f)
	// Increase buffer size for large lines
	buf := make([]byte, 0, parser.BufInitSize)
	scanner.Buffer(buf, parser.BufMaxSize)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var raw claudeRawMessage
		if unmarshalErr := json.Unmarshal(line, &raw); unmarshalErr != nil {
			// Skip malformed lines, don't fail entire file
			continue
		}

		// Skip non-message lines (e.g., file-history-snapshot)
		if raw.Type != claude.RoleUser && raw.Type != claude.RoleAssistant {
			continue
		}

		if raw.SessionID == "" {
			continue
		}

		sessionMsgs[raw.SessionID] = append(sessionMsgs[raw.SessionID], raw)
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, errParser.ScanFile(scanErr)
	}

	// Convert to sessions
	var sessions []*entity.Session
	for sessionID, msgs := range sessionMsgs {
		s := p.buildSession(sessionID, msgs, path)
		if s != nil {
			sessions = append(sessions, s)
		}
	}

	// Sort sessions by start time
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].StartTime.Before(sessions[j].StartTime)
	})

	return sessions, nil
}

// ParseLine parses a single JSONL line into a Message.
//
// Unmarshals the JSON, extracts the message content, and converts it to the
// common Message type. Non-message lines (e.g., file-history-snapshot) are
// skipped and return nil.
//
// Parameters:
//   - line: Raw JSONL line bytes to parse
//
// Returns:
//   - *entity.Message: The parsed message, or nil if the line should be skipped
//   - string: The session ID this message belongs to
//   - error: Non-nil if JSON unmarshaling fails
func (p *ClaudeCode) ParseLine(line []byte) (*entity.Message, string, error) {
	if len(line) == 0 {
		return nil, "", nil
	}

	var raw claudeRawMessage
	if unmarshalErr := json.Unmarshal(line, &raw); unmarshalErr != nil {
		return nil, "", errParser.Unmarshal(unmarshalErr)
	}

	// Skip non-message lines
	if raw.Type != claude.RoleUser && raw.Type != claude.RoleAssistant {
		return nil, "", nil
	}

	msg := p.convertMessage(raw)
	return &msg, raw.SessionID, nil
}

// Ensure ClaudeCode implements Session.
var _ Session = (*ClaudeCode)(nil)
