//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// parseIndex attempts to parse a string as a positive integer index.
//
// Parameters:
//   - s: String to parse as an integer
//
// Returns:
//   - int: Parsed positive integer
//   - error: Non-nil if parsing fails or index is not positive
func parseIndex(s string) (int, error) {
	var idx int
	_, scanErr := fmt.Sscanf(s, "%d", &idx)
	if scanErr != nil {
		return 0, scanErr
	}
	if idx < 1 {
		return 0, errIndexNotPositive()
	}
	return idx, nil
}

// parseJsonlTranscript parses a .jsonl file and returns formatted Markdown.
//
// Reads a Claude Code JSONL transcript and converts it to readable Markdown
// with message headers, timestamps, and formatted content blocks.
//
// Parameters:
//   - path: Path to the JSONL transcript file
//
// Returns:
//   - string: Markdown-formatted transcript
//   - error: Non-nil if the file cannot be opened or read
func parseJsonlTranscript(path string) (result string, retErr error) {
	file, openErr := os.Open(filepath.Clean(path))
	if openErr != nil {
		return "", openErr
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && retErr == nil {
			retErr = closeErr
		}
	}()

	nl := config.NewlineLF
	var sb strings.Builder
	sb.WriteString(config.SessionHeadingTranscript + nl + nl)
	sb.WriteString(
		fmt.Sprintf(config.MetadataSource+" %s"+nl+nl, filepath.Base(path)),
	)
	sb.WriteString(config.Separator + nl + nl)

	scanner := bufio.NewScanner(file)
	// Increase buffer size for large lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024) // 10MB max line size

	messageCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var entry transcriptEntry
		if unmarshalErr := json.Unmarshal(
			[]byte(line), &entry,
		); unmarshalErr != nil {
			continue
		}

		// Skip non-message entries
		if entry.Type != config.RoleUser &&
			entry.Type != config.RoleAssistant {
			continue
		}

		messageCount++
		formatted := formatTranscriptEntry(entry)
		if formatted != "" {
			sb.WriteString(formatted)
			sb.WriteString(nl + config.Separator + nl + nl)
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return "", errReadingFile(scanErr)
	}

	sb.WriteString(
		fmt.Sprintf(config.TplSessionTotalMessages+nl, messageCount),
	)

	return sb.String(), nil
}

// parseSessionFile extracts metadata from a session file.
//
// Parses a session Markdown file to extract topic, date, type, and summary
// information. Handles both "# Session: topic" and "# topic" header formats.
//
// Parameters:
//   - path: Path to the session file
//
// Returns:
//   - sessionInfo: Parsed session metadata
//   - error: Non-nil if the file cannot be read
func parseSessionFile(path string) (sessionInfo, error) {
	content, readErr := os.ReadFile(filepath.Clean(path))
	if readErr != nil {
		return sessionInfo{}, readErr
	}

	nl := config.NewlineLF
	contentStr := string(content)
	info := sessionInfo{}

	// Extract topic from first line (# Session: topic)
	if strings.HasPrefix(contentStr, config.SessionHeadingPrefix) {
		lineEnd := strings.Index(contentStr, nl)
		if lineEnd != -1 {
			info.Topic = strings.TrimSpace(
				contentStr[len(config.SessionHeadingPrefix):lineEnd],
			)
		}
	} else if strings.HasPrefix(contentStr, config.HeadingLevelOneStart) {
		// Alternative format: # Topic
		lineEnd := strings.Index(contentStr, nl)
		if lineEnd != -1 {
			info.Topic = strings.TrimSpace(
				contentStr[len(config.HeadingLevelOneStart):lineEnd],
			)
		}
	}

	// Extract date
	if idx := strings.Index(contentStr, config.MetadataDate); idx != -1 {
		lineEnd := strings.Index(contentStr[idx:], nl)
		if lineEnd != -1 {
			info.Date = strings.TrimSpace(
				contentStr[idx+len(config.MetadataDate) : idx+lineEnd],
			)
		}
	}

	// Extract type
	if idx := strings.Index(contentStr, config.MetadataType); idx != -1 {
		lineEnd := strings.Index(contentStr[idx:], nl)
		if lineEnd != -1 {
			info.Type = strings.TrimSpace(
				contentStr[idx+len(config.MetadataType) : idx+lineEnd],
			)
		}
	}

	// Extract summary (first non-empty line after ## Summary)
	if idx := strings.Index(contentStr, config.RecallHeadingSummary); idx != -1 {
		afterSummary := contentStr[idx+len(config.RecallHeadingSummary):]
		lines := strings.Split(afterSummary, nl)
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" &&
				!strings.HasPrefix(line, config.PrefixHeading) &&
				!strings.HasPrefix(line, config.Separator) &&
				!strings.HasPrefix(line, config.PrefixBracket) {
				info.Summary = line
				break
			}
		}
	}

	return info, nil
}
