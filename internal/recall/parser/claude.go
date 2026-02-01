//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
)

// ClaudeCodeParser parses Claude Code JSONL session files.
//
// Claude Code stores sessions as JSONL files where each line is a
// self-contained JSON object representing a message. Messages are
// linked via parentUuid and grouped by sessionId.
type ClaudeCodeParser struct{}

// NewClaudeCodeParser creates a new Claude Code session parser.
//
// Returns:
//   - *ClaudeCodeParser: A parser instance for Claude Code JSONL files
func NewClaudeCodeParser() *ClaudeCodeParser {
	return &ClaudeCodeParser{}
}

// Tool returns the tool identifier for this parser.
//
// Returns:
//   - string: The identifier "claude-code"
func (p *ClaudeCodeParser) Tool() string {
	return "claude-code"
}

// CanParse returns true if the file appears to be a Claude Code session file.
//
// Checks if the file has a .jsonl extension and contains Claude Code message
// format (sessionId and slug fields in the first few lines).
//
// Parameters:
//   - path: File path to check
//
// Returns:
//   - bool: True if this parser can handle the file
func (p *ClaudeCodeParser) CanParse(path string) bool {
	// Check extension
	if !strings.HasSuffix(path, ".jsonl") {
		return false
	}

	// Peek at the first few lines to detect the Claude Code format
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err.Error())
		}
	}(file)

	scanner := bufio.NewScanner(file)
	// Check the first N lines - slug may not appear until later in the file
	// (early lines can be file-history-snapshot or messages without the slug)
	for i := 0; i < config.ParserPeekLines && scanner.Scan(); i++ {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var raw claudeRawMessage
		if err := json.Unmarshal(line, &raw); err != nil {
			continue
		}

		// Claude Code messages have sessionId and slug fields
		if raw.SessionID != "" && raw.Slug != "" {
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
//   - []*Session: All sessions found in the file, sorted by start time
//   - error: Non-nil if the file cannot be opened or read
func (p *ClaudeCodeParser) ParseFile(path string) ([]*Session, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err.Error())
		}
	}(file)

	// Group messages by session ID
	sessionMsgs := make(map[string][]claudeRawMessage)

	scanner := bufio.NewScanner(file)
	// Increase buffer size for large lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // 1MB max line size

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var raw claudeRawMessage
		if err := json.Unmarshal(line, &raw); err != nil {
			// Skip malformed lines, don't fail entire file
			continue
		}

		// Skip non-message lines (e.g., file-history-snapshot)
		if raw.Type != config.RoleUser && raw.Type != config.RoleAssistant {
			continue
		}

		if raw.SessionID == "" {
			continue
		}

		sessionMsgs[raw.SessionID] = append(sessionMsgs[raw.SessionID], raw)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan file: %w", err)
	}

	// Convert to sessions
	var sessions []*Session
	for sessionID, msgs := range sessionMsgs {
		session := p.buildSession(sessionID, msgs, path)
		if session != nil {
			sessions = append(sessions, session)
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
//   - *Message: The parsed message, or nil if the line should be skipped
//   - string: The session ID this message belongs to
//   - error: Non-nil if JSON unmarshaling fails
func (p *ClaudeCodeParser) ParseLine(line []byte) (*Message, string, error) {
	if len(line) == 0 {
		return nil, "", nil
	}

	var raw claudeRawMessage
	if err := json.Unmarshal(line, &raw); err != nil {
		return nil, "", fmt.Errorf("unmarshal: %w", err)
	}

	// Skip non-message lines
	if raw.Type != config.RoleUser && raw.Type != config.RoleAssistant {
		return nil, "", nil
	}

	msg := p.convertMessage(raw)
	return &msg, raw.SessionID, nil
}

// Claude Code-specific raw types for parsing JSONL

type claudeRawMessage struct {
	UUID        string           `json:"uuid"`
	ParentUUID  *string          `json:"parentUuid"`
	SessionID   string           `json:"sessionId"`
	RequestID   string           `json:"requestId,omitempty"`
	Timestamp   time.Time        `json:"timestamp"`
	Type        string           `json:"type"` // "user", "assistant", or other
	UserType    string           `json:"userType,omitempty"`
	IsSidechain bool             `json:"isSidechain,omitempty"`
	CWD         string           `json:"cwd"`
	GitBranch   string           `json:"gitBranch,omitempty"`
	Version     string           `json:"version"`
	Slug        string           `json:"slug"`
	Message     claudeRawContent `json:"message"`
}

type claudeRawContent struct {
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Model        string          `json:"model,omitempty"`
	Role         string          `json:"role"`
	Content      json.RawMessage `json:"content"` // Can be string or []claudeRawBlock
	StopReason   *string         `json:"stop_reason,omitempty"`
	StopSequence *string         `json:"stop_sequence,omitempty"`
	Usage        *claudeRawUsage `json:"usage,omitempty"`
}

type claudeRawBlock struct {
	Type      string          `json:"type"`
	Text      string          `json:"text,omitempty"`
	Thinking  string          `json:"thinking,omitempty"`
	Signature string          `json:"signature,omitempty"`
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
	ToolUseID string          `json:"tool_use_id,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
	IsError   bool            `json:"is_error,omitempty"`
}

type claudeRawUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
}

// Ensure ClaudeCodeParser implements SessionParser
var _ SessionParser = (*ClaudeCodeParser)(nil)
