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
	"runtime"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	cfgCopilot "github.com/ActiveMemory/ctx/internal/config/copilot"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	"github.com/ActiveMemory/ctx/internal/log/warn"
)

// CopilotCLI parses GitHub Copilot CLI session files.
//
// Copilot CLI stores sessions as JSONL files in ~/.copilot/sessions/
// (or $COPILOT_HOME/sessions/). Each file contains one session with
// JSONL-formatted messages similar to Claude Code's format.
type CopilotCLI struct{}

// NewCopilotCLI creates a new Copilot CLI session parser.
//
// Returns:
//   - *CopilotCLI: a new parser instance
func NewCopilotCLI() *CopilotCLI {
	return &CopilotCLI{}
}

// Tool returns the tool identifier for this parser.
//
// Returns:
//   - string: the Copilot CLI tool identifier
func (p *CopilotCLI) Tool() string {
	return session.ToolCopilotCLI
}

// Matches returns true if the file appears to be a Copilot CLI session file.
//
// Checks if the file has a .jsonl extension and lives in a Copilot CLI
// session directory (under ~/.copilot/ or $COPILOT_HOME), and the first
// line contains a valid Copilot CLI message with a role or type field.
//
// Parameters:
//   - path: file path to check
//
// Returns:
//   - bool: true if the file is a Copilot CLI session
func (p *CopilotCLI) Matches(path string) bool {
	if !strings.HasSuffix(path, file.ExtJSONL) {
		return false
	}

	// Must be under a .copilot directory (not chatSessions, which is VS Code)
	dir := filepath.Dir(path)
	if strings.Contains(dir, cfgCopilot.DirChatSessions) {
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
	buf := make([]byte, 0, cfgCopilot.ScanBufInit)
	scanner.Buffer(buf, cfgCopilot.ScanBufMatchMax)

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
// Each file represents one session. Messages are parsed line by line from
// the JSONL file and assembled into a single Session entity.
//
// Parameters:
//   - path: path to the JSONL session file
//
// Returns:
//   - []*entity.Session: the parsed sessions (at most one for Copilot CLI)
//   - error: any error encountered during parsing
func (p *CopilotCLI) ParseFile(path string) ([]*entity.Session, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errParser.OpenFile(err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			warn.Warn("copilot-cli: close %s: %v", path, closeErr)
		}
	}()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, cfgCopilot.ScanBufInit)
	scanner.Buffer(buf, cfgCopilot.ScanBufMax)

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
		return nil, errParser.ScanFile(scanErr)
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

// ParseLine is not meaningful for Copilot CLI sessions since each file
// represents a complete session. Returns nil for all lines.
//
// Parameters:
//   - line: the raw line bytes (unused)
//
// Returns:
//   - *entity.Message: always nil
//   - string: always empty
//   - error: always nil
func (p *CopilotCLI) ParseLine(_ []byte) (*entity.Message, string, error) {
	return nil, "", nil
}

// buildSession converts raw Copilot CLI messages into a Session entity.
//
// Iterates through all messages to extract metadata (CWD, model, timestamps)
// and assemble a complete session with turn counts and preview text.
//
// Parameters:
//   - msgs: raw messages parsed from the JSONL file
//   - sourcePath: path to the JSONL source file
//
// Returns:
//   - *entity.Session: the built session, or nil if msgs is empty
func (p *CopilotCLI) buildSession(
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
				if len(preview) > session.PreviewMaxLen {
					preview = preview[:session.PreviewMaxLen] + token.Ellipsis
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
// may be stored. Checks ~/.copilot/sessions and ~/.copilot/history, and
// on Windows also checks LOCALAPPDATA. Respects $COPILOT_HOME env var.
//
// Returns:
//   - []string: paths to session directories found on the system
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
	candidates := []string{cfgCopilot.DirSessions, cfgCopilot.DirHistory}
	for _, sub := range candidates {
		dir := filepath.Join(copilotHome, sub)
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			dirs = append(dirs, dir)
		}
	}

	// On Windows, also check under LOCALAPPDATA
	if runtime.GOOS == env.OSWindows {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData != "" {
			for _, sub := range candidates {
				dir := filepath.Join(localAppData, cfgCopilot.CLIAppName, sub)
				if info, err := os.Stat(dir); err == nil && info.IsDir() {
					dirs = append(dirs, dir)
				}
			}
		}
	}

	return dirs
}

// Ensure CopilotCLI implements Session.
var _ Session = (*CopilotCLI)(nil)
