//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// buildSession converts raw Copilot CLI messages into a Session
// entity.
//
// Iterates through all messages to extract metadata (CWD, model,
// timestamps) and assemble a complete session with turn counts
// and preview text.
//
// Parameters:
//   - msgs: raw messages parsed from the JSONL file
//   - sourcePath: path to the JSONL source file
//
// Returns:
//   - *entity.Session: the built session, or nil if empty
func (p *CopilotCLI) buildSession(
	msgs []copilotCLIRawMessage, sourcePath string,
) *entity.Session {
	if len(msgs) == 0 {
		return nil
	}

	sess := &entity.Session{
		ID: filepath.Base(
			strings.TrimSuffix(sourcePath, file.ExtJSONL),
		),
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
					preview = preview[:session.PreviewMaxLen] +
						token.Ellipsis
				}
				sess.FirstUserMsg = preview
			}
		}

		sess.Messages = append(sess.Messages, entityMsg)
	}

	if !sess.StartTime.IsZero() && !sess.EndTime.IsZero() {
		sess.Duration = sess.EndTime.Sub(
			sess.StartTime,
		)
	}

	return sess
}
