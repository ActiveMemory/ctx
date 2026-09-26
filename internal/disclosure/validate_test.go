//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure_test

import (
	"errors"
	"strings"
	"testing"

	readTpl "github.com/ActiveMemory/ctx/internal/assets/read/template"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
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

	// M4/C11: the per-kind prefix means a bare "## " section stranded
	// below ## Themes is the convention-kind form of the same violation.
	conventionSectionBelowThemes = "# Conventions\n\n" +
		"## Themes\n\n- naming — g → [naming](conventions/naming.md)\n\n" +
		"## Misplaced Section\n\nprose below the themes region.\n"

	// M4/C3: with no timestamp, two sections sharing a title are
	// indistinguishable to a plan.
	conventionDuplicateTitle = "# Conventions\n\n<!-- guide -->\n\n" +
		"## Naming\n\nfirst.\n\n" +
		"## Naming\n\nsecond.\n\n" +
		"## Themes\n\n- naming — g → [naming](conventions/naming.md)\n"

	// A duplicate title in an entry root is legitimate: the timestamp
	// still separates the two entries.
	entryDuplicateTitle = "# Learnings\n\n<!-- guide -->\n\n" +
		"## [2026-07-15-120000] Same Title\n\nfirst.\n\n" +
		"## [2026-07-16-090000] Same Title\n\nsecond.\n\n" +
		"## Themes\n\n- a — g → [a](learnings/a.md)\n"

	// A date-only header (legacy or hand-written) is not an entry start
	// for the block parser, so without the guard it is silently folded
	// into the entry above it — and moved into that entry's theme. The
	// offending heading is on line 9.
	dateOnlyHeader = "# Learnings\n\n<!-- guide -->\n\n" +
		"## [2026-07-15-120000] a staged entry\n\n**Context**: x.\n\n" +
		"## [2026-09-26] Legacy\n\n**Context**: y.\n"

	// Conventions carry no timestamp: "## [" is ordinary title text.
	conventionBracketTitle = "# Conventions\n\n<!-- guide -->\n\n" +
		"## [Draft] Naming\n\nprose.\n"
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
		// M4/T10 (C11): entry-below-themes generalized to the per-kind prefix.
		{
			"convention section below themes", conventionSectionBelowThemes,
			disclosure.KindConvention, errDisc.ErrEntryBelowThemes,
		},
		// M4/T11 (C3): title-only identity must be unique.
		{
			"convention duplicate title", conventionDuplicateTitle,
			disclosure.KindConvention, errDisc.ErrDuplicateStagedTitle,
		},
		// The same duplicate title is fine for an entry kind: the
		// timestamp still distinguishes the two entries.
		{
			"entry duplicate title is allowed", entryDuplicateTitle,
			disclosure.KindLearning, nil,
		},
		// M4/T09: an un-migrated convention root is the first-run case.
		{
			"convention un-migrated", conventionUnmigrated,
			disclosure.KindConvention, nil,
		},
		// A malformed header is the precise form of unparsable staging.
		{
			"date-only header after an entry", dateOnlyHeader,
			disclosure.KindLearning, errDisc.ErrStagingUnparsable,
		},
		{
			"convention bracketed title", conventionBracketTitle,
			disclosure.KindConvention, nil,
		},
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

// DECISIONS.md's shipped template carries a "## [YYYY-MM-DD] Decision
// Title" example inside its <!-- DECISION FORMATS --> comment. That is
// documentation, not an entry: the malformed-header guard must skip it.
func TestValidate_DecisionTemplateExampleIgnored(t *testing.T) {
	tpl, tplErr := readTpl.Template(cfgCtx.Decision)
	if tplErr != nil {
		t.Fatalf("read template: %v", tplErr)
	}
	if !strings.Contains(string(tpl), "## [YYYY-MM-DD] Decision Title") {
		t.Fatal("template lost its commented date-only example; " +
			"this test no longer guards anything")
	}
	content := string(tpl) +
		"\n## [2026-09-26-120000] A real decision\n\n**Status**: Accepted\n"

	err := disclosure.Validate(disclosure.Parse(content, disclosure.KindDecision))
	if err != nil {
		t.Errorf("Validate = %v, want nil (commented example is not an entry)", err)
	}
}

// The malformed-header refusal names the 1-based line and the heading,
// on LF and CRLF files alike (the carriage return is not reported), and
// stays matchable as ErrStagingUnparsable.
func TestValidate_MalformedEntryHeader(t *testing.T) {
	const wantLine, wantHeading = 9, "## [2026-09-26] Legacy"
	for name, content := range map[string]string{
		"LF":   dateOnlyHeader,
		"CRLF": strings.ReplaceAll(dateOnlyHeader, "\n", "\r\n"),
	} {
		t.Run(name, func(t *testing.T) {
			err := disclosure.Validate(
				disclosure.Parse(content, disclosure.KindLearning),
			)
			mErr, ok := errors.AsType[*errDisc.MalformedEntryHeaderError](err)
			if !ok {
				t.Fatalf("Validate = %v, want *MalformedEntryHeaderError", err)
			}
			if mErr.Line != wantLine || mErr.Heading != wantHeading {
				t.Errorf("got line %d heading %q, want line %d heading %q",
					mErr.Line, mErr.Heading, wantLine, wantHeading)
			}
			if !errors.Is(err, errDisc.ErrStagingUnparsable) {
				t.Errorf("errors.Is(%v, ErrStagingUnparsable) = false", err)
			}
			msg := err.Error()
			if !strings.Contains(msg, "line 9") ||
				!strings.Contains(msg, wantHeading) {
				t.Errorf("message %q must name line 9 and the heading", msg)
			}
		})
	}
}
