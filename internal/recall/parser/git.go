//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"os"
	"os/exec"
	"strings"
)

// gitRemote returns the git remote origin URL for a directory.
//
// Runs `git remote get-url origin` in the given directory.
// Returns an empty string if the directory does not exist, is not a git
// repository, or has no remote named "origin".
//
// Errors are intentionally swallowed — this is a best-effort enrichment
// helper. Callers treat "" as "unknown" and proceed without a remote URL.
//
// Parameters:
//   - dir: Directory path to query for git remote
//
// Returns:
//   - string: Remote URL, or empty string on any error
func gitRemote(dir string) string {
	if dir == "" {
		return ""
	}

	if _, statErr := os.Stat(dir); statErr != nil {
		return ""
	}

	cmd := exec.Command("git", "-C", dir, "remote", "get-url", "origin")
	output, cmdErr := cmd.Output()
	if cmdErr != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}
