//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
)

// HookInput represents the JSON payload that Claude Code sends to hook
// commands via stdin.
type HookInput struct {
	SessionID string    `json:"session_id"`
	ToolInput ToolInput `json:"tool_input"`
}

// ToolInput contains the tool-specific fields from a Claude Code hook
// invocation. For Bash hooks, Command holds the shell command.
type ToolInput struct {
	Command string `json:"command"`
}

// Response is the JSON output format for Claude Code hooks.
// Using structured JSON ensures the agent processes the output as a directive
// rather than treating it as ignorable plain text.
type Response struct {
	Output *ResponseOutput `json:"hookSpecificOutput,omitempty"`
}

// ResponseOutput carries event-specific fields inside a Response.
type ResponseOutput struct {
	HookEventName     string `json:"hookEventName"`
	AdditionalContext string `json:"additionalContext,omitempty"`
}

// BlockResponse is the JSON output for blocked commands.
type BlockResponse struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// TokenInfo holds token usage and model information extracted from a
// session's JSONL file.
type TokenInfo struct {
	Tokens int    // Total input tokens (input + cache_creation + cache_read)
	Model  string // Model ID from the last assistant message, or ""
}

// usageData represents the minimal usage fields from a Claude Code JSONL
// assistant message. Only the fields needed for token counting are included.
type usageData struct {
	InputTokens              int `json:"input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

// jsonlMessage represents the minimal structure of a Claude Code JSONL line
// needed to extract usage and model data from assistant messages.
type jsonlMessage struct {
	Type    string `json:"type"`
	Message struct {
		Role  string    `json:"role"`
		Model string    `json:"model"`
		Usage usageData `json:"usage"`
	} `json:"message"`
}

// FormatContext builds a JSON Response with additionalContext for the
// given hook event. This is the standard way for non-blocking hooks to inject
// directives that the agent will actually process (plain text gets ignored).
//
// Parameters:
//   - event: Hook event name
//   - context: Additional context string
//
// Returns:
//   - string: JSON-encoded hook response
func FormatContext(event, context string) string {
	resp := Response{
		Output: &ResponseOutput{
			HookEventName:     event,
			AdditionalContext: context,
		},
	}
	data, _ := json.Marshal(resp)
	return string(data)
}

// SessionStats holds the fields written to the per-session stats JSONL file.
type SessionStats struct {
	Timestamp  string `json:"ts"`
	Prompt     int    `json:"prompt"`
	Tokens     int    `json:"tokens"`
	Pct        int    `json:"pct"`
	WindowSize int    `json:"window"`
	Model      string `json:"model,omitempty"`
	Event      string `json:"event"`
}

// ReadInput reads and parses the JSON hook input from r.
// Returns a zero-value HookInput on any error (graceful degradation).
//
// Guards against blocking forever on stdin:
//   - Terminal (character device): returns immediately
//   - Pipe/file with no EOF within 2s: times out and returns zero value
//
// Both cases are harmless — hooks degrade gracefully with zero input.
//
// Parameters:
//   - r: Reader to read hook input from
//
// Returns:
//   - HookInput: Parsed input or zero value
func ReadInput(r io.Reader) HookInput {
	if f, ok := r.(*os.File); ok {
		if fi, readErr := f.Stat(); readErr == nil && fi.Mode()&os.ModeCharDevice != 0 {
			return HookInput{}
		}
	}

	type readResult struct {
		data []byte
		err  error
	}
	ch := make(chan readResult, 1)
	go func() {
		data, readErr := io.ReadAll(r)
		ch <- readResult{data, readErr}
	}()

	var input HookInput
	select {
	case res := <-ch:
		if res.err == nil {
			_ = json.Unmarshal(res.data, &input)
		}
	case <-time.After(2 * time.Second):
	}
	return input
}

// ReadSessionID reads the session ID from stdin JSON, returning the
// fallback "unknown" if stdin is empty or unparseable.
//
// Parameters:
//   - stdin: File to read input from
//
// Returns:
//   - string: Session ID or config.IDSessionUnknown
func ReadSessionID(stdin *os.File) string {
	input := ReadInput(stdin)
	if input.SessionID == "" {
		return cfgSession.IDUnknown
	}
	return input.SessionID
}

// WriteSessionStats appends a JSONL line to .context/state/stats-{sessionID}.jsonl.
// The file is designed for `tail -f` monitoring of token usage across prompts.
// Best-effort: errors are silently ignored.
//
// Parameters:
//   - sessionID: Session identifier
//   - stats: Stats entry to write
func WriteSessionStats(sessionID string, stats SessionStats) {
	path := filepath.Join(core.StateDir(), "stats-"+sessionID+".jsonl")
	data, marshalErr := json.Marshal(stats)
	if marshalErr != nil {
		return
	}
	data = append(data, '\n')

	f, openErr := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) //nolint:gosec // state dir path
	if openErr != nil {
		return
	}
	defer func() { _ = f.Close() }()
	_, _ = f.Write(data)
}
