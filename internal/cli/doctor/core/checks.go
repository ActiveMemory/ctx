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

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

// CheckContextInitialized verifies that a .context/ directory exists.
//
// Parameters:
//   - report: Report to append the result to
func CheckContextInitialized(report *Report) {
	if context.Exists("") {
		report.Results = append(report.Results, Result{
			Name:     "context_initialized",
			Category: "Structure",
			Status:   StatusOK,
			Message:  "Context initialized (.context/)",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "context_initialized",
			Category: "Structure",
			Status:   StatusError,
			Message:  "Context not initialized — run ctx init",
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
	for _, f := range config.FilesRequired {
		path := filepath.Join(dir, f)
		if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
			missing = append(missing, f)
		}
	}

	total := len(config.FilesRequired)
	present := total - len(missing)

	if len(missing) == 0 {
		report.Results = append(report.Results, Result{
			Name:     "required_files",
			Category: "Structure",
			Status:   StatusOK,
			Message:  fmt.Sprintf("Required files present (%d/%d)", present, total),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "required_files",
			Category: "Structure",
			Status:   StatusError,
			Message:  fmt.Sprintf("Missing required files (%d/%d): %s", present, total, strings.Join(missing, ", ")),
		})
	}
}

// CheckCtxrcValidation validates the .ctxrc file for unknown fields or parse errors.
//
// Parameters:
//   - report: Report to append the result to
func CheckCtxrcValidation(report *Report) {
	data, readErr := os.ReadFile(config.FileContextRC) //nolint:gosec // project-local config file
	if readErr != nil {
		// No .ctxrc is fine — defaults are used.
		report.Results = append(report.Results, Result{
			Name:     "ctxrc_validation",
			Category: "Structure",
			Status:   StatusOK,
			Message:  "No .ctxrc file (using defaults)",
		})
		return
	}

	warnings, validateErr := rc.Validate(data)
	if validateErr != nil {
		report.Results = append(report.Results, Result{
			Name:     "ctxrc_validation",
			Category: "Structure",
			Status:   StatusError,
			Message:  fmt.Sprintf(".ctxrc parse error: %v", validateErr),
		})
		return
	}

	if len(warnings) > 0 {
		report.Results = append(report.Results, Result{
			Name:     "ctxrc_validation",
			Category: "Structure",
			Status:   StatusWarning,
			Message:  fmt.Sprintf(".ctxrc has unknown fields: %s", strings.Join(warnings, "; ")),
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     "ctxrc_validation",
		Category: "Structure",
		Status:   StatusOK,
		Message:  ".ctxrc valid",
	})
}

// CheckDrift detects stale paths or missing files referenced in context.
//
// Parameters:
//   - report: Report to append the result to
func CheckDrift(report *Report) {
	if !context.Exists("") {
		return // skip drift check if not initialized
	}

	ctx, loadErr := context.Load("")
	if loadErr != nil {
		report.Results = append(report.Results, Result{
			Name:     "drift",
			Category: "Quality",
			Status:   StatusWarning,
			Message:  fmt.Sprintf("Could not load context for drift check: %v", loadErr),
		})
		return
	}

	driftReport := drift.Detect(ctx)
	warnCount := len(driftReport.Warnings)
	violCount := len(driftReport.Violations)

	if warnCount == 0 && violCount == 0 {
		report.Results = append(report.Results, Result{
			Name:     "drift",
			Category: "Quality",
			Status:   StatusOK,
			Message:  "No drift detected",
		})
		return
	}

	var parts []string
	if violCount > 0 {
		parts = append(parts, fmt.Sprintf("%d violations", violCount))
	}
	if warnCount > 0 {
		parts = append(parts, fmt.Sprintf("%d warnings", warnCount))
	}

	status := StatusWarning
	if violCount > 0 {
		status = StatusError
	}

	report.Results = append(report.Results, Result{
		Name:     "drift",
		Category: "Quality",
		Status:   status,
		Message:  fmt.Sprintf("Drift: %s — run ctx drift for details", strings.Join(parts, ", ")),
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
			Name:     "plugin_installed",
			Category: "Plugin",
			Status:   StatusInfo,
			Message:  "ctx plugin not installed",
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     "plugin_installed",
		Category: "Plugin",
		Status:   StatusOK,
		Message:  "ctx plugin installed",
	})

	globalEnabled := initialize.PluginEnabledGlobally()
	localEnabled := initialize.PluginEnabledLocally()

	if globalEnabled {
		report.Results = append(report.Results, Result{
			Name:     "plugin_enabled_global",
			Category: "Plugin",
			Status:   StatusOK,
			Message:  "Plugin enabled globally (~/.claude/settings.json)",
		})
	}

	if localEnabled {
		report.Results = append(report.Results, Result{
			Name:     "plugin_enabled_local",
			Category: "Plugin",
			Status:   StatusOK,
			Message:  "Plugin enabled locally (.claude/settings.local.json)",
		})
	}

	if !globalEnabled && !localEnabled {
		report.Results = append(report.Results, Result{
			Name:     "plugin_enabled",
			Category: "Plugin",
			Status:   StatusWarning,
			Message: "Plugin installed but not enabled — run 'ctx init' to auto-enable, " +
				"or add {\"enabledPlugins\": {\"" + config.PluginID +
				"\": true}} to ~/.claude/settings.json",
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
			Name:     "event_logging",
			Category: "Hooks",
			Status:   StatusOK,
			Message:  "Event logging enabled",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "event_logging",
			Category: "Hooks",
			Status:   StatusInfo,
			Message:  "Event logging disabled (enable with event_log: true in .ctxrc)",
		})
	}
}

// CheckWebhook checks whether a webhook notification endpoint is configured.
//
// Parameters:
//   - report: Report to append the result to
func CheckWebhook(report *Report) {
	dir := rc.ContextDir()
	encPath := filepath.Join(dir, ".notify.enc")
	if _, statErr := os.Stat(encPath); statErr == nil {
		report.Results = append(report.Results, Result{
			Name:     "webhook",
			Category: "Hooks",
			Status:   StatusOK,
			Message:  "Webhook configured",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "webhook",
			Category: "Hooks",
			Status:   StatusInfo,
			Message:  "No webhook configured (optional — use ctx notify setup)",
		})
	}
}

// CheckReminders checks for pending reminders in the context directory.
//
// Parameters:
//   - report: Report to append the result to
func CheckReminders(report *Report) {
	dir := rc.ContextDir()
	remindersPath := filepath.Join(dir, "reminders.json")
	data, readErr := os.ReadFile(remindersPath) //nolint:gosec // project-local path
	if readErr != nil {
		report.Results = append(report.Results, Result{
			Name:     "reminders",
			Category: "State",
			Status:   StatusOK,
			Message:  "No pending reminders",
		})
		return
	}

	var reminders []any
	if unmarshalErr := json.Unmarshal(data, &reminders); unmarshalErr != nil {
		report.Results = append(report.Results, Result{
			Name:     "reminders",
			Category: "State",
			Status:   StatusOK,
			Message:  "No pending reminders",
		})
		return
	}

	count := len(reminders)
	if count == 0 {
		report.Results = append(report.Results, Result{
			Name:     "reminders",
			Category: "State",
			Status:   StatusOK,
			Message:  "No pending reminders",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "reminders",
			Category: "State",
			Status:   StatusInfo,
			Message:  fmt.Sprintf("%d pending reminders", count),
		})
	}
}

// CheckTaskCompletion analyzes task completion ratio and suggests archiving.
//
// Parameters:
//   - report: Report to append the result to
func CheckTaskCompletion(report *Report) {
	dir := rc.ContextDir()
	tasksPath := filepath.Join(dir, config.FileTask)
	data, readErr := os.ReadFile(tasksPath) //nolint:gosec // project-local path
	if readErr != nil {
		return // no tasks file, skip
	}

	matches := config.RegExTaskMultiline.FindAllStringSubmatch(string(data), -1)
	var completed, pending int
	for _, m := range matches {
		if len(m) > 2 && m[2] == config.MarkTaskComplete {
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
	msg := fmt.Sprintf("Tasks: %d/%d completed (%d%%)", completed, total, ratio)

	if ratio >= 80 && completed > 5 {
		report.Results = append(report.Results, Result{
			Name:     "task_completion",
			Category: "State",
			Status:   StatusWarning,
			Message:  msg + " — consider archiving with ctx tasks archive",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "task_completion",
			Category: "State",
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
	// Only count files in FileReadOrder — these are the files actually
	// loaded into agent context. Other .md files (DETAILED_DESIGN.md,
	// map-tracking, etc.) exist on disk but aren't injected.
	indexed := make(map[string]bool, len(config.FileReadOrder))
	for _, f := range config.FileReadOrder {
		indexed[f] = true
	}

	var totalTokens int
	ctx, loadErr := context.Load("")
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
			tokens := context.EstimateTokens(f.Content)
			totalTokens += tokens
			breakdown = append(breakdown, fileTokens{name: f.Name, tokens: tokens})
		}
	}

	window := rc.ContextWindow()
	msg := fmt.Sprintf("Context size: ~%d tokens (window: %d)", totalTokens, window)

	warnThreshold := window / 5 // 20% of context window
	if totalTokens > warnThreshold {
		report.Results = append(report.Results, Result{
			Name:     "context_size",
			Category: "Size",
			Status:   StatusWarning,
			Message:  msg + " — consider ctx compact",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "context_size",
			Category: "Size",
			Status:   StatusOK,
			Message:  msg,
		})
	}

	// Add per-file breakdown as info results.
	for _, ft := range breakdown {
		report.Results = append(report.Results, Result{
			Name:     "context_file_" + ft.name,
			Category: "Size",
			Status:   StatusInfo,
			Message:  fmt.Sprintf("%-22s ~%d tokens", ft.name, ft.tokens),
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

	events, queryErr := eventlog.Query(eventlog.QueryOpts{Last: 1})
	if queryErr != nil || len(events) == 0 {
		report.Results = append(report.Results, Result{
			Name:     "recent_events",
			Category: "Events",
			Status:   StatusInfo,
			Message:  "No events in log",
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     "recent_events",
		Category: "Events",
		Status:   StatusOK,
		Message:  fmt.Sprintf("Last event: %s", events[len(events)-1].Timestamp),
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
		msg := fmt.Sprintf("Memory %d%% (%s / %s GB)",
			pct,
			sysinfo.FormatGiB(snap.Memory.UsedBytes),
			sysinfo.FormatGiB(snap.Memory.TotalBytes))
		report.Results = append(report.Results, Result{
			Name:     "resource_memory",
			Category: "Resources",
			Status:   SeverityToStatus(sevMap["memory"]),
			Message:  msg,
		})
	}

	// Swap (only when swap is configured).
	if snap.Memory.Supported && snap.Memory.SwapTotalBytes > 0 {
		pct := ResourcePct(snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes)
		msg := fmt.Sprintf("Swap %d%% (%s / %s GB)",
			pct,
			sysinfo.FormatGiB(snap.Memory.SwapUsedBytes),
			sysinfo.FormatGiB(snap.Memory.SwapTotalBytes))
		report.Results = append(report.Results, Result{
			Name:     "resource_swap",
			Category: "Resources",
			Status:   SeverityToStatus(sevMap["swap"]),
			Message:  msg,
		})
	}

	// Disk.
	if snap.Disk.Supported && snap.Disk.TotalBytes > 0 {
		pct := ResourcePct(snap.Disk.UsedBytes, snap.Disk.TotalBytes)
		msg := fmt.Sprintf("Disk %d%% (%s / %s GB)",
			pct,
			sysinfo.FormatGiB(snap.Disk.UsedBytes),
			sysinfo.FormatGiB(snap.Disk.TotalBytes))
		report.Results = append(report.Results, Result{
			Name:     "resource_disk",
			Category: "Resources",
			Status:   SeverityToStatus(sevMap["disk"]),
			Message:  msg,
		})
	}

	// Load (1-minute average relative to CPU count).
	if snap.Load.Supported && snap.Load.NumCPU > 0 {
		ratio := snap.Load.Load1 / float64(snap.Load.NumCPU)
		msg := fmt.Sprintf("Load %.2fx (%.1f / %d CPUs)",
			ratio, snap.Load.Load1, snap.Load.NumCPU)
		report.Results = append(report.Results, Result{
			Name:     "resource_load",
			Category: "Resources",
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
