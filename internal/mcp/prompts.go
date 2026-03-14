//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context"
)

// promptDefs defines all available MCP prompts.
var promptDefs = []Prompt{
	{
		Name:        "ctx-session-start",
		Description: "Initialize a new session: loads full context and provides orientation",
	},
	{
		Name:        "ctx-add-decision",
		Description: "Record an architectural decision with context, rationale, and consequences",
		Arguments: []PromptArgument{
			{Name: "content", Description: "Decision title", Required: true},
			{Name: "context", Description: "Background context for the decision", Required: true},
			{Name: "rationale", Description: "Why this decision was made", Required: true},
			{Name: "consequences", Description: "Impact of the decision", Required: true},
		},
	},
	{
		Name:        "ctx-add-learning",
		Description: "Record a lesson learned with context, lesson, and application",
		Arguments: []PromptArgument{
			{Name: "content", Description: "Learning title", Required: true},
			{Name: "context", Description: "Background context", Required: true},
			{Name: "lesson", Description: "What was learned", Required: true},
			{Name: "application", Description: "How to apply this lesson", Required: true},
		},
	},
	{
		Name:        "ctx-reflect",
		Description: "Review the current session and capture outstanding learnings and decisions",
	},
	{
		Name:        "ctx-checkpoint",
		Description: "Summarize session progress and persist important context before ending",
	},
}

// handlePromptsList returns all available MCP prompts.
func (s *Server) handlePromptsList(req Request) *Response {
	return s.ok(req.ID, PromptListResult{Prompts: promptDefs})
}

// handlePromptsGet returns the content of a requested prompt.
func (s *Server) handlePromptsGet(req Request) *Response {
	var params GetPromptParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.error(req.ID, errCodeInvalidArg, "invalid params")
	}

	switch params.Name {
	case "ctx-session-start":
		return s.promptSessionStart(req.ID)
	case "ctx-add-decision":
		return s.promptAddDecision(req.ID, params.Arguments)
	case "ctx-add-learning":
		return s.promptAddLearning(req.ID, params.Arguments)
	case "ctx-reflect":
		return s.promptReflect(req.ID)
	case "ctx-checkpoint":
		return s.promptCheckpoint(req.ID)
	default:
		return s.error(req.ID, errCodeNotFound,
			fmt.Sprintf("unknown prompt: %s", params.Name))
	}
}

// promptSessionStart loads context and provides a session orientation.
func (s *Server) promptSessionStart(id json.RawMessage) *Response {
	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.error(id, errCodeInternal,
			fmt.Sprintf("failed to load context: %v", err))
	}

	var sb strings.Builder
	sb.WriteString("You are starting a new session. Read the following context files carefully.")
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	for _, fileName := range ctxCfg.ReadOrder {
		f := ctx.File(fileName)
		if f == nil || f.IsEmpty {
			continue
		}
		fmt.Fprintf(&sb, "## %s%s%s%s%s",
			fileName, token.NewlineLF, token.NewlineLF,
			string(f.Content), token.NewlineLF)
	}

	sb.WriteString(token.NewlineLF)
	sb.WriteString("Remember this context throughout the session. ")
	sb.WriteString("Use ctx_add to record decisions and learnings as you work. ")
	sb.WriteString("At session end, use ctx-checkpoint to capture outstanding context.")

	return s.ok(id, GetPromptResult{
		Description: "Session initialization with full context load",
		Messages: []PromptMessage{
			{
				Role:    "user",
				Content: ToolContent{Type: "text", Text: sb.String()},
			},
		},
	})
}

