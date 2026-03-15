//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import "fmt"

// NoSiteConfig returns an error when the zensical config file is missing.
//
// Parameters:
//   - dir: directory where the config was expected.
//
// Returns:
//   - error: "no zensical.toml found in <dir>"
func NoSiteConfig(dir string) error {
	return fmt.Errorf("no zensical.toml found in %s", dir)
}

// ZensicalNotFound returns an error when zensical is not installed.
//
// Returns:
//   - error: includes installation instructions
func ZensicalNotFound() error {
	return fmt.Errorf(
		"zensical not found. Install with: pipx install zensical (requires Python >= 3.10)",
	)
}
