//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package query

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/entity"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/journal/parser"
)

// FindSessions returns sessions for the current project, or all projects if
// allProjects is true.
//
// Parameters:
//   - allProjects: when true, scan all projects instead
//     of just the current one.
//
// Returns:
//   - []*entity.Session: matching sessions sorted by start time.
//   - error: non-nil if the working directory or session scan fails.
func FindSessions(allProjects bool) ([]*entity.Session, error) {
	if allProjects {
		return parser.FindSessions()
	}
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return nil, errFs.WorkingDirectory(cwdErr)
	}
	return parser.FindSessionsForCWD(cwd)
}
