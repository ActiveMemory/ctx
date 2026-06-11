//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
)

// DreamEnabled reports whether the dream is turned on. Returns false
// (opt-in default) when the dream section is absent. The auto-trigger
// gate honors this before a pass.
//
// Returns:
//   - bool: true only when dream.enabled is set true in .ctxrc
func DreamEnabled() bool {
	d := RC().Dream
	return d != nil && d.Enabled
}

// DreamMode returns the configured execution mode, defaulting to
// discipline (the only mode built in v1) when unset.
//
// Returns:
//   - string: the dream mode ("discipline" by default)
func DreamMode() string {
	d := RC().Dream
	if d == nil || d.Mode == "" {
		return cfgDream.ModeDiscipline
	}
	return d.Mode
}

// DreamMax returns the per-pass ceiling on ideas/ files, defaulting to
// cfgDream.DefaultMax when unset or non-positive.
//
// Returns:
//   - int: the file ceiling for a pass
func DreamMax() int {
	d := RC().Dream
	if d == nil || d.Max <= 0 {
		return cfgDream.DefaultMax
	}
	return d.Max
}

// DreamCadence returns the configured cron schedule string, or empty
// when unset (no cron installed).
//
// Returns:
//   - string: the cron cadence, or "" when unconfigured
func DreamCadence() string {
	d := RC().Dream
	if d == nil {
		return ""
	}
	return d.Cadence
}

// DreamQuietMinutes returns the activity quiet window the trigger gate
// honors, defaulting to cfgDream.DefaultQuietMinutes when unset or
// non-positive.
//
// Returns:
//   - int: the quiet window in minutes
func DreamQuietMinutes() int {
	d := RC().Dream
	if d == nil || d.QuietMinutes <= 0 {
		return cfgDream.DefaultQuietMinutes
	}
	return d.QuietMinutes
}

// DreamModel returns the executor model override, or empty when the
// session default model should be used.
//
// Returns:
//   - string: the model override, or "" for the session default
func DreamModel() string {
	d := RC().Dream
	if d == nil {
		return ""
	}
	return d.Model
}

// DreamBudget returns the step/token budget for a pass, defaulting to
// cfgDream.DefaultBudget when unset or non-positive.
//
// Returns:
//   - int: the pass budget
func DreamBudget() int {
	d := RC().Dream
	if d == nil || d.Budget <= 0 {
		return cfgDream.DefaultBudget
	}
	return d.Budget
}

// DreamExecutor returns the configured executor command template, or
// empty when the reference claude -p invocation should be used.
//
// Returns:
//   - string: the executor command template, or "" for the default
func DreamExecutor() string {
	d := RC().Dream
	if d == nil {
		return ""
	}
	return d.Executor
}
