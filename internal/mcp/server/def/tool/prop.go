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

// EntryAttrProps returns the five entry-attribute properties shared by
// tools that accept structured context entries (add, watch-update).
// contextKey selects the DescKey for the "context" field since
// each tool describes it differently.
//
// Parameters:
//   - contextKey: DescKey selecting the description for the context field
//
// Returns:
//   - map[string]proto.Property: property map with context, rationale,
//     consequence, lesson, and application fields
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
	}
}
