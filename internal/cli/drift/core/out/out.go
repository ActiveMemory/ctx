//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package out

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/drift/core/sanitize"
	cfgDrift "github.com/ActiveMemory/ctx/internal/config/drift"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/drift"
	errDrift "github.com/ActiveMemory/ctx/internal/err/drift"
	writeDrift "github.com/ActiveMemory/ctx/internal/write/drift"
)

// JSONOutput represents the JSON structure for
// machine-readable drift output.
//
// Fields:
//   - Timestamp: RFC3339-formatted UTC time
//   - Status: Overall drift status
//   - Warnings: Issues that should be addressed
//   - Violations: Constitution violations
//   - Passed: Names of checks that passed
type JSONOutput struct {
	Timestamp  string               `json:"timestamp"`
	Status     cfgDrift.StatusType  `json:"status"`
	Warnings   []drift.Issue        `json:"warnings"`
	Violations []drift.Issue        `json:"violations"`
	Passed     []cfgDrift.CheckName `json:"passed"`
}

// DriftText writes the drift report as formatted text with
// icons.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - report: Drift detection report to display
//
// Returns:
//   - error: Non-nil if violations were detected
func DriftText(
	cmd *cobra.Command, report *drift.Report,
) error {
	writeDrift.ReportHeader(cmd)

	// Violations
	if len(report.Violations) > 0 {
		writeDrift.ViolationsHeading(
			cmd, len(report.Violations),
		)
		for _, v := range report.Violations {
			line := fmt.Sprintf(
				desc.Text(text.DescKeyDriftViolationLine),
				v.File, v.Message,
			)
			if v.Line > 0 {
				line = fmt.Sprintf(
					desc.Text(
						text.DescKeyDriftViolationLineLoc,
					),
					v.File, v.Line, v.Message,
				)
			}
			if v.Rule != "" {
				line += fmt.Sprintf(
					desc.Text(text.DescKeyDriftViolationRule),
					v.Rule,
				)
			}
			writeDrift.ViolationLine(cmd, line)
		}
		writeDrift.BlankLine(cmd)
	}

	// Warnings
	if len(report.Warnings) > 0 {
		writeDrift.WarningsHeading(
			cmd, len(report.Warnings),
		)

		var pathRefs []drift.Issue
		var staleness []drift.Issue
		var other []drift.Issue

		for _, w := range report.Warnings {
			switch w.Type {
			case cfgDrift.IssueDeadPath:
				pathRefs = append(pathRefs, w)
			case cfgDrift.IssueStaleness, cfgDrift.IssueStaleAge:
				staleness = append(staleness, w)
			default:
				other = append(other, w)
			}
		}

		if len(pathRefs) > 0 {
			items := make([]string, len(pathRefs))
			for i, w := range pathRefs {
				items[i] = fmt.Sprintf(
					desc.Text(text.DescKeyDriftPathRefLine),
					w.File, w.Line, w.Path,
				)
			}
			writeDrift.PathRefsBlock(cmd, items)
		}

		if len(staleness) > 0 {
			items := make([]string, len(staleness))
			for i, w := range staleness {
				items[i] = fmt.Sprintf(
					desc.Text(
						text.DescKeyDriftStalenessLine,
					),
					w.File, w.Message,
				)
			}
			writeDrift.StalenessBlock(cmd, items)
		}

		if len(other) > 0 {
			items := make([]string, len(other))
			for i, w := range other {
				items[i] = fmt.Sprintf(
					desc.Text(text.DescKeyDriftOtherLine),
					w.File, w.Message,
				)
			}
			writeDrift.OtherBlock(cmd, items)
		}
	}

	// Passed
	if len(report.Passed) > 0 {
		writeDrift.PassedHeading(cmd, len(report.Passed))
		for _, p := range report.Passed {
			writeDrift.PassedLine(
				cmd, sanitize.FormatCheckName(p),
			)
		}
		writeDrift.BlankLine(cmd)
	}

	// Summary
	status := report.Status()
	switch status {
	case cfgDrift.StatusViolation:
		writeDrift.StatusViolation(cmd)
		return errDrift.Violations()
	case cfgDrift.StatusWarning:
		writeDrift.StatusWarning(cmd)
	default:
		writeDrift.StatusOK(cmd)
	}

	return nil
}

// DriftJSON writes the drift report as pretty-printed JSON.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - report: Drift detection report to serialize
//
// Returns:
//   - error: Non-nil if JSON encoding fails
func DriftJSON(
	cmd *cobra.Command, report *drift.Report,
) error {
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
