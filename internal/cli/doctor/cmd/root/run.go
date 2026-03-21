//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/doctor/core"
	"github.com/ActiveMemory/ctx/internal/config/stats"
)

// Run executes the doctor command logic, running all health checks and
// producing either JSON or human-readable output.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - jsonOutput: If true, output as JSON
//
// Returns:
//   - error: Non-nil if output formatting fails
func Run(cmd *cobra.Command, jsonOutput bool) error {
	report := &core.Report{}

	core.CheckContextInitialized(report)
	core.CheckRequiredFiles(report)
	core.CheckCtxrcValidation(report)
	core.CheckDrift(report)
	core.CheckPluginEnablement(report)
	core.CheckEventLogging(report)
	core.CheckWebhook(report)
	core.CheckReminders(report)
	core.CheckTaskCompletion(report)
	core.CheckContextTokenSize(report)
	core.CheckSystemResources(report)
	core.CheckRecentEventActivity(report)

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
		return core.OutputJSON(cmd, report)
	}
	return core.OutputHuman(cmd, report)
}
