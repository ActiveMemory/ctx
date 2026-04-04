//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"errors"
	"os"
	"path/filepath"

	errSkill "github.com/ActiveMemory/ctx/internal/err/skill"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Remove deletes the skill directory for the given name from skillsDir.
// Returns an error if the skill does not exist.
func Remove(skillsDir, name string) error {
	dir := filepath.Join(skillsDir, name)

	info, statErr := ctxIo.SafeStat(dir)
	if statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return errSkill.NotFound(name)
		}
		return errSkill.Remove(name, statErr)
	}

	if !info.IsDir() {
		return errSkill.NotValidDir(name)
	}

	if removeErr := os.RemoveAll(dir); removeErr != nil {
		return errSkill.Remove(name, removeErr)
	}
	return nil
}
