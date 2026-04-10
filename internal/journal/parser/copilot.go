//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
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

	f, scanner, scanErr := openScanner(
		path, cfgCopilot.ScanBufMatchMax,
	)
	if scanErr != nil {
		return false
	}
	defer func() { _ = f.Close() }()

	if !scanner.Scan() {
		return false
	}

	var line copilotRawLine
	if unmarshalLineErr := json.Unmarshal(
		scanner.Bytes(), &line,
	); unmarshalLineErr != nil {
		return false
	}

	if line.Kind != copilotKindSnapshot {
		return false
	}

	var sess copilotRawSession
	if unmarshalSessionErr := json.Unmarshal(
		line.V, &sess,
	); unmarshalSessionErr != nil {
		return false
	}

	return sess.SessionID != "" && sess.Version > 0
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
	f, scanner, scanErr := openScanner(
		path, cfgCopilot.ScanBufMax,
	)
	if scanErr != nil {
		return nil, errParser.OpenFile(scanErr)
	}
	defer func() { _ = f.Close() }()

	var sess *copilotRawSession

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
			sess = &s

		case copilotKindScalarPatch:
			// Scalar property patch — apply to session
			if sess != nil {
				p.applyScalarPatch(sess, line.K, line.V)
			}

		case copilotKindObjectPatch:
			// Array/object patch — apply to session
			if sess != nil {
				p.applyPatch(sess, line.K, line.V)
			}
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return nil, errParser.ScanFile(scanErr)
	}

	if sess == nil {
		return nil, nil
	}

	// Resolve workspace folder from workspace.json next to chatSessions/
	cwd := p.resolveWorkspaceCWD(path)

	result := p.buildSession(sess, path, cwd)
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
