//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/schema"
)

// buildRecordTypes constructs field expectations per record
// type from config constants.
//
// Returns:
//   - map[string]RecordSchema: record type to field schema
func buildRecordTypes() map[string]RecordSchema {
	msgSchema := RecordSchema{
		Required: cfgSchema.RequiredFields,
		Optional: cfgSchema.OptionalFields,
	}

	rt := map[string]RecordSchema{
		cfgSchema.RecordUser:      msgSchema,
		cfgSchema.RecordAssistant: msgSchema,
	}

	for _, t := range cfgSchema.MetadataRecordTypes {
		rt[t] = RecordSchema{}
	}
	for _, t := range cfgSchema.InfraRecordTypes {
		rt[t] = RecordSchema{}
	}

	return rt
}

// buildBlockTypes constructs the known block type map from
// config constants.
//
// Returns:
//   - map[string]BlockKind: block type to parse kind
func buildBlockTypes() map[string]BlockKind {
	bt := make(map[string]BlockKind)
	for _, t := range cfgSchema.ParsedBlockTypes {
		bt[t] = BlockParsed
	}
	for _, t := range cfgSchema.KnownBlockTypes {
		bt[t] = BlockKnown
	}
	return bt
}
