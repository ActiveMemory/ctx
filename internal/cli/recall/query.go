//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"os"

	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// findSessions returns sessions for the current project, or all projects if
// allProjects is true.
//
// Parameters:
//   - allProjects: when true, scan all projects instead of just the current one.
//
// Returns:
//   - []*parser.Session: matching sessions sorted by start time.
//   - error: non-nil if the working directory or session scan fails.
func findSessions(allProjects bool) ([]*parser.Session, error) {
	if allProjects {
		return parser.FindSessions()
	}
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return nil, ctxerr.WorkingDirectory(cwdErr)
	}
	return parser.FindSessionsForCWD(cwd)
}
