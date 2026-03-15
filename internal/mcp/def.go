//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp"
)

// toolDefs defines all available MCP tools.
var toolDefs = []Tool{
	{
		Name: mcp.MCPToolStatus,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolStatusDesc),
		InputSchema: InputSchema{Type: mcp.MCPSchemaObject},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: mcp.MCPToolAdd,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolAddDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				cli.AttrType: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropType),
					Enum: []string{
						"task", "decision",
						"learning", "convention",
					},
				},
				mcp.MCPFieldContent: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropContent),
				},
				mcp.MCPFieldPriority: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropPriority),
					Enum: []string{"high", "medium", "low"},
				},
				cli.AttrContext: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropContext),
				},
				cli.AttrRationale: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropRationale),
				},
				cli.AttrConsequences: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropConseq),
				},
				cli.AttrLesson: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropLesson),
				},
				cli.AttrApplication: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropApplication),
				},
			},
			Required: []string{cli.AttrType, mcp.MCPFieldContent},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name: mcp.MCPToolComplete,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolCompleteDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				mcp.MCPFieldQuery: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropQuery),
				},
			},
			Required: []string{mcp.MCPFieldQuery},
		},
		Annotations: &ToolAnnotations{IdempotentHint: true},
	},
	{
		Name: mcp.MCPToolDrift,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolDriftDesc),
		InputSchema: InputSchema{Type: mcp.MCPSchemaObject},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: mcp.MCPToolRecall,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolRecallDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				mcp.MCPFieldLimit: {
					Type: mcp.MCPSchemaNumber,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropLimit),
				},
				mcp.MCPFieldSince: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropSince),
				},
			},
		},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: mcp.MCPToolWatchUpdate,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolWatchUpdateDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				cli.AttrType: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropEntryType),
				},
				mcp.MCPFieldContent: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropMainContent),
				},
				cli.AttrContext: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropCtxBg),
				},
				cli.AttrRationale: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropRationale),
				},
				cli.AttrConsequences: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropConseq),
				},
				cli.AttrLesson: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropLesson),
				},
				cli.AttrApplication: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropApplication),
				},
			},
			Required: []string{cli.AttrType, mcp.MCPFieldContent},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name: mcp.MCPToolCompact,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolCompactDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				mcp.MCPFieldArchive: {
					Type: mcp.MCPSchemaBoolean,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropArchive),
				},
			},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name: mcp.MCPToolNext,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolNextDesc),
		InputSchema: InputSchema{Type: mcp.MCPSchemaObject},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: mcp.MCPToolCheckTaskCompletion,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolCheckTaskDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				mcp.MCPFieldRecentAction: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropRecentAct),
				},
			},
		},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name: mcp.MCPToolSessionEvent,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolSessionDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				cli.AttrType: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropEventType),
				},
				mcp.MCPFieldCaller: {
					Type: mcp.MCPSchemaString,
					Description: assets.TextDesc(
						assets.TextDescKeyMCPToolPropCaller),
				},
			},
			Required: []string{cli.AttrType},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name: mcp.MCPToolRemind,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPToolRemindDesc),
		InputSchema: InputSchema{Type: mcp.MCPSchemaObject},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
}
