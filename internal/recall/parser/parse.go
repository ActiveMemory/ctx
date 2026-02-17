//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"encoding/json"
	"path/filepath"
	"sort"

	"github.com/ActiveMemory/ctx/internal/config"
)

// buildSession constructs a Session from raw Claude Code messages.
//
// Parameters:
//   - id: Session ID to use
//   - rawMsgs: Raw messages belonging to this session
//   - sourcePath: Path to the source JSONL file
//
// Returns:
//   - *Session: Constructed session with messages, stats, and metadata
func (p *ClaudeCodeParser) buildSession(
	id string, rawMsgs []claudeRawMessage, sourcePath string,
) *Session {
	if len(rawMsgs) == 0 {
		return nil
	}

	// Sort by timestamp
	sort.Slice(rawMsgs, func(i, j int) bool {
		return rawMsgs[i].Timestamp.Before(rawMsgs[j].Timestamp)
	})

	first := rawMsgs[0]
	last := rawMsgs[len(rawMsgs)-1]

	session := &Session{
		ID:         id,
		Slug:       first.Slug,
		Tool:       config.ToolClaudeCode,
		SourceFile: sourcePath,
		CWD:        first.CWD,
		Project:    filepath.Base(first.CWD),
		GitBranch:  first.GitBranch,
		StartTime:  first.Timestamp,
		EndTime:    last.Timestamp,
		Duration:   last.Timestamp.Sub(first.Timestamp),
	}

	// Convert messages and accumulate stats
	for _, raw := range rawMsgs {
		msg := p.convertMessage(raw)
		session.Messages = append(session.Messages, msg)

		if msg.BelongsToUser() {
			session.TurnCount++
			if session.FirstUserMsg == "" && msg.Text != "" {
				// Truncate preview
				preview := msg.Text
				if len(preview) > 100 {
					preview = preview[:100] + "..."
				}
				session.FirstUserMsg = preview
			}
		}

		session.TotalTokensIn += msg.TokensIn
		session.TotalTokensOut += msg.TokensOut

		// Check for errors in tool results
		for _, tr := range msg.ToolResults {
			if tr.IsError {
				session.HasErrors = true
			}
		}

		// Track model
		if raw.Message.Model != "" && session.Model == "" {
			session.Model = raw.Message.Model
		}
	}

	session.TotalTokens = session.TotalTokensIn + session.TotalTokensOut

	return session
}

// convertMessage converts a Claude Code raw message to the common Message type.
//
// Parameters:
//   - raw: Raw message from JSONL parsing
//
// Returns:
//   - Message: Normalized message with text, tool uses,
//     and tool results extracted
func (p *ClaudeCodeParser) convertMessage(raw claudeRawMessage) Message {
	msg := Message{
		ID:        raw.UUID,
		Timestamp: raw.Timestamp,
		Role:      raw.Type,
	}

	if raw.Message.Usage != nil {
		msg.TokensIn = raw.Message.Usage.InputTokens
		msg.TokensOut = raw.Message.Usage.OutputTokens
	}

	// Parse content - can be a string or array of blocks
	blocks := p.parseContentBlocks(raw.Message.Content)

	// Extract content from blocks
	for _, block := range blocks {
		switch block.Type {
		case config.ClaudeBlockText:
			if msg.Text != "" {
				msg.Text += config.NewlineLF
			}
			msg.Text += block.Text

		case config.ClaudeBlockThinking:
			if msg.Thinking != "" {
				msg.Thinking += config.NewlineLF
			}
			msg.Thinking += block.Thinking

		case config.ClaudeBlockToolUse:
			inputStr := ""
			if block.Input != nil {
				inputStr = string(block.Input)
			}
			msg.ToolUses = append(msg.ToolUses, ToolUse{
				ID:    block.ID,
				Name:  block.Name,
				Input: inputStr,
			})

		case config.ClaudeBlockToolResult:
			contentStr := ""
			if block.Content != nil {
				// Try to unmarshal as JSON string first (handles escaping)
				var unescaped string
				if err := json.Unmarshal(block.Content, &unescaped); err == nil {
					contentStr = unescaped
				} else {
					// Fallback to raw bytes
					contentStr = string(block.Content)
				}
			}
			msg.ToolResults = append(msg.ToolResults, ToolResult{
				ToolUseID: block.ToolUseID,
				Content:   contentStr,
				IsError:   block.IsError,
			})
		}
	}

	return msg
}

// parseContentBlocks parses the content field, which can be a string or array.
//
// Parameters:
//   - content: Raw JSON content that may be a string or array of blocks
//
// Returns:
//   - []claudeRawBlock: Parsed content blocks
//     (text, thinking, tool_use, tool_result)
func (p *ClaudeCodeParser) parseContentBlocks(
	content json.RawMessage,
) []claudeRawBlock {
	if len(content) == 0 {
		return nil
	}

	// Try parsing as the array of blocks first
	var blocks []claudeRawBlock
	if err := json.Unmarshal(content, &blocks); err == nil {
		return blocks
	}

	// Try parsing as a simple string
	var text string
	if err := json.Unmarshal(content, &text); err == nil && text != "" {
		return []claudeRawBlock{{Type: config.ClaudeBlockText, Text: text}}
	}

	return nil
}
