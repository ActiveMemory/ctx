//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package permissions

import "fmt"

// errSettingsNotFound returns an error when settings.local.json is missing.
func errSettingsNotFound() error {
	return fmt.Errorf("no .claude/settings.local.json found")
}

// errGoldenNotFound returns an error when settings.golden.json is missing.
func errGoldenNotFound() error {
	return fmt.Errorf("no .claude/settings.golden.json found — run 'ctx permissions snapshot' first")
}

// errReadFile wraps a file read failure.
func errReadFile(path string, err error) error {
	return fmt.Errorf("failed to read %s: %w", path, err)
}

// errWriteFile wraps a file write failure.
func errWriteFile(path string, err error) error {
	return fmt.Errorf("failed to write %s: %w", path, err)
}

// errParseSettings wraps a JSON parse failure for a settings file.
func errParseSettings(path string, err error) error {
	return fmt.Errorf("failed to parse %s: %w", path, err)
}
