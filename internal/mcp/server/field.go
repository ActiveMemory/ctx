//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/entry"
)

// applyOptionalFields copies optional entry fields from MCP args
// to the params struct.
func applyOptionalFields(
	params *entry.Params,
	args map[string]interface{},
) {
	if v, ok := args[field.Priority].(string); ok {
		params.Priority = v
	}
	if v, ok := args[cli.AttrContext].(string); ok {
		params.Context = v
	}
	if v, ok := args[cli.AttrRationale].(string); ok {
		params.Rationale = v
	}
	if v, ok := args[cli.AttrConsequences].(string); ok {
		params.Consequences = v
	}
	if v, ok := args[cli.AttrLesson].(string); ok {
		params.Lesson = v
	}
	if v, ok := args[cli.AttrApplication].(string); ok {
		params.Application = v
	}
}
