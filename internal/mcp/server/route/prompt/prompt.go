//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	"github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context/load"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	promptdef "github.com/ActiveMemory/ctx/internal/mcp/server/def/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
	"github.com/ActiveMemory/ctx/internal/mcp/server/stat"
)

// sessionStart loads context and provides session orientation.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - contextDir: path to the .context/ directory
//
// Returns:
//   - *proto.Response: rendered session start prompt with context files
func sessionStart(
	id json.RawMessage, contextDir string,
) *proto.Response {
	ctx, loadErr := load.Do(contextDir)
	if loadErr != nil {
		return out.ErrResponse(id, proto.ErrCodeInternal,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPLoadContext), loadErr))
	}

	var sb strings.Builder
	sb.WriteString(
		assets.TextDesc(assets.TextDescKeyMCPPromptSessionStartHeader),
	)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	for _, fileName := range ctxCfg.ReadOrder {
		f := ctx.File(fileName)
		if f == nil || f.IsEmpty {
			continue
		}
		_, _ = fmt.Fprintf(
			&sb,
			assets.TextDesc(assets.TextDescKeyMCPPromptSectionFormat),
			fileName, string(f.Content),
		)
	}

	sb.WriteString(token.NewlineLF)
	sb.WriteString(
		assets.TextDesc(assets.TextDescKeyMCPPromptSessionStartFooter),
	)

	return out.OkResponse(id, proto.GetPromptResult{
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptSessionStartResultD,
		),
		Messages: []proto.PromptMessage{
			{
				Role: prompt.RoleUser,
				Content: proto.ToolContent{
					Type: mime.ContentTypeText,
					Text: sb.String(),
				},
			},
		},
	})
}

// checkpoint summarizes progress and prepares for session end.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - toolCalls: number of tool calls in the session
//   - addsPerformed: map of entry type to add count
//   - pending: number of pending updates
//
// Returns:
//   - *proto.Response: checkpoint prompt with session stats
func checkpoint(
	id json.RawMessage, toolCalls int,
	addsPerformed map[string]int, pending int,
) *proto.Response {
	adds := stat.TotalAdds(addsPerformed)

	var sb strings.Builder
	sb.WriteString(
		assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointHeader),
	)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	_, _ = fmt.Fprintf(
		&sb,
		assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointStatsFormat),
		toolCalls, adds, pending,
	)

	sb.WriteString(token.NewlineLF)
	sb.WriteString(
		assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointSteps),
	)

	return out.OkResponse(id, proto.GetPromptResult{
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptCheckpointResultD,
		),
		Messages: []proto.PromptMessage{
			{
				Role: prompt.RoleUser,
				Content: proto.ToolContent{
					Type: mime.ContentTypeText,
					Text: sb.String(),
				},
			},
		},
	})
}

// addDecision formats a decision for recording.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: prompt arguments (content, context, rationale,
//     consequence)
//
// Returns:
//   - *proto.Response: formatted decision prompt
func addDecision(
	id json.RawMessage, args map[string]string,
) *proto.Response {
	return buildEntry(id, promptdef.EntryPromptSpec{
		KeyHeader:  assets.TextDescKeyMCPPromptAddDecisionHeader,
		KeyFooter:  assets.TextDescKeyMCPPromptAddDecisionFooter,
		FieldFmtK:  assets.TextDescKeyMCPPromptAddDecisionFieldFmt,
		KeyResultD: assets.TextDescKeyMCPPromptAddDecisionResultD,
		Fields: []promptdef.EntryField{
			{KeyLabel: assets.TextDescKeyMCPPromptLabelDecision,
				Value: args[field.Content]},
			{KeyLabel: assets.TextDescKeyMCPPromptLabelContext,
				Value: args[cli.AttrContext]},
			{KeyLabel: assets.TextDescKeyMCPPromptLabelRationale,
				Value: args[cli.AttrRationale]},
			{KeyLabel: assets.TextDescKeyMCPPromptLabelConsequence,
				Value: args[cli.AttrConsequence]},
		},
	})
}

// addLearning formats a learning for recording.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: prompt arguments (content, context, lesson,
//     application)
//
// Returns:
//   - *proto.Response: formatted learning prompt
func addLearning(
	id json.RawMessage, args map[string]string,
) *proto.Response {
	return buildEntry(id, promptdef.EntryPromptSpec{
		KeyHeader:  assets.TextDescKeyMCPPromptAddLearningHeader,
		KeyFooter:  assets.TextDescKeyMCPPromptAddLearningFooter,
		FieldFmtK:  assets.TextDescKeyMCPPromptAddLearningFieldFmt,
		KeyResultD: assets.TextDescKeyMCPPromptAddLearningResultD,
		Fields: []promptdef.EntryField{
			{KeyLabel: assets.TextDescKeyMCPPromptLabelLearning,
				Value: args[field.Content]},
			{KeyLabel: assets.TextDescKeyMCPPromptLabelContext,
				Value: args[cli.AttrContext]},
			{KeyLabel: assets.TextDescKeyMCPPromptLabelLesson,
				Value: args[cli.AttrLesson]},
			{KeyLabel: assets.TextDescKeyMCPPromptLabelApplication,
				Value: args[cli.AttrApplication]},
		},
	})
}

// reflect reviews the current session for outstanding items.
//
// Parameters:
//   - id: JSON-RPC request ID
//
// Returns:
//   - *proto.Response: reflection prompt text
func reflect(id json.RawMessage) *proto.Response {
	return out.OkResponse(id, proto.GetPromptResult{
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptReflectResultD),
		Messages: []proto.PromptMessage{
			{
				Role: prompt.RoleUser,
				Content: proto.ToolContent{
					Type: mime.ContentTypeText,
					Text: assets.TextDesc(
						assets.TextDescKeyMCPPromptReflectBody,
					),
				},
			},
		},
	})
}
