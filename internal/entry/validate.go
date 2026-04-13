//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/entity"
	errAdd "github.com/ActiveMemory/ctx/internal/err/add"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Validate checks that required fields are present for the given entry type.
//
// Parameters:
//   - params: Entry parameters to validate
//   - examplesFn: Function returning example text for a given type
//     (pass nil to omit examples from error messages)
//
// Returns:
//   - error: Non-nil with details about missing fields, nil if valid
func Validate(params entity.EntryParams, examplesFn func(string) string) error {
	if params.Content == "" {
		examples := ""
		if examplesFn != nil {
			examples = examplesFn(params.Type)
		}
		return errAdd.NoContentProvided(params.Type, examples)
	}

	// Provenance is required for task, decision, and learning
	// unless relaxed per-project via .ctxrc provenance_required.
	var provenance [][2]string
	if rc.ProvenanceSessionRequired() {
		provenance = append(provenance,
			[2]string{flag.PrefixLong + flag.SessionID, params.SessionID})
	}
	if rc.ProvenanceBranchRequired() {
		provenance = append(provenance,
			[2]string{flag.PrefixLong + flag.Branch, params.Branch})
	}
	if rc.ProvenanceCommitRequired() {
		provenance = append(provenance,
			[2]string{flag.PrefixLong + flag.Commit, params.Commit})
	}

	var extra [][2]string
	switch params.Type {
	case entry.Task:
		if params.Section == "" {
			return errAdd.SectionRequired()
		}

	case entry.Decision:
		extra = [][2]string{
			{flag.PrefixLong + flag.Context, params.Context},
			{flag.PrefixLong + flag.Rationale, params.Rationale},
			{flag.PrefixLong + flag.Consequence, params.Consequence},
		}

	case entry.Learning:
		extra = [][2]string{
			{flag.PrefixLong + flag.Context, params.Context},
			{flag.PrefixLong + flag.Lesson, params.Lesson},
			{flag.PrefixLong + flag.Application, params.Application},
		}
	}

	if params.Type == entry.Task ||
		params.Type == entry.Decision ||
		params.Type == entry.Learning {
		if m := checkRequired(
			append(provenance, extra...),
		); len(m) > 0 {
			return errAdd.MissingFields(params.Type, m)
		}
	}

	return nil
}
