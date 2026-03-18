//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	promptCfg "github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Defs defines all available MCP prompts.
var Defs = []proto.Prompt{
	{
		Name: promptCfg.SessionStart,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPPromptSessionStartDesc),
	},
	{
		Name: promptCfg.AddDecision,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPPromptAddDecisionDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: assets.TextDesc(embed.TextDescKeyMCPPromptArgDecisionTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: assets.TextDesc(embed.TextDescKeyMCPPromptArgDecisionCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrRationale,
				Description: assets.TextDesc(embed.TextDescKeyMCPPromptArgDecisionRat),
				Required:    true,
			},
			{
				Name:        cli.AttrConsequence,
				Description: assets.TextDesc(embed.TextDescKeyMCPPromptArgDecisionConseq),
				Required:    true,
			},
		},
	},
	{
		Name: promptCfg.AddLearning,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPPromptAddLearningDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: assets.TextDesc(embed.TextDescKeyMCPPromptArgLearningTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: assets.TextDesc(embed.TextDescKeyMCPPromptArgLearningCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrLesson,
				Description: assets.TextDesc(embed.TextDescKeyMCPPromptArgLearningLesson),
				Required:    true,
			},
			{
				Name:        cli.AttrApplication,
				Description: assets.TextDesc(embed.TextDescKeyMCPPromptArgLearningApp),
				Required:    true,
			},
		},
	},
	{
		Name: promptCfg.Reflect,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPPromptReflectDesc),
	},
	{
		Name: promptCfg.Checkpoint,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPPromptCheckpointDesc),
	},
}
