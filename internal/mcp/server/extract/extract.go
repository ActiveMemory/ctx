//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package extract

import (
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/entity"
	errMcp "github.com/ActiveMemory/ctx/internal/err/mcp"
	"github.com/ActiveMemory/ctx/internal/sanitize"
)

// EntryArgs extracts required type/content from MCP args.
//
// Validates that both fields are present and that content does not
// exceed MaxContentLen.
//
// Parameters:
//   - args: MCP tool arguments
//
// Returns:
//   - string: extracted entry type
//   - string: extracted content string
//   - error: non-nil if type or content is missing, or content too long
func EntryArgs(
	args map[string]interface{},
) (string, string, error) {
	entryType, _ := args[cli.AttrType].(string)
	content, _ := args[field.Content].(string)

	if entryType == "" || content == "" {
		return "", "", errMcp.TypeContentRequired()
	}

	// MCP-SAN.1: Enforce input length limits.
	if len(content) > cfg.MaxContentLen {
		return "", "", errMcp.InputTooLong(
			field.Content, cfg.MaxContentLen,
		)
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

// SanitizedOpts builds EntryOpts with content sanitization applied
// to all text fields.
//
// Parameters:
//   - args: MCP tool arguments with optional entry fields
//
// Returns:
//   - entity.EntryOpts: sanitized options struct
func SanitizedOpts(
	args map[string]interface{},
) entity.EntryOpts {
	opts := Opts(args)
	opts.Context = sanitize.Content(opts.Context)
	opts.Rationale = sanitize.Content(opts.Rationale)
	opts.Consequence = sanitize.Content(opts.Consequence)
	opts.Lesson = sanitize.Content(opts.Lesson)
	opts.Application = sanitize.Content(opts.Application)
	opts.SessionID = sanitize.SessionID(opts.SessionID)
	return opts
}
