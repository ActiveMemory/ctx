//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/doctor"
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
	cmd.Println(assets.TextDesc(assets.TextDescKeyDoctorOutputHeader))
	cmd.Println(assets.TextDesc(assets.TextDescKeyDoctorOutputSeparator))
	cmd.Println()

	// Group by category.
	categories := []string{
		doctor.CategoryStructure,
		doctor.CategoryQuality,
		doctor.CategoryPlugin,
		doctor.CategoryHooks,
		doctor.CategoryState,
		doctor.CategorySize,
		doctor.CategoryResources,
		doctor.CategoryEvents,
	}
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
			cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyDoctorOutputResultLine), icon, r.Message))
		}
		cmd.Println()
	}

	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyDoctorOutputSummary), report.Warnings, report.Errors))
	return nil
}

// statusIcon returns a unicode icon for the given status string.
//
// Parameters:
//   - status: One of StatusOK, StatusWarning, StatusError, or StatusInfo
//
// Returns:
//   - string: A single unicode character representing the status
func statusIcon(status string) string {
	switch status {
	case StatusOK:
		return "✓"
	case StatusWarning:
		return "⚠"
	case StatusError:
		return "✗"
	case StatusInfo:
		return "○"
	default:
		return "?"
	}
}
