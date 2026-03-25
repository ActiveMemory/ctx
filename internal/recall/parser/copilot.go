//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// osWindows is the runtime.GOOS value for Windows.
const osWindows = "windows"

// copilotKeyRequests is the key path segment for request arrays.
const copilotKeyRequests = "requests"

// CopilotParser parses VS Code Copilot Chat JSONL session files.
//
// Copilot Chat stores sessions as JSONL files in VS Code's workspaceStorage
// directory. Each file contains one session. The first line is a full session
// snapshot (kind=0), subsequent lines are incremental patches (kind=1, kind=2).
type CopilotParser struct{}

// NewCopilotParser creates a new Copilot Chat session parser.
func NewCopilotParser() *CopilotParser {
	return &CopilotParser{}
}

// Tool returns the tool identifier for this parser.
func (p *CopilotParser) Tool() string {
	return session.ToolCopilot
}

// Matches returns true if the file appears to be a Copilot Chat session file.
//
// Checks if the file has a .jsonl extension and lives in a chatSessions
// directory, and the first line contains a Copilot session snapshot.
func (p *CopilotParser) Matches(path string) bool {
	if !strings.HasSuffix(path, file.ExtJSONL) {
		return false
	}

	// Copilot sessions live in chatSessions/ directories
	if !strings.Contains(filepath.Dir(path), "chatSessions") {
		return false
	}

	file, openErr := os.Open(filepath.Clean(path))
	if openErr != nil {
		return false
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	if !scanner.Scan() {
		return false
	}

	var line copilotRawLine
	if err := json.Unmarshal(scanner.Bytes(), &line); err != nil {
		return false
	}

	// kind=0 is the full session snapshot
	if line.Kind != 0 {
		return false
	}

	var session copilotRawSession
	if err := json.Unmarshal(line.V, &session); err != nil {
		return false
	}

	return session.SessionID != "" && session.Version > 0
}

// ParseFile reads a Copilot Chat JSONL file and returns the session.
//
// Reconstructs the session by reading the initial snapshot (kind=0) and
// applying incremental patches (kind=1 for scalar, kind=2 for array/object).
func (p *CopilotParser) ParseFile(path string) ([]*entity.Session, error) {
	file, openErr := os.Open(filepath.Clean(path))
	if openErr != nil {
		return nil, fmt.Errorf("open file: %w", openErr)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 4*1024*1024) // 4MB — Copilot lines can be very large

	var session *copilotRawSession

	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		if len(lineBytes) == 0 {
			continue
		}

		var line copilotRawLine
		if err := json.Unmarshal(lineBytes, &line); err != nil {
			continue
		}

		switch line.Kind {
		case 0:
			// Full session snapshot
			var s copilotRawSession
			if err := json.Unmarshal(line.V, &s); err != nil {
				return nil, fmt.Errorf("parse session snapshot: %w", err)
			}
			session = &s

		case 1:
			// Scalar property patch — apply to session
			if session != nil {
				p.applyScalarPatch(session, line.K, line.V)
			}

		case 2:
			// Array/object patch — apply to session
			if session != nil {
				p.applyPatch(session, line.K, line.V)
			}
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, fmt.Errorf("scan file: %w", scanErr)
	}

	if session == nil {
		return nil, nil
	}

	// Resolve workspace folder from workspace.json next to chatSessions/
	cwd := p.resolveWorkspaceCWD(path)

	result := p.buildSession(session, path, cwd)
	if result == nil {
		return nil, nil
	}

	return []*entity.Session{result}, nil
}

// ParseLine is not meaningful for Copilot sessions since they use patches.
// Returns nil for all lines.
func (p *CopilotParser) ParseLine(_ []byte) (*entity.Message, string, error) {
	return nil, "", nil
}

// applyScalarPatch applies a kind=1 scalar patch to the session.
// These update individual properties like result, modelState, followups.
func (p *CopilotParser) applyScalarPatch(
	session *copilotRawSession, keys []json.RawMessage, value json.RawMessage,
) {
	path := p.parseKeyPath(keys)
	if len(path) < 2 {
		return
	}

	// Handle requests.<N>.result patches — these contain token counts
	if path[0] == copilotKeyRequests && len(path) == 3 && path[2] == "result" {
		idx, err := strconv.Atoi(path[1])
		if err != nil || idx < 0 || idx >= len(session.Requests) {
			return
		}
		var result copilotRawResult
		if err := json.Unmarshal(value, &result); err == nil {
			session.Requests[idx].Result = &result
		}
	}
}

// applyPatch applies a kind=2 array/object patch to the session.
func (p *CopilotParser) applyPatch(
	session *copilotRawSession, keys []json.RawMessage, value json.RawMessage,
) {
	path := p.parseKeyPath(keys)
	if len(path) == 0 {
		return
	}

	switch {
	case len(path) == 1 && path[0] == copilotKeyRequests:
		// New request(s) appended
		var requests []copilotRawRequest
		if err := json.Unmarshal(value, &requests); err == nil {
			session.Requests = append(session.Requests, requests...)
		}

	case len(path) == 3 && path[0] == copilotKeyRequests && path[2] == "response":
		// Response update for a specific request
		idx, err := strconv.Atoi(path[1])
		if err != nil || idx < 0 || idx >= len(session.Requests) {
			return
		}
		var items []copilotRawRespItem
		if err := json.Unmarshal(value, &items); err == nil {
			session.Requests[idx].Response = items
		}
	}
}

// parseKeyPath converts the K array from JSONL into string path segments.
func (p *CopilotParser) parseKeyPath(keys []json.RawMessage) []string {
	path := make([]string, 0, len(keys))
	for _, k := range keys {
		var s string
		if err := json.Unmarshal(k, &s); err == nil {
			path = append(path, s)
			continue
		}
		var n int
		if err := json.Unmarshal(k, &n); err == nil {
			path = append(path, strconv.Itoa(n))
			continue
		}
	}
	return path
}

// buildSession converts a reconstructed copilotRawSession into a Session.
func (p *CopilotParser) buildSession(
	raw *copilotRawSession, sourcePath string, cwd string,
) *entity.Session {
	if len(raw.Requests) == 0 {
		return nil
	}

	session := &entity.Session{
		ID:         raw.SessionID,
		Tool:       session.ToolCopilot,
		SourceFile: sourcePath,
		CWD:        cwd,
		Project:    filepath.Base(cwd),
		StartTime:  time.UnixMilli(raw.CreationDate),
	}

	if raw.CustomTitle != "" {
		session.Slug = raw.CustomTitle
	}

	for _, req := range raw.Requests {
		// User message
		userMsg := entity.Message{
			ID:        req.RequestID,
			Timestamp: time.UnixMilli(req.Timestamp),
			Role:      claude.RoleUser,
			Text:      req.Message.Text,
		}

		if req.Result != nil {
			userMsg.TokensIn = req.Result.Metadata.PromptTokens
		}

		session.Messages = append(session.Messages, userMsg)
		session.TurnCount++

		if session.FirstUserMsg == "" && userMsg.Text != "" {
			preview := userMsg.Text
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			session.FirstUserMsg = preview
		}

		// Assistant response
		assistantMsg := p.buildAssistantMessage(req)
		if assistantMsg != nil {
			session.Messages = append(session.Messages, *assistantMsg)

			if session.Model == "" && req.ModelID != "" {
				session.Model = req.ModelID
			}
		}

		// Accumulate tokens
		if req.Result != nil {
			session.TotalTokensIn += req.Result.Metadata.PromptTokens
			session.TotalTokensOut += req.Result.Metadata.OutputTokens
		}
	}

	session.TotalTokens = session.TotalTokensIn + session.TotalTokensOut

	// Set end time from last request
	if last := raw.Requests[len(raw.Requests)-1]; last.Result != nil {
		session.EndTime = time.UnixMilli(last.Timestamp).Add(
			time.Duration(last.Result.Timings.TotalElapsed) * time.Millisecond,
		)
	} else {
		session.EndTime = time.UnixMilli(
			raw.Requests[len(raw.Requests)-1].Timestamp,
		)
	}
	session.Duration = session.EndTime.Sub(session.StartTime)

	return session
}

// buildAssistantMessage extracts the assistant response from a request.
func (p *CopilotParser) buildAssistantMessage(
	req copilotRawRequest,
) *entity.Message {
	if len(req.Response) == 0 {
		return nil
	}

	msg := &entity.Message{
		ID:        req.RequestID + "-response",
		Timestamp: time.UnixMilli(req.Timestamp),
		Role:      claude.RoleAssistant,
	}

	if req.Result != nil {
		msg.TokensOut = req.Result.Metadata.OutputTokens
	}

	for _, item := range req.Response {
		switch item.Kind {
		case "thinking":
			var text string
			if err := json.Unmarshal(item.Value, &text); err == nil {
				if msg.Thinking != "" {
					msg.Thinking += token.NewlineLF
				}
				msg.Thinking += text
			}

		case "toolInvocationSerialized":
			tu := p.parseToolInvocation(item)
			if tu != nil {
				msg.ToolUses = append(msg.ToolUses, *tu)
			}

		case "":
			// Plain markdown text (no kind field)
			var text string
			if err := json.Unmarshal(item.Value, &text); err == nil {
				text = strings.TrimSpace(text)
				if text != "" {
					if msg.Text != "" {
						msg.Text += token.NewlineLF
					}
					msg.Text += text
				}
			}

			// Skip: codeblockUri, inlineReference, progressTaskSerialized,
			//        textEditGroup, undoStop, mcpServersStarting
		}
	}

	// Check for tool errors
	for _, tr := range msg.ToolResults {
		if tr.IsError {
			return msg // HasErrors is set at session level
		}
	}

	return msg
}

// parseToolInvocation extracts a ToolUse from a toolInvocationSerialized item.
func (p *CopilotParser) parseToolInvocation(item copilotRawRespItem) *entity.ToolUse {
	toolID := item.ToolID
	if toolID == "" {
		return nil
	}

	// Extract the tool name from toolId (e.g., "copilot_readFile" -> "readFile")
	name := toolID
	if idx := strings.LastIndex(toolID, "_"); idx >= 0 {
		name = toolID[idx+1:]
	}

	// Use invocationMessage as the input description
	inputStr := ""
	if item.InvocationMessage != nil {
		// InvocationMessage can be a string or object with value field
		var simple string
		if err := json.Unmarshal(item.InvocationMessage, &simple); err == nil {
			inputStr = simple
		} else {
			var obj struct {
				Value string `json:"value"`
			}
			if err := json.Unmarshal(item.InvocationMessage, &obj); err == nil {
				inputStr = obj.Value
			}
		}
	}

	return &entity.ToolUse{
		ID:    item.ToolCallID,
		Name:  name,
		Input: inputStr,
	}
}

// resolveWorkspaceCWD reads workspace.json from the workspaceStorage
// directory to determine the workspace folder path.
func (p *CopilotParser) resolveWorkspaceCWD(sessionPath string) string {
	// sessionPath is like: .../workspaceStorage/<hash>/chatSessions/<id>.jsonl
	// workspace.json is at: .../workspaceStorage/<hash>/workspace.json
	chatDir := filepath.Dir(sessionPath) // chatSessions/
	storageDir := filepath.Dir(chatDir)  // <hash>/
	wsFile := filepath.Join(storageDir, "workspace.json")

	data, err := os.ReadFile(filepath.Clean(wsFile))
	if err != nil {
		return ""
	}

	var ws copilotRawWorkspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return ""
	}

	return fileURIToPath(ws.Folder)
}

