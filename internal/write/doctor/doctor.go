//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package doctor

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/doctor"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/stats"
)

// JSON prints pre-marshaled JSON data to the command's output stream.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - data: Pre-marshaled JSON string
func JSON(cmd *cobra.Command, data string) {
	if cmd == nil {
		return
	}
	cmd.Println(data)
}

// Report writes the doctor report in a human-readable format grouped
// by category. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - results: Individual check results to display
//   - warnings: Total warning count
//   - errors: Total error count
func Report(
	cmd *cobra.Command, results []ResultItem, warnings, errors int,
) {
	if cmd == nil {
		return
	}
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
	grouped := make(map[string][]ResultItem)
	for _, r := range results {
		grouped[r.Category] = append(grouped[r.Category], r)
	}

	for _, cat := range categories {
		items, ok := grouped[cat]
		if !ok {
			continue
		}
		cmd.Println(cat)
		for _, r := range items {
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
			warnings, errors,
		),
	)
}
