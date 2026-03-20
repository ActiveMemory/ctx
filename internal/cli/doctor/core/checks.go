//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/crypto"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/doctor"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/reminder"
	"github.com/ActiveMemory/ctx/internal/context/load"
	"github.com/ActiveMemory/ctx/internal/context/token"
	"github.com/ActiveMemory/ctx/internal/context/validate"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/log"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

// CheckContextInitialized verifies that a .context/ directory exists.
//
// Parameters:
//   - report: Report to append the result to
func CheckContextInitialized(report *Report) {
	if validate.Exists("") {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckContextInit,
			Category: doctor.CategoryStructure,
			Status:   StatusOK,
			Message:  desc.TextDesc(text.DescKeyDoctorContextInitializedOk),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckContextInit,
			Category: doctor.CategoryStructure,
			Status:   StatusError,
			Message:  desc.TextDesc(text.DescKeyDoctorContextInitializedError),
		})
	}
}

// CheckRequiredFiles verifies that all required context files are present.
//
// Parameters:
//   - report: Report to append the result to
func CheckRequiredFiles(report *Report) {
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
			Status:   StatusOK,
			Message:  fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorRequiredFilesOk), present, total),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckRequiredFiles,
			Category: doctor.CategoryStructure,
			Status:   StatusError,
			Message:  fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorRequiredFilesError), present, total, strings.Join(missing, ", ")),
		})
	}
}

// CheckCtxrcValidation validates the .ctxrc file for unknown fields or parse errors.
//
// Parameters:
//   - report: Report to append the result to
func CheckCtxrcValidation(report *Report) {
	data, readErr := io.SafeReadUserFile(file.CtxRC)
	if readErr != nil {
		// No .ctxrc is fine — defaults are used.
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckCtxrcValidation,
			Category: doctor.CategoryStructure,
			Status:   StatusOK,
			Message:  desc.TextDesc(text.DescKeyDoctorCtxrcValidationOkNoFile),
		})
		return
	}

	warnings, validateErr := rc.Validate(data)
	if validateErr != nil {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckCtxrcValidation,
			Category: doctor.CategoryStructure,
			Status:   StatusError,
			Message:  fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorCtxrcValidationError), validateErr),
		})
		return
	}

	if len(warnings) > 0 {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckCtxrcValidation,
			Category: doctor.CategoryStructure,
			Status:   StatusWarning,
			Message:  fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorCtxrcValidationWarning), strings.Join(warnings, "; ")),
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     doctor.CheckCtxrcValidation,
		Category: doctor.CategoryStructure,
		Status:   StatusOK,
		Message:  desc.TextDesc(text.DescKeyDoctorCtxrcValidationOk),
	})
}

// CheckDrift detects stale paths or missing files referenced in context.
//
// Parameters:
//   - report: Report to append the result to
func CheckDrift(report *Report) {
	if !validate.Exists("") {
		return // skip drift check if not initialized
	}

	ctx, loadErr := load.Do("")
	if loadErr != nil {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckDrift,
			Category: doctor.CategoryQuality,
			Status:   StatusWarning,
			Message:  fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorDriftWarningLoad), loadErr),
		})
		return
	}

	driftReport := drift.Detect(ctx)
	warnCount := len(driftReport.Warnings)
	violCount := len(driftReport.Violations)

	if warnCount == 0 && violCount == 0 {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckDrift,
			Category: doctor.CategoryQuality,
			Status:   StatusOK,
			Message:  desc.TextDesc(text.DescKeyDoctorDriftOk),
		})
		return
	}

	var parts []string
	if violCount > 0 {
		parts = append(parts, fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorDriftViolations), violCount))
	}
	if warnCount > 0 {
		parts = append(parts, fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorDriftWarnings), warnCount))
	}

	status := StatusWarning
	if violCount > 0 {
		status = StatusError
	}

	report.Results = append(report.Results, Result{
		Name:     doctor.CheckDrift,
		Category: doctor.CategoryQuality,
		Status:   status,
		Message:  fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorDriftDetected), strings.Join(parts, ", ")),
	})
}

