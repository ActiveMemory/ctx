//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"fmt"
	"strings"
)

// ErrNoContent returns a simple error when no content source is available.
//
// Returns:
//   - error: "no content provided"
func ErrNoContent() error {
	return fmt.Errorf("no content provided")
}

// ErrMissingDecision returns an error with usage help for incomplete decisions.
//
// Parameters:
//   - missing: List of missing required flag names (e.g., "--context")
//
// Returns:
//   - error: Formatted error with ADR format requirements and example
func ErrMissingDecision(missing []string) error {
	return fmt.Errorf(`decisions require complete ADR format

Missing required flags: %s

Usage:
  ctx add decision "Decision title" \
    --context "What prompted this decision" \
    --rationale "Why this choice over alternatives" \
    --consequences "What changes as a result"

Example:
  ctx add decision "Use PostgreSQL for primary database" \
    --context "Need a reliable database for production workloads" \
    --rationale "PostgreSQL offers ACID compliance, JSON support, and team familiarity" \
    --consequences "Team needs PostgreSQL training; must set up replication"`,
		strings.Join(missing, ", "))
}

// ErrMissingLearning returns an error with usage help for incomplete learnings.
//
// Parameters:
//   - missing: List of missing required flag names (e.g., "--lesson")
//
// Returns:
//   - error: Formatted error with learning format requirements and example
func ErrMissingLearning(missing []string) error {
	return fmt.Errorf(`learnings require complete format

Missing required flags: %s

Usage:
  ctx add learning "Learning title" \
    --context "What prompted this learning" \
    --lesson "The key insight" \
    --application "How to apply this going forward"

Example:
  ctx add learning "Go embed requires files in same package" \
    --context "Tried to embed files from parent directory, got compile error" \
    --lesson "go:embed only works with files in same or child directories" \
    --application "Keep embedded files in internal/templates/, not project root"`,
		strings.Join(missing, ", "))
}

// ErrNoContentProvided returns an error with usage help when content is missing.
//
// Parameters:
//   - fType: Entry type (e.g., "decision", "task") for contextual examples
//   - examples: Type-specific example text
//
// Returns:
//   - error: Formatted error showing input methods and type-specific examples
func ErrNoContentProvided(fType, examples string) error {
	return fmt.Errorf(`no content provided

Usage:
  ctx add %s "your content here"
  ctx add %s --file /path/to/content.md
  echo "content" | ctx add %s

Examples:
%s`, fType, fType, fType, examples)
}

// ErrFileRead wraps a file read failure with the file path.
//
// Parameters:
//   - path: File path that failed to read
//   - cause: Underlying error from the read operation
//
// Returns:
//   - error: "failed to read <path>: <cause>"
func ErrFileRead(path string, cause error) error {
	return fmt.Errorf("failed to read %s: %w", path, cause)
}

// ErrFileWriteAdd wraps a file write failure with the file path.
//
// Parameters:
//   - path: File path that failed to write
//   - cause: Underlying error from the write operation
//
// Returns:
//   - error: "failed to write <path>: <cause>"
func ErrFileWriteAdd(path string, cause error) error {
	return fmt.Errorf("failed to write %s: %w", path, cause)
}

// ErrStdinRead wraps a failure to read from standard input.
//
// Parameters:
//   - cause: Underlying error from the stdin read
//
// Returns:
//   - error: "failed to read from stdin: <cause>"
func ErrStdinRead(cause error) error {
	return fmt.Errorf("failed to read from stdin: %w", cause)
}

// ErrIndexUpdate wraps a failure to update the index in a context file.
//
// Parameters:
//   - path: File path where the index update failed
//   - cause: Underlying error from the write operation
//
// Returns:
//   - error: "failed to update index in <path>: <cause>"
func ErrIndexUpdate(path string, cause error) error {
	return fmt.Errorf("failed to update index in %s: %w", path, cause)
}

// ErrUnknownType returns an error for an unrecognized entry type.
//
// Parameters:
//   - fType: The unrecognized type string
//
// Returns:
//   - error: Formatted error listing valid types
func ErrUnknownType(fType string) error {
	return fmt.Errorf(
		"unknown type %q. Valid types: decision, task, learning, convention",
		fType,
	)
}

// ErrFileNotFound returns an error when a context file does not exist.
//
// Parameters:
//   - path: File path that was not found
//
// Returns:
//   - error: Formatted error suggesting "ctx init"
func ErrFileNotFound(path string) error {
	return fmt.Errorf(
		"context file %s not found. Run 'ctx init' first", path,
	)
}

// ErrMissingFields returns a validation error for missing required fields.
//
// Parameters:
//   - entryType: The entry type (e.g., "decision", "learning")
//   - missing: List of missing field names
//
// Returns:
//   - error: Formatted error listing the missing fields
func ErrMissingFields(entryType string, missing []string) error {
	return fmt.Errorf(
		"%s missing required fields: %s", entryType, strings.Join(missing, ", "),
	)
}
