//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package drift provides functionality for detecting stale or invalid context.
package drift

import "github.com/ActiveMemory/ctx/internal/entity"

// Status returns the overall status of the report.
//
// Returns:
//   - StatusType: StatusViolation if any violations, StatusWarning if only
//     warnings, StatusOk otherwise
func (r *Report) Status() StatusType {
	if len(r.Violations) > 0 {
		return StatusViolation
	}
	if len(r.Warnings) > 0 {
		return StatusWarning
	}
	return StatusOk
}

// Detect runs all drift detection checks on the given context.
//
// Performs multiple validation checks including path references, staleness
// indicators, constitution compliance, and required file presence.
//
// Parameters:
//   - ctx: Loaded context containing files to check
//
// Returns:
//   - *Report: Drift report with warnings, violations, and passed checks
func Detect(ctx *entity.Context) *Report {
	report := &Report{
		Warnings:   []Issue{},
		Violations: []Issue{},
		Passed:     []CheckName{},
	}

	// Check path references in context files
	checkPathReferences(ctx, report)

	// Check for staleness indicators
	checkStaleness(ctx, report)

	// Check constitution rules (basic heuristics)
	checkConstitution(ctx, report)

	// Check for empty required files
	checkRequiredFiles(ctx, report)

	// Check for files not modified recently
	checkFileAge(ctx, report)

	// Check for excessive entry counts in knowledge files
	checkEntryCount(ctx, report)

	// Check for undocumented internal packages
	checkMissingPackages(ctx, report)

	// Check context file comment headers against templates
	checkTemplateHeaders(ctx, report)

	return report
}
