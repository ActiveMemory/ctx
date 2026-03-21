//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/drift"
	errdrift "github.com/ActiveMemory/ctx/internal/err/drift"
)

// OutputDriftText writes the drift report as formatted text with icons.
//
// Output is grouped into violations, warnings (by type), and passed checks.
// Includes a summary status line at the end.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - report: Drift detection report to display
//
// Returns:
//   - error: Non-nil if violations were detected
func OutputDriftText(cmd *cobra.Command, report *drift.Report) error {
	cmd.Println(desc.TextDesc(text.DescKeyDriftReportHeading))
	cmd.Println(desc.TextDesc(text.DescKeyDriftReportSeparator))
	cmd.Println()

	// Violations
	if len(report.Violations) > 0 {
		cmd.Println(fmt.Sprintf(
			desc.TextDesc(
				text.DescKeyDriftViolationsHeading), len(report.Violations)),
		)
		cmd.Println()
		for _, v := range report.Violations {
			line := fmt.Sprintf(
				desc.TextDesc(text.DescKeyDriftViolationLine), v.File, v.Message,
			)
			if v.Line > 0 {
				line = fmt.Sprintf(
					desc.TextDesc(text.DescKeyDriftViolationLineLoc),
					v.File, v.Line, v.Message,
				)
			}
			if v.Rule != "" {
				line += fmt.Sprintf(
					desc.TextDesc(text.DescKeyDriftViolationRule), v.Rule,
				)
			}
			cmd.Println(line)
		}
		cmd.Println()
	}

	// Warnings
	if len(report.Warnings) > 0 {
		cmd.Println(fmt.Sprintf(
			desc.TextDesc(text.DescKeyDriftWarningsHeading), len(report.Warnings)))
		cmd.Println()

		// Group by type
		var pathRefs []drift.Issue
		var staleness []drift.Issue
		var other []drift.Issue

		for _, w := range report.Warnings {
			switch w.Type {
			case drift.IssueDeadPath:
				pathRefs = append(pathRefs, w)
			case drift.IssueStaleness, drift.IssueStaleAge:
				staleness = append(staleness, w)
			default:
				other = append(other, w)
			}
		}

		if len(pathRefs) > 0 {
			cmd.Println(desc.TextDesc(text.DescKeyDriftPathRefsLabel))
			for _, w := range pathRefs {
				cmd.Println(fmt.Sprintf(
					desc.TextDesc(text.DescKeyDriftPathRefLine), w.File, w.Line, w.Path))
			}
			cmd.Println()
		}

		if len(staleness) > 0 {
			cmd.Println(desc.TextDesc(text.DescKeyDriftStalenessLabel))
			for _, w := range staleness {
				cmd.Println(fmt.Sprintf(
					desc.TextDesc(text.DescKeyDriftStalenessLine), w.File, w.Message))
			}
			cmd.Println()
		}

		if len(other) > 0 {
			cmd.Println(desc.TextDesc(text.DescKeyDriftOtherLabel))
			for _, w := range other {
				cmd.Println(fmt.Sprintf(
					desc.TextDesc(text.DescKeyDriftOtherLine), w.File, w.Message))
			}
			cmd.Println()
		}
	}

	// Passed
	if len(report.Passed) > 0 {
		cmd.Println(fmt.Sprintf(
			desc.TextDesc(text.DescKeyDriftPassedHeading), len(report.Passed)))
		for _, p := range report.Passed {
			cmd.Println(fmt.Sprintf(
				desc.TextDesc(text.DescKeyDriftPassedLine), FormatCheckName(p)))
		}
		cmd.Println()
	}

	// Summary
	status := report.Status()
	switch status {
	case drift.StatusViolation:
		cmd.Println()
		cmd.Println(desc.TextDesc(text.DescKeyDriftStatusViolation))
		return errdrift.Violations()
	case drift.StatusWarning:
		cmd.Println()
		cmd.Println(desc.TextDesc(text.DescKeyDriftStatusWarning))
	default:
		cmd.Println()
		cmd.Println(desc.TextDesc(text.DescKeyDriftStatusOK))
	}

	return nil
}

// OutputDriftJSON writes the drift report as pretty-printed JSON.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - report: Drift detection report to serialize
//
// Returns:
//   - error: Non-nil if JSON encoding fails
func OutputDriftJSON(cmd *cobra.Command, report *drift.Report) error {
	output := JSONOutput{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Status:     report.Status(),
		Warnings:   report.Warnings,
		Violations: report.Violations,
		Passed:     report.Passed,
	}

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}
