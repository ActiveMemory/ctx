//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package proto

import (
	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp/schema"
)

// mergeProps copies all entries from src into dst, returning dst.
func mergeProps(dst, src map[string]Property) map[string]Property {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// entryAttrProps returns the five entry-attribute properties shared by
// tools that accept structured context entries (add, watch-update).
// contextKey selects the TextDescKey for the "context" field since
// each tool describes it differently.
func entryAttrProps(contextKey string) map[string]Property {
	return map[string]Property{
		cli.AttrContext: {
			Type:        schema.String,
			Description: assets.TextDesc(contextKey),
		},
		cli.AttrRationale: {
			Type: schema.String,
			Description: assets.TextDesc(
				assets.TextDescKeyMCPToolPropRationale),
		},
		cli.AttrConsequence: {
			Type: schema.String,
			Description: assets.TextDesc(
				assets.TextDescKeyMCPToolPropConseq),
		},
		cli.AttrLesson: {
			Type: schema.String,
			Description: assets.TextDesc(
				assets.TextDescKeyMCPToolPropLesson),
		},
		cli.AttrApplication: {
			Type: schema.String,
			Description: assets.TextDesc(
				assets.TextDescKeyMCPToolPropApplication),
		},
	}
}
