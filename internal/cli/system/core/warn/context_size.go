//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package warn

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	"github.com/ActiveMemory/ctx/internal/io"

	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// TokenUsageLine formats a context window usage line for display.
// Shows an icon (normal or warning), token count, percentage, and window size.
//
// Parameters:
//   - tokens: number of tokens used
//   - pct: percentage of the context window used
//   - windowSize: total context window size
//
// Returns:
//   - string: formatted usage line (e.g., "⏱ Context window: ~12k tokens (~60% of 200k)")
func TokenUsageLine(tokens, pct, windowSize int) string {
	icon := desc.Text(text.DescKeyCheckContextSizeTokenNormal)
	suffix := ""
	if pct >= stats.ContextWindowThresholdPct {
		icon = desc.Text(text.DescKeyCheckContextSizeTokenLow)
		suffix = desc.Text(text.DescKeyCheckContextSizeRunningLowSuffix)
	}
	return fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeTokenUsage),
		icon, core.FormatTokenCount(tokens), pct, core.FormatWindowSize(windowSize), suffix)
}

// OversizeNudgeContent checks for an injection-oversize flag file and returns
// the raw nudge content if present. Deletes the flag after reading (one-shot).
//
// Returns:
//   - string: raw oversize nudge content, or empty string if no flag
func OversizeNudgeContent() string {
	baseDir := filepath.Join(rc.ContextDir(), dir.State)
	flagPath := filepath.Join(baseDir, stats.ContextSizeInjectionOversizeFlag)
	data, readErr := io.SafeReadFile(baseDir, stats.ContextSizeInjectionOversizeFlag)
	if readErr != nil {
		return ""
	}

	tokenCount := ExtractOversizeTokens(data)
	fallback := fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeOversizeFallback), tokenCount)
	content := core.LoadMessage(hook.CheckContextSize, hook.VariantOversize,
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

// EmitCheckpoint builds the standard checkpoint box with optional token usage.
//
// Parameters:
//   - logFile: absolute path to the log file
//   - sessionID: session identifier
//   - count: current prompt count
//   - tokens: token usage count
//   - pct: context window usage percentage
//   - windowSize: total context window size
//
// Returns:
//   - string: formatted nudge box, or empty string if silenced
func EmitCheckpoint(logFile, sessionID string, count, tokens, pct, windowSize int) string {
	fallback := desc.Text(text.DescKeyCheckContextSizeCheckpointFallback)
	content := core.LoadMessage(hook.CheckContextSize, hook.VariantCheckpoint, nil, fallback)
	if content == "" {
		core.LogMessage(logFile, sessionID, fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeSilencedCheckpointLog), count))
		return ""
	}
	// Append optional token usage and oversize nudge to content
	if tokens > 0 {
		content += token.NewlineLF + TokenUsageLine(tokens, pct, windowSize)
	}
	if extra := OversizeNudgeContent(); extra != "" {
		content += token.NewlineLF + extra
	}
	box := core.NudgeBox(
		desc.Text(text.DescKeyCheckContextSizeRelayPrefix),
		fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeCheckpointBoxTitle), count),
		content)
	core.LogMessage(logFile, sessionID, fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeCheckpointLogFormat), count, tokens, pct))
	ref := notify.NewTemplateRef(hook.CheckContextSize, hook.VariantCheckpoint, nil)
	checkpointMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckContextSize,
		fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeCheckpointRelayFormat), count))
	core.NudgeAndRelay(checkpointMsg, sessionID, ref)
	return box
}

// EmitWindowWarning builds an independent context window warning (>80%).
//
// Parameters:
//   - logFile: absolute path to the log file
//   - sessionID: session identifier
//   - count: current prompt count
//   - tokens: token usage count
//   - pct: context window usage percentage
//
// Returns:
//   - string: formatted nudge box, or empty string if silenced
func EmitWindowWarning(logFile, sessionID string, count, tokens, pct int) string {
	fallback := fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeWindowFallback), pct, core.FormatTokenCount(tokens))
	content := core.LoadMessage(hook.CheckContextSize, hook.VariantWindow,
		map[string]any{tpl.VarPercentage: pct, tpl.VarTokenCount: core.FormatTokenCount(tokens)}, fallback)
	if content == "" {
		core.LogMessage(logFile, sessionID, fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeSilencedWindowLog), count, pct))
		return ""
	}
	box := core.NudgeBox(
		desc.Text(text.DescKeyCheckContextSizeRelayPrefix),
		desc.Text(text.DescKeyCheckContextSizeWindowBoxTitle),
		content)
	core.LogMessage(logFile, sessionID, fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeWindowLogFormat), count, tokens, pct))
	ref := notify.NewTemplateRef(hook.CheckContextSize, hook.VariantWindow,
		map[string]any{tpl.VarPercentage: pct, tpl.VarTokenCount: core.FormatTokenCount(tokens)})
	windowMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckContextSize,
		fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeWindowRelayFormat), pct))
	core.NudgeAndRelay(windowMsg, sessionID, ref)
	return box
}

// EmitBillingWarning builds a one-shot warning when token usage crosses the
// billing_token_warn threshold.
//
// Parameters:
//   - logFile: absolute path to the log file
//   - sessionID: session identifier
//   - count: current prompt count
//   - tokens: token usage count
//   - threshold: billing token warning threshold
//
// Returns:
//   - string: formatted nudge box, or empty string if silenced or already fired
func EmitBillingWarning(logFile, sessionID string, count, tokens, threshold int) string {
	// One-shot guard: skip if already warned this session.
	warnedFile := filepath.Join(core.StateDir(), stats.ContextSizeBillingWarnedPrefix+sessionID)
	if _, statErr := os.Stat(warnedFile); statErr == nil {
		return "" // already fired
	}

	fallback := fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeBillingFallback),
		core.FormatTokenCount(tokens), core.FormatTokenCount(threshold))
	content := core.LoadMessage(hook.CheckContextSize, hook.VariantBilling,
		map[string]any{tpl.VarTokenCount: core.FormatTokenCount(tokens), tpl.VarThreshold: core.FormatTokenCount(threshold)}, fallback)
	if content == "" {
		core.LogMessage(logFile, sessionID, fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeSilencedBillingLog), count, tokens, threshold))
		io.TouchFile(warnedFile) // silenced counts as fired
		return ""
	}

	box := core.NudgeBox(
		desc.Text(text.DescKeyCheckContextSizeBillingRelayPrefix),
		desc.Text(text.DescKeyCheckContextSizeBillingBoxTitle),
		content)

	io.TouchFile(warnedFile) // one-shot: mark as fired
	core.LogMessage(logFile, sessionID, fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeBillingLogFormat), count, tokens, threshold))
	ref := notify.NewTemplateRef(hook.CheckContextSize, hook.VariantBilling,
		map[string]any{tpl.VarTokenCount: core.FormatTokenCount(tokens), tpl.VarThreshold: core.FormatTokenCount(threshold)})
	billingMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckContextSize,
		fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeBillingRelayFormat),
			core.FormatTokenCount(tokens), core.FormatTokenCount(threshold)))
	core.NudgeAndRelay(billingMsg, sessionID, ref)
	return box
}
