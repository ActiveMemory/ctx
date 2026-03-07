//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import "fmt"

// ErrDirNotFound returns an error when the serve directory does not exist.
//
// Parameters:
//   - dir: Directory path that was not found
//
// Returns:
//   - error: Formatted error with the missing path
func ErrDirNotFound(dir string) error {
	return fmt.Errorf("directory not found: %s", dir)
}

// ErrNotDir returns an error when the path exists but is not a directory.
//
// Parameters:
//   - path: Path that is not a directory
//
// Returns:
//   - error: Formatted error with the path
func ErrNotDir(path string) error {
	return fmt.Errorf("not a directory: %s", path)
}

// ErrNoSiteConfig returns an error when the zensical config file is missing.
//
// Parameters:
//   - dir: Directory where the config was expected
//
// Returns:
//   - error: Formatted error with the directory path
func ErrNoSiteConfig(dir string) error {
	return fmt.Errorf("no zensical.toml found in %s", dir)
}

// ErrZensicalNotFound returns an error when zensical is not installed.
//
// Returns:
//   - error: Formatted error with install instructions
func ErrZensicalNotFound() error {
	return fmt.Errorf("zensical not found. Install with: pipx install zensical (requires Python >= 3.10)")
}
