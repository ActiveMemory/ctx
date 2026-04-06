//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package extract

import (
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/entity"
	errMcp "github.com/ActiveMemory/ctx/internal/err/mcp"
)

// EntryArgs extracts required type/content from MCP args.
//
// Parameters:
//   - args: MCP tool arguments
//
// Returns:
//   - string: extracted entry type
//   - string: extracted content string
//   - error: non-nil if type or content is missing
func EntryArgs(
	args map[string]interface{},
) (string, string, error) {
	entryType, _ := args[cli.AttrType].(string)
	content, _ := args[field.Content].(string)

	if entryType == "" || content == "" {
		return "", "", errMcp.TypeContentRequired()
	}

	return entryType, content, nil
}

// Opts builds EntryOpts from MCP tool arguments.
//
// Parameters:
//   - args: MCP tool arguments with optional entry fields
//
// Returns:
//   - entity.EntryOpts: populated options struct
func Opts(args map[string]interface{}) entity.EntryOpts {
	opts := entity.EntryOpts{}
	if v, ok := args[field.Priority].(string); ok {
		opts.Priority = v
	}
	if v, ok := args[field.Section].(string); ok {
		opts.Section = v
	}
	if v, ok := args[cli.AttrContext].(string); ok {
		opts.Context = v
	}
	if v, ok := args[cli.AttrRationale].(string); ok {
		opts.Rationale = v
	}
	if v, ok := args[cli.AttrConsequence].(string); ok {
		opts.Consequence = v
	}
	if v, ok := args[cli.AttrLesson].(string); ok {
		opts.Lesson = v
	}
	if v, ok := args[cli.AttrApplication].(string); ok {
		opts.Application = v
	}
	if v, ok := args[field.SessionID].(string); ok {
		opts.SessionID = v
	}
	if v, ok := args[field.Branch].(string); ok {
		opts.Branch = v
	}
	if v, ok := args[field.Commit].(string); ok {
		opts.Commit = v
	}
	return opts
}
