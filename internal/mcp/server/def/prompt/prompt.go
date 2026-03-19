//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	promptCfg "github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Defs defines all available MCP prompts.
var Defs = []proto.Prompt{
	{
		Name: promptCfg.SessionStart,
		Description: desc.TextDesc(
			text.TextDescKeyMCPPromptSessionStartDesc),
	},
	{
		Name: promptCfg.AddDecision,
		Description: desc.TextDesc(
			text.TextDescKeyMCPPromptAddDecisionDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: desc.TextDesc(text.TextDescKeyMCPPromptArgDecisionTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: desc.TextDesc(text.TextDescKeyMCPPromptArgDecisionCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrRationale,
				Description: desc.TextDesc(text.TextDescKeyMCPPromptArgDecisionRat),
				Required:    true,
			},
			{
				Name:        cli.AttrConsequence,
				Description: desc.TextDesc(text.TextDescKeyMCPPromptArgDecisionConseq),
				Required:    true,
			},
		},
	},
	{
		Name: promptCfg.AddLearning,
		Description: desc.TextDesc(
			text.TextDescKeyMCPPromptAddLearningDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: desc.TextDesc(text.TextDescKeyMCPPromptArgLearningTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: desc.TextDesc(text.TextDescKeyMCPPromptArgLearningCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrLesson,
				Description: desc.TextDesc(text.TextDescKeyMCPPromptArgLearningLesson),
				Required:    true,
			},
			{
				Name:        cli.AttrApplication,
				Description: desc.TextDesc(text.TextDescKeyMCPPromptArgLearningApp),
				Required:    true,
			},
		},
	},
	{
		Name: promptCfg.Reflect,
		Description: desc.TextDesc(
			text.TextDescKeyMCPPromptReflectDesc),
	},
	{
		Name: promptCfg.Checkpoint,
		Description: desc.TextDesc(
			text.TextDescKeyMCPPromptCheckpointDesc),
	},
}
