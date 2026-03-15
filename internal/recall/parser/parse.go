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
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/session"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/rc"
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

	s := &Session{
		ID:         id,
		Slug:       first.Slug,
		Tool:       session.ToolClaudeCode,
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
		s.Messages = append(s.Messages, msg)

		if msg.BelongsToUser() {
			s.TurnCount++
			if s.FirstUserMsg == "" && msg.Text != "" {
				// Truncate preview
				preview := msg.Text
				if len(preview) > 100 {
					preview = preview[:100] + token.Ellipsis
				}
				s.FirstUserMsg = preview
			}
		}

		s.TotalTokensIn += msg.TokensIn
		s.TotalTokensOut += msg.TokensOut

		// Check for errors in tool results
		for _, tr := range msg.ToolResults {
			if tr.IsError {
				s.HasErrors = true
			}
		}

		// Track model
		if raw.Message.Model != "" && s.Model == "" {
			s.Model = raw.Message.Model
		}
	}

	s.TotalTokens = s.TotalTokensIn + s.TotalTokensOut

	return s
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
		case claude.BlockText:
			if msg.Text != "" {
				msg.Text += token.NewlineLF
			}
			msg.Text += block.Text

		case claude.BlockThinking:
			if msg.Thinking != "" {
				msg.Thinking += token.NewlineLF
			}
			msg.Thinking += block.Thinking

		case claude.BlockToolUse:
			inputStr := ""
			if block.Input != nil {
				inputStr = string(block.Input)
			}
			msg.ToolUses = append(msg.ToolUses, ToolUse{
				ID:    block.ID,
				Name:  block.Name,
				Input: inputStr,
			})

		case claude.BlockToolResult:
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
		return []claudeRawBlock{{Type: claude.BlockText, Text: text}}
	}

	return nil
}

// parseMarkdownSession extracts a Session from Markdown content.
//
// Parameters:
//   - content: Raw Markdown content
//   - sourcePath: Path to the source file
//
// Returns:
//   - *Session: The parsed session, or nil if no session header found
func (p *MarkdownSessionParser) parseMarkdownSession(
	content string, sourcePath string,
) *Session {
	lines := strings.Split(content, token.NewlineLF)

	var headerLine string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if sessionHeader(trimmed) {
			headerLine = trimmed
			break
		}
	}

	if headerLine == "" {
		return nil
	}

	date, topic := parseSessionHeader(headerLine)

	// Derive a session ID from the filename (stable, OS-agnostic)
	base := filepath.Base(sourcePath)
	sessionID := strings.TrimSuffix(base, file.ExtMarkdown)

	// Parse date from header or fall back to file modification time
	startTime := parseSessionDate(date)
	if startTime.IsZero() {
		info, err := os.Stat(sourcePath)
		if err == nil {
			startTime = info.ModTime()
		} else {
			startTime = time.Now()
		}
	}

	// Extract sections
	sections := extractSections(lines)

	// Build messages from sections
	var messages []Message
	turnCount := 0

	// The session summary itself is treated as an assistant message
	var bodyParts []string
	for _, sec := range sections {
		if sec.body != "" {
			bodyParts = append(bodyParts, token.HeadingLevelTwoStart+
				sec.heading+token.NewlineLF+sec.body,
			)
		}
	}

	if len(bodyParts) > 0 {
		messages = append(messages, Message{
			ID:        sessionID + "-summary",
			Timestamp: startTime,
			Role:      claude.RoleAssistant,
			Text:      strings.Join(bodyParts, token.NewlineLF+token.NewlineLF),
		})
	}

	// The topic acts as the initial user message
	if topic != "" {
		turnCount = 1
		messages = append([]Message{{
			ID:        sessionID + "-topic",
			Timestamp: startTime,
			Role:      claude.RoleUser,
			Text:      topic,
		}}, messages...)
	}

	cwd := ""
	project := ""
	// Try to infer the project from the path
	// (look for .context/sessions/ pattern)
	d := filepath.Dir(sourcePath)
	if filepath.Base(d) == dir.Sessions {
		contextDir := filepath.Dir(d)
		if filepath.Base(contextDir) == dir.Context {
			projectDir := filepath.Dir(contextDir)
			project = filepath.Base(projectDir)
			cwd = projectDir
		}
	}

	return &Session{
		ID:           sessionID,
		Slug:         sessionID,
		Tool:         session.ToolMarkdown,
		SourceFile:   sourcePath,
		CWD:          cwd,
		Project:      project,
		StartTime:    startTime,
		EndTime:      startTime,
		Duration:     0,
		Messages:     messages,
		TurnCount:    turnCount,
		FirstUserMsg: topic,
	}
}

// parseSessionHeader extracts the date and topic from a session header line.
//
// Parameters:
//   - line: The full header line (e.g., "# Session: 2026-01-15 - Fix API")
//
// Returns:
//   - string: The date portion (e.g., "2026-01-15")
//   - string: The topic portion (e.g., "Fix API")
func parseSessionHeader(line string) (string, string) {
	// Remove "# " prefix
	rest := strings.TrimPrefix(line, token.HeadingLevelOneStart)

	// Remove session prefix if present (e.g., "Session: ")
	for _, prefix := range rc.SessionPrefixes() {
		rest = strings.TrimPrefix(rest, prefix+token.Space)
		rest = strings.TrimPrefix(rest, prefix)
	}

	rest = strings.TrimSpace(rest)

	// Split on em dash or hyphen topic separator
	for _, sep := range token.TopicSeparators {
		if idx := strings.Index(rest, sep); idx >= 0 {
			return strings.TrimSpace(
					rest[:idx]),
				strings.TrimSpace(rest[idx+len(sep):])
		}
	}

	// No separator found: treat the entire rest as a topic
	return "", rest
}

// parseSessionDate parses a date string into a time.Time.
//
// Supports YYYY-MM-DD format.
//
// Parameters:
//   - dateStr: Date string to parse
//
// Returns:
//   - time.Time: Parsed time, or zero value on failure
func parseSessionDate(dateStr string) time.Time {
	t, err := time.Parse(cfgTime.DateFormat, dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}

// extractSections extracts H2 sections from Markdown lines.
//
// Sections are returned in document order to ensure deterministic output.
//
// Parameters:
//   - lines: All lines of the Markdown file
//
// Returns:
//   - []section: Sections in document order
func extractSections(lines []string) []section {
	var sections []section
	var currentHeading string
	var currentBody []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, token.HeadingLevelTwoStart) {
			// Save previous section
			if currentHeading != "" {
				sections = append(sections, section{
					heading: currentHeading,
					body: strings.TrimSpace(
						strings.Join(currentBody, token.NewlineLF),
					),
				})
			}
			currentHeading = strings.TrimPrefix(trimmed, token.HeadingLevelTwoStart)
			currentBody = nil
			continue
		}

		if currentHeading != "" {
			currentBody = append(currentBody, line)
		}
	}

	// Save the last section
	if currentHeading != "" {
		sections = append(sections, section{
			heading: currentHeading,
			body: strings.TrimSpace(
				strings.Join(currentBody, token.NewlineLF),
			),
		})
	}

	return sections
}
