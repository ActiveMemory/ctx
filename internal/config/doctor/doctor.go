//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package doctor

// Doctor check name constants: used as Result.Name values.
const (
	// CheckContextInit identifies the context initialization check.
	CheckContextInit = "context_initialized"
	// CheckRequiredFiles identifies the required files check.
	CheckRequiredFiles = "required_files"
	// CheckCtxrcValidation identifies the .ctxrc validation check.
	CheckCtxrcValidation = "ctxrc_validation"
	// CheckDrift identifies the drift detection check.
	CheckDrift = "drift"
	// CheckPluginInstalled identifies the plugin installation check.
	CheckPluginInstalled = "plugin_installed"
	// CheckPluginEnabledGlobal identifies the global plugin enablement check.
	CheckPluginEnabledGlobal = "plugin_enabled_global"
	// CheckPluginEnabledLocal identifies the local plugin enablement check.
	CheckPluginEnabledLocal = "plugin_enabled_local"
	// CheckPluginEnabled identifies the plugin enablement check
	// (when neither scope is active).
	CheckPluginEnabled = "plugin_enabled"
	// CheckEventLogging identifies the event logging check.
	CheckEventLogging = "event_logging"
	// CheckWebhook identifies the webhook configuration check.
	CheckWebhook = "webhook"
	// CheckReminders identifies the pending reminders check.
	CheckReminders = "reminders"
	// CheckTaskCompletion identifies the task completion check.
	CheckTaskCompletion = "task_completion"
	// CheckContextSize identifies the context token size check.
	CheckContextSize = "context_size"
	// CheckContextFilePrefix is the prefix for per-file context size results.
	CheckContextFilePrefix = "context_file_"
	// CheckRecentEvents identifies the recent event log check.
	CheckRecentEvents = "recent_events"
	// CheckResourceMemory identifies the memory resource check.
	CheckResourceMemory = "resource_memory"
	// CheckResourceSwap identifies the swap resource check.
	CheckResourceSwap = "resource_swap"
	// CheckResourceDisk identifies the disk resource check.
	CheckResourceDisk = "resource_disk"
	// CheckResourceLoad identifies the load resource check.
	CheckResourceLoad = "resource_load"
)

// Doctor category constants: used as Result.Category values.
const (
	// CategoryStructure groups context directory and file checks.
	CategoryStructure = "Structure"
	// CategoryQuality groups drift and content quality checks.
	CategoryQuality = "Quality"
	// CategoryPlugin groups plugin installation and enablement checks.
	CategoryPlugin = "Plugin"
	// CategoryHooks groups hook configuration checks.
	CategoryHooks = "Hooks"
	// CategoryState groups runtime state checks.
	CategoryState = "State"
	// CategorySize groups token size and budget checks.
	CategorySize = "Size"
	// CategoryResources groups system resource checks.
	CategoryResources = "Resources"
	// CategoryEvents groups event log checks.
	CategoryEvents = "Events"
)

// Thresholds for doctor health checks.
const (
	// TaskCompletionWarnPct is the completed-task ratio that triggers a warning.
	TaskCompletionWarnPct = 80
	// TaskCompletionMinCount is the minimum completed tasks before the ratio check applies.
	TaskCompletionMinCount = 5
	// ContextSizeWarnPct is the percentage of context window usage that triggers a warning.
	ContextSizeWarnPct = 20
)
