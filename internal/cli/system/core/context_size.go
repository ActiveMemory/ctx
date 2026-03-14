//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// TokenUsageLine formats a context window usage line for display.
// Shows an icon (normal or warning), token count, percentage, and window size.
//
// Parameters:
//   - tokens: number of tokens used
//   - pct: percentage of context window used
//   - windowSize: total context window size
//
// Returns:
//   - string: formatted usage line (e.g., "⏱ Context window: ~12k tokens (~60% of 200k)")
func TokenUsageLine(tokens, pct, windowSize int) string {
	icon := assets.TextDesc(assets.TextDescKeyCheckContextSizeTokenNormal)
	suffix := ""
	if pct >= stats.ContextWindowThresholdPct {
		icon = assets.TextDesc(assets.TextDescKeyCheckContextSizeTokenLow)
		suffix = assets.TextDesc(assets.TextDescKeyCheckContextSizeRunningLowSuffix)
	}
	return fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeTokenUsage),
		icon, FormatTokenCount(tokens), pct, FormatWindowSize(windowSize), suffix)
}

// OversizeNudgeContent checks for an injection-oversize flag file and returns
// the raw nudge content if present. Deletes the flag after reading (one-shot).
//
// Returns:
//   - string: raw oversize nudge content, or empty string if no flag
func OversizeNudgeContent() string {
	baseDir := filepath.Join(rc.ContextDir(), dir.State)
	flagPath := filepath.Join(baseDir, stats.ContextSizeInjectionOversizeFlag)
	data, readErr := validation.SafeReadFile(baseDir, stats.ContextSizeInjectionOversizeFlag)
	if readErr != nil {
		return ""
	}

	tokenCount := ExtractOversizeTokens(data)
	fallback := fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeOversizeFallback), tokenCount)
	content := LoadMessage(hook.CheckContextSize, hook.VariantOversize,
		map[string]any{tpl.VarTokenCount: tokenCount}, fallback)
	if content == "" {
		_ = os.Remove(flagPath) // silenced, still consume the flag
		return ""
	}

	_ = os.Remove(flagPath) // one-shot: consumed
	return content
}

// ExtractOversizeTokens parses the token count from an injection-oversize flag file.
//
// Parameters:
//   - data: raw bytes from the flag file
//
// Returns:
//   - int: parsed token count, or 0 if not found
func ExtractOversizeTokens(data []byte) int {
	m := regex.OversizeTokens.FindSubmatch(data)
	if m == nil {
		return 0
	}
	n, parseErr := strconv.Atoi(string(m[1]))
	if parseErr != nil {
		return 0
	}
	return n
}

// EmitCheckpoint emits the standard checkpoint box with optional token usage.
//
// Parameters:
//   - cmd: Cobra command for output
//   - logFile: absolute path to the log file
//   - sessionID: session identifier
//   - count: current prompt count
//   - tokens: token usage count
//   - pct: context window usage percentage
//   - windowSize: total context window size
func EmitCheckpoint(cmd *cobra.Command, logFile, sessionID string, count, tokens, pct, windowSize int) {
	fallback := assets.TextDesc(assets.TextDescKeyCheckContextSizeCheckpointFallback)
	content := LoadMessage(hook.CheckContextSize, hook.VariantCheckpoint, nil, fallback)
	if content == "" {
		LogMessage(logFile, sessionID, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeSilencedCheckpointLog), count))
		return
	}
	// Append optional token usage and oversize nudge to content
	if tokens > 0 {
		content += token.NewlineLF + TokenUsageLine(tokens, pct, windowSize)
	}
	if extra := OversizeNudgeContent(); extra != "" {
		content += token.NewlineLF + extra
	}
	cmd.Println(NudgeBox(
		assets.TextDesc(assets.TextDescKeyCheckContextSizeRelayPrefix),
		fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeCheckpointBoxTitle), count),
		content))
	cmd.Println()
	LogMessage(logFile, sessionID, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeCheckpointLogFormat), count, tokens, pct))
	ref := notify.NewTemplateRef(hook.CheckContextSize, hook.VariantCheckpoint, nil)
	checkpointMsg := hook.CheckContextSize + ": " + fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeCheckpointRelayFormat), count)
	NudgeAndRelay(checkpointMsg, sessionID, ref)
}

