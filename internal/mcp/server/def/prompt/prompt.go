//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	promptCfg "github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Defs defines all available MCP prompts.
var Defs = []proto.Prompt{
	{
		Name: promptCfg.SessionStart,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptSessionStartDesc),
	},
	{
		Name: promptCfg.AddDecision,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptAddDecisionDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrRationale,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionRat),
				Required:    true,
			},
			{
				Name:        cli.AttrConsequence,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionConseq),
				Required:    true,
			},
		},
	},
	{
		Name: promptCfg.AddLearning,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptAddLearningDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrLesson,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningLesson),
				Required:    true,
			},
			{
				Name:        cli.AttrApplication,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningApp),
				Required:    true,
			},
		},
	},
	{
		Name: promptCfg.Reflect,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptReflectDesc),
	},
	{
		Name: promptCfg.Checkpoint,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptCheckpointDesc),
	},
}
