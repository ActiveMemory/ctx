//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package build

import (
	coreEntry "github.com/ActiveMemory/ctx/internal/cli/add/core/entry"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/validate"
)

// requireSection enforces the convention add-path's section contract.
//
// A convention is a bullet living under an H2 section, and the section is
// what the digesting pass folds into a theme file. There is deliberately
// no default: a catch-all section is the path of least resistance for an
// agent that has not decided where a convention belongs, and everything
// lands there. Choosing the section is the thinking; the CLI refuses to
// do it for the caller.
//
// Placeholder values ("TBD", "n/a", "pending", …) are refused for the
// same reason — they are the catch-all wearing a different hat.
// [validate.RejectPlaceholder] covers both: an empty value yields the
// flag-empty error, a placeholder the flag-placeholder error.
//
// Other nouns are unaffected: tasks resolve their own section default,
// and decisions/learnings do not group by section at all.
//
// Parameters:
//   - noun: the add subcommand's entry type
//   - section: the caller-supplied --section value
//
// Returns:
//   - error: a flag validation error, or nil when acceptable
func requireSection(noun, section string) error {
	if !coreEntry.FileTypeIsConvention(noun) {
		return nil
	}
	return validate.RejectPlaceholder(cFlag.Section, section)
}
