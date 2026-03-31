//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

// PercentMultiplier is the multiplier for converting ratios to percentages.
const PercentMultiplier = 100

// Context size hook configuration.
const (
	// ContextSizeCounterPrefix is the state file prefix for
	// per-session prompt counters.
	ContextSizeCounterPrefix = "context-check-"
	// ContextSizeLogFile is the log file name within .context/logs/.
	ContextSizeLogFile = "check-context-size.log"
	// ContextCheckpointPct is the context window usage percentage that
	// triggers a one-shot checkpoint nudge. Fires once per session to
	// encourage persisting progress (decisions, learnings, task updates).
	ContextCheckpointPct = 60
	// ContextWindowWarnPct is the context window usage percentage that
	// triggers a recurring urgent warning. Fires on every prompt at or
	// above this threshold to signal imminent context compaction.
	ContextWindowWarnPct = 90
	// ContextCheckpointNudgedPrefix is the state file prefix for the
	// one-shot checkpoint guard. Prevents the 60% nudge from repeating.
	ContextCheckpointNudgedPrefix = "checkpoint-nudged-"
	// ContextSizeBillingWarnedPrefix is the state file prefix
	// for the one-shot billing warning guard.
	ContextSizeBillingWarnedPrefix = "billing-warned-"
	// ContextSizeInjectionOversizeFlag is the state file name
	// for the injection-oversize one-shot flag.
	ContextSizeInjectionOversizeFlag = "injection-oversize"
	// JsonlPathCachePrefix is the state file prefix for cached JSONL file paths.
	JsonlPathCachePrefix = "jsonl-path-"
	// ContextSizeOversizeSepLen is the separator length for the
	// oversize flag file header.
	ContextSizeOversizeSepLen = 35

	// ViolationSpecMissing is the score for a missing Spec: trailer.
	ViolationSpecMissing = 3
	// ViolationSignoffMissing is the score for a missing Signed-off-by: trailer.
	ViolationSignoffMissing = 1
	// ViolationTaskRefMissing is the score for no task reference in the message.
	ViolationTaskRefMissing = 1
	// ViolationSingleLine is the score for a single-line commit message.
	ViolationSingleLine = 1
	// ViolationNoTasksChanged is the score for source changes
	// without TASKS.md update.
	ViolationNoTasksChanged = 1
	// ViolationThresholdNudge is the minimum score to emit a nudge.
	ViolationThresholdNudge = 2
	// ViolationThresholdWarn is the minimum score to emit a warning.
	ViolationThresholdWarn = 4
)
