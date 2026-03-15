//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import "fmt"

// CtxNotInPath returns an error indicating that ctx was not found in PATH.
//
// Returns:
//   - error: "ctx not found in PATH"
func CtxNotInPath() error {
	return fmt.Errorf("ctx not found in PATH")
}

// WorkingDirectory wraps a failure to determine the working directory.
//
// Parameters:
//   - cause: the underlying error from os.Getwd.
//
// Returns:
//   - error: "failed to get working directory: <cause>"
func WorkingDirectory(cause error) error {
	return fmt.Errorf("failed to get working directory: %w", cause)
}

// DriftViolations returns an error when drift detection found violations.
//
// Returns:
//   - error: "drift detection found violations"
func DriftViolations() error {
	return fmt.Errorf("drift detection found violations")
}

// FlagRequired returns an error for a missing required flag.
//
// Parameters:
//   - name: the flag name.
//
// Returns:
//   - error: "required flag \"<name>\" not set"
func FlagRequired(name string) error {
	return fmt.Errorf("required flag %q not set", name)
}

// ArgRequired returns an error for a missing required argument.
//
// Parameters:
//   - name: the argument name.
//
// Returns:
//   - error: "<name> argument is required"
func ArgRequired(name string) error {
	return fmt.Errorf("%s argument is required", name)
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
	return fmt.Errorf("failed to parse %s: %w", path, cause)
}
