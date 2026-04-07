//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"testing"
)

func TestDefaultSchema(t *testing.T) {
	s := Default()

	if s.Version == "" {
		t.Fatal("schema version is empty")
	}
	if s.CCVersionRange == "" {
		t.Fatal("CC version range is empty")
	}

	// Must have user and assistant record types.
	for _, rt := range []string{"user", "assistant"} {
		rs, ok := s.RecordTypes[rt]
		if !ok {
			t.Fatalf("missing record type: %s", rt)
		}
		if len(rs.Required) == 0 {
			t.Fatalf("record type %s has no required fields", rt)
		}
	}

	// Must have core block types.
	for _, bt := range []string{"text", "thinking", "tool_use", "tool_result"} {
		if _, ok := s.BlockTypes[bt]; !ok {
			t.Fatalf("missing block type: %s", bt)
		}
		if s.BlockTypes[bt] != BlockParsed {
			t.Fatalf("block type %s should be BlockParsed", bt)
		}
	}
}

func TestKnownField(t *testing.T) {
	s := Default()

	// Required fields are known.
	if !s.KnownField("user", "uuid") {
		t.Fatal("uuid should be known for user")
	}

	// Optional fields are known.
	if !s.KnownField("user", "gitBranch") {
		t.Fatal("gitBranch should be known for user")
	}

	// Unknown fields are not known.
	if s.KnownField("user", "fakeField") {
		t.Fatal("fakeField should not be known")
	}

	// Unknown record type returns false.
	if s.KnownField("bogus", "uuid") {
		t.Fatal("bogus record type should return false")
	}
}

func TestKnownRecordType(t *testing.T) {
	s := Default()

	for _, rt := range []string{
		"user", "assistant", "progress", "file-history-snapshot",
		"last-prompt", "attachment", "system",
	} {
		if !s.KnownRecordType(rt) {
			t.Fatalf("record type %s should be known", rt)
		}
	}

	if s.KnownRecordType("imaginary-type") {
		t.Fatal("imaginary-type should not be known")
	}
}

func TestKnownBlockType(t *testing.T) {
	s := Default()

	if !s.KnownBlockType("text") {
		t.Fatal("text should be known")
	}
	if !s.KnownBlockType("mcp_tool_use") {
		t.Fatal("mcp_tool_use should be known")
	}
	if s.KnownBlockType("alien_block") {
		t.Fatal("alien_block should not be known")
	}
}
