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
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// MergeProps copies all entries from src into dst, returning dst.
//
// Parameters:
//   - dst: target property map
//   - src: source property map to merge from
//
// Returns:
//   - map[string]proto.Property: the dst map with src entries added
func MergeProps(dst, src map[string]proto.Property) map[string]proto.Property {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// EntryAttrProps returns the entry-attribute properties shared by
// tools that accept structured context entries (add, watch-update).
// contextKey selects the DescKey for the "context" field since
// each tool describes it differently. The task placement and
// provenance fields are declared here because entry validation
// requires them, and a client can only send what the schema
// advertises.
//
// Parameters:
//   - contextKey: DescKey selecting the description for the context field
//
// Returns:
//   - map[string]proto.Property: property map with context, rationale,
//     consequence, lesson, application, section, priority,
//     session_id, branch, and commit fields
func EntryAttrProps(contextKey string) map[string]proto.Property {
	return map[string]proto.Property{
		cli.AttrContext: {
			Type:        schema.String,
			Description: desc.Text(contextKey),
		},
		cli.AttrRationale: {
			Type: schema.String,
			Description: desc.Text(
				text.DescKeyMCPToolPropRationale),
		},
		cli.AttrConsequence: {
			Type: schema.String,
			Description: desc.Text(
				text.DescKeyMCPToolPropConseq),
		},
		cli.AttrLesson: {
			Type: schema.String,
			Description: desc.Text(
				text.DescKeyMCPToolPropLesson),
		},
		cli.AttrApplication: {
			Type: schema.String,
			Description: desc.Text(
				text.DescKeyMCPToolPropApplication),
		},
		field.Section: {
			Type: schema.String,
			Description: desc.Text(
				text.DescKeyMCPToolPropSection),
		},
		field.Priority: {
			Type: schema.String,
			Description: desc.Text(
				text.DescKeyMCPToolPropPriority),
			Enum: entry.Priorities,
		},
		field.SessionID: {
			Type: schema.String,
			Description: desc.Text(
				text.DescKeyMCPToolPropSessionID),
		},
		field.Branch: {
			Type: schema.String,
			Description: desc.Text(
				text.DescKeyMCPToolPropBranch),
		},
		field.Commit: {
			Type: schema.String,
			Description: desc.Text(
				text.DescKeyMCPToolPropCommit),
		},
	}
}
