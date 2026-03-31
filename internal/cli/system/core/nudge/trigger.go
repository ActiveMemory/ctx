//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

import (
	"github.com/ActiveMemory/ctx/internal/config/event"
	"github.com/ActiveMemory/ctx/internal/config/stats"
)

// TriggerResult is the outcome of evaluating percentage-based triggers.
//
// Fields:
//   - Event: trigger event type (silent, checkpoint, or window warning)
//   - Checkpoint: true when the 60% one-shot should fire
//   - Window: true when the 90% recurring warning should fire
type TriggerResult struct {
	Event      string
	Checkpoint bool
	Window     bool
}

// EvaluateTrigger determines which nudge (if any) should fire based on
// context window usage percentage and whether the checkpoint has already
// fired this session.
//
// Parameters:
//   - pct: current context window usage percentage (0-100)
//   - checkpointFired: true if the one-shot checkpoint already fired
//
// Returns:
//   - TriggerResult: which trigger(s) matched
func EvaluateTrigger(pct int, checkpointFired bool) TriggerResult {
	r := TriggerResult{Event: event.Silent}

	if pct >= stats.ContextWindowWarnPct {
		r.Event = event.WindowWarning
		r.Window = true
		return r
	}

	if pct >= stats.ContextCheckpointPct && !checkpointFired {
		r.Event = event.Checkpoint
		r.Checkpoint = true
		return r
	}

	return r
}
