package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidateBoundary checks that dir resolves to a path within the current
// working directory. Returns an error if the resolved path escapes the
// project root.
func ValidateBoundary(dir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("validate boundary: %w", err)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("validate boundary: %w", err)
	}

	// Resolve symlinks in both paths so traversal via symlinked parents
	// is caught.
	resolvedCwd, err := filepath.EvalSymlinks(cwd)
	if err != nil {
		return fmt.Errorf("validate boundary: %w", err)
	}

	resolvedDir, err := filepath.EvalSymlinks(absDir)
	if err != nil {
		// If the target doesn't exist yet (e.g. before init), fall back
		// to the absolute path for the prefix check.
		resolvedDir = filepath.Clean(absDir)
	}

	// Ensure the resolved dir is equal to or nested under the project root.
	// Append os.PathSeparator to avoid "/foo/bar" matching "/foo/b".
	root := resolvedCwd + string(os.PathSeparator)
	if resolvedDir != resolvedCwd && !strings.HasPrefix(resolvedDir, root) {
		return fmt.Errorf("context directory %q resolves outside project root %q", dir, resolvedCwd)
	}

	return nil
}

// CheckSymlinks checks whether dir itself or any of its immediate children
// are symlinks. Returns an error describing the first symlink found.
func CheckSymlinks(dir string) error {
	// Check the directory itself.
	info, err := os.Lstat(dir)
	if err != nil {
		// Non-existent dir is not our concern — let the caller handle it.
		return nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("context directory %q is a symlink", dir)
	}

	// Check immediate children.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		child := filepath.Join(dir, entry.Name())
		ci, err := os.Lstat(child)
		if err != nil {
			continue
		}
		if ci.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("context file %q is a symlink", child)
		}
	}

	return nil
}
