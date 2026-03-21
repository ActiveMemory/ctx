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

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
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
				desc.Text(text.DescKeyMCPLoadContext), loadErr))
	}

	var sb strings.Builder
	sb.WriteString(
		desc.Text(text.DescKeyMCPPromptSessionStartHeader),
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
			desc.Text(text.DescKeyMCPPromptSectionFormat),
			fileName, string(f.Content),
		)
	}

	sb.WriteString(token.NewlineLF)
	sb.WriteString(
		desc.Text(text.DescKeyMCPPromptSessionStartFooter),
	)

	return out.OkResponse(id, proto.GetPromptResult{
		Description: desc.Text(
			text.DescKeyMCPPromptSessionStartResultD,
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
		desc.Text(text.DescKeyMCPPromptCheckpointHeader),
	)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	_, _ = fmt.Fprintf(
		&sb,
		desc.Text(text.DescKeyMCPPromptCheckpointStatsFormat),
		toolCalls, adds, pending,
	)

	sb.WriteString(token.NewlineLF)
	sb.WriteString(
		desc.Text(text.DescKeyMCPPromptCheckpointSteps),
	)

	return out.OkResponse(id, proto.GetPromptResult{
		Description: desc.Text(
			text.DescKeyMCPPromptCheckpointResultD,
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
		KeyHeader:  text.DescKeyMCPPromptAddDecisionHeader,
		KeyFooter:  text.DescKeyMCPPromptAddDecisionFooter,
		FieldFmtK:  text.DescKeyMCPPromptAddDecisionFieldFmt,
		KeyResultD: text.DescKeyMCPPromptAddDecisionResultD,
		Fields: []promptdef.EntryField{
			{KeyLabel: text.DescKeyMCPPromptLabelDecision,
				Value: args[field.Content]},
			{KeyLabel: text.DescKeyMCPPromptLabelContext,
				Value: args[cli.AttrContext]},
			{KeyLabel: text.DescKeyMCPPromptLabelRationale,
				Value: args[cli.AttrRationale]},
			{KeyLabel: text.DescKeyMCPPromptLabelConsequence,
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
		KeyHeader:  text.DescKeyMCPPromptAddLearningHeader,
		KeyFooter:  text.DescKeyMCPPromptAddLearningFooter,
		FieldFmtK:  text.DescKeyMCPPromptAddLearningFieldFmt,
		KeyResultD: text.DescKeyMCPPromptAddLearningResultD,
		Fields: []promptdef.EntryField{
			{KeyLabel: text.DescKeyMCPPromptLabelLearning,
				Value: args[field.Content]},
			{KeyLabel: text.DescKeyMCPPromptLabelContext,
				Value: args[cli.AttrContext]},
			{KeyLabel: text.DescKeyMCPPromptLabelLesson,
				Value: args[cli.AttrLesson]},
			{KeyLabel: text.DescKeyMCPPromptLabelApplication,
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
		Description: desc.Text(
			text.DescKeyMCPPromptReflectResultD),
		Messages: []proto.PromptMessage{
			{
				Role: prompt.RoleUser,
				Content: proto.ToolContent{
					Type: mime.ContentTypeText,
					Text: desc.Text(
						text.DescKeyMCPPromptReflectBody,
					),
				},
			},
		},
	})
}
