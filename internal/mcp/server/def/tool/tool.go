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
	cfgMcpTool "github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Defs returns all available MCP tool definitions.
//
// This is a function (not a package-level var) because desc.Text()
// reads from lookup maps that are populated by lookup.Init() in main().
// Package-level vars are initialized before main(), so desc.Text()
// would return empty strings.
//
// Returns:
//   - []proto.Tool: Complete set of MCP tool definitions
func Defs() []proto.Tool {
	return []proto.Tool{
		{
			Name: cfgMcpTool.Status,
			Description: desc.Text(
				text.DescKeyMCPToolStatusDesc),
			InputSchema: proto.InputSchema{Type: schema.Object},
			Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
		},
		{
			Name: cfgMcpTool.Add,
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
							entry.Task, entry.Decision,
							entry.Learning, entry.Convention,
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
					field.Section: {
						Type: schema.String,
						Description: desc.Text(
							text.DescKeyMCPToolPropSection),
					},
				}, EntryAttrProps(
					text.DescKeyMCPToolPropContext)),
				Required: []string{cli.AttrType, field.Content},
			},
			Annotations: &proto.ToolAnnotations{},
		},
		{
			Name: cfgMcpTool.Complete,
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
			Name: cfgMcpTool.Drift,
			Description: desc.Text(
				text.DescKeyMCPToolDriftDesc),
			InputSchema: proto.InputSchema{Type: schema.Object},
			Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
		},
		{
			Name: cfgMcpTool.JournalSource,
			Description: desc.Text(
				text.DescKeyMCPToolJournalSourceDesc),
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
			Name: cfgMcpTool.WatchUpdate,
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
			Name: cfgMcpTool.Compact,
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
			Name: cfgMcpTool.Next,
			Description: desc.Text(
				text.DescKeyMCPToolNextDesc),
			InputSchema: proto.InputSchema{Type: schema.Object},
			Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
		},
		{
			Name: cfgMcpTool.CheckTaskCompletion,
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
			Name: cfgMcpTool.SessionEvent,
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
			Name: cfgMcpTool.Remind,
			Description: desc.Text(
				text.DescKeyMCPToolRemindDesc),
			InputSchema: proto.InputSchema{Type: schema.Object},
			Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
		},
		{
			Name: cfgMcpTool.SteeringGet,
			Description: desc.Text(
				text.DescKeyMCPToolSteeringGetDesc),
			InputSchema: proto.InputSchema{
				Type: schema.Object,
				Properties: map[string]proto.Property{
					field.Prompt: {
						Type: schema.String,
						Description: desc.Text(
							text.DescKeyMCPToolPropPrompt),
					},
				},
			},
			Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
		},
		{
			Name: cfgMcpTool.Search,
			Description: desc.Text(
				text.DescKeyMCPToolSearchDesc),
			InputSchema: proto.InputSchema{
				Type: schema.Object,
				Properties: map[string]proto.Property{
					field.Query: {
						Type: schema.String,
						Description: desc.Text(
							text.DescKeyMCPToolPropSearchQuery),
					},
				},
				Required: []string{field.Query},
			},
			Annotations: &proto.ToolAnnotations{ReadOnlyHint: true},
		},
		{
			Name: cfgMcpTool.SessionStart,
			Description: desc.Text(
				text.DescKeyMCPToolSessionStartDesc),
			InputSchema: proto.InputSchema{Type: schema.Object},
			Annotations: &proto.ToolAnnotations{},
		},
		{
			Name: cfgMcpTool.SessionEnd,
			Description: desc.Text(
				text.DescKeyMCPToolSessionEndDesc),
			InputSchema: proto.InputSchema{
				Type: schema.Object,
				Properties: map[string]proto.Property{
					field.Summary: {
						Type: schema.String,
						Description: desc.Text(
							text.DescKeyMCPToolPropSummary),
					},
				},
			},
			Annotations: &proto.ToolAnnotations{},
		},
	}
}
