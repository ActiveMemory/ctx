//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/log"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	hook2 "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"

	"github.com/ActiveMemory/ctx/internal/notify"
)

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
	content := message.LoadMessage(hook.CheckContextSize, hook.VariantCheckpoint, nil, fallback)
	if content == "" {
		log.Message(logFile, sessionID, fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeSilencedCheckpointLog), count))
		return ""
	}
	// Append optional token usage and oversize nudge to content
	if tokens > 0 {
		content += token.NewlineLF + TokenUsageLine(tokens, pct, windowSize)
	}
	if extra := oversizeNudgeContent(); extra != "" {
		content += token.NewlineLF + extra
	}
	box := message.NudgeBox(
		desc.Text(text.DescKeyCheckContextSizeRelayPrefix),
		fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeCheckpointBoxTitle), count),
		content)
	log.Message(logFile, sessionID, fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeCheckpointLogFormat), count, tokens, pct))
	ref := notify.NewTemplateRef(hook.CheckContextSize, hook.VariantCheckpoint, nil)
	checkpointMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckContextSize,
		fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeCheckpointRelayFormat), count))
	NudgeAndRelay(checkpointMsg, sessionID, ref)
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
	fallback := fmt.Sprintf(
		desc.Text(text.DescKeyCheckContextSizeWindowFallback),
		pct, hook2.FormatTokenCount(tokens),
	)
	content := message.LoadMessage(hook.CheckContextSize, hook.VariantWindow,
		map[string]any{
			stats.VarPercentage: pct,
			stats.VarTokenCount: hook2.FormatTokenCount(tokens),
		}, fallback)
	if content == "" {
		log.Message(
			logFile, sessionID,
			fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeSilencedWindowLog),
				count, pct,
			),
		)
		return ""
	}
	box := message.NudgeBox(
		desc.Text(text.DescKeyCheckContextSizeRelayPrefix),
		desc.Text(text.DescKeyCheckContextSizeWindowBoxTitle),
		content)
	log.Message(
		logFile, sessionID,
		fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeWindowLogFormat),
			count, tokens, pct,
		),
	)
	ref := notify.NewTemplateRef(hook.CheckContextSize, hook.VariantWindow,
		map[string]any{
			stats.VarPercentage: pct,
			stats.VarTokenCount: hook2.FormatTokenCount(tokens),
		},
	)
	windowMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckContextSize,
		fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeWindowRelayFormat), pct))
	NudgeAndRelay(windowMsg, sessionID, ref)
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
	warnedFile := filepath.Join(
		state.StateDir(), stats.ContextSizeBillingWarnedPrefix+sessionID,
	)
	if _, statErr := os.Stat(warnedFile); statErr == nil {
		return "" // already fired
	}

	fallback := fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeBillingFallback),
		hook2.FormatTokenCount(tokens), hook2.FormatTokenCount(threshold))
	content := message.LoadMessage(hook.CheckContextSize, hook.VariantBilling,
		map[string]any{
			stats.VarTokenCount: hook2.FormatTokenCount(tokens),
			stats.VarThreshold:  hook2.FormatTokenCount(threshold),
		}, fallback)
	if content == "" {
		log.Message(
			logFile, sessionID,
			fmt.Sprintf(
				desc.Text(text.DescKeyCheckContextSizeSilencedBillingLog),
				count, tokens, threshold,
			),
		)
		io.TouchFile(warnedFile) // silenced counts as fired
		return ""
	}

	box := message.NudgeBox(
		desc.Text(text.DescKeyCheckContextSizeBillingRelayPrefix),
		desc.Text(text.DescKeyCheckContextSizeBillingBoxTitle),
		content)

	io.TouchFile(warnedFile) // one-shot: mark as fired
	log.Message(
		logFile, sessionID, fmt.Sprintf(
			desc.Text(text.DescKeyCheckContextSizeBillingLogFormat),
			count, tokens, threshold),
	)
	ref := notify.NewTemplateRef(
		hook.CheckContextSize, hook.VariantBilling,
		map[string]any{
			stats.VarTokenCount: hook2.FormatTokenCount(tokens),
			stats.VarThreshold:  hook2.FormatTokenCount(threshold),
		},
	)
	billingMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckContextSize,
		fmt.Sprintf(
			desc.Text(text.DescKeyCheckContextSizeBillingRelayFormat),
			hook2.FormatTokenCount(tokens), hook2.FormatTokenCount(threshold),
		),
	)
	NudgeAndRelay(billingMsg, sessionID, ref)
	return box
}
