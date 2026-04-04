//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	initCore "github.com/ActiveMemory/ctx/internal/cli/initialize/core/plugin"
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/crypto"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/doctor"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/reminder"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	cfgSysinfo "github.com/ActiveMemory/ctx/internal/config/sysinfo"
	cfgToken "github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context/load"
	"github.com/ActiveMemory/ctx/internal/context/token"
	"github.com/ActiveMemory/ctx/internal/context/validate"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/log/event"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

// ContextInitialized verifies that a .context/ directory
// exists.
//
// Parameters:
//   - report: Report to append the result to
func ContextInitialized(report *Report) {
	if validate.Exists("") {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckContextInit,
			Category: doctor.CategoryStructure,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorContextInitializedOk),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckContextInit,
			Category: doctor.CategoryStructure,
			Status:   stats.StatusError,
			Message:  desc.Text(text.DescKeyDoctorContextInitializedError),
		})
	}
}

// RequiredFiles verifies that all required context files are
// present.
//
// Parameters:
//   - report: Report to append the result to
func RequiredFiles(report *Report) {
	dir := rc.ContextDir()
	var missing []string
	for _, f := range ctx.FilesRequired {
		path := filepath.Join(dir, f)
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			missing = append(missing, f)
		}
	}

	total := len(ctx.FilesRequired)
	present := total - len(missing)

	if len(missing) == 0 {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckRequiredFiles,
			Category: doctor.CategoryStructure,
			Status:   stats.StatusOK,
			Message: fmt.Sprintf(
				desc.Text(text.DescKeyDoctorRequiredFilesOk),
				present, total,
			),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckRequiredFiles,
			Category: doctor.CategoryStructure,
			Status:   stats.StatusError,
			Message: fmt.Sprintf(
				desc.Text(text.DescKeyDoctorRequiredFilesError),
				present, total,
				strings.Join(missing, cfgToken.CommaSpace),
			),
		})
	}
}

// CtxrcValidation validates the .ctxrc file for unknown
// fields or parse errors.
//
// Parameters:
//   - report: Report to append the result to
func CtxrcValidation(report *Report) {
	data, readErr := io.SafeReadUserFile(file.CtxRC)
	if readErr != nil {
		// No .ctxrc is fine - defaults are used.
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckCtxrcValidation,
			Category: doctor.CategoryStructure,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorCtxrcValidationOkNoFile),
		})
		return
	}

	warnings, validateErr := rc.Validate(data)
	if validateErr != nil {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckCtxrcValidation,
			Category: doctor.CategoryStructure,
			Status:   stats.StatusError,
			Message: fmt.Sprintf(
				desc.Text(text.DescKeyDoctorCtxrcValidationError),
				validateErr,
			),
		})
		return
	}

	if len(warnings) > 0 {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckCtxrcValidation,
			Category: doctor.CategoryStructure,
			Status:   stats.StatusWarning,
			Message: fmt.Sprintf(
				desc.Text(
					text.DescKeyDoctorCtxrcValidationWarning),
				strings.Join(
					warnings, cfgToken.SemicolonSpace,
				),
			),
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     doctor.CheckCtxrcValidation,
		Category: doctor.CategoryStructure,
		Status:   stats.StatusOK,
		Message:  desc.Text(text.DescKeyDoctorCtxrcValidationOk),
	})
}

// Drift detects stale paths or missing files referenced in
// context.
//
// Parameters:
//   - report: Report to append the result to
func Drift(report *Report) {
	if !validate.Exists("") {
		return // skip drift check if not initialized
	}

	c, loadErr := load.Do("")
	if loadErr != nil {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckDrift,
			Category: doctor.CategoryQuality,
			Status:   stats.StatusWarning,
			Message: fmt.Sprintf(
				desc.Text(text.DescKeyDoctorDriftWarningLoad),
				loadErr,
			),
		})
		return
	}

	driftReport := drift.Detect(c)
	warnCount := len(driftReport.Warnings)
	violCount := len(driftReport.Violations)

	if warnCount == 0 && violCount == 0 {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckDrift,
			Category: doctor.CategoryQuality,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorDriftOk),
		})
		return
	}

	var parts []string
	if violCount > 0 {
		parts = append(
			parts,
			fmt.Sprintf(
				desc.Text(text.DescKeyDoctorDriftViolations),
				violCount,
			),
		)
	}
	if warnCount > 0 {
		parts = append(
			parts,
			fmt.Sprintf(
				desc.Text(text.DescKeyDoctorDriftWarnings),
				warnCount,
			),
		)
	}

	status := stats.StatusWarning
	if violCount > 0 {
		status = stats.StatusError
	}

	report.Results = append(report.Results, Result{
		Name:     doctor.CheckDrift,
		Category: doctor.CategoryQuality,
		Status:   status,
		Message: fmt.Sprintf(
			desc.Text(text.DescKeyDoctorDriftDetected),
			strings.Join(parts, cfgToken.CommaSpace),
		),
	})
}

