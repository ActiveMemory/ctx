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
	cfgMcpPrompt "github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Defs defines all available MCP prompts.
var Defs = []proto.Prompt{
	{
		Name: cfgMcpPrompt.SessionStart,
		Description: desc.Text(
			text.DescKeyMCPPromptSessionStartDesc),
	},
	{
		Name: cfgMcpPrompt.AddDecision,
		Description: desc.Text(
			text.DescKeyMCPPromptAddDecisionDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: desc.Text(text.DescKeyMCPPromptArgDecisionTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: desc.Text(text.DescKeyMCPPromptArgDecisionCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrRationale,
				Description: desc.Text(text.DescKeyMCPPromptArgDecisionRat),
				Required:    true,
			},
			{
				Name:        cli.AttrConsequence,
				Description: desc.Text(text.DescKeyMCPPromptArgDecisionConseq),
				Required:    true,
			},
		},
	},
	{
		Name: cfgMcpPrompt.AddLearning,
		Description: desc.Text(
			text.DescKeyMCPPromptAddLearningDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: desc.Text(text.DescKeyMCPPromptArgLearningTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: desc.Text(text.DescKeyMCPPromptArgLearningCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrLesson,
				Description: desc.Text(text.DescKeyMCPPromptArgLearningLesson),
				Required:    true,
			},
			{
				Name:        cli.AttrApplication,
				Description: desc.Text(text.DescKeyMCPPromptArgLearningApp),
				Required:    true,
			},
		},
	},
	{
		Name: cfgMcpPrompt.Reflect,
		Description: desc.Text(
			text.DescKeyMCPPromptReflectDesc),
	},
	{
		Name: cfgMcpPrompt.Checkpoint,
		Description: desc.Text(
			text.DescKeyMCPPromptCheckpointDesc),
	},
}
