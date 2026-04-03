//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/claude"
	cfgCopilot "github.com/ActiveMemory/ctx/internal/config/copilot"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// buildSession converts a reconstructed copilotRawSession into a Session.
//
// Parameters:
//   - raw: the reconstructed raw session data
//   - sourcePath: path to the JSONL source file
//   - cwd: resolved workspace directory
//
// Returns:
//   - *entity.Session: the built session, or nil if the session has no requests
func (p *Copilot) buildSession(
	raw *copilotRawSession, sourcePath string, cwd string,
) *entity.Session {
	if len(raw.Requests) == 0 {
		return nil
	}

	sess := &entity.Session{
		ID:         raw.SessionID,
		Tool:       session.ToolCopilot,
		SourceFile: sourcePath,
		CWD:        cwd,
		Project:    filepath.Base(cwd),
		StartTime:  time.UnixMilli(raw.CreationDate),
	}

	if raw.CustomTitle != "" {
		sess.Slug = raw.CustomTitle
	}

	for _, req := range raw.Requests {
		// User message
		userMsg := entity.Message{
			ID:        req.RequestID,
			Timestamp: time.UnixMilli(req.Timestamp),
			Role:      claude.RoleUser,
			Text:      req.Message.Text,
		}

		if req.Result != nil {
			userMsg.TokensIn = req.Result.Metadata.PromptTokens
		}

		sess.Messages = append(sess.Messages, userMsg)
		sess.TurnCount++

		if sess.FirstUserMsg == "" && userMsg.Text != "" {
			preview := userMsg.Text
			if len(preview) > session.PreviewMaxLen {
				preview = preview[:session.PreviewMaxLen] + token.Ellipsis
			}
			sess.FirstUserMsg = preview
		}

		// Assistant response
		assistantMsg := p.buildAssistantMessage(req)
		if assistantMsg != nil {
			sess.Messages = append(sess.Messages, *assistantMsg)

			if sess.Model == "" && req.ModelID != "" {
				sess.Model = req.ModelID
			}
		}

		// Accumulate tokens
		if req.Result != nil {
			sess.TotalTokensIn += req.Result.Metadata.PromptTokens
			sess.TotalTokensOut += req.Result.Metadata.OutputTokens
		}
	}

	sess.TotalTokens = sess.TotalTokensIn + sess.TotalTokensOut

	// Set end time from last request
	if last := raw.Requests[len(raw.Requests)-1]; last.Result != nil {
		sess.EndTime = time.UnixMilli(last.Timestamp).Add(
			time.Duration(last.Result.Timings.TotalElapsed) * time.Millisecond,
		)
	} else {
		sess.EndTime = time.UnixMilli(
			raw.Requests[len(raw.Requests)-1].Timestamp,
		)
	}
	sess.Duration = sess.EndTime.Sub(sess.StartTime)

	return sess
}

// buildAssistantMessage extracts the assistant response from a request.
//
// Parameters:
//   - req: the raw request containing response items
//
// Returns:
//   - *entity.Message: the assistant message, or nil if
//     the request has no response
func (p *Copilot) buildAssistantMessage(
	req copilotRawRequest,
) *entity.Message {
	if len(req.Response) == 0 {
		return nil
	}

	msg := &entity.Message{
		ID:        req.RequestID + cfgCopilot.ResponseSuffix,
		Timestamp: time.UnixMilli(req.Timestamp),
		Role:      claude.RoleAssistant,
	}

	if req.Result != nil {
		msg.TokensOut = req.Result.Metadata.OutputTokens
	}

	for _, item := range req.Response {
		switch item.Kind {
		case cfgCopilot.RespKindThinking:
			var text string
			if unmarshalThinkErr := json.Unmarshal(
				item.Value, &text,
			); unmarshalThinkErr == nil {
				if msg.Thinking != "" {
					msg.Thinking += token.NewlineLF
				}
				msg.Thinking += text
			}

		case cfgCopilot.RespKindToolInvoke:
			tu := p.parseToolInvocation(item)
			if tu != nil {
				msg.ToolUses = append(msg.ToolUses, *tu)
			}

		case "":
			// Plain markdown text (no kind field)
			var text string
			if unmarshalTextErr := json.Unmarshal(
				item.Value, &text,
			); unmarshalTextErr == nil {
				text = strings.TrimSpace(text)
				if text != "" {
					if msg.Text != "" {
						msg.Text += token.NewlineLF
					}
					msg.Text += text
				}
			}

			// Skip: codeblockUri, inlineReference, progressTaskSerialized,
			//        textEditGroup, undoStop, mcpServersStarting
		}
	}

	// Check for tool errors
	for _, tr := range msg.ToolResults {
		if tr.IsError {
			return msg // HasErrors is set at session level
		}
	}

	return msg
}

// parseToolInvocation extracts a ToolUse from a toolInvocationSerialized item.
//
// Parameters:
//   - item: the raw response item containing tool invocation data
//
// Returns:
//   - *entity.ToolUse: the parsed tool use, or nil if the item has no tool ID
func (p *Copilot) parseToolInvocation(item copilotRawRespItem) *entity.ToolUse {
	toolID := item.ToolID
	if toolID == "" {
		return nil
	}

	// Extract the tool name from toolId (e.g., "copilot_readFile" -> "readFile")
	name := toolID
	if idx := strings.LastIndex(toolID, cfgCopilot.ToolIDSeparator); idx >= 0 {
		name = toolID[idx+1:]
	}

	// Use invocationMessage as the input description
	inputStr := ""
	if item.InvocationMessage != nil {
		// InvocationMessage can be a string or object with value field
		var simple string
		if unmarshalStrErr := json.Unmarshal(
			item.InvocationMessage, &simple,
		); unmarshalStrErr == nil {
			inputStr = simple
		} else {
			var obj struct {
				Value string `json:"value"`
			}
			if unmarshalObjErr := json.Unmarshal(
				item.InvocationMessage, &obj,
			); unmarshalObjErr == nil {
				inputStr = obj.Value
			}
		}
	}

	return &entity.ToolUse{
		ID:    item.ToolCallID,
		Name:  name,
		Input: inputStr,
	}
}