// CompanionConfig reports whether companion tool checks
// are enabled or suppressed in .ctxrc.
//
// Parameters:
//   - report: Report to append the result to
func CompanionConfig(report *Report) {
	if rc.CompanionCheck() {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckCompanionConfig,
			Category: doctor.CategoryPlugin,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorCompanionConfigOk),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckCompanionConfig,
			Category: doctor.CategoryPlugin,
			Status:   stats.StatusInfo,
			Message:  desc.Text(text.DescKeyDoctorCompanionConfigInfo),
		})
	}
}

// PluginEnablement checks whether the ctx plugin is
// installed and enabled.
//
// Parameters:
//   - report: Report to append the result to
func PluginEnablement(report *Report) {
	installed := initCore.Installed()
	if !installed {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckPluginInstalled,
			Category: doctor.CategoryPlugin,
			Status:   stats.StatusInfo,
			Message:  desc.Text(text.DescKeyDoctorPluginInstalledInfo),
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     doctor.CheckPluginInstalled,
		Category: doctor.CategoryPlugin,
		Status:   stats.StatusOK,
		Message:  desc.Text(text.DescKeyDoctorPluginInstalledOk),
	})

	globalEnabled := initCore.EnabledGlobally()
	localEnabled := initCore.EnabledLocally()

	if globalEnabled {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckPluginEnabledGlobal,
			Category: doctor.CategoryPlugin,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorPluginEnabledGlobalOk),
		})
	}

	if localEnabled {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckPluginEnabledLocal,
			Category: doctor.CategoryPlugin,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorPluginEnabledLocalOk),
		})
	}

	if !globalEnabled && !localEnabled {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckPluginEnabled,
			Category: doctor.CategoryPlugin,
			Status:   stats.StatusWarning,
			Message: fmt.Sprintf(
				desc.Text(
					text.DescKeyDoctorPluginEnabledWarning,
				), claude.PluginID,
			),
		})
	}
}

// EventLogging checks whether event logging is enabled.
//
// Parameters:
//   - report: Report to append the result to
func EventLogging(report *Report) {
	if rc.EventLog() {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckEventLogging,
			Category: doctor.CategoryHooks,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorEventLoggingOk),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckEventLogging,
			Category: doctor.CategoryHooks,
			Status:   stats.StatusInfo,
			Message:  desc.Text(text.DescKeyDoctorEventLoggingInfo),
		})
	}
}

// Webhook checks whether a webhook notification endpoint
// is configured.
//
// Parameters:
//   - report: Report to append the result to
func Webhook(report *Report) {
	dir := rc.ContextDir()
	encPath := filepath.Join(dir, crypto.NotifyEnc)
	if _, statErr := os.Stat(encPath); statErr == nil {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckWebhook,
			Category: doctor.CategoryHooks,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorWebhookOk),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckWebhook,
			Category: doctor.CategoryHooks,
			Status:   stats.StatusInfo,
			Message:  desc.Text(text.DescKeyDoctorWebhookInfo),
		})
	}
}

