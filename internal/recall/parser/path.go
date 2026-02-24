//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"path/filepath"
	"strings"
)

// getPathRelativeToHome returns the path relative to the user's home directory.
// Returns an empty string if the path is not under a home directory.
func getPathRelativeToHome(path string) string {
	if path == "" {
		return ""
	}

	// Handle common home directory patterns
	// /home/username/... -> strip /home/username
	// /Users/username/... -> strip /Users/username (macOS)
	parts := strings.Split(path, string(filepath.Separator))

	for i, part := range parts {
		if part == "home" || part == "Users" {
			// Next part is username, rest is relative path
			if i+2 < len(parts) {
				return filepath.Join(parts[i+2:]...)
			}
			return ""
		}
	}

	return ""
}
