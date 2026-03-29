//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	errAdd "github.com/ActiveMemory/ctx/internal/err/add"

	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/flag"
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
func Validate(params Params, examplesFn func(string) string) error {
	if params.Content == "" {
		examples := ""
		if examplesFn != nil {
			examples = examplesFn(params.Type)
		}
		return errAdd.NoContentProvided(params.Type, examples)
	}

	switch params.Type {
	case entry.Decision:
		if m := checkRequired([][2]string{
			{flag.PrefixLong + flag.Context, params.Context},
			{flag.PrefixLong + flag.Rationale, params.Rationale},
			{flag.PrefixLong + flag.Consequence, params.Consequence},
		}); len(m) > 0 {
			return errAdd.MissingFields(entry.Decision, m)
		}

	case entry.Learning:
		if m := checkRequired([][2]string{
			{flag.PrefixLong + flag.Context, params.Context},
			{flag.PrefixLong + flag.Lesson, params.Lesson},
			{flag.PrefixLong + flag.Application, params.Application},
		}); len(m) > 0 {
			return errAdd.MissingFields(entry.Learning, m)
		}
	}

	return nil
}
