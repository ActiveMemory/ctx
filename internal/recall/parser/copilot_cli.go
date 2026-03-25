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
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// CopilotCLIParser parses GitHub Copilot CLI session files.
//
// Copilot CLI stores sessions as JSONL files in ~/.copilot/sessions/
// (or $COPILOT_HOME/sessions/). Each file contains one session with
// JSONL-formatted messages similar to Claude Code's format.
type CopilotCLIParser struct{}

// NewCopilotCLIParser creates a new Copilot CLI session parser.
func NewCopilotCLIParser() *CopilotCLIParser {
	return &CopilotCLIParser{}
}

// Tool returns the tool identifier for this parser.
func (p *CopilotCLIParser) Tool() string {
	return session.ToolCopilotCLI
}

// Matches returns true if the file appears to be a Copilot CLI session file.
//
// Checks if the file has a .jsonl extension and lives in a Copilot CLI
// session directory (under ~/.copilot/ or $COPILOT_HOME).
func (p *CopilotCLIParser) Matches(path string) bool {
	if !strings.HasSuffix(path, file.ExtJSONL) {
		return false
	}

	// Must be under a .copilot directory (not chatSessions, which is VS Code)
	dir := filepath.Dir(path)
	if strings.Contains(dir, "chatSessions") {
		return false
	}

	// Check if this is under a .copilot directory
	if !strings.Contains(path, cfgHook.DirCopilotHome) {
		return false
	}

	// Verify the first line looks like a Copilot CLI session
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	if !scanner.Scan() {
		return false
	}

	var msg copilotCLIRawMessage
	if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
		return false
	}

	// A Copilot CLI session line must have a role and type
	return msg.Role != "" || msg.Type != ""
}

// ParseFile reads a Copilot CLI JSONL session file and returns sessions.
//
// Each file represents one session. Messages are parsed line by line.
func (p *CopilotCLIParser) ParseFile(path string) ([]*entity.Session, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 4*1024*1024) // 4MB per line

	var messages []copilotCLIRawMessage

	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		if len(lineBytes) == 0 {
			continue
		}

		var msg copilotCLIRawMessage
		if err := json.Unmarshal(lineBytes, &msg); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, fmt.Errorf("scan file: %w", scanErr)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	result := p.buildSession(messages, path)
	if result == nil {
		return nil, nil
	}

	return []*entity.Session{result}, nil
}

// ParseLine is not used for Copilot CLI sessions since each file
// represents a complete session. Returns nil.
func (p *CopilotCLIParser) ParseLine(_ []byte) (*entity.Message, string, error) {
	return nil, "", nil
}

// buildSession converts raw messages into a Session entity.
func (p *CopilotCLIParser) buildSession(
	msgs []copilotCLIRawMessage, sourcePath string,
) *entity.Session {
	if len(msgs) == 0 {
		return nil
	}

	sess := &entity.Session{
		ID:         filepath.Base(strings.TrimSuffix(sourcePath, file.ExtJSONL)),
		Tool:       session.ToolCopilotCLI,
		SourceFile: sourcePath,
	}

	for _, msg := range msgs {
		// Extract CWD from first message that has it
		if sess.CWD == "" && msg.CWD != "" {
			sess.CWD = msg.CWD
			sess.Project = filepath.Base(msg.CWD)
		}

		// Extract session ID if present
		if msg.SessionID != "" {
			sess.ID = msg.SessionID
		}

		// Extract model
		if sess.Model == "" && msg.Model != "" {
			sess.Model = msg.Model
		}

		// Set timestamps
		if !msg.Timestamp.IsZero() {
			if sess.StartTime.IsZero() {
				sess.StartTime = msg.Timestamp
			}
			sess.EndTime = msg.Timestamp
		}

		// Build entity message
		entityMsg := entity.Message{
			ID:        msg.ID,
			Timestamp: msg.Timestamp,
			Role:      msg.Role,
			Text:      msg.Text,
		}

		if msg.Role == claude.RoleUser {
			sess.TurnCount++
			if sess.FirstUserMsg == "" && msg.Text != "" {
				preview := msg.Text
				if len(preview) > 100 {
					preview = preview[:100] + "..."
				}
				sess.FirstUserMsg = preview
			}
		}

		sess.Messages = append(sess.Messages, entityMsg)
	}

	if !sess.StartTime.IsZero() && !sess.EndTime.IsZero() {
		sess.Duration = sess.EndTime.Sub(sess.StartTime)
	}

	return sess
}

// CopilotCLISessionDirs returns the directories where Copilot CLI sessions
// may be stored. Respects $COPILOT_HOME env var.
func CopilotCLISessionDirs() []string {
	var dirs []string

	copilotHome := os.Getenv(cfgHook.EnvCopilotHome)
	if copilotHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		copilotHome = filepath.Join(home, cfgHook.DirCopilotHome)
	}

	// Check common session subdirectories
	candidates := []string{"sessions", "history"}
	for _, sub := range candidates {
		dir := filepath.Join(copilotHome, sub)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			dirs = append(dirs, dir)
		}
	}

	// On Windows, also check under LOCALAPPDATA
	if runtime.GOOS == osWindows {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			for _, sub := range candidates {
				dir := filepath.Join(localAppData, "GitHub Copilot CLI", sub)
				if info, err := os.Stat(dir); err == nil && info.IsDir() {
					dirs = append(dirs, dir)
				}
			}
		}
	}

	return dirs
}

// Ensure CopilotCLIParser implements SessionParser.
var _ SessionParser = (*CopilotCLIParser)(nil)
