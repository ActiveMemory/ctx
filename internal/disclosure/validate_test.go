//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"errors"
	"testing"

	"github.com/ActiveMemory/ctx/internal/disclosure"
	errDisc "github.com/ActiveMemory/ctx/internal/err/disclosure"
)

const (
	twoThemes = "# Learnings\n\n## Themes\n\n- a — g → [a](learnings/a.md)\n\n" +
		"## Themes\n\n- b — g → [b](learnings/b.md)\n"

	entryBelowThemes = "# Learnings\n\n## Themes\n\n- a — g → [a](learnings/a.md)\n\n" +
		"## [2026-07-15-120000] misplaced entry\n\n**Context**: below themes.\n"

	unparsableStaging = "# Learnings\n\n<!-- guide -->\n\n" +
		"## [not-a-real-timestamp] malformed\n\nbody.\n\n" +
		"## Themes\n\n- a — g → [a](learnings/a.md)\n"
)

// T06: the Validate precondition returns the named sentinel for each
// malformed shape, and nil for the two valid shapes (well-formed and
// not-yet-migrated).
func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		content string
		kind    disclosure.Kind
		want    error // nil = must pass
	}{
		{"well-formed populated", entryMigratedPopulated, disclosure.KindLearning, nil},
		{"un-migrated (zero themes)", entryUnmigrated, disclosure.KindLearning, nil},
		{"well-formed empty staging", entryMigratedEmpty, disclosure.KindLearning, nil},
		{"convention well-formed", conventionMigrated, disclosure.KindConvention, nil},
		{"two ## Themes", twoThemes, disclosure.KindLearning, errDisc.ErrMultipleThemes},
		{"entry below themes", entryBelowThemes, disclosure.KindLearning, errDisc.ErrEntryBelowThemes},
		{"unparsable staging", unparsableStaging, disclosure.KindLearning, errDisc.ErrStagingUnparsable},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := disclosure.Validate(disclosure.Parse(tc.content, tc.kind))
			switch {
			case tc.want == nil && got != nil:
				t.Errorf("Validate = %v, want nil (valid shape)", got)
			case tc.want != nil && !errors.Is(got, tc.want):
				t.Errorf("Validate = %v, want %v", got, tc.want)
			}
		})
	}
}