// Reminders checks for pending reminders in the context
// directory.
//
// Parameters:
//   - report: Report to append the result to
func Reminders(report *Report) {
	dir := rc.ContextDir()
	remindersPath := filepath.Join(dir, reminder.File)
	data, readErr := io.SafeReadUserFile(remindersPath)
	if readErr != nil {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckReminders,
			Category: doctor.CategoryState,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorRemindersOk),
		})
		return
	}

	var reminders []any
	if unmarshalErr := json.Unmarshal(
		data, &reminders,
	); unmarshalErr != nil {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckReminders,
			Category: doctor.CategoryState,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorRemindersOk),
		})
		return
	}

	count := len(reminders)
	if count == 0 {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckReminders,
			Category: doctor.CategoryState,
			Status:   stats.StatusOK,
			Message:  desc.Text(text.DescKeyDoctorRemindersOk),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckReminders,
			Category: doctor.CategoryState,
			Status:   stats.StatusInfo,
			Message: fmt.Sprintf(
				desc.Text(text.DescKeyDoctorRemindersInfo),
				count,
			),
		})
	}
}

// TaskCompletion analyzes the task completion ratio and
// suggests archiving.
//
// Parameters:
//   - report: Report to append the result to
func TaskCompletion(report *Report) {
	dir := rc.ContextDir()
	tasksPath := filepath.Join(dir, ctx.Task)
	data, readErr := io.SafeReadUserFile(tasksPath)
	if readErr != nil {
		return // no tasks file, skip
	}

	matches := regex.TaskMultiline.FindAllStringSubmatch(
		string(data), -1,
	)
	var completed, pending int
	for _, m := range matches {
		if len(m) > 2 && m[2] == marker.MarkTaskComplete {
			completed++
		} else {
			pending++
		}
	}
	total := completed + pending

	if total == 0 {
		return // no tasks to report on
	}

	ratio := completed * stats.PercentMultiplier / total
	msg := fmt.Sprintf(desc.Text(
		text.DescKeyDoctorTaskCompletionFormat),
		completed, total, ratio,
	)

	aboveWarn := ratio >= doctor.TaskCompletionWarnPct
	aboveMin := completed > doctor.TaskCompletionMinCount
	if aboveWarn && aboveMin {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckTaskCompletion,
			Category: doctor.CategoryState,
			Status:   stats.StatusWarning,
			Message: msg + desc.Text(
				text.DescKeyDoctorTaskCompletionWarningSuffix,
			),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckTaskCompletion,
			Category: doctor.CategoryState,
			Status:   stats.StatusOK,
			Message:  msg,
		})
	}
}

// ContextTokenSize estimates context token usage and
// reports per-file breakdown.
//
// Parameters:
//   - report: Report to append the result to
func ContextTokenSize(report *Report) {
	indexed := make(
		map[string]bool, len(ctx.ReadOrder),
	)
	for _, f := range ctx.ReadOrder {
		indexed[f] = true
	}

	var totalTokens int
	c, loadErr := load.Do("")
	if loadErr != nil {
		return
	}

	type fileTokens struct {
		name   string
		tokens int
	}
	var breakdown []fileTokens

	for _, f := range c.Files {
		if indexed[f.Name] {
			tokens := token.Estimate(f.Content)
			totalTokens += tokens
			breakdown = append(
				breakdown,
				fileTokens{name: f.Name, tokens: tokens},
			)
		}
	}

	window := rc.ContextWindow()
	msg := fmt.Sprintf(
		desc.Text(text.DescKeyDoctorContextSizeFormat),
		totalTokens, window,
	)

	warnThreshold := window * doctor.ContextSizeWarnPct /
		stats.PercentMultiplier
	if totalTokens > warnThreshold {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckContextSize,
			Category: doctor.CategorySize,
			Status:   stats.StatusWarning,
			Message: msg + desc.Text(
				text.DescKeyDoctorContextSizeWarningSuffix,
			),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckContextSize,
			Category: doctor.CategorySize,
			Status:   stats.StatusOK,
			Message:  msg,
		})
	}

	for _, ft := range breakdown {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckContextFilePrefix + ft.name,
			Category: doctor.CategorySize,
			Status:   stats.StatusInfo,
			Message: fmt.Sprintf(
				desc.Text(text.DescKeyDoctorContextFileFormat),
				ft.name, ft.tokens,
			),
		})
	}
}

