//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/doctor/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/doctor/core/output"
	"github.com/ActiveMemory/ctx/internal/config/doctor"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
)

// Run executes the doctor command logic, running all health
// checks and producing either JSON or human-readable output.
//
// Context-dependent checks that fail with
// [errCtx.ErrDirNotDeclared] emit exactly one "did not run
// (cascade)" line; later dependent checks are silently skipped
// so the report shows one loud entry instead of N copies of the
// same message. Non-dependent checks (companion config, plugin,
// system resources, etc.) continue to run regardless.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - jsonOutput: If true, output as JSON
//
// Returns:
//   - error: Non-nil if output formatting fails
func Run(cmd *cobra.Command, jsonOutput bool) error {
	report := &check.Report{}

	entries := []check.Entry{
		{
			Name:     doctor.CheckContextInit,
			Category: doctor.CategoryStructure,
			Fn:       check.ContextInitialized,
		},
		{
			Name:     doctor.CheckRequiredFiles,
			Category: doctor.CategoryStructure,
			Fn:       check.RequiredFiles,
		},
		{
			Name:     doctor.CheckCtxrcValidation,
			Category: doctor.CategoryStructure,
			Fn:       check.CtxrcValidation,
		},
		{
			Name:     doctor.CheckDrift,
			Category: doctor.CategoryQuality,
			Fn:       check.Drift,
		},
		{
			Name:     doctor.CheckPluginInstalled,
			Category: doctor.CategoryPlugin,
			Fn:       check.PluginEnablement,
		},
		{
			Name:     doctor.CheckCompanionConfig,
			Category: doctor.CategoryPlugin,
			Fn:       check.CompanionConfig,
		},
		{
			Name:     doctor.CheckEventLogging,
			Category: doctor.CategoryHooks,
			Fn:       check.EventLogging,
		},
		{
			Name:     doctor.CheckWebhook,
			Category: doctor.CategoryHooks,
			Fn:       check.Webhook,
		},
		{
			Name:     doctor.CheckReminders,
			Category: doctor.CategoryState,
			Fn:       check.Reminders,
		},
		{
			Name:     doctor.CheckTaskCompletion,
			Category: doctor.CategoryState,
			Fn:       check.TaskCompletion,
		},
		{
			Name:     doctor.CheckContextSize,
			Category: doctor.CategorySize,
			Fn:       check.ContextTokenSize,
		},
		{
			Name:     doctor.CheckResourceMemory,
			Category: doctor.CategoryResources,
			Fn:       check.SystemResources,
		},
		{
			Name:     doctor.CheckRecentEvents,
			Category: doctor.CategoryEvents,
			Fn:       check.RecentEventActivity,
		},
	}

	// Track whether a context-dependent check has already
	// failed due to errCtx.ErrDirNotDeclared. Subsequent
	// dependent failures with the same root cause are folded
	// into a single diagnostic.
	ctxCascadeAnnounced := false

	for _, entry := range entries {
		err := entry.Fn(report)
		if err == nil {
			continue
		}
		if errors.Is(err, errCtx.ErrDirNotDeclared) {
			if ctxCascadeAnnounced {
				// Already reported once; skip silently.
				continue
			}
			ctxCascadeAnnounced = true
			report.Results = append(report.Results, check.Result{
				Name:     entry.Name,
				Category: entry.Category,
				Status:   stats.StatusError,
				Message: fmt.Sprintf(desc.Text(
					text.DescKeyDoctorCheckDidNotRunCascade,
				), err),
			})
			continue
		}
		// Non-cascade error: attribute to the specific check.
		report.Results = append(report.Results, check.Result{
			Name:     entry.Name,
			Category: entry.Category,
			Status:   stats.StatusError,
			Message: fmt.Sprintf(
				desc.Text(text.DescKeyDoctorCheckDidNotRun), err,
			),
		})
	}

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
