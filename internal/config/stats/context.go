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
	// ContextSizeBillingWarnedPrefix is the state file prefix for the one-shot billing warning guard.
	ContextSizeBillingWarnedPrefix = "billing-warned-"
	// ContextSizeInjectionOversizeFlag is the state file name for the injection-oversize one-shot flag.
	ContextSizeInjectionOversizeFlag = "injection-oversize"
	// JsonlPathCachePrefix is the state file prefix for cached JSONL file paths.
	JsonlPathCachePrefix = "jsonl-path-"
	// ContextSizeOversizeSepLen is the separator length for the oversize flag file header.
	ContextSizeOversizeSepLen = 35
)
