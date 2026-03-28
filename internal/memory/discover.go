//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/memory"
	errMemory "github.com/ActiveMemory/ctx/internal/err/memory"
)

// DiscoverPath locates Claude Code's auto memory file for the
// given project root. The path is derived from how Claude Code encodes
// project directories: absolute path with "/" replaced by "-", prefixed
// with "-".
//
// Returns the resolved path if the file exists, or an error if auto
// memory has not been created yet.
//
// Parameters:
//   - projectRoot: Project root directory to derive the memory path from
func DiscoverPath(projectRoot string) (string, error) {
	abs, absErr := filepath.Abs(projectRoot)
	if absErr != nil {
		return "", errMemory.DiscoverResolveRoot(absErr)
	}

	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", errMemory.DiscoverResolveHome(homeErr)
	}

	slug := ProjectSlug(abs)
	memPath := filepath.Join(home, dir.Claude, dir.Projects, slug, dir.Memory, memory.Source)

	if _, statErr := os.Stat(memPath); statErr != nil {
		return "", errMemory.NoDiscovery(memPath)
	}
	return memPath, nil
}

// ProjectSlug encodes an absolute project path into the Claude Code
// project directory slug format: "/" replaced by "-", prefixed with "-".
//
// Example: /home/jose/WORKSPACE/ctx → -home-jose-WORKSPACE-ctx
//
// Parameters:
//   - absPath: Absolute project path to encode
//
// Returns:
//   - string: Slug-encoded path with dashes replacing separators
func ProjectSlug(absPath string) string {
	// Strip leading "/" then replace remaining "/" with "-", prefix with "-"
	return "-" + strings.ReplaceAll(absPath[1:], "/", "-")
}
