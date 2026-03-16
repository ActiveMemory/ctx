//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// PromptDefs defines all available MCP prompts.
var PromptDefs = []proto.Prompt{
	{
		Name: prompt.SessionStart,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptSessionStartDesc),
	},
	{
		Name: prompt.AddDecision,
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
				Name:        cli.AttrConsequences,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionConseq),
				Required:    true,
			},
		},
	},
	{
		Name: prompt.AddLearning,
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
		Name: prompt.Reflect,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptReflectDesc),
	},
	{
		Name: prompt.Checkpoint,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptCheckpointDesc),
	},
}