// RecentEventActivity reports the most recent event log
// entry.
//
// Parameters:
//   - report: Report to append the result to
func RecentEventActivity(report *Report) {
	if !rc.EventLog() {
		return // skip if logging disabled
	}

	events, queryErr := event.Query(
		entity.EventQueryOpts{Last: 1},
	)
	if queryErr != nil || len(events) == 0 {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckRecentEvents,
			Category: doctor.CategoryEvents,
			Status:   stats.StatusInfo,
			Message:  desc.Text(text.DescKeyDoctorRecentEventsInfo),
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     doctor.CheckRecentEvents,
		Category: doctor.CategoryEvents,
		Status:   stats.StatusOK,
		Message: fmt.Sprintf(
			desc.Text(text.DescKeyDoctorRecentEventsOk),
			events[len(events)-1].Timestamp,
		),
	})
}

// SystemResources collects and evaluates system resource
// metrics.
//
// Parameters:
//   - report: Report to append the result to
func SystemResources(report *Report) {
	snap := sysinfo.Collect()
	AddResourceResults(report, snap)
}

// AddResourceResults appends per-metric resource results to
// the report. Extracted for testability with constructed
// Snapshot values.
//
// Parameters:
//   - report: Report to append the results to
//   - snap: System resource snapshot to evaluate
func AddResourceResults(
	report *Report, snap sysinfo.Snapshot,
) {
	alerts := sysinfo.Evaluate(snap)

	sevMap := make(
		map[string]sysinfo.Severity, len(alerts),
	)
	for _, a := range alerts {
		sevMap[a.Resource] = a.Severity
	}

	byteChecks := []struct {
		supported bool
		used      uint64
		total     uint64
		fmtKey    string
		checkName string
		resource  string
	}{
		{
			snap.Memory.Supported,
			snap.Memory.UsedBytes,
			snap.Memory.TotalBytes,
			text.DescKeyDoctorResourceMemoryFormat,
			doctor.CheckResourceMemory,
			cfgSysinfo.ResourceMemory,
		},
		{
			snap.Memory.Supported,
			snap.Memory.SwapUsedBytes,
			snap.Memory.SwapTotalBytes,
			text.DescKeyDoctorResourceSwapFormat,
			doctor.CheckResourceSwap,
			cfgSysinfo.ResourceSwap,
		},
		{
			snap.Disk.Supported,
			snap.Disk.UsedBytes,
			snap.Disk.TotalBytes,
			text.DescKeyDoctorResourceDiskFormat,
			doctor.CheckResourceDisk,
			cfgSysinfo.ResourceDisk,
		},
	}
	for _, bc := range byteChecks {
		if !bc.supported || bc.total == 0 {
			continue
		}
		pct := ResourcePct(bc.used, bc.total)
		msg := fmt.Sprintf(desc.Text(bc.fmtKey),
			pct,
			sysinfo.FormatGiB(bc.used),
			sysinfo.FormatGiB(bc.total))
		report.Results = append(report.Results, Result{
			Name:     bc.checkName,
			Category: doctor.CategoryResources,
			Status:   SeverityToStatus(sevMap[bc.resource]),
			Message:  msg,
		})
	}

	if snap.Load.Supported && snap.Load.NumCPU > 0 {
		ratio := snap.Load.Load1 /
			float64(snap.Load.NumCPU)
		msg := fmt.Sprintf(
			desc.Text(text.DescKeyDoctorResourceLoadFormat),
			ratio, snap.Load.Load1, snap.Load.NumCPU)
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckResourceLoad,
			Category: doctor.CategoryResources,
			Status: SeverityToStatus(
				sevMap[cfgSysinfo.ResourceLoad],
			),
			Message: msg,
		})
	}
}

// SeverityToStatus converts a sysinfo.Severity to a doctor
// status string.
//
// Parameters:
//   - sev: Severity level from system resource evaluation
//
// Returns:
//   - string: Corresponding status constant
func SeverityToStatus(sev sysinfo.Severity) string {
	switch sev {
	case sysinfo.SeverityWarning:
		return stats.StatusWarning
	case sysinfo.SeverityDanger:
		return stats.StatusError
	default:
		return stats.StatusOK
	}
}

// ResourcePct calculates the percentage of used versus
// total.
//
// Parameters:
//   - used: Used amount
//   - total: Total capacity
//
// Returns:
//   - int: Percentage (0 if the total is 0)
func ResourcePct(used, total uint64) int {
	if total == 0 {
		return 0
	}
	return int(
		float64(used) / float64(total) *
			stats.PercentMultiplier,
	)
}
