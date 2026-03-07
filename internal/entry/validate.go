//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/write"
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
		return write.ErrNoContentProvided(params.Type, examples)
	}

	switch config.UserInputToEntry(params.Type) {
	case config.EntryDecision:
		if m := checkRequired([][2]string{
			{config.FlagPrefixLong + config.FlagContext, params.Context},
			{config.FlagPrefixLong + config.FlagRationale, params.Rationale},
			{config.FlagPrefixLong + config.FlagConsequences, params.Consequences},
		}); len(m) > 0 {
			return write.ErrMissingFields(config.EntryDecision, m)
		}

	case config.EntryLearning:
		if m := checkRequired([][2]string{
			{config.FlagPrefixLong + config.FlagContext, params.Context},
			{config.FlagPrefixLong + config.FlagLesson, params.Lesson},
			{config.FlagPrefixLong + config.FlagApplication, params.Application},
		}); len(m) > 0 {
			return write.ErrMissingFields(config.EntryLearning, m)
		}
	}

	return nil
}

// checkRequired returns the names of any fields whose values are empty.
func checkRequired(fields [][2]string) []string {
	var missing []string
	for _, f := range fields {
		if f[1] == "" {
			missing = append(missing, f[0])
		}
	}
	return missing
}
