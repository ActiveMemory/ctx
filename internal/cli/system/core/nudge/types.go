//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

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
