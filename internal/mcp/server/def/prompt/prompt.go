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
			text.DescKeyMCPPromptSessionStartDesc),
	},
	{
		Name: promptCfg.AddDecision,
		Description: desc.TextDesc(
			text.DescKeyMCPPromptAddDecisionDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: desc.TextDesc(text.DescKeyMCPPromptArgDecisionTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: desc.TextDesc(text.DescKeyMCPPromptArgDecisionCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrRationale,
				Description: desc.TextDesc(text.DescKeyMCPPromptArgDecisionRat),
				Required:    true,
			},
			{
				Name:        cli.AttrConsequence,
				Description: desc.TextDesc(text.DescKeyMCPPromptArgDecisionConseq),
				Required:    true,
			},
		},
	},
	{
		Name: promptCfg.AddLearning,
		Description: desc.TextDesc(
			text.DescKeyMCPPromptAddLearningDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: desc.TextDesc(text.DescKeyMCPPromptArgLearningTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: desc.TextDesc(text.DescKeyMCPPromptArgLearningCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrLesson,
				Description: desc.TextDesc(text.DescKeyMCPPromptArgLearningLesson),
				Required:    true,
			},
			{
				Name:        cli.AttrApplication,
				Description: desc.TextDesc(text.DescKeyMCPPromptArgLearningApp),
				Required:    true,
			},
		},
	},
	{
		Name: promptCfg.Reflect,
		Description: desc.TextDesc(
			text.DescKeyMCPPromptReflectDesc),
	},
	{
		Name: promptCfg.Checkpoint,
		Description: desc.TextDesc(
			text.DescKeyMCPPromptCheckpointDesc),
	},
}