// EmitWindowWarning emits an independent context window warning (>80%).
//
// Parameters:
//   - cmd: Cobra command for output
//   - logFile: absolute path to the log file
//   - sessionID: session identifier
//   - count: current prompt count
//   - tokens: token usage count
//   - pct: context window usage percentage
func EmitWindowWarning(cmd *cobra.Command, logFile, sessionID string, count, tokens, pct int) {
	fallback := fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeWindowFallback), pct, FormatTokenCount(tokens))
	content := LoadMessage(hook.CheckContextSize, hook.VariantWindow,
		map[string]any{tpl.VarPercentage: pct, tpl.VarTokenCount: FormatTokenCount(tokens)}, fallback)
	if content == "" {
		LogMessage(logFile, sessionID, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeSilencedWindowLog), count, pct))
		return
	}
	cmd.Println(NudgeBox(
		assets.TextDesc(assets.TextDescKeyCheckContextSizeRelayPrefix),
		assets.TextDesc(assets.TextDescKeyCheckContextSizeWindowBoxTitle),
		content))
	cmd.Println()
	LogMessage(logFile, sessionID, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeWindowLogFormat), count, tokens, pct))
	ref := notify.NewTemplateRef(hook.CheckContextSize, hook.VariantWindow,
		map[string]any{tpl.VarPercentage: pct, tpl.VarTokenCount: FormatTokenCount(tokens)})
	windowMsg := hook.CheckContextSize + ": " + fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeWindowRelayFormat), pct)
	NudgeAndRelay(windowMsg, sessionID, ref)
}

// EmitBillingWarning emits a one-shot warning when token usage crosses the
// billing_token_warn threshold.
//
// Parameters:
//   - cmd: Cobra command for output
//   - logFile: absolute path to the log file
//   - sessionID: session identifier
//   - count: current prompt count
//   - tokens: token usage count
//   - threshold: billing token warning threshold
func EmitBillingWarning(cmd *cobra.Command, logFile, sessionID string, count, tokens, threshold int) {
	// One-shot guard: skip if already warned this session.
	warnedFile := filepath.Join(StateDir(), stats.ContextSizeBillingWarnedPrefix+sessionID)
	if _, statErr := os.Stat(warnedFile); statErr == nil {
		return // already fired
	}

	fallback := fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeBillingFallback),
		FormatTokenCount(tokens), FormatTokenCount(threshold))
	content := LoadMessage(hook.CheckContextSize, hook.VariantBilling,
		map[string]any{tpl.VarTokenCount: FormatTokenCount(tokens), tpl.VarThreshold: FormatTokenCount(threshold)}, fallback)
	if content == "" {
		LogMessage(logFile, sessionID, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeSilencedBillingLog), count, tokens, threshold))
		TouchFile(warnedFile) // silenced counts as fired
		return
	}

	cmd.Println(NudgeBox(
		assets.TextDesc(assets.TextDescKeyCheckContextSizeBillingRelayPrefix),
		assets.TextDesc(assets.TextDescKeyCheckContextSizeBillingBoxTitle),
		content))
	cmd.Println()

	TouchFile(warnedFile) // one-shot: mark as fired
	LogMessage(logFile, sessionID, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeBillingLogFormat), count, tokens, threshold))
	ref := notify.NewTemplateRef(hook.CheckContextSize, hook.VariantBilling,
		map[string]any{tpl.VarTokenCount: FormatTokenCount(tokens), tpl.VarThreshold: FormatTokenCount(threshold)})
	billingMsg := hook.CheckContextSize + ": " + fmt.Sprintf(assets.TextDesc(assets.TextDescKeyCheckContextSizeBillingRelayFormat),
		FormatTokenCount(tokens), FormatTokenCount(threshold))
	NudgeAndRelay(billingMsg, sessionID, ref)
}
