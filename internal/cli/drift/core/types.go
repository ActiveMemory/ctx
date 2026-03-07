//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import "github.com/ActiveMemory/ctx/internal/drift"

// FixResult tracks fixes applied during drift fix.
//
// Fields:
//   - Fixed: Number of issues successfully fixed
//   - Skipped: Number of issues skipped (not auto-fixable)
//   - Errors: Error messages from failed fix attempts
type FixResult struct {
	Fixed   int
	Skipped int
	Errors  []string
}

// JSONOutput represents the JSON structure for machine-readable drift output.
//
// Fields:
//   - Timestamp: RFC3339-formatted UTC time when the report was generated
//   - Status: Overall drift status ("ok", "warning", or "violation")
//   - Warnings: Issues that should be addressed but don't block
//   - Violations: Constitution violations that must be fixed
//   - Passed: Names of checks that passed successfully
type JSONOutput struct {
	Timestamp  string            `json:"timestamp"`
	Status     drift.StatusType  `json:"status"`
	Warnings   []drift.Issue     `json:"warnings"`
	Violations []drift.Issue     `json:"violations"`
	Passed     []drift.CheckName `json:"passed"`
}
