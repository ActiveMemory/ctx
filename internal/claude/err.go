//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import "fmt"

// errSkillList wraps a failure to list embedded skill directories.
//
// Parameters:
//   - err: Underlying error from the list operation
//
// Returns:
//   - error: Wrapped error with format "failed to list skills: <cause>"
func errSkillList(err error) error {
	return fmt.Errorf("failed to list skills: %w", err)
}

// errSkillRead wraps a failure to read a skill's SKILL.md content.
//
// Parameters:
//   - name: Skill directory name that failed to read
//   - err: Underlying error from the read operation
//
// Returns:
//   - error: Wrapped error with format "failed to read skill <name>: <cause>"
func errSkillRead(name string, err error) error {
	return fmt.Errorf("failed to read skill %s: %w", name, err)
}