// CheckPluginEnablement checks whether the ctx plugin is installed and enabled.
//
// Parameters:
//   - report: Report to append the result to
func CheckPluginEnablement(report *Report) {
	installed := initialize.PluginInstalled()
	if !installed {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckPluginInstalled,
			Category: doctor.CategoryPlugin,
			Status:   StatusInfo,
			Message:  desc.TextDesc(text.DescKeyDoctorPluginInstalledInfo),
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     doctor.CheckPluginInstalled,
		Category: doctor.CategoryPlugin,
		Status:   StatusOK,
		Message:  desc.TextDesc(text.DescKeyDoctorPluginInstalledOk),
	})

	globalEnabled := initialize.PluginEnabledGlobally()
	localEnabled := initialize.PluginEnabledLocally()

	if globalEnabled {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckPluginEnabledGlobal,
			Category: doctor.CategoryPlugin,
			Status:   StatusOK,
			Message:  desc.TextDesc(text.DescKeyDoctorPluginEnabledGlobalOk),
		})
	}

	if localEnabled {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckPluginEnabledLocal,
			Category: doctor.CategoryPlugin,
			Status:   StatusOK,
			Message:  desc.TextDesc(text.DescKeyDoctorPluginEnabledLocalOk),
		})
	}

	if !globalEnabled && !localEnabled {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckPluginEnabled,
			Category: doctor.CategoryPlugin,
			Status:   StatusWarning,
			Message:  fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorPluginEnabledWarning), claude.PluginID),
		})
	}
}

// CheckEventLogging checks whether event logging is enabled.
//
// Parameters:
//   - report: Report to append the result to
func CheckEventLogging(report *Report) {
	if rc.EventLog() {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckEventLogging,
			Category: doctor.CategoryHooks,
			Status:   StatusOK,
			Message:  desc.TextDesc(text.DescKeyDoctorEventLoggingOk),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckEventLogging,
			Category: doctor.CategoryHooks,
			Status:   StatusInfo,
			Message:  desc.TextDesc(text.DescKeyDoctorEventLoggingInfo),
		})
	}
}

// CheckWebhook checks whether a webhook notification endpoint is configured.
//
// Parameters:
//   - report: Report to append the result to
func CheckWebhook(report *Report) {
	dir := rc.ContextDir()
	encPath := filepath.Join(dir, crypto.NotifyEnc)
	if _, statErr := os.Stat(encPath); statErr == nil {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckWebhook,
			Category: doctor.CategoryHooks,
			Status:   StatusOK,
			Message:  desc.TextDesc(text.DescKeyDoctorWebhookOk),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckWebhook,
			Category: doctor.CategoryHooks,
			Status:   StatusInfo,
			Message:  desc.TextDesc(text.DescKeyDoctorWebhookInfo),
		})
	}
}

// CheckReminders checks for pending reminders in the context directory.
//
// Parameters:
//   - report: Report to append the result to
func CheckReminders(report *Report) {
	dir := rc.ContextDir()
	remindersPath := filepath.Join(dir, reminder.Reminders)
	data, readErr := io.SafeReadUserFile(remindersPath)
	if readErr != nil {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckReminders,
			Category: doctor.CategoryState,
			Status:   StatusOK,
			Message:  desc.TextDesc(text.DescKeyDoctorRemindersOk),
		})
		return
	}

	var reminders []any
	if unmarshalErr := json.Unmarshal(data, &reminders); unmarshalErr != nil {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckReminders,
			Category: doctor.CategoryState,
			Status:   StatusOK,
			Message:  desc.TextDesc(text.DescKeyDoctorRemindersOk),
		})
		return
	}

	count := len(reminders)
	if count == 0 {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckReminders,
			Category: doctor.CategoryState,
			Status:   StatusOK,
			Message:  desc.TextDesc(text.DescKeyDoctorRemindersOk),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckReminders,
			Category: doctor.CategoryState,
			Status:   StatusInfo,
			Message:  fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorRemindersInfo), count),
		})
	}
}

