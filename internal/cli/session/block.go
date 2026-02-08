//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// formatTextBlock appends a text content block to sb.
//
// Parameters:
//   - sb: Builder to append to
//   - blockMap: Parsed content block
func formatTextBlock(sb *strings.Builder, blockMap contentBlock) {
	if text, hasText := blockMap[config.ClaudeFieldText].(string); hasText {
		sb.WriteString(text)
		sb.WriteString(config.NewlineLF)
	}
}

// formatThinkingBlock appends a collapsible thinking block to sb.
//
// Parameters:
//   - sb: Builder to append to
//   - blockMap: Parsed content block
func formatThinkingBlock(sb *strings.Builder, blockMap contentBlock) {
	nl := config.NewlineLF
	if thinking, hasThinking := blockMap[config.ClaudeFieldThinking].(string); hasThinking {
		sb.WriteString(config.TplSessionThinkingOpen + nl + nl)
		sb.WriteString(thinking)
		sb.WriteString(nl + config.TplRecallDetailsClose + nl + nl)
	}
}

// formatToolUseBlock appends a tool invocation block with parameters to sb.
//
// Parameters:
//   - sb: Builder to append to
//   - blockMap: Parsed content block
func formatToolUseBlock(sb *strings.Builder, blockMap contentBlock) {
	nl := config.NewlineLF
	name, _ := blockMap[config.ClaudeFieldName].(string)
	fmt.Fprintf(sb, config.TplSessionToolName+nl, name)
	if input, hasInput := blockMap[config.ClaudeFieldInput].(contentBlock); hasInput {
		for k, v := range input {
			vStr := fmt.Sprintf("%v", v)
			if len(vStr) > config.SessionMaxParamLen {
				vStr = vStr[:config.SessionMaxParamLen] + config.Ellipsis
			}
			fmt.Fprintf(sb, config.TplSessionToolParam+nl, k, vStr)
		}
	}
	sb.WriteString(nl)
}

// formatToolResultBlock appends a tool result block with optional truncation to sb.
//
// Parameters:
//   - sb: Builder to append to
//   - blockMap: Parsed content block
func formatToolResultBlock(sb *strings.Builder, blockMap contentBlock) {
	nl := config.NewlineLF
	sb.WriteString(config.TplSessionToolResult + nl)
	if result, hasResult := blockMap[config.ClaudeFieldContent].(string); hasResult {
		if len(result) > config.SessionMaxResultLen {
			result = result[:config.SessionMaxResultLen] + config.TplSessionTruncated
		}
		sb.WriteString(config.CodeFence + nl)
		sb.WriteString(result)
		sb.WriteString(nl + config.CodeFence + nl + nl)
	}
}
