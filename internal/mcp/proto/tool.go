//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package proto

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
)

// ToolDefs defines all available MCP tools.
var ToolDefs = []Tool{
	{
		Name: tool.Status,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolStatusDesc),
		InputSchema: InputSchema{Type: schema.Object},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: tool.Add,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolAddDesc),
		InputSchema: InputSchema{
			Type: schema.Object,
			Properties: mergeProps(map[string]Property{
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
			}, entryAttrProps(
				assets.TextDescKeyMCPToolPropContext)),
			Required: []string{cli.AttrType, field.Content},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name: tool.Complete,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolCompleteDesc),
		InputSchema: InputSchema{
			Type: schema.Object,
			Properties: map[string]Property{
				field.Query: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropQuery),
				},
			},
			Required: []string{field.Query},
		},
		Annotations: &ToolAnnotations{IdempotentHint: true},
	},
	{
		Name: tool.Drift,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolDriftDesc),
		InputSchema: InputSchema{Type: schema.Object},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: tool.Recall,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolRecallDesc),
		InputSchema: InputSchema{
			Type: schema.Object,
			Properties: map[string]Property{
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
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: tool.WatchUpdate,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolWatchUpdateDesc),
		InputSchema: InputSchema{
			Type: schema.Object,
			Properties: mergeProps(map[string]Property{
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
			}, entryAttrProps(
				assets.TextDescKeyMCPToolPropCtxBg)),
			Required: []string{cli.AttrType, field.Content},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name: tool.Compact,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolCompactDesc),
		InputSchema: InputSchema{
			Type: schema.Object,
			Properties: map[string]Property{
				field.Archive: {
					Type: schema.Boolean,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropArchive),
				},
			},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name: tool.Next,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolNextDesc),
		InputSchema: InputSchema{Type: schema.Object},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: tool.CheckTaskCompletion,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolCheckTaskDesc),
		InputSchema: InputSchema{
			Type: schema.Object,
			Properties: map[string]Property{
				field.RecentAction: {
					Type: schema.String,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropRecentAct),
				},
			},
		},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: tool.SessionEvent,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolSessionDesc),
		InputSchema: InputSchema{
			Type: schema.Object,
			Properties: map[string]Property{
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
		Annotations: &ToolAnnotations{},
	},
	{
		Name: tool.Remind,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolRemindDesc),
		InputSchema: InputSchema{Type: schema.Object},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
}
