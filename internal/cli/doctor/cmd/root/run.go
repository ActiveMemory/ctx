//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/doctor/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/doctor/core/output"
	"github.com/ActiveMemory/ctx/internal/config/stats"
)

// Run executes the doctor command logic, running all health
// checks and producing either JSON or human-readable output.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - jsonOutput: If true, output as JSON
//
// Returns:
//   - error: Non-nil if output formatting fails
func Run(cmd *cobra.Command, jsonOutput bool) error {
	report := &check.Report{}

	check.ContextInitialized(report)
	check.RequiredFiles(report)
	check.CtxrcValidation(report)
	check.Drift(report)
	check.PluginEnablement(report)
	check.CompanionConfig(report)
	check.EventLogging(report)
	check.Webhook(report)
	check.Reminders(report)
	check.TaskCompletion(report)
	check.ContextTokenSize(report)
	check.SystemResources(report)
	check.RecentEventActivity(report)

	// Count warnings and errors.
	for _, r := range report.Results {
		switch r.Status {
		case stats.StatusWarning:
			report.Warnings++
		case stats.StatusError:
			report.Errors++
		}
	}

	if jsonOutput {
		return output.JSON(cmd, report)
	}
	return output.Human(cmd, report)
}
