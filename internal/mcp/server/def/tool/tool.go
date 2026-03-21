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
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	toolCfg "github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Defs defines all available MCP tools.
var Defs = []proto.Tool{
	{
		Name: toolCfg.Status,
		Description: desc.Text(
			text.DescKeyMCPToolStatusDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.Add,
		Description: desc.Text(
			text.DescKeyMCPToolAddDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: MergeProps(map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: desc.Text(
						text.DescKeyMCPToolPropType),
					Enum: []string{
						"task", "decision",
						"learning", "convention",
					},
				},
				field.Content: {
					Type: schema.String,
					Description: desc.Text(
						text.DescKeyMCPToolPropContent),
				},
				field.Priority: {
					Type: schema.String,
					Description: desc.Text(
						text.DescKeyMCPToolPropPriority),
					Enum: entry.Priorities,
				},
			}, EntryAttrProps(
				text.DescKeyMCPToolPropContext)),
			Required: []string{cli.AttrType, field.Content},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Complete,
		Description: desc.Text(
			text.DescKeyMCPToolCompleteDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Query: {
					Type: schema.String,
					Description: desc.Text(
						text.DescKeyMCPToolPropQuery),
				},
			},
			Required: []string{field.Query},
		},
		Annotations: &proto.ToolAnnotations{IdempotentHint: true},
	},
	{
		Name: toolCfg.Drift,
		Description: desc.Text(
			text.DescKeyMCPToolDriftDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.Recall,
		Description: desc.Text(
			text.DescKeyMCPToolRecallDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Limit: {
					Type: schema.Number,
					Description: desc.Text(
						text.DescKeyMCPToolPropLimit),
				},
				field.Since: {
					Type: schema.String,
					Description: desc.Text(
						text.DescKeyMCPToolPropSince),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.WatchUpdate,
		Description: desc.Text(
			text.DescKeyMCPToolWatchUpdateDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: MergeProps(map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: desc.Text(
						text.DescKeyMCPToolPropEntryType),
				},
				field.Content: {
					Type: schema.String,
					Description: desc.Text(
						text.DescKeyMCPToolPropMainContent),
				},
			}, EntryAttrProps(
				text.DescKeyMCPToolPropCtxBg)),
			Required: []string{cli.AttrType, field.Content},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Compact,
		Description: desc.Text(
			text.DescKeyMCPToolCompactDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.Archive: {
					Type: schema.Boolean,
					Description: desc.Text(
						text.DescKeyMCPToolPropArchive),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Next,
		Description: desc.Text(
			text.DescKeyMCPToolNextDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.CheckTaskCompletion,
		Description: desc.Text(
			text.DescKeyMCPToolCheckTaskDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				field.RecentAction: {
					Type: schema.String,
					Description: desc.Text(
						text.DescKeyMCPToolPropRecentAct),
				},
			},
		},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: toolCfg.SessionEvent,
		Description: desc.Text(
			text.DescKeyMCPToolSessionDesc),
		InputSchema: proto.InputSchema{
			Type: schema.Object,
			Properties: map[string]proto.Property{
				cli.AttrType: {
					Type: schema.String,
					Description: desc.Text(
						text.DescKeyMCPToolPropEventType),
				},
				field.Caller: {
					Type: schema.String,
					Description: desc.Text(
						text.DescKeyMCPToolPropCaller),
				},
			},
			Required: []string{cli.AttrType},
		},
		Annotations: &proto.ToolAnnotations{},
	},
	{
		Name: toolCfg.Remind,
		Description: desc.Text(
			text.DescKeyMCPToolRemindDesc),
		InputSchema: proto.InputSchema{Type: schema.Object},
		Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
	},
}
