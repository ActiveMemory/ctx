//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"os"
	"path/filepath"
	"strings"

	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
)

// walkForContextDir walks upward from the current working directory
// looking for an existing directory whose basename matches name.
//
// When a candidate is found above CWD, it is validated against the
// git root (if any). If the candidate falls outside the git root,
// it belongs to a different project and is discarded — the git root
// is used as the anchor instead.
//
// Absolute configured names skip the walk entirely. When no matching
// directory is found upward, returns the context directory anchored
// to the git root (if found) or filepath.Join(cwd, name) as an
// absolute path so that ctx init can create a fresh context directory
// at the current location.
//
// Parameters:
//   - name: Configured context directory name (may be relative or absolute)
//
// Returns:
//   - string: Absolute path to the resolved context directory
func walkForContextDir(name string) string {
	if filepath.IsAbs(name) {
		return name
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return name
	}

	// Walk upward looking for an existing context directory.
	var candidate string
	cur := cwd
	for {
		path := filepath.Join(cur, name)
		if info, statErr := os.Stat(path); statErr == nil && info.IsDir() {
			candidate = path
			break
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}

	gitRoot := findGitRoot(cwd)

	// No candidate found — anchor to git root or CWD.
	if candidate == "" {
		if gitRoot != "" {
			return filepath.Join(gitRoot, name)
		}
		return filepath.Join(cwd, name)
	}

	// Candidate found in CWD itself — always valid.
	candidateParent := filepath.Dir(candidate)
	if candidateParent == cwd {
		return candidate
	}

	// Candidate found above CWD — validate against git root.
	if gitRoot == "" {
		// No git root to confirm ownership; don't trust the ancestor.
		return filepath.Join(cwd, name)
	}

	// Check whether the candidate is within the git root.
	// Append separator to avoid "/foo/bar" matching "/foo/b".
	root := gitRoot + string(os.PathSeparator)
	if candidateParent == gitRoot || strings.HasPrefix(candidateParent, root) {
		return candidate
	}

	// Candidate is outside the git root — belongs to a different project.
	// Anchor to the git root instead.
	return filepath.Join(gitRoot, name)
}

// findGitRoot walks upward from start looking for a .git entry
// (directory or file, to support worktrees). Returns the parent
// directory of the .git entry, or "" if none is found.
//
// Parameters:
//   - start: Directory to start searching from
//
// Returns:
//   - string: Absolute path to the git root, or "" if not found
func findGitRoot(start string) string {
	cur := start
	for {
		gitPath := filepath.Join(cur, cfgGit.DotDir)
		if _, statErr := os.Stat(gitPath); statErr == nil {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return ""
		}
		cur = parent
	}
}