// CheckTaskCompletion analyzes task completion ratio and suggests archiving.
//
// Parameters:
//   - report: Report to append the result to
func CheckTaskCompletion(report *Report) {
	dir := rc.ContextDir()
	tasksPath := filepath.Join(dir, ctx.Task)
	data, readErr := io.SafeReadUserFile(tasksPath)
	if readErr != nil {
		return // no tasks file, skip
	}

	matches := regex.TaskMultiline.FindAllStringSubmatch(string(data), -1)
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

	ratio := completed * 100 / total
	msg := fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorTaskCompletionFormat), completed, total, ratio)

	if ratio >= 80 && completed > 5 {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckTaskCompletion,
			Category: doctor.CategoryState,
			Status:   StatusWarning,
			Message:  msg + desc.TextDesc(text.DescKeyDoctorTaskCompletionWarningSuffix),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckTaskCompletion,
			Category: doctor.CategoryState,
			Status:   StatusOK,
			Message:  msg,
		})
	}
}

// CheckContextTokenSize estimates context token usage and reports per-file breakdown.
//
// Parameters:
//   - report: Report to append the result to
func CheckContextTokenSize(report *Report) {
	// Only count files in ReadOrder — these are the files actually
	// loaded into agent context. Other .md files (DETAILED_DESIGN.md,
	// map-tracking, etc.) exist on disk but aren't injected.
	indexed := make(map[string]bool, len(ctx.ReadOrder))
	for _, f := range ctx.ReadOrder {
		indexed[f] = true
	}

	var totalTokens int
	ctx, loadErr := load.Do("")
	if loadErr != nil {
		return
	}

	// Collect per-file token counts for breakdown.
	type fileTokens struct {
		name   string
		tokens int
	}
	var breakdown []fileTokens

	for _, f := range ctx.Files {
		if indexed[f.Name] {
			tokens := token.EstimateTokens(f.Content)
			totalTokens += tokens
			breakdown = append(breakdown, fileTokens{name: f.Name, tokens: tokens})
		}
	}

	window := rc.ContextWindow()
	msg := fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorContextSizeFormat), totalTokens, window)

	warnThreshold := window / 5 // 20% of context window
	if totalTokens > warnThreshold {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckContextSize,
			Category: doctor.CategorySize,
			Status:   StatusWarning,
			Message:  msg + desc.TextDesc(text.DescKeyDoctorContextSizeWarningSuffix),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckContextSize,
			Category: doctor.CategorySize,
			Status:   StatusOK,
			Message:  msg,
		})
	}

	// Add per-file breakdown as info results.
	for _, ft := range breakdown {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckContextFilePrefix + ft.name,
			Category: doctor.CategorySize,
			Status:   StatusInfo,
			Message:  fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorContextFileFormat), ft.name, ft.tokens),
		})
	}
}

// CheckRecentEventActivity reports the most recent event log entry.
//
// Parameters:
//   - report: Report to append the result to
func CheckRecentEventActivity(report *Report) {
	if !rc.EventLog() {
		return // skip if logging disabled
	}

	events, queryErr := log.Query(log.QueryOpts{Last: 1})
	if queryErr != nil || len(events) == 0 {
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckRecentEvents,
			Category: doctor.CategoryEvents,
			Status:   StatusInfo,
			Message:  desc.TextDesc(text.DescKeyDoctorRecentEventsInfo),
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     doctor.CheckRecentEvents,
		Category: doctor.CategoryEvents,
		Status:   StatusOK,
		Message:  fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorRecentEventsOk), events[len(events)-1].Timestamp),
	})
}

