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
	// ContextSizeCounterPrefix is the state file prefix for per-session prompt counters.
	ContextSizeCounterPrefix = "context-check-"
	// ContextSizeLogFile is the log file name within .context/logs/.
	ContextSizeLogFile = "check-context-size.log"
	// ContextWindowThresholdPct is the percentage of context window usage
	// that triggers an independent warning, regardless of prompt count.
	ContextWindowThresholdPct = 80
	// ContextCheckpointMinPct is the minimum context window usage percentage
	// below which counter-based checkpoint nudges are suppressed. This
	// eliminates noise on large context windows (e.g., 1M) where prompt
	// count is a poor proxy for session depth.
	ContextCheckpointMinPct = 20
	// ContextSizeBillingWarnedPrefix is the state file prefix for the one-shot billing warning guard.
	ContextSizeBillingWarnedPrefix = "billing-warned-"
	// ContextSizeInjectionOversizeFlag is the state file name for the injection-oversize one-shot flag.
	ContextSizeInjectionOversizeFlag = "injection-oversize"
	// JsonlPathCachePrefix is the state file prefix for cached JSONL file paths.
	JsonlPathCachePrefix = "jsonl-path-"
	// ContextSizeOversizeSepLen is the separator length for the oversize flag file header.
	ContextSizeOversizeSepLen = 35

	// CheckpointLateThreshold is the prompt count above which the late
	// checkpoint frequency kicks in.
	CheckpointLateThreshold = 30
	// CheckpointLateInterval is how often (in prompts) checkpoints fire
	// after the late threshold.
	CheckpointLateInterval = 3
	// CheckpointEarlyThreshold is the prompt count above which the early
	// checkpoint frequency kicks in.
	CheckpointEarlyThreshold = 15
	// CheckpointEarlyInterval is how often (in prompts) checkpoints fire
	// during the early window (between early and late thresholds).
	CheckpointEarlyInterval = 5

	// ViolationSpecMissing is the score for a missing Spec: trailer.
	ViolationSpecMissing = 3
	// ViolationSignoffMissing is the score for a missing Signed-off-by: trailer.
	ViolationSignoffMissing = 1
	// ViolationTaskRefMissing is the score for no task reference in the message.
	ViolationTaskRefMissing = 1
	// ViolationSingleLine is the score for a single-line commit message.
	ViolationSingleLine = 1
	// ViolationNoTasksChanged is the score for source changes without TASKS.md update.
	ViolationNoTasksChanged = 1
	// ViolationThresholdNudge is the minimum score to emit a nudge.
	ViolationThresholdNudge = 2
	// ViolationThresholdWarn is the minimum score to emit a warning.
	ViolationThresholdWarn = 4
)
