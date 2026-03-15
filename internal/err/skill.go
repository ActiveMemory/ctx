//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// SkillList wraps a failure to list embedded skill directories.
//
// Parameters:
//   - cause: the underlying error from the list operation
//
// Returns:
//   - error: "failed to list skills: <cause>"
func SkillList(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrSkillList), cause)
}

// SkillRead wraps a failure to read a skill's content.
//
// Parameters:
//   - name: Skill directory name that failed to read
//   - cause: the underlying error from the read operation
//
// Returns:
//   - error: "failed to read skill <name>: <cause>"
func SkillRead(name string, cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrSkillRead), name, cause,
	)
}