// CheckSystemResources collects and evaluates system resource metrics.
//
// Parameters:
//   - report: Report to append the result to
func CheckSystemResources(report *Report) {
	snap := sysinfo.Collect(".")
	AddResourceResults(report, snap)
}

// AddResourceResults appends per-metric resource results to the report.
// Extracted for testability with constructed Snapshot values.
//
// Parameters:
//   - report: Report to append the results to
//   - snap: System resource snapshot to evaluate
func AddResourceResults(report *Report, snap sysinfo.Snapshot) {
	alerts := sysinfo.Evaluate(snap)

	// Build severity lookup by resource name.
	sevMap := make(map[string]sysinfo.Severity, len(alerts))
	for _, a := range alerts {
		sevMap[a.Resource] = a.Severity
	}

	// Memory.
	if snap.Memory.Supported && snap.Memory.TotalBytes > 0 {
		pct := ResourcePct(snap.Memory.UsedBytes, snap.Memory.TotalBytes)
		msg := fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorResourceMemoryFormat),
			pct,
			sysinfo.FormatGiB(snap.Memory.UsedBytes),
			sysinfo.FormatGiB(snap.Memory.TotalBytes))
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckResourceMemory,
			Category: doctor.CategoryResources,
			Status:   SeverityToStatus(sevMap["memory"]),
			Message:  msg,
		})
	}

	// Swap (only when swap is configured).
	if snap.Memory.Supported && snap.Memory.SwapTotalBytes > 0 {
		pct := ResourcePct(snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes)
		msg := fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorResourceSwapFormat),
			pct,
			sysinfo.FormatGiB(snap.Memory.SwapUsedBytes),
			sysinfo.FormatGiB(snap.Memory.SwapTotalBytes))
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckResourceSwap,
			Category: doctor.CategoryResources,
			Status:   SeverityToStatus(sevMap["swap"]),
			Message:  msg,
		})
	}

	// Disk.
	if snap.Disk.Supported && snap.Disk.TotalBytes > 0 {
		pct := ResourcePct(snap.Disk.UsedBytes, snap.Disk.TotalBytes)
		msg := fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorResourceDiskFormat),
			pct,
			sysinfo.FormatGiB(snap.Disk.UsedBytes),
			sysinfo.FormatGiB(snap.Disk.TotalBytes))
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckResourceDisk,
			Category: doctor.CategoryResources,
			Status:   SeverityToStatus(sevMap["disk"]),
			Message:  msg,
		})
	}

	// Do (1-minute average relative to CPU count).
	if snap.Load.Supported && snap.Load.NumCPU > 0 {
		ratio := snap.Load.Load1 / float64(snap.Load.NumCPU)
		msg := fmt.Sprintf(desc.TextDesc(text.DescKeyDoctorResourceLoadFormat),
			ratio, snap.Load.Load1, snap.Load.NumCPU)
		report.Results = append(report.Results, Result{
			Name:     doctor.CheckResourceLoad,
			Category: doctor.CategoryResources,
			Status:   SeverityToStatus(sevMap["load"]),
			Message:  msg,
		})
	}
}

// SeverityToStatus converts a sysinfo.Severity to a doctor status string.
//
// Parameters:
//   - sev: Severity level from system resource evaluation
//
// Returns:
//   - string: Corresponding status constant (StatusOK, StatusWarning, StatusError)
func SeverityToStatus(sev sysinfo.Severity) string {
	switch sev {
	case sysinfo.SeverityWarning:
		return StatusWarning
	case sysinfo.SeverityDanger:
		return StatusError
	default:
		return StatusOK
	}
}

// ResourcePct calculates the percentage of used vs total.
//
// Parameters:
//   - used: Used amount
//   - total: Total capacity
//
// Returns:
//   - int: Percentage (0 if total is 0)
func ResourcePct(used, total uint64) int {
	if total == 0 {
		return 0
	}
	return int(float64(used) / float64(total) * 100)
}
