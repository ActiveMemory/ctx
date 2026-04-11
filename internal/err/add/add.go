//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// NoContent returns an error when no content source is available.
//
// Returns:
//   - error: "no content provided"
func NoContent() error {
	return errors.New(desc.Text(text.DescKeyErrAddNoContent))
}

// NoContentProvided returns an error with usage help when content is missing.
//
// Parameters:
//   - fType: Entry type (e.g., "decision", "task") for contextual examples
//   - examples: Type-specific example text
//
// Returns:
//   - error: Formatted error showing input methods and type-specific examples
func NoContentProvided(fType, examples string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrAddNoContentProvided),
		fType, fType, fType, examples,
	)
}

// IndexUpdate wraps a failure to update the index in a context file.
//
// Parameters:
//   - path: File path where the index update failed
//   - cause: Underlying error from the write operation
//
// Returns:
//   - error: "failed to update index in <path>: <cause>"
func IndexUpdate(path string, cause error) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrAddIndexUpdate), path, cause,
	)
}

// UnknownType returns an error for an unrecognized entry type.
//
// Parameters:
//   - fType: The unrecognized type string
//
// Returns:
//   - error: Formatted error listing valid types
func UnknownType(fType string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrAddUnknownType), fType,
	)
}

// FileNotFound returns an error when a context file does not exist.
//
// Parameters:
//   - path: File path that was not found
//
// Returns:
//   - error: Formatted error suggesting "ctx init"
func FileNotFound(path string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrAddFileNotFound), path,
	)
}

// SectionRequired returns a validation error when a task is added without
// the --section flag.
//
// Returns:
//   - error: Formatted error explaining that --section is mandatory for tasks
func SectionRequired() error {
	return errors.New(desc.Text(text.DescKeyErrAddSectionRequired))
}

// MissingFields returns a validation error for missing required fields.
//
// Parameters:
//   - entryType: The entry type (e.g., "decision", "learning")
//   - missing: List of missing field names
//
// Returns:
//   - error: Formatted error listing the missing fields
func MissingFields(entryType string, missing []string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrAddMissingFields),
		entryType, strings.Join(missing, token.CommaSpace),
	)
}