// fileURIToPath converts a file:// URI to a local file path.
// Example: "file:///g%3A/GitProjects/ctx" -> "G:\GitProjects\ctx" (Windows)
//
//	"file:///home/user/project" -> "/home/user/project" (Unix)
func fileURIToPath(uri string) string {
	if uri == "" {
		return ""
	}

	parsed, err := url.Parse(uri)
	if err != nil {
		return ""
	}

	if parsed.Scheme != "file" {
		return ""
	}

	path := parsed.Path

	// URL-decode the path (e.g., %3A -> :)
	decoded, err := url.PathUnescape(path)
	if err != nil {
		decoded = path
	}

	// On Windows, file URIs have /G:/... — strip the leading slash
	if runtime.GOOS == osWindows && len(decoded) > 2 && decoded[0] == '/' {
		decoded = decoded[1:]
	}

	return filepath.FromSlash(decoded)
}

// CopilotSessionDirs returns the directories where Copilot Chat sessions
// are stored. Checks both VS Code stable and Insiders paths.
func CopilotSessionDirs() []string {
	var dirs []string

	appData := os.Getenv("APPDATA")
	if runtime.GOOS != osWindows {
		// On macOS/Linux, VS Code stores data in different locations
		home, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		switch runtime.GOOS {
		case "darwin":
			appData = filepath.Join(home, "Library", "Application Support")
		default: // Linux
			appData = filepath.Join(home, ".config")
		}
	}

	if appData == "" {
		return nil
	}

	// Check both Code stable and Code Insiders
	variants := []string{"Code", "Code - Insiders"}
	for _, variant := range variants {
		wsDir := filepath.Join(appData, variant, "User", "workspaceStorage")
		if info, err := os.Stat(wsDir); err == nil && info.IsDir() {
			// Scan each workspace for chatSessions/ subdirectory
			entries, err := os.ReadDir(wsDir)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				chatDir := filepath.Join(wsDir, entry.Name(), "chatSessions")
				if info, err := os.Stat(chatDir); err == nil && info.IsDir() {
					dirs = append(dirs, chatDir)
				}
			}
		}
	}

	return dirs
}

// Ensure CopilotParser implements SessionParser.
var _ SessionParser = (*CopilotParser)(nil)
