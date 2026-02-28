//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package messages provides metadata for hook message templates.
//
// The registry is a static list of all hook message entries, their
// categories (customizable vs ctx-specific), descriptions, and
// available template variables. It is the metadata layer over the
// embedded FS and changes only when hooks are added or removed.
package messages

// HookMessageInfo describes a single hook message template entry.
type HookMessageInfo struct {
	// Hook is the hook directory name (e.g., "qa-reminder").
	Hook string

	// Variant is the template file stem (e.g., "gate").
	Variant string

	// Category is "customizable" or "ctx-specific".
	Category string

	// Description is a one-line human description of this message.
	Description string

	// TemplateVars lists available Go template variables (e.g., "PromptsSinceNudge").
	TemplateVars []string
}

// CategoryCustomizable marks messages intended for project-specific customization.
const CategoryCustomizable = "customizable"

// CategoryCtxSpecific marks messages specific to ctx's own development workflow.
const CategoryCtxSpecific = "ctx-specific"

// Registry returns the static list of all hook message entries.
//
// Returns:
//   - []HookMessageInfo: All 27 entries sorted by hook then variant
func Registry() []HookMessageInfo {
	return []HookMessageInfo{
		// block-dangerous-commands: ctx-specific block responses
		{
			Hook:        "block-dangerous-commands",
			Variant:     "cp-to-bin",
			Category:    CategoryCtxSpecific,
			Description: "Block copying binaries to bin directories",
		},
		{
			Hook:        "block-dangerous-commands",
			Variant:     "install-to-local-bin",
			Category:    CategoryCtxSpecific,
			Description: "Block copying binaries to ~/.local/bin",
		},
		{
			Hook:        "block-dangerous-commands",
			Variant:     "mid-git-push",
			Category:    CategoryCtxSpecific,
			Description: "Block git push without user approval",
		},
		{
			Hook:        "block-dangerous-commands",
			Variant:     "mid-sudo",
			Category:    CategoryCtxSpecific,
			Description: "Block sudo usage",
		},

		// block-non-path-ctx: ctx-specific block responses
		{
			Hook:        "block-non-path-ctx",
			Variant:     "absolute-path",
			Category:    CategoryCtxSpecific,
			Description: "Block absolute path ctx invocation",
		},
		{
			Hook:        "block-non-path-ctx",
			Variant:     "dot-slash",
			Category:    CategoryCtxSpecific,
			Description: "Block ./ctx or ./dist/ctx invocation",
		},
		{
			Hook:        "block-non-path-ctx",
			Variant:     "go-run",
			Category:    CategoryCtxSpecific,
			Description: "Block go run ./cmd/ctx invocation",
		},

		// check-backup-age: customizable
		{
			Hook:         "check-backup-age",
			Variant:      "warning",
			Category:     CategoryCustomizable,
			Description:  "Backup staleness warning",
			TemplateVars: []string{"Warnings"},
		},

		// check-ceremonies: customizable
		{
			Hook:        "check-ceremonies",
			Variant:     "both",
			Category:    CategoryCustomizable,
			Description: "Both ceremonies missing nudge",
		},
		{
			Hook:        "check-ceremonies",
			Variant:     "remember",
			Category:    CategoryCustomizable,
			Description: "Start-of-session ceremony nudge",
		},
		{
			Hook:        "check-ceremonies",
			Variant:     "wrapup",
			Category:    CategoryCustomizable,
			Description: "End-of-session ceremony nudge",
		},

		// check-context-size: customizable
		{
			Hook:        "check-context-size",
			Variant:     "checkpoint",
			Category:    CategoryCustomizable,
			Description: "Context capacity warning",
		},
		{
			Hook:         "check-context-size",
			Variant:      "oversize",
			Category:     CategoryCustomizable,
			Description:  "Injection oversize nudge",
			TemplateVars: []string{"TokenCount"},
		},
		{
			Hook:         "check-context-size",
			Variant:      "window",
			Category:     CategoryCustomizable,
			Description:  "Context window usage warning (>80%)",
			TemplateVars: []string{"TokenCount", "Percentage"},
		},

		// check-journal: customizable
		{
			Hook:         "check-journal",
			Variant:      "both",
			Category:     CategoryCustomizable,
			Description:  "Unexported sessions and unenriched entries",
			TemplateVars: []string{"UnexportedCount", "UnenrichedCount"},
		},
		{
			Hook:         "check-journal",
			Variant:      "unenriched",
			Category:     CategoryCustomizable,
			Description:  "Unenriched journal entries",
			TemplateVars: []string{"UnenrichedCount"},
		},
		{
			Hook:         "check-journal",
			Variant:      "unexported",
			Category:     CategoryCustomizable,
			Description:  "Unexported sessions reminder",
			TemplateVars: []string{"UnexportedCount"},
		},

		// check-knowledge: customizable
		{
			Hook:         "check-knowledge",
			Variant:      "warning",
			Category:     CategoryCustomizable,
			Description:  "Knowledge file growth warning",
			TemplateVars: []string{"FileWarnings"},
		},

		// check-map-staleness: customizable
		{
			Hook:         "check-map-staleness",
			Variant:      "stale",
			Category:     CategoryCustomizable,
			Description:  "Architecture map staleness nudge",
			TemplateVars: []string{"LastRefreshDate", "ModuleCount"},
		},

		// check-persistence: customizable
		{
			Hook:         "check-persistence",
			Variant:      "nudge",
			Category:     CategoryCustomizable,
			Description:  "Context persistence nudge",
			TemplateVars: []string{"PromptsSinceNudge"},
		},

		// check-reminders: ctx-specific
		{
			Hook:         "check-reminders",
			Variant:      "reminders",
			Category:     CategoryCtxSpecific,
			Description:  "Pending reminders relay",
			TemplateVars: []string{"ReminderList"},
		},

		// check-resources: ctx-specific
		{
			Hook:         "check-resources",
			Variant:      "alert",
			Category:     CategoryCtxSpecific,
			Description:  "System resource pressure alert",
			TemplateVars: []string{"AlertMessages"},
		},

		// check-version: ctx-specific
		{
			Hook:         "check-version",
			Variant:      "key-rotation",
			Category:     CategoryCtxSpecific,
			Description:  "Encryption key rotation nudge",
			TemplateVars: []string{"KeyAgeDays"},
		},
		{
			Hook:         "check-version",
			Variant:      "mismatch",
			Category:     CategoryCtxSpecific,
			Description:  "Binary/plugin version mismatch",
			TemplateVars: []string{"BinaryVersion", "PluginVersion"},
		},

		// post-commit: customizable
		{
			Hook:        "post-commit",
			Variant:     "nudge",
			Category:    CategoryCustomizable,
			Description: "Post-commit context capture nudge",
		},

		// qa-reminder: customizable
		{
			Hook:        "qa-reminder",
			Variant:     "gate",
			Category:    CategoryCustomizable,
			Description: "Pre-commit QA gate instructions",
		},

		// specs-nudge: customizable
		{
			Hook:        "specs-nudge",
			Variant:     "nudge",
			Category:    CategoryCustomizable,
			Description: "Plan-to-specs directory nudge",
		},
	}
}

// Lookup returns the HookMessageInfo for the given hook and variant,
// or nil if not found.
//
// Parameters:
//   - hook: Hook directory name (e.g., "qa-reminder")
//   - variant: Template file stem (e.g., "gate")
//
// Returns:
//   - *HookMessageInfo: The matching entry, or nil
func Lookup(hook, variant string) *HookMessageInfo {
	for _, info := range Registry() {
		if info.Hook == hook && info.Variant == variant {
			return &info
		}
	}
	return nil
}

// Hooks returns a deduplicated list of hook names in the registry.
//
// Returns:
//   - []string: Hook names in alphabetical order
func Hooks() []string {
	seen := make(map[string]bool)
	var hooks []string
	for _, info := range Registry() {
		if !seen[info.Hook] {
			seen[info.Hook] = true
			hooks = append(hooks, info.Hook)
		}
	}
	return hooks
}

// Variants returns the variant names for a given hook.
//
// Parameters:
//   - hook: Hook directory name
//
// Returns:
//   - []string: Variant names for the hook, or nil if hook not found
func Variants(hook string) []string {
	var variants []string
	for _, info := range Registry() {
		if info.Hook == hook {
			variants = append(variants, info.Variant)
		}
	}
	return variants
}
