//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import "github.com/ActiveMemory/ctx/internal/entity"

// Session defines the interface for tool-specific session parsers.
//
// Each AI tool (Claude Code, Aider, Cursor) implements this interface
// to parse its specific format into the common Session type.
type Session interface {
	// ParseFile reads a session file and returns all sessions found.
	// A single file may contain multiple sessions (grouped by session ID).
	ParseFile(path string) ([]*entity.Session, error)

	// ParseLine parses a single line from a session file.
	// Returns nil if the line should be skipped (e.g., non-message lines).
	// Returns: message, sessionID, error.
	ParseLine(line []byte) (*entity.Message, string, error)

	// Matches returns true if this parser can handle the given file.
	// Implementations may check file extension, peek at content, etc.
	Matches(path string) bool

	// Tool returns the tool identifier (e.g., "claude-code", "aider").
	Tool() string
}
