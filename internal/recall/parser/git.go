//   /    Context:                     https://ctx.ist
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
// Returns an empty string if not a git repo or no remote configured.
func gitRemote(dir string) string {
	if dir == "" {
		return ""
	}

	// Check if the directory exists
	if _, err := os.Stat(dir); err != nil {
		return ""
	}

	// Try to get git remote
	cmd := exec.Command("git", "-C", dir, "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}
