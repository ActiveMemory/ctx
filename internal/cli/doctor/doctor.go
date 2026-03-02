//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package doctor provides the "ctx doctor" command for structural
// health checks across context, hooks, and configuration.
package doctor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/initialize"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/sysinfo"
)

// Status constants for check results.
const (
	statusOK      = "ok"
	statusWarning = "warning"
	statusError   = "error"
	statusInfo    = "info"
)

// Result represents a single check outcome.
type Result struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Status   string `json:"status"` // "ok", "warning", "error", "info"
	Message  string `json:"message"`
}

// Report is the complete doctor output.
type Report struct {
	Results  []Result `json:"results"`
	Warnings int      `json:"warnings"`
	Errors   int      `json:"errors"`
}

// Cmd returns the "ctx doctor" command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Structural health check",
		Long: `Run mechanical health checks across context, hooks, and configuration.

Checks:
  - Context initialized and required files present
  - Drift detected (stale paths, missing files)
  - Plugin installed and enabled
  - Event logging status
  - Webhook configured
  - Pending reminders
  - Task completion ratio
  - Context token size
  - System resources (memory, swap, disk, load)

Use --json for machine-readable output.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			jsonOut, _ := cmd.Flags().GetBool("json")
			return runDoctor(cmd, jsonOut)
		},
	}
	cmd.Flags().BoolP("json", "j", false, "Machine-readable JSON output")
	return cmd
}

func runDoctor(cmd *cobra.Command, jsonOutput bool) error {
	report := &Report{}

	checkContextInitialized(report)
	checkRequiredFiles(report)
	checkDrift(report)
	checkPluginEnablement(report)
	checkEventLogging(report)
	checkWebhook(report)
	checkReminders(report)
	checkTaskCompletion(report)
	checkContextTokenSize(report)
	checkSystemResources(report)
	checkRecentEventActivity(report)

	// Count warnings and errors.
	for _, r := range report.Results {
		switch r.Status {
		case statusWarning:
			report.Warnings++
		case statusError:
			report.Errors++
		}
	}

	if jsonOutput {
		return outputDoctorJSON(cmd, report)
	}
	return outputDoctorHuman(cmd, report)
}

func checkContextInitialized(report *Report) {
	if context.Exists("") {
		report.Results = append(report.Results, Result{
			Name:     "context_initialized",
			Category: "Structure",
			Status:   statusOK,
			Message:  "Context initialized (.context/)",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "context_initialized",
			Category: "Structure",
			Status:   statusError,
			Message:  "Context not initialized — run ctx init",
		})
	}
}

func checkRequiredFiles(report *Report) {
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
			Status:   statusOK,
			Message:  fmt.Sprintf("Required files present (%d/%d)", present, total),
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "required_files",
			Category: "Structure",
			Status:   statusError,
			Message:  fmt.Sprintf("Missing required files (%d/%d): %s", present, total, strings.Join(missing, ", ")),
		})
	}
}

func checkDrift(report *Report) {
	if !context.Exists("") {
		return // skip drift check if not initialized
	}

	ctx, loadErr := context.Load("")
	if loadErr != nil {
		report.Results = append(report.Results, Result{
			Name:     "drift",
			Category: "Quality",
			Status:   statusWarning,
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
			Status:   statusOK,
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

	status := "warning"
	if violCount > 0 {
		status = "error"
	}

	report.Results = append(report.Results, Result{
		Name:     "drift",
		Category: "Quality",
		Status:   status,
		Message:  fmt.Sprintf("Drift: %s — run ctx drift for details", strings.Join(parts, ", ")),
	})
}

func checkPluginEnablement(report *Report) {
	installed := initialize.PluginInstalled()
	if !installed {
		report.Results = append(report.Results, Result{
			Name:     "plugin_installed",
			Category: "Plugin",
			Status:   statusInfo,
			Message:  "ctx plugin not installed",
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     "plugin_installed",
		Category: "Plugin",
		Status:   statusOK,
		Message:  "ctx plugin installed",
	})

	globalEnabled := initialize.PluginEnabledGlobally()
	localEnabled := initialize.PluginEnabledLocally()

	if globalEnabled {
		report.Results = append(report.Results, Result{
			Name:     "plugin_enabled_global",
			Category: "Plugin",
			Status:   statusOK,
			Message:  "Plugin enabled globally (~/.claude/settings.json)",
		})
	}

	if localEnabled {
		report.Results = append(report.Results, Result{
			Name:     "plugin_enabled_local",
			Category: "Plugin",
			Status:   statusOK,
			Message:  "Plugin enabled locally (.claude/settings.local.json)",
		})
	}

	if !globalEnabled && !localEnabled {
		report.Results = append(report.Results, Result{
			Name:     "plugin_enabled",
			Category: "Plugin",
			Status:   statusWarning,
			Message: "Plugin installed but not enabled — run 'ctx init' to auto-enable, " +
				"or add {\"enabledPlugins\": {\"" + config.PluginID +
				"\": true}} to ~/.claude/settings.json",
		})
	}
}

func checkEventLogging(report *Report) {
	if rc.EventLog() {
		report.Results = append(report.Results, Result{
			Name:     "event_logging",
			Category: "Hooks",
			Status:   statusOK,
			Message:  "Event logging enabled",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "event_logging",
			Category: "Hooks",
			Status:   statusInfo,
			Message:  "Event logging disabled (enable with event_log: true in .ctxrc)",
		})
	}
}

func checkWebhook(report *Report) {
	dir := rc.ContextDir()
	encPath := filepath.Join(dir, ".notify.enc")
	if _, statErr := os.Stat(encPath); statErr == nil {
		report.Results = append(report.Results, Result{
			Name:     "webhook",
			Category: "Hooks",
			Status:   statusOK,
			Message:  "Webhook configured",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "webhook",
			Category: "Hooks",
			Status:   statusInfo,
			Message:  "No webhook configured (optional — use ctx notify setup)",
		})
	}
}

func checkReminders(report *Report) {
	dir := rc.ContextDir()
	remindersPath := filepath.Join(dir, "reminders.json")
	data, readErr := os.ReadFile(remindersPath) //nolint:gosec // project-local path
	if readErr != nil {
		report.Results = append(report.Results, Result{
			Name:     "reminders",
			Category: "State",
			Status:   statusOK,
			Message:  "No pending reminders",
		})
		return
	}

	var reminders []any
	if unmarshalErr := json.Unmarshal(data, &reminders); unmarshalErr != nil {
		report.Results = append(report.Results, Result{
			Name:     "reminders",
			Category: "State",
			Status:   statusOK,
			Message:  "No pending reminders",
		})
		return
	}

	count := len(reminders)
	if count == 0 {
		report.Results = append(report.Results, Result{
			Name:     "reminders",
			Category: "State",
			Status:   statusOK,
			Message:  "No pending reminders",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "reminders",
			Category: "State",
			Status:   statusInfo,
			Message:  fmt.Sprintf("%d pending reminders", count),
		})
	}
}

func checkTaskCompletion(report *Report) {
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
			Status:   statusWarning,
			Message:  msg + " — consider archiving with ctx tasks archive",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "task_completion",
			Category: "State",
			Status:   statusOK,
			Message:  msg,
		})
	}
}

func checkContextTokenSize(report *Report) {
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
			Status:   statusWarning,
			Message:  msg + " — consider ctx compact",
		})
	} else {
		report.Results = append(report.Results, Result{
			Name:     "context_size",
			Category: "Size",
			Status:   statusOK,
			Message:  msg,
		})
	}

	// Add per-file breakdown as info results.
	for _, ft := range breakdown {
		report.Results = append(report.Results, Result{
			Name:     "context_file_" + ft.name,
			Category: "Size",
			Status:   statusInfo,
			Message:  fmt.Sprintf("%-22s ~%d tokens", ft.name, ft.tokens),
		})
	}
}

func checkRecentEventActivity(report *Report) {
	if !rc.EventLog() {
		return // skip if logging disabled
	}

	events, queryErr := eventlog.Query(eventlog.QueryOpts{Last: 1})
	if queryErr != nil || len(events) == 0 {
		report.Results = append(report.Results, Result{
			Name:     "recent_events",
			Category: "Events",
			Status:   statusInfo,
			Message:  "No events in log",
		})
		return
	}

	report.Results = append(report.Results, Result{
		Name:     "recent_events",
		Category: "Events",
		Status:   statusOK,
		Message:  fmt.Sprintf("Last event: %s", events[len(events)-1].Timestamp),
	})
}

func checkSystemResources(report *Report) {
	snap := sysinfo.Collect(".")
	addResourceResults(report, snap)
}

// addResourceResults appends per-metric resource results to the report.
// Extracted for testability with constructed Snapshot values.
func addResourceResults(report *Report, snap sysinfo.Snapshot) {
	alerts := sysinfo.Evaluate(snap)

	// Build severity lookup by resource name.
	sevMap := make(map[string]sysinfo.Severity, len(alerts))
	for _, a := range alerts {
		sevMap[a.Resource] = a.Severity
	}

	// Memory.
	if snap.Memory.Supported && snap.Memory.TotalBytes > 0 {
		pct := resourcePct(snap.Memory.UsedBytes, snap.Memory.TotalBytes)
		msg := fmt.Sprintf("Memory %d%% (%s / %s GB)",
			pct,
			sysinfo.FormatGiB(snap.Memory.UsedBytes),
			sysinfo.FormatGiB(snap.Memory.TotalBytes))
		report.Results = append(report.Results, Result{
			Name:     "resource_memory",
			Category: "Resources",
			Status:   severityToStatus(sevMap["memory"]),
			Message:  msg,
		})
	}

	// Swap (only when swap is configured).
	if snap.Memory.Supported && snap.Memory.SwapTotalBytes > 0 {
		pct := resourcePct(snap.Memory.SwapUsedBytes, snap.Memory.SwapTotalBytes)
		msg := fmt.Sprintf("Swap %d%% (%s / %s GB)",
			pct,
			sysinfo.FormatGiB(snap.Memory.SwapUsedBytes),
			sysinfo.FormatGiB(snap.Memory.SwapTotalBytes))
		report.Results = append(report.Results, Result{
			Name:     "resource_swap",
			Category: "Resources",
			Status:   severityToStatus(sevMap["swap"]),
			Message:  msg,
		})
	}

	// Disk.
	if snap.Disk.Supported && snap.Disk.TotalBytes > 0 {
		pct := resourcePct(snap.Disk.UsedBytes, snap.Disk.TotalBytes)
		msg := fmt.Sprintf("Disk %d%% (%s / %s GB)",
			pct,
			sysinfo.FormatGiB(snap.Disk.UsedBytes),
			sysinfo.FormatGiB(snap.Disk.TotalBytes))
		report.Results = append(report.Results, Result{
			Name:     "resource_disk",
			Category: "Resources",
			Status:   severityToStatus(sevMap["disk"]),
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
			Status:   severityToStatus(sevMap["load"]),
			Message:  msg,
		})
	}
}

func severityToStatus(sev sysinfo.Severity) string {
	switch sev {
	case sysinfo.SeverityWarning:
		return statusWarning
	case sysinfo.SeverityDanger:
		return statusError
	default:
		return statusOK
	}
}

func resourcePct(used, total uint64) int {
	if total == 0 {
		return 0
	}
	return int(float64(used) / float64(total) * 100)
}

func outputDoctorJSON(cmd *cobra.Command, report *Report) error {
	data, marshalErr := json.MarshalIndent(report, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	cmd.Println(string(data))
	return nil
}

func outputDoctorHuman(cmd *cobra.Command, report *Report) error {
	cmd.Println("ctx doctor")
	cmd.Println("==========")
	cmd.Println()

	// Group by category.
	categories := []string{"Structure", "Quality", "Plugin", "Hooks", "State", "Size", "Resources", "Events"}
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
			cmd.Printf("  %s %s\n", icon, r.Message)
		}
		cmd.Println()
	}

	cmd.Printf("Summary: %d warnings, %d errors\n", report.Warnings, report.Errors)
	return nil
}

func statusIcon(status string) string {
	switch status {
	case "ok":
		return "\u2713" // ✓
	case "warning":
		return "\u26a0" // ⚠
	case "error":
		return "\u2717" // ✗
	case "info":
		return "\u25cb" // ○
	default:
		return "?"
	}
}