// promptAddDecision formats a decision for recording.
func (s *Server) promptAddDecision(
	id json.RawMessage, args map[string]string,
) *Response {
	content := args["content"]
	ctx := args["context"]
	rationale := args["rationale"]
	consequences := args["consequences"]

	var sb strings.Builder
	sb.WriteString("Record this architectural decision using ctx_add:")
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	fmt.Fprintf(&sb, "- **Decision**: %s%s", content, token.NewlineLF)
	fmt.Fprintf(&sb, "- **Context**: %s%s", ctx, token.NewlineLF)
	fmt.Fprintf(&sb, "- **Rationale**: %s%s", rationale, token.NewlineLF)
	fmt.Fprintf(&sb, "- **Consequences**: %s%s", consequences, token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	sb.WriteString("Call ctx_add with type=\"decision\" and all fields above.")

	return s.ok(id, GetPromptResult{
		Description: "Record an architectural decision",
		Messages: []PromptMessage{
			{
				Role:    "user",
				Content: ToolContent{Type: "text", Text: sb.String()},
			},
		},
	})
}

// promptAddLearning formats a learning for recording.
func (s *Server) promptAddLearning(
	id json.RawMessage, args map[string]string,
) *Response {
	content := args["content"]
	ctx := args["context"]
	lesson := args["lesson"]
	application := args["application"]

	var sb strings.Builder
	sb.WriteString("Record this learning using ctx_add:")
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	fmt.Fprintf(&sb, "- **Learning**: %s%s", content, token.NewlineLF)
	fmt.Fprintf(&sb, "- **Context**: %s%s", ctx, token.NewlineLF)
	fmt.Fprintf(&sb, "- **Lesson**: %s%s", lesson, token.NewlineLF)
	fmt.Fprintf(&sb, "- **Application**: %s%s", application, token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	sb.WriteString("Call ctx_add with type=\"learning\" and all fields above.")

	return s.ok(id, GetPromptResult{
		Description: "Record a lesson learned",
		Messages: []PromptMessage{
			{
				Role:    "user",
				Content: ToolContent{Type: "text", Text: sb.String()},
			},
		},
	})
}

// promptReflect reviews the current session for outstanding items.
func (s *Server) promptReflect(id json.RawMessage) *Response {
	var sb strings.Builder
	sb.WriteString("Reflect on this session and identify:")
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	sb.WriteString("1. **Decisions made** — Record each with ctx_add type=\"decision\"")
	sb.WriteString(token.NewlineLF)
	sb.WriteString("2. **Lessons learned** — Record each with ctx_add type=\"learning\"")
	sb.WriteString(token.NewlineLF)
	sb.WriteString("3. **Tasks completed** — Mark done with ctx_complete")
	sb.WriteString(token.NewlineLF)
	sb.WriteString("4. **New tasks identified** — Add with ctx_add type=\"task\"")
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	sb.WriteString("Review what was discussed and changed. ")
	sb.WriteString("Don't let important context slip away.")

	return s.ok(id, GetPromptResult{
		Description: "Review session for outstanding learnings and decisions",
		Messages: []PromptMessage{
			{
				Role:    "user",
				Content: ToolContent{Type: "text", Text: sb.String()},
			},
		},
	})
}

// promptCheckpoint summarizes progress and prepares for session end.
func (s *Server) promptCheckpoint(id json.RawMessage) *Response {
	pending := s.session.pendingCount()
	adds := totalAdds(s.session.addsPerformed)

	var sb strings.Builder
	sb.WriteString("Session checkpoint. Before ending this session:")
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	fmt.Fprintf(&sb, "- Tool calls this session: %d%s", s.session.toolCalls, token.NewlineLF)
	fmt.Fprintf(&sb, "- Entries added: %d%s", adds, token.NewlineLF)
	fmt.Fprintf(&sb, "- Pending updates: %d%s", pending, token.NewlineLF)

	sb.WriteString(token.NewlineLF)
	sb.WriteString("1. Check ctx_status for current context state")
	sb.WriteString(token.NewlineLF)
	sb.WriteString("2. Record any remaining decisions or learnings")
	sb.WriteString(token.NewlineLF)
	sb.WriteString("3. Mark completed tasks with ctx_complete")
	sb.WriteString(token.NewlineLF)
	sb.WriteString("4. Run ctx_compact if needed")
	sb.WriteString(token.NewlineLF)
	sb.WriteString("5. Call ctx_session_event type=\"end\" when done")

	return s.ok(id, GetPromptResult{
		Description: "Summarize session progress and persist context",
		Messages: []PromptMessage{
			{
				Role:    "user",
				Content: ToolContent{Type: "text", Text: sb.String()},
			},
		},
	})
}
