//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package budget

import (
	"github.com/ActiveMemory/ctx/internal/cli/agent/core/extract"
	"github.com/ActiveMemory/ctx/internal/config/agent"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/index"
)

// ParseEntryBlocks parses a context file into entry blocks.
//
// Parameters:
//   - ctx: Loaded context
//   - fileName: Name of the file to parse (e.g., config.Decision)
//
// Returns:
//   - []index.EntryBlock: Parsed entry blocks; nil if the file is not found
func ParseEntryBlocks(ctx *entity.Context, fileName string) []index.EntryBlock {
	if f := ctx.File(fileName); f != nil {
		return index.ParseEntryBlocks(string(f.Content))
	}
	return nil
}

// ExtractAllConventions extracts all bullet items from CONVENTIONS.md
// (not limited to 5 like the old implementation).
//
// Parameters:
//   - ctx: Loaded context containing the files
//
// Returns:
//   - []string: All convention bullet items; nil if the file is not found
func ExtractAllConventions(ctx *entity.Context) []string {
	if f := ctx.File(cfgCtx.Convention); f != nil {
		return extract.BulletItems(string(f.Content), agent.BulletItemLimit)
	}
	return nil
}
