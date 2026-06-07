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

// TestKnownField_PostV1FieldDrift pins the optional fields
// added to OptionalFields after the initial 1.0.0 schema —
// guards against silent regression of the drift-fix that
// landed on 2026-05-23. Each field name here corresponds to
// JSONL data observed in user-submitted journals from
// Claude Code versions beyond the 2.1.92 range covered by
// 1.0.0. If a future refactor drops one of these from
// OptionalFields, this test fires immediately and the
// schema-drift CLI starts complaining about it again.
func TestKnownField_PostV1FieldDrift(t *testing.T) {
	s := Default()
	for _, field := range []string{
		"interruptedMessageId",
		"attributionPlugin",
		"attributionSkill",
		"attributionMcpServer",
		"attributionMcpTool",
		"promptSource",
		"apiErrorStatus",
		"errorDetails",
	} {
		t.Run(field, func(t *testing.T) {
			for _, rt := range []string{"user", "assistant"} {
				if !s.KnownField(rt, field) {
					t.Errorf(
						"%q should be known for %q (post-1.0 drift fix)",
						field, rt,
					)
				}
			}
		})
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
