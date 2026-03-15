//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// CtxNotInPath returns an error indicating that ctx was not found in PATH.
//
// Returns:
//   - error: "ctx not found in PATH"
func CtxNotInPath() error {
	return errors.New(
		assets.TextDesc(assets.TextDescKeyErrValidationCtxNotInPath),
	)
}

// WorkingDirectory wraps a failure to determine the working directory.
//
// Parameters:
//   - cause: the underlying error from os.Getwd.
//
// Returns:
//   - error: "failed to get working directory: <cause>"
func WorkingDirectory(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidationWorkingDirectory), cause,
	)
}

// DriftViolations returns an error when drift detection found violations.
//
// Returns:
//   - error: "drift detection found violations"
func DriftViolations() error {
	return errors.New(
		assets.TextDesc(assets.TextDescKeyErrValidationDriftViolations),
	)
}

// FlagRequired returns an error for a missing required flag.
//
// Parameters:
//   - name: the flag name.
//
// Returns:
//   - error: "required flag \"<name>\" not set"
func FlagRequired(name string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidationFlagRequired), name,
	)
}

// ArgRequired returns an error for a missing required argument.
//
// Parameters:
//   - name: the argument name.
//
// Returns:
//   - error: "<name> argument is required"
func ArgRequired(name string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidationArgRequired), name,
	)
}

// ContextOutsideRoot returns an error when .context/ resolves outside the project root.
//
// Parameters:
//   - dir: the context directory path
//   - root: the project root path
//
// Returns:
//   - error: "context directory <dir> resolves outside project root <root>"
func ContextOutsideRoot(dir, root string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidateContextOutsideRoot), dir, root,
	)
}

// ContextDirSymlink returns an error when .context/ is a symlink.
//
// Parameters:
//   - dir: the context directory path
//
// Returns:
//   - error: "context directory <dir> is a symlink"
func ContextDirSymlink(dir string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidateContextDirSymlink), dir,
	)
}

// ContextFileSymlink returns an error when a file inside .context/ is a symlink.
//
// Parameters:
//   - file: the symlinked file path
//
// Returns:
//   - error: "context file <file> is a symlink"
func ContextFileSymlink(file string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidateContextFileSymlink), file,
	)
}

// InvalidSelection returns an error for an invalid interactive selection.
//
// Parameters:
//   - input: the user's input
//   - max: the maximum valid selection number
//
// Returns:
//   - error: "invalid selection: <input> (expected 1-<max>)"
func InvalidSelection(input string, max int) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidateInvalidSelection), input, max,
	)
}

// UnknownDocument returns an error for an unrecognized document alias.
//
// Parameters:
//   - alias: the unrecognized alias
//
// Returns:
//   - error: "unknown document <alias> (available: manifesto, about, invariants)"
func UnknownDocument(alias string) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidateUnknownDocument), alias,
	)
}

// ParseFile wraps a failure to parse a file.
//
// Parameters:
//   - path: file path that could not be parsed
//   - cause: the underlying parse error
//
// Returns:
//   - error: "failed to parse %s: <cause>"
func ParseFile(path string, cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrValidationParseFile), path, cause,
	)
}
