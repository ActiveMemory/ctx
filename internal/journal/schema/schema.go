//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/schema"
)

// Default returns the embedded schema derived from empirical
// analysis of CC versions 2.1.2 through 2.1.92.
//
// Returns:
//   - *Schema: the default schema definition
func Default() *Schema {
	return &Schema{
		Version:        cfgSchema.Version,
		CCVersionRange: cfgSchema.CCVersionRange,
		RecordTypes:    buildRecordTypes(),
		BlockTypes:     buildBlockTypes(),
	}
}

// KnownField reports whether name is a required or optional
// field for the given record type.
//
// Parameters:
//   - recordType: the record type to check
//   - name: the field name to look up
//
// Returns:
//   - bool: true if the field is declared for the type
func (s *Schema) KnownField(
	recordType, name string,
) bool {
	rs, ok := s.RecordTypes[recordType]
	if !ok {
		return false
	}
	for _, f := range rs.Required {
		if f == name {
			return true
		}
	}
	for _, f := range rs.Optional {
		if f == name {
			return true
		}
	}
	return false
}

// KnownRecordType reports whether the record type is
// recognized.
//
// Parameters:
//   - recordType: the type value to check
//
// Returns:
//   - bool: true if the record type is in the schema
func (s *Schema) KnownRecordType(recordType string) bool {
	_, ok := s.RecordTypes[recordType]
	return ok
}

// KnownBlockType reports whether the block type is recognized.
//
// Parameters:
//   - blockType: the content block type to check
//
// Returns:
//   - bool: true if the block type is in the schema
func (s *Schema) KnownBlockType(blockType string) bool {
	_, ok := s.BlockTypes[blockType]
	return ok
}
