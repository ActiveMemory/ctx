//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import cfgDrift "github.com/ActiveMemory/ctx/internal/config/drift"

// Issue represents a detected drift issue.
//
// Issues are categorized by type and may reference specific files, lines,
// or paths in the codebase.
//
// Fields:
//   - File: Context file where the issue was detected (e.g., "ARCHITECTURE.md")
//   - Line: Line number in the file, if applicable
//   - Type: Issue category (e.g., "dead_path", "staleness", "missing_file")
//   - Message: Human-readable description of the issue
//   - Path: Referenced path that caused the issue, if applicable
//   - Rule: Constitution rule that was violated, if applicable
type Issue struct {
	File    string             `json:"file"`
	Line    int                `json:"line,omitempty"`
	Type    cfgDrift.IssueType `json:"type"`
	Message string             `json:"message"`
	Path    string             `json:"path,omitempty"`
	Rule    string             `json:"rule,omitempty"`
}

// Report represents the complete drift detection report.
//
// Contains categorized issues and a list of checks that passed.
//
// Fields:
//   - Warnings: Non-critical issues that should be addressed
//   - Violations: Critical issues that indicate constitution violations
//   - Passed: Names of checks that are completed without issues
type Report struct {
	Warnings   []Issue              `json:"warnings"`
	Violations []Issue              `json:"violations"`
	Passed     []cfgDrift.CheckName `json:"passed"`
}
