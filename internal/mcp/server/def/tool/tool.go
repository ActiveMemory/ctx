//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tool

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	toolCfg "github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Defs defines all available MCP tools.
var Defs = []proto.Tool{
	{
		Name: toolCfg.Status,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolStatusDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.Add,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolAddDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: MergeProps(map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropType),
					Enum: []string{
						"task", "decision",
						"learning", "convention",
					},
				},
				field.Content: {
					Type: schema.String,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropContent),
				},
				field.Priority: {
					Type: schema.String,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropPriority),
					Enum: []string{"high", "medium", "low"},
				},
			}, EntryAttrProps(
				text.TextDescKeyMCPToolPropContext)),
			Required: []string{cli.AttrType, field.Content},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Complete,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolCompleteDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Query: {
					Type: schema.String,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropQuery),
				},
			},
			Required: []string{field.Query},
		},
		Annotations: &proto.ToolAnnotations{IdempotentHint: true},
	},
	{
		Name: toolCfg.Drift,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolDriftDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.Recall,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolRecallDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Limit: {
					Type: schema.Number,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropLimit),
				},
				field.Since: {
					Type: schema.String,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropSince),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.WatchUpdate,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolWatchUpdateDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: MergeProps(map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropEntryType),
				},
				field.Content: {
					Type: schema.String,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropMainContent),
				},
			}, EntryAttrProps(
				text.TextDescKeyMCPToolPropCtxBg)),
			Required: []string{cli.AttrType, field.Content},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Compact,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolCompactDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Archive: {
					Type: schema.Boolean,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropArchive),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Next,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolNextDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.CheckTaskCompletion,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolCheckTaskDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.RecentAction: {
					Type: schema.String,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropRecentAct),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.SessionEvent,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolSessionDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropEventType),
				},
				field.Caller: {
					Type: schema.String,
					Description: desc.TextDesc(
						text.TextDescKeyMCPToolPropCaller),
				},
			},
			Required: []string{cli.AttrType},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Remind,
		Description: desc.TextDesc(
			text.TextDescKeyMCPToolRemindDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
}
