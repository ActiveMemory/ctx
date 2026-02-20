//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
)

// MarkdownSessionParser parses Markdown session files written by AI agents.
//
// This parser handles the tool-agnostic session format used by non-Claude
// tools (Copilot, Cursor, Aider, etc.) where the AI agent saves session
// summaries as structured Markdown in .context/sessions/.
//
// Expected format:
//
//	# Session: YYYY-MM-DD — Topic
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
type MarkdownSessionParser struct{}

// NewMarkdownSessionParser creates a new Markdown session parser.
//
// Returns:
//   - *MarkdownSessionParser: A parser instance for Markdown session files
func NewMarkdownSessionParser() *MarkdownSessionParser {
	return &MarkdownSessionParser{}
}

// Tool returns the tool identifier for this parser.
//
// Returns:
//   - string: The identifier "markdown"
func (p *MarkdownSessionParser) Tool() string {
	return config.ToolMarkdown
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
func (p *MarkdownSessionParser) Matches(path string) bool {
	if !strings.HasSuffix(path, config.ExtMarkdown) {
		return false
	}

	// Skip README.md files
	base := filepath.Base(path)
	if strings.EqualFold(base, config.FilenameReadme) {
		return false
	}

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return false
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for i := 0; i < config.ParserPeekLines && scanner.Scan(); i++ {
		line := strings.TrimSpace(scanner.Text())
		if isSessionHeader(line) {
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
//   - []*Session: A single-element slice with the parsed session
//   - error: Non-nil if the file cannot be opened or read
func (p *MarkdownSessionParser) ParseFile(path string) ([]*Session, error) {
	content, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	session := p.parseMarkdownSession(string(content), path)
	if session == nil {
		return nil, nil
	}

	return []*Session{session}, nil
}

// ParseLine is not applicable for Markdown files (they are not line-oriented).
//
// Parameters:
//   - line: Ignored
//
// Returns:
//   - nil, "", nil always
func (p *MarkdownSessionParser) ParseLine(_ []byte) (*Message, string, error) {
	return nil, "", nil
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
	lines := strings.Split(content, config.NewlineLF)

	var headerLine string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if isSessionHeader(trimmed) {
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
	sessionID := strings.TrimSuffix(base, config.ExtMarkdown)

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
			bodyParts = append(bodyParts, "## "+sec.heading+config.NewlineLF+sec.body)
		}
	}

	if len(bodyParts) > 0 {
		messages = append(messages, Message{
			ID:        sessionID + "-summary",
			Timestamp: startTime,
			Role:      config.RoleAssistant,
			Text:      strings.Join(bodyParts, config.NewlineLF+config.NewlineLF),
		})
	}

	// The topic acts as the initial user message
	if topic != "" {
		turnCount = 1
		messages = append([]Message{{
			ID:        sessionID + "-topic",
			Timestamp: startTime,
			Role:      config.RoleUser,
			Text:      topic,
		}}, messages...)
	}

	cwd := ""
	project := ""
	// Try to infer project from the path (look for .context/sessions/ pattern)
	dir := filepath.Dir(sourcePath)
	if filepath.Base(dir) == "sessions" {
		contextDir := filepath.Dir(dir)
		if filepath.Base(contextDir) == config.DirContext {
			projectDir := filepath.Dir(contextDir)
			project = filepath.Base(projectDir)
			cwd = projectDir
		}
	}

	return &Session{
		ID:           sessionID,
		Slug:         sessionID,
		Tool:         config.ToolMarkdown,
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

// isSessionHeader checks if a line is a session header.
//
// Recognized formats:
//   - "# Session: YYYY-MM-DD — Topic"
//   - "# Session: YYYY-MM-DD - Topic"
//   - "# Oturum: YYYY-MM-DD — Topic" (Turkish)
//   - "# YYYY-MM-DD — Topic"
//   - "# YYYY-MM-DD - Topic"
//
// Parameters:
//   - line: Trimmed line to check
//
// Returns:
//   - bool: True if the line matches a session header pattern
func isSessionHeader(line string) bool {
	if !strings.HasPrefix(line, "# ") {
		return false
	}

	rest := line[2:]

	// Check for "Session:" or "Oturum:" prefix
	for _, prefix := range []string{"Session:", "Oturum:"} {
		if strings.HasPrefix(rest, prefix) {
			return true
		}
	}

	// Check for direct date pattern (YYYY-MM-DD)
	if len(rest) >= 10 && rest[4] == '-' && rest[7] == '-' {
		return true
	}

	return false
}

// parseSessionHeader extracts the date and topic from a session header line.
//
// Parameters:
//   - line: The full header line (e.g., "# Session: 2026-01-15 — Fix API")
//
// Returns:
//   - string: The date portion (e.g., "2026-01-15")
//   - string: The topic portion (e.g., "Fix API")
func parseSessionHeader(line string) (string, string) {
	// Remove "# " prefix
	rest := strings.TrimPrefix(line, "# ")

	// Remove "Session: " or "Oturum: " prefix if present
	for _, prefix := range []string{"Session: ", "Oturum: ", "Session:", "Oturum:"} {
		rest = strings.TrimPrefix(rest, prefix)
	}

	rest = strings.TrimSpace(rest)

	// Split on " — " (em dash) or " - " (hyphen)
	for _, sep := range []string{" \u2014 ", " - "} {
		if idx := strings.Index(rest, sep); idx >= 0 {
			return strings.TrimSpace(rest[:idx]), strings.TrimSpace(rest[idx+len(sep):])
		}
	}

	// No separator found — treat entire rest as topic
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
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return t
}

// section holds a heading and its body content in document order.
type section struct {
	heading string
	body    string
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
		if strings.HasPrefix(trimmed, "## ") {
			// Save previous section
			if currentHeading != "" {
				sections = append(sections, section{
					heading: currentHeading,
					body: strings.TrimSpace(
						strings.Join(currentBody, config.NewlineLF),
					),
				})
			}
			currentHeading = strings.TrimPrefix(trimmed, "## ")
			currentBody = nil
			continue
		}

		if currentHeading != "" {
			currentBody = append(currentBody, line)
		}
	}

	// Save last section
	if currentHeading != "" {
		sections = append(sections, section{
			heading: currentHeading,
			body: strings.TrimSpace(
				strings.Join(currentBody, config.NewlineLF),
			),
		})
	}

	return sections
}
