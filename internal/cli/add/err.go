//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"fmt"
	"strings"
)

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
