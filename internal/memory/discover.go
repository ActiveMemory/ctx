//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DiscoverMemoryPath locates Claude Code's auto memory file for the
// given project root. The path is derived from how Claude Code encodes
// project directories: absolute path with "/" replaced by "-", prefixed
// with "-".
//
// Returns the resolved path if the file exists, or an error if auto
// memory has not been created yet.
func DiscoverMemoryPath(projectRoot string) (string, error) {
	abs, absErr := filepath.Abs(projectRoot)
	if absErr != nil {
		return "", fmt.Errorf("resolving project root: %w", absErr)
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", fmt.Errorf("resolving home directory: %w", homeErr)
	}

	slug := ProjectSlug(abs)
	memPath := filepath.Join(home, ".claude", "projects", slug, "memory", "MEMORY.md")

	if _, statErr := os.Stat(memPath); statErr != nil {
		return "", fmt.Errorf("no auto memory found at %s", memPath)
	}
	return memPath, nil
}

// ProjectSlug encodes an absolute project path into the Claude Code
// project directory slug format: "/" replaced by "-", prefixed with "-".
//
// Example: /home/jose/WORKSPACE/ctx → -home-jose-WORKSPACE-ctx
func ProjectSlug(absPath string) string {
	// Strip leading "/" then replace remaining "/" with "-", prefix with "-"
	return "-" + strings.ReplaceAll(absPath[1:], "/", "-")
}
