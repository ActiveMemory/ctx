//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
)

// formatTranscriptEntry formats a single transcript entry as Markdown.
//
// Converts a transcript entry into readable Markdown with role headers,
// timestamps, and formatted content blocks. Handles text, thinking (as
// collapsible details), tool_use (with parameters), and tool_result blocks.
//
// Parameters:
//   - entry: Transcript entry to format
//
// Returns:
//   - string: Markdown-formatted representation of the entry
func formatTranscriptEntry(entry transcriptEntry) string {
	var sb strings.Builder
	nl := config.NewlineLF

	// Header with role and timestamp
	role := strings.ToUpper(entry.Message.Role[:1]) + entry.Message.Role[1:]
	sb.WriteString(fmt.Sprintf(config.TplSessionRoleHeading+nl+nl, role))

	if entry.Timestamp != "" {
		// Parse and format timestamp
		if t, parseErr := time.Parse(
			time.RFC3339, entry.Timestamp,
		); parseErr == nil {
			sb.WriteString(
				fmt.Sprintf(
					config.TplSessionTimestamp+nl+nl, t.Format("2006-01-02 15:04:05"),
				),
			)
		}
	}

	// Handle content
	switch content := entry.Message.Content.(type) {
	case string:
		sb.WriteString(content)
		sb.WriteString(nl)
	case []any:
		for _, block := range content {
			blockMap, isMap := block.(contentBlock)
			if !isMap {
				continue
			}

			blockType, _ := blockMap[config.ClaudeFieldType].(string)
			switch blockType {
			case config.ClaudeBlockText:
				formatTextBlock(&sb, blockMap)
			case config.ClaudeBlockThinking:
				formatThinkingBlock(&sb, blockMap)
			case config.ClaudeBlockToolUse:
				formatToolUseBlock(&sb, blockMap)
			case config.ClaudeBlockToolResult:
				formatToolResultBlock(&sb, blockMap)
			}
		}
	}

	return sb.String()
}
