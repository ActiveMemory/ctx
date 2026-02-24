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

// errNoContent returns a simple error when no content source is available.
// Callers that know the entry type should use errNoContentProvided(fType)
// instead for a richer message.
//
// Returns:
//   - error: "no content provided"
func errNoContent() error {
	return fmt.Errorf("no content provided")
}

// errMissingDecision returns an error with usage help for incomplete decisions.
//
// Parameters:
//   - missing: List of missing required flag names (e.g., "--context")
//
// Returns:
//   - error: Formatted error with ADR format requirements and example
func errMissingDecision(missing []string) error {
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

// errMissingLearning returns an error with usage help for incomplete learnings.
//
// Parameters:
//   - missing: List of missing required flag names (e.g., "--lesson")
//
// Returns:
//   - error: Formatted error with learning format requirements and example
func errMissingLearning(missing []string) error {
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

// errNoContentProvided returns an error with usage help when content is missing.
//
// Parameters:
//   - fType: Entry type (e.g., "decision", "task") for contextual examples
//
// Returns:
//   - error: Formatted error showing input methods and type-specific examples
func errNoContentProvided(fType string) error {
	examples := examplesForType(fType)
	return fmt.Errorf(`no content provided

Usage:
  ctx add %s "your content here"
  ctx add %s --file /path/to/content.md
  echo "content" | ctx add %s

Examples:
%s`, fType, fType, fType, examples)
}

// errFileRead wraps a file read failure with the file path.
//
// Parameters:
//   - path: File path that failed to read
//   - err: Underlying error from the read operation
//
// Returns:
//   - error: Wrapped error with format "failed to read <path>: <cause>"
func errFileRead(path string, err error) error {
	return fmt.Errorf("failed to read %s: %w", path, err)
}

// errFileWrite wraps a file write failure with the file path.
//
// Parameters:
//   - path: File path that failed to write
//   - err: Underlying error from the write operation
//
// Returns:
//   - error: Wrapped error with format "failed to write <path>: <cause>"
func errFileWrite(path string, err error) error {
	return fmt.Errorf("failed to write %s: %w", path, err)
}

// errStdinRead wraps a failure to read from standard input.
//
// Parameters:
//   - err: Underlying error from the stdin read
//
// Returns:
//   - error: Wrapped error with format "failed to read from stdin: <cause>"
func errStdinRead(err error) error {
	return fmt.Errorf("failed to read from stdin: %w", err)
}

// errIndexUpdate wraps a failure to update the index in a context file.
//
// Parameters:
//   - path: File path where the index update failed
//   - err: Underlying error from the write operation
//
// Returns:
//   - error: Wrapped error with format "failed to update index in <path>: <cause>"
func errIndexUpdate(path string, err error) error {
	return fmt.Errorf("failed to update index in %s: %w", path, err)
}

// errUnknownType returns an error for an unrecognized entry type.
//
// Parameters:
//   - fType: The unrecognized type string
//
// Returns:
//   - error: Formatted error listing valid types
func errUnknownType(fType string) error {
	return fmt.Errorf(
		"unknown type %q. Valid types: decision, task, learning, convention",
		fType,
	)
}

// errFileNotFound returns an error when a context file does not exist.
//
// Parameters:
//   - path: File path that was not found
//
// Returns:
//   - error: Formatted error suggesting "ctx init"
func errFileNotFound(path string) error {
	return fmt.Errorf(
		"context file %s not found. Run 'ctx init' first", path,
	)
}

// errMissingFields returns a validation error for missing required fields.
//
// Parameters:
//   - entryType: The entry type (e.g., "decision", "learning")
//   - missing: List of missing field names
//
// Returns:
//   - error: Formatted error listing the missing fields
func errMissingFields(entryType string, missing []string) error {
	return fmt.Errorf(
		"%s missing required fields: %s", entryType, strings.Join(missing, ", "),
	)
}
