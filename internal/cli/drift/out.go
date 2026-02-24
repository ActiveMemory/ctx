//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/drift"
)

// outputDriftText writes the drift report as formatted text with colors.
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
func outputDriftText(cmd *cobra.Command, report *drift.Report) error {
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	cmd.Println(cyan("Drift Detection Report"))
	cmd.Println(cyan("======================"))
	cmd.Println()

	// Violations
	if len(report.Violations) > 0 {
		cmd.Println(fmt.Sprintf(
			"%s VIOLATIONS (%d)", red("❌"), len(report.Violations)),
		)
		cmd.Println()
		for _, v := range report.Violations {
			line := fmt.Sprintf("  - %s: %s", v.File, v.Message)
			if v.Line > 0 {
				line = fmt.Sprintf("  - %s:%d %s", v.File, v.Line, v.Message)
			}
			if v.Rule != "" {
				line += fmt.Sprintf(" (rule: %s)", v.Rule)
			}
			cmd.Println(line)
		}
		cmd.Println()
	}

	// Warnings
	if len(report.Warnings) > 0 {
		cmd.Println(
			fmt.Sprintf("%s WARNINGS (%d)", yellow("⚠️ "), len(report.Warnings)),
		)
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
			cmd.Println("  Path References:")
			for _, w := range pathRefs {
				cmd.Println(fmt.Sprintf(
					"  - %s:%d references '%s' (not found)", w.File, w.Line, w.Path,
				))
			}
			cmd.Println()
		}

		if len(staleness) > 0 {
			cmd.Println("  Staleness:")
			for _, w := range staleness {
				cmd.Println(fmt.Sprintf("  - %s %s", w.File, w.Message))
			}
			cmd.Println()
		}

		if len(other) > 0 {
			cmd.Println("  Other:")
			for _, w := range other {
				cmd.Println(fmt.Sprintf("  - %s: %s", w.File, w.Message))
			}
			cmd.Println()
		}
	}

	// Passed
	if len(report.Passed) > 0 {
		cmd.Println(fmt.Sprintf("%s PASSED (%d)", green("✅"), len(report.Passed)))
		for _, p := range report.Passed {
			cmd.Println(fmt.Sprintf("  - %s", formatCheckName(p)))
		}
		cmd.Println()
	}

	// Summary
	status := report.Status()
	switch status {
	case drift.StatusViolation:
		cmd.Println()
		cmd.Println(fmt.Sprintf(
			"Status: %s — Constitution violations detected", red("VIOLATION")),
		)
		return errViolationsFound()
	case drift.StatusWarning:
		cmd.Println()
		cmd.Println(fmt.Sprintf(
			"Status: %s — Issues detected that should be addressed",
			yellow("WARNING")),
		)
	default:
		cmd.Println()
		cmd.Println(fmt.Sprintf("Status: %s — No drift detected", green("OK")))
	}

	return nil
}

// outputDriftJSON writes the drift report as pretty-printed JSON.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - report: Drift detection report to serialize
//
// Returns:
//   - error: Non-nil if JSON encoding fails
func outputDriftJSON(cmd *cobra.Command, report *drift.Report) error {
	output := JsonOutput{
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
