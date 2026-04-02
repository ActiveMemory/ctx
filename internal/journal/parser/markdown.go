//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"bufio"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/parser"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	errParser "github.com/ActiveMemory/ctx/internal/err/parser"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// MarkdownSession parses Markdown session files written by AI agents.
//
// This parser handles the tool-agnostic session format used by non-Claude
// tools (Copilot, Cursor, Aider, etc.) where the AI agent saves session
// summaries as structured Markdown in .context/sessions/.
//
// Expected format:
//
//	# Session: YYYY-MM-DD - Topic
//
//	## What Was Done
//	- ...
//
//	## Decisions
//	- ...
//
//	## Learnings
//	- ...
//
//	## Next Steps
//	- ...
type MarkdownSession struct{}

// NewMarkdownSession creates a new Markdown session parser.
//
// Returns:
//   - *MarkdownSession: A parser instance for Markdown session files
func NewMarkdownSession() *MarkdownSession {
	return &MarkdownSession{}
}

// Tool returns the tool identifier for this parser.
//
// Returns:
//   - string: The identifier "markdown"
func (p *MarkdownSession) Tool() string {
	return session.ToolMarkdown
}

// Matches returns true if the file appears to be a Markdown session file.
//
// Checks if the file has a .md extension and contains a session header
// in one of the recognized formats.
//
// Parameters:
//   - path: File path to check
//
// Returns:
//   - bool: True if this parser can handle the file
func (p *MarkdownSession) Matches(path string) bool {
	if !strings.HasSuffix(path, file.ExtMarkdown) {
		return false
	}

	// Skip README.md files
	base := filepath.Base(path)
	if strings.EqualFold(base, file.Readme) {
		return false
	}

	f, openErr := ctxIo.SafeOpenUserFile(path)
	if openErr != nil {
		return false
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			ctxLog.Warn(warn.Close, path, closeErr)
		}
	}()

	scanner := bufio.NewScanner(f)
	for i := 0; i < parser.LinesToPeek && scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		if sessionHeader(line) {
			return true
		}
	}

	return false
}

// ParseFile reads a Markdown session file and returns all sessions.
//
// Each file is treated as a single session. Metadata is extracted from
// the H1 header (date, topic) and H2 sections (content).
//
// Parameters:
//   - path: Path to the Markdown file to parse
//
// Returns:
//   - []*entity.Session: A single-element slice with the parsed session
//   - error: Non-nil if the file cannot be opened or read
func (p *MarkdownSession) ParseFile(path string) ([]*entity.Session, error) {
	content, readErr := ctxIo.SafeReadUserFile(filepath.Clean(path))
	if readErr != nil {
		return nil, errParser.ReadFile(readErr)
	}

	s := p.parseMarkdownSession(string(content), path)
	if s == nil {
		return nil, nil
	}

	return []*entity.Session{s}, nil
}

// ParseLine is not applicable for Markdown files (they are not line-oriented).
//
// Parameters:
//   - line: Ignored
//
// Returns:
//   - nil, "", nil always
func (p *MarkdownSession) ParseLine(_ []byte) (*entity.Message, string, error) {
	return nil, "", nil
}
