//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// List wraps a failure to list embedded skill directories.
//
// Parameters:
//   - cause: the underlying error from the list operation
//
// Returns:
//   - error: "failed to list skills: <cause>"
func List(cause error) error {
	return fmt.Errorf(desc.Text(text.DescKeyErrSkillList), cause)
}

// Read wraps a failure to read a skill's content.
//
// Parameters:
//   - name: Skill directory name that failed to read
//   - cause: the underlying error from the read operation
//
// Returns:
//   - error: "failed to read skill <name>: <cause>"
func Read(name string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrSkillRead), name, cause,
	)
}
