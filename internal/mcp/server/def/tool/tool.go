//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tool

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	toolCfg "github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Defs defines all available MCP tools.
var Defs = []proto.Tool{
	{
		Name: toolCfg.Status,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolStatusDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.Add,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolAddDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: MergeProps(map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropType),
					Enum: []string{
						"task", "decision",
						"learning", "convention",
					},
				},
				field.Content: {
					Type: schema.String,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropContent),
				},
				field.Priority: {
					Type: schema.String,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropPriority),
					Enum: []string{"high", "medium", "low"},
				},
			}, EntryAttrProps(
				embed.TextDescKeyMCPToolPropContext)),
			Required: []string{cli.AttrType, field.Content},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Complete,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolCompleteDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Query: {
					Type: schema.String,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropQuery),
				},
			},
			Required: []string{field.Query},
		},
		Annotations: &proto.ToolAnnotations{IdempotentHint: true},
	},
	{
		Name: toolCfg.Drift,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolDriftDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.Recall,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolRecallDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Limit: {
					Type: schema.Number,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropLimit),
				},
				field.Since: {
					Type: schema.String,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropSince),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.WatchUpdate,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolWatchUpdateDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: MergeProps(map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropEntryType),
				},
				field.Content: {
					Type: schema.String,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropMainContent),
				},
			}, EntryAttrProps(
				embed.TextDescKeyMCPToolPropCtxBg)),
			Required: []string{cli.AttrType, field.Content},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Compact,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolCompactDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Archive: {
					Type: schema.Boolean,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropArchive),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Next,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolNextDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.CheckTaskCompletion,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolCheckTaskDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.RecentAction: {
					Type: schema.String,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropRecentAct),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.SessionEvent,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolSessionDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropEventType),
				},
				field.Caller: {
					Type: schema.String,
					Description: assets.TextDesc(
						embed.TextDescKeyMCPToolPropCaller),
				},
			},
			Required: []string{cli.AttrType},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Remind,
		Description: assets.TextDesc(
			embed.TextDescKeyMCPToolRemindDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
}
