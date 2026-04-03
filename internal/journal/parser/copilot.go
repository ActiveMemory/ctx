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

	cfgCopilot "github.com/ActiveMemory/ctx/internal/config/copilot"
	"github.com/ActiveMemory/ctx/internal/config/env"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/entity"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Ensure Copilot implements Session.
var _ Session = (*Copilot)(nil)

// Copilot parses VS Code Copilot Chat JSONL session files.
//
// Copilot Chat stores sessions as JSONL files in VS Code's workspaceStorage
// directory. Each file contains one session. The first line is a full session
// snapshot (kind=0), subsequent lines are incremental patches (kind=1, kind=2).
type Copilot struct{}

// NewCopilot creates a new Copilot Chat session parser.
//
// Returns:
//   - *Copilot: a new parser instance
func NewCopilot() *Copilot {
	return &Copilot{}
}

// Tool returns the tool identifier for this parser.
//
// Returns:
//   - string: the Copilot tool identifier
func (p *Copilot) Tool() string {
	return session.ToolCopilot
}

// Matches returns true if the file appears to be a Copilot Chat session file.
//
// Checks if the file has a .jsonl extension and lives in a chatSessions
// directory, and the first line contains a Copilot session snapshot.
//
// Parameters:
//   - path: file path to check
//
// Returns:
//   - bool: true if the file is a Copilot Chat session
func (p *Copilot) Matches(path string) bool {
	if !strings.HasSuffix(path, file.ExtJSONL) {
		return false
	}

	// Copilot sessions live in chatSessions/ directories
	if !strings.Contains(filepath.Dir(path), cfgCopilot.DirChatSessions) {
		return false
	}

	f, openErr := io.SafeOpenUserFile(path)
	if openErr != nil {
		return false
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, cfgCopilot.ScanBufInit)
	scanner.Buffer(buf, cfgCopilot.ScanBufMatchMax)

	if !scanner.Scan() {
		return false
	}

	var line copilotRawLine
	if unmarshalLineErr := json.Unmarshal(
		scanner.Bytes(), &line,
	); unmarshalLineErr != nil {
		return false
	}

	// kind=0 is the full session snapshot
	if line.Kind != copilotKindSnapshot {
		return false
	}

	var session copilotRawSession
	if unmarshalSessionErr := json.Unmarshal(
		line.V, &session,
	); unmarshalSessionErr != nil {
		return false
	}

	return session.SessionID != "" && session.Version > 0
}

// ParseFile reads a Copilot Chat JSONL file and returns the session.
//
// Reconstructs the session by reading the initial snapshot (kind=0) and
// applying incremental patches (kind=1 for scalar, kind=2 for array/object).
//
// Parameters:
//   - path: path to the JSONL session file
//
// Returns:
//   - []*entity.Session: the parsed sessions (at most one for Copilot)
//   - error: any error encountered during parsing
func (p *Copilot) ParseFile(path string) ([]*entity.Session, error) {
	f, openErr := io.SafeOpenUserFile(path)
	if openErr != nil {
		return nil, errParser.OpenFile(openErr)
	}
	defer func() { _ = f.Close() }()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, cfgCopilot.ScanBufInit)
	scanner.Buffer(buf, cfgCopilot.ScanBufMax)

	var session *copilotRawSession

	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		if len(lineBytes) == 0 {
			continue
		}

		var line copilotRawLine
		if unmarshalLineErr := json.Unmarshal(
			lineBytes, &line,
		); unmarshalLineErr != nil {
			continue
		}

		switch line.Kind {
		case copilotKindSnapshot:
			// Full session snapshot
			var s copilotRawSession
			if unmarshalSnapErr := json.Unmarshal(
				line.V, &s,
			); unmarshalSnapErr != nil {
				return nil, errParser.Unmarshal(unmarshalSnapErr)
			}
			session = &s

		case copilotKindScalarPatch:
			// Scalar property patch — apply to session
			if session != nil {
				p.applyScalarPatch(session, line.K, line.V)
			}

		case copilotKindObjectPatch:
			// Array/object patch — apply to session
			if session != nil {
				p.applyPatch(session, line.K, line.V)
			}
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, errParser.ScanFile(scanErr)
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
//
// Parameters:
//   - line: the raw line bytes (unused)
//
// Returns:
//   - *entity.Message: always nil
//   - string: always empty
//   - error: always nil
func (p *Copilot) ParseLine(_ []byte) (*entity.Message, string, error) {
	return nil, "", nil
}

// CopilotSessionDirs returns the directories where Copilot Chat sessions
// are stored. Checks both VS Code stable and Insiders paths.
//
// Returns:
//   - []string: paths to chatSessions directories found on the system
func CopilotSessionDirs() []string {
	var dirs []string

	appData := os.Getenv(cfgCopilot.EnvAppData)
	if runtime.GOOS != env.OSWindows {
		// On macOS/Linux, VS Code stores data in different locations
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return nil
		}
		switch runtime.GOOS {
		case cfgCopilot.OSDarwin:
			appData = filepath.Join(
				home, cfgCopilot.DirLibrary,
				cfgCopilot.DirAppSupport)
		default: // Linux
			appData = filepath.Join(home, cfgCopilot.DirDotConfig)
		}
	}

	if appData == "" {
		return nil
	}

	// Check both Code stable and Code Insiders
	variants := []string{cfgCopilot.AppCode, cfgCopilot.AppCodeInsiders}
	for _, variant := range variants {
		wsDir := filepath.Join(
			appData, variant,
			cfgCopilot.DirUser, cfgCopilot.DirWorkspace)
		if info, statErr := io.SafeStat(wsDir); statErr == nil && info.IsDir() {
			// Scan each workspace for chatSessions/ subdirectory
			entries, readDirErr := os.ReadDir(wsDir)
			if readDirErr != nil {
				continue
			}
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				chatDir := filepath.Join(
					wsDir, entry.Name(),
					cfgCopilot.DirChatSessions,
				)
				info, statChatErr := io.SafeStat(chatDir)
				if statChatErr == nil && info.IsDir() {
					dirs = append(dirs, chatDir)
				}
			}
		}
	}

	return dirs
}
