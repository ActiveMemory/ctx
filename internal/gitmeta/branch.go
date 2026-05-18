//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gitmeta

import (
	"strings"

	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	cfgGitmeta "github.com/ActiveMemory/ctx/internal/config/gitmeta"
	execGit "github.com/ActiveMemory/ctx/internal/exec/git"
)

// resolveBranchOrDetached returns the current symbolic ref
// name, or the literal "detached" when HEAD is not on a
// branch. Failures (binary not on PATH, repo absent) collapse
// to "detached"; the caller treats this as best-effort
// metadata for provenance lines.
//
// Parameters:
//   - projectRoot: absolute path to the project root.
//
// Returns:
//   - string: branch name, or "detached".
func resolveBranchOrDetached(projectRoot string) string {
	out, runErr := execGit.Run(
		cfgGit.FlagChangeDir, projectRoot,
		cfgGit.RevParse, cfgGit.FlagShowCurrent,
	)
	if runErr != nil {
		return cfgGitmeta.BranchDetached
	}
	b := strings.TrimSpace(string(out))
	if b == "" {
		return cfgGitmeta.BranchDetached
	}
	return b
}
