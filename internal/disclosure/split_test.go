//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"testing"

	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	"github.com/ActiveMemory/ctx/internal/disclosure"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
)

// Two-entry staging with a "---" separator and blank lines between
// entries — exactly the whitespace ParseEntryBlocks would trim.
const stagingTwo = "## [2026-01-01-000000] Alpha\n\n**Context**: a.\n\n---\n\n" +
	"## [2026-01-02-000000] Beta\n\n**Context**: b.\n"

const (
	spanAlpha = "## [2026-01-01-000000] Alpha\n\n**Context**: a.\n\n---\n\n"
	spanBeta  = "## [2026-01-02-000000] Beta\n\n**Context**: b.\n"
)

func idFor(ts, title string) string {
	return ts + cfgDisc.IDSeparator + title
}

// T02: the byte-cut is lossless — a moved entry's span is verbatim
// (including the trailing blank lines and "---" separator that
// ParseEntryBlocks trims), and the remaining staging is byte-exact.
func TestSplitStaging(t *testing.T) {
	idA := idFor("2026-01-01-000000", "Alpha")
	idB := idFor("2026-01-02-000000", "Beta")

	t.Run("move first", func(t *testing.T) {
		moved, remaining, err := disclosure.SplitStaging(stagingTwo, []string{idA})
		if err != nil {
			t.Fatalf("SplitStaging err = %v", err)
		}
		if moved[idA] != spanAlpha {
			t.Errorf("moved[idA] = %q, want %q", moved[idA], spanAlpha)
		}
		if remaining != spanBeta {
			t.Errorf("remaining = %q, want %q", remaining, spanBeta)
		}
	})

	t.Run("move second", func(t *testing.T) {
		moved, remaining, err := disclosure.SplitStaging(stagingTwo, []string{idB})
		if err != nil {
			t.Fatalf("SplitStaging err = %v", err)
		}
		if moved[idB] != spanBeta {
			t.Errorf("moved[idB] = %q, want %q", moved[idB], spanBeta)
		}
		if remaining != spanAlpha {
			t.Errorf("remaining = %q, want %q", remaining, spanAlpha)
		}
	})

	t.Run("move both leaves empty staging", func(t *testing.T) {
		moved, remaining, err := disclosure.SplitStaging(stagingTwo, []string{idA, idB})
		if err != nil {
			t.Fatalf("SplitStaging err = %v", err)
		}
		if len(moved) != 2 {
			t.Errorf("moved has %d entries, want 2", len(moved))
		}
		if remaining != "" {
			t.Errorf("remaining = %q, want empty", remaining)
		}
	})

	t.Run("move none is identity", func(t *testing.T) {
		moved, remaining, err := disclosure.SplitStaging(stagingTwo, nil)
		if err != nil {
			t.Fatalf("SplitStaging err = %v", err)
		}
		if len(moved) != 0 {
			t.Errorf("moved has %d entries, want 0", len(moved))
		}
		if remaining != stagingTwo {
			t.Errorf("remaining = %q, want the input unchanged", remaining)
		}
	})

	t.Run("unknown id", func(t *testing.T) {
		_, _, err := disclosure.SplitStaging(stagingTwo,
			[]string{idFor("2026-09-09-090909", "Ghost")})
		if err != errDisc.ErrEntryNotInStaging {
			t.Errorf("err = %v, want ErrEntryNotInStaging", err)
		}
	})

	t.Run("duplicate id", func(t *testing.T) {
		_, _, err := disclosure.SplitStaging(stagingTwo, []string{idA, idA})
		if err != errDisc.ErrEntryAssignedTwice {
			t.Errorf("err = %v, want ErrEntryAssignedTwice", err)
		}
	})
}

// T02: FlattenPlan concatenates entry IDs in assignment order and refuses
// an assignment that moves nothing.
func TestFlattenPlan(t *testing.T) {
	idA := idFor("2026-01-01-000000", "Alpha")
	idB := idFor("2026-01-02-000000", "Beta")

	t.Run("ordered concat", func(t *testing.T) {
		ids, err := disclosure.FlattenPlan(disclosure.Plan{
			Kind: "learning",
			Assignments: []disclosure.Assignment{
				{Theme: "t1", Slug: "t1", Gist: "g1", Entries: []string{idA}},
				{Theme: "t2", Slug: "t2", Gist: "g2", Entries: []string{idB}},
			},
		})
		if err != nil {
			t.Fatalf("FlattenPlan err = %v", err)
		}
		if len(ids) != 2 || ids[0] != idA || ids[1] != idB {
			t.Errorf("ids = %v, want [%s %s]", ids, idA, idB)
		}
	})

	t.Run("empty assignment", func(t *testing.T) {
		_, err := disclosure.FlattenPlan(disclosure.Plan{
			Assignments: []disclosure.Assignment{
				{Theme: "t1", Slug: "t1", Gist: "g1", Entries: nil},
			},
		})
		if err != errDisc.ErrEmptyAssignment {
			t.Errorf("err = %v, want ErrEmptyAssignment", err)
		}
	})
}
