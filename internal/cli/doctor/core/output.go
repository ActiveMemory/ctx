//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"

	"github.com/spf13/cobra"

	writeDoctor "github.com/ActiveMemory/ctx/internal/write/doctor"
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
	writeDoctor.JSON(cmd, string(data))
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
	items := make([]writeDoctor.ResultItem, len(report.Results))
	for i, r := range report.Results {
		items[i] = writeDoctor.ResultItem{
			Category: r.Category,
			Status:   r.Status,
			Message:  r.Message,
		}
	}
	writeDoctor.Report(cmd, items, report.Warnings, report.Errors)
	return nil
}
