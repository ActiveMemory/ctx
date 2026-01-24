//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2025-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	driftJSON bool
	driftFix  bool
)

// DriftCmd returns the drift command.
func DriftCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drift",
		Short: "Detect stale or invalid context",
		Long: `Run drift detection to find stale paths, broken references, and constitution violations.

Checks performed:
  - Path references in ARCHITECTURE.md and CONVENTIONS.md exist
  - Staleness indicators (many completed tasks)
  - Constitution rule violations (potential secrets)
  - Required files are present

Use --json for machine-readable output.`,
		RunE: runDrift,
	}

	cmd.Flags().BoolVar(&driftJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&driftFix, "fix", false, "Auto-fix simple issues (not yet implemented)")

	return cmd
}

// DriftJSONOutput represents the JSON output format.
type DriftJSONOutput struct {
	Timestamp  string        `json:"timestamp"`
	Status     string        `json:"status"`
	Warnings   []drift.Issue `json:"warnings"`
	Violations []drift.Issue `json:"violations"`
	Passed     []string      `json:"passed"`
}

func runDrift(cmd *cobra.Command, _ []string) error {
	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	report := drift.Detect(ctx)

	if driftJSON {
		return outputDriftJSON(cmd, report)
	}

	return outputDriftText(cmd, report)
}

func outputDriftJSON(cmd *cobra.Command, report *drift.Report) error {
	output := DriftJSONOutput{
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
		cmd.Printf("%s VIOLATIONS (%d)\n\n", red("❌"), len(report.Violations))
		for _, v := range report.Violations {
			if v.Line > 0 {
				cmd.Printf("  - %s:%d %s", v.File, v.Line, v.Message)
			} else {
				cmd.Printf("  - %s: %s", v.File, v.Message)
			}
			if v.Rule != "" {
				cmd.Printf(" (rule: %s)", v.Rule)
			}
			cmd.Println()
		}
		cmd.Println()
	}

	// Warnings
	if len(report.Warnings) > 0 {
		cmd.Printf("%s WARNINGS (%d)\n\n", yellow("⚠️ "), len(report.Warnings))

		// Group by type
		var pathRefs []drift.Issue
		var staleness []drift.Issue
		var other []drift.Issue

		for _, w := range report.Warnings {
			switch w.Type {
			case "dead_path":
				pathRefs = append(pathRefs, w)
			case "staleness":
				staleness = append(staleness, w)
			default:
				other = append(other, w)
			}
		}

		if len(pathRefs) > 0 {
			cmd.Println("  Path References:")
			for _, w := range pathRefs {
				cmd.Printf("  - %s:%d references '%s' (not found)\n", w.File, w.Line, w.Path)
			}
			cmd.Println()
		}

		if len(staleness) > 0 {
			cmd.Println("  Staleness:")
			for _, w := range staleness {
				cmd.Printf("  - %s %s\n", w.File, w.Message)
			}
			cmd.Println()
		}

		if len(other) > 0 {
			cmd.Println("  Other:")
			for _, w := range other {
				cmd.Printf("  - %s: %s\n", w.File, w.Message)
			}
			cmd.Println()
		}
	}

	// Passed
	if len(report.Passed) > 0 {
		cmd.Printf("%s PASSED (%d)\n", green("✅"), len(report.Passed))
		for _, p := range report.Passed {
			cmd.Printf("  - %s\n", formatCheckName(p))
		}
		cmd.Println()
	}

	// Summary
	status := report.Status()
	switch status {
	case "violation":
		cmd.Printf("\nStatus: %s — Constitution violations detected\n", red("VIOLATION"))
		return fmt.Errorf("drift detection found violations")
	case "warning":
		cmd.Printf("\nStatus: %s — Issues detected that should be addressed\n", yellow("WARNING"))
	default:
		cmd.Printf("\nStatus: %s — No drift detected\n", green("OK"))
	}

	return nil
}

func formatCheckName(name string) string {
	switch name {
	case "path_references":
		return "Path references are valid"
	case "staleness_check":
		return "No staleness indicators"
	case "constitution_check":
		return "Constitution rules respected"
	case "required_files":
		return "All required files present"
	default:
		return name
	}
}
