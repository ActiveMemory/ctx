//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/doctor"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/stats"
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
	cmd.Println(desc.Text(text.DescKeyDoctorOutputHeader))
	cmd.Println(desc.Text(text.DescKeyDoctorOutputSeparator))
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
			icon := stats.StatusIcon(r.Status)
			cmd.Println(fmt.Sprintf(
				desc.Text(text.DescKeyDoctorOutputResultLine), icon, r.Message),
			)
		}
		cmd.Println()
	}

	cmd.Println(
		fmt.Sprintf(
			desc.Text(text.DescKeyDoctorOutputSummary),
			report.Warnings, report.Errors,
		),
	)
	return nil
}
