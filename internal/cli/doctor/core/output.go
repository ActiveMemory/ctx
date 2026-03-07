//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// OutputJSON writes the report as indented JSON to the command's output stream.
//
// Parameters:
//   - cmd: Cobra command providing the output writer
//   - report: Doctor report to serialize
//
// Returns:
//   - error: Non-nil if JSON marshaling fails
func OutputJSON(cmd *cobra.Command, report *Report) error {
	data, marshalErr := json.MarshalIndent(report, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	cmd.Println(string(data))
	return nil
}

// OutputHuman writes the report in a human-readable format grouped by category.
//
// Parameters:
//   - cmd: Cobra command providing the output writer
//   - report: Doctor report to display
//
// Returns:
//   - error: Always nil (satisfies interface)
func OutputHuman(cmd *cobra.Command, report *Report) error {
	cmd.Println("ctx doctor")
	cmd.Println("==========")
	cmd.Println()

	// Group by category.
	categories := []string{"Structure", "Quality", "Plugin", "Hooks", "State", "Size", "Resources", "Events"}
	grouped := make(map[string][]Result)
	for _, r := range report.Results {
		grouped[r.Category] = append(grouped[r.Category], r)
	}

	for _, cat := range categories {
		results, ok := grouped[cat]
		if !ok {
			continue
		}
		cmd.Println(cat)
		for _, r := range results {
			icon := statusIcon(r.Status)
			cmd.Println(fmt.Sprintf("  %s %s", icon, r.Message))
		}
		cmd.Println()
	}

	cmd.Println(fmt.Sprintf("Summary: %d warnings, %d errors", report.Warnings, report.Errors))
	return nil
}

// statusIcon returns a unicode icon for the given status string.
func statusIcon(status string) string {
	switch status {
	case StatusOK:
		return "\u2713" // check mark
	case StatusWarning:
		return "\u26a0" // warning sign
	case StatusError:
		return "\u2717" // ballot x
	case StatusInfo:
		return "\u25cb" // white circle
	default:
		return "?"
	}
}
