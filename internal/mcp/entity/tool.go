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
	"github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// ToolDefs defines all available MCP tools.
var ToolDefs = []proto.Tool{
	{
		Name: tool.Status,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolStatusDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: tool.Add,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolAddDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: MergeProps(map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropType),
					Enum: []string{
						"task", "decision",
						"learning", "convention",
					},
				},
				field.Content: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropContent),
				},
				field.Priority: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropPriority),
					Enum: []string{"high", "medium", "low"},
				},
			}, EntryAttrProps(
				assets.TextDescKeyMCPToolPropContext)),
			Required: []string{cli.AttrType, field.Content},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: tool.Complete,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolCompleteDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Query: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropQuery),
				},
			},
			Required: []string{field.Query},
		},
		Annotations: &proto.ToolAnnotations{IdempotentHint: true},
	},
	{
		Name: tool.Drift,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolDriftDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: tool.Recall,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolRecallDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Limit: {
					Type: schema.Number,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropLimit),
				},
				field.Since: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropSince),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: tool.WatchUpdate,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolWatchUpdateDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: MergeProps(map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropEntryType),
				},
				field.Content: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropMainContent),
				},
			}, EntryAttrProps(
				assets.TextDescKeyMCPToolPropCtxBg)),
			Required: []string{cli.AttrType, field.Content},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: tool.Compact,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolCompactDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Archive: {
					Type: schema.Boolean,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropArchive),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: tool.Next,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolNextDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: tool.CheckTaskCompletion,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolCheckTaskDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.RecentAction: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropRecentAct),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: tool.SessionEvent,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolSessionDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropEventType),
				},
				field.Caller: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropCaller),
				},
			},
			Required: []string{cli.AttrType},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: tool.Remind,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolRemindDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
}
