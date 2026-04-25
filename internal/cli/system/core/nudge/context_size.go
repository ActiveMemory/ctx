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
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/io"
)

// EmitCheckpoint builds the standard checkpoint box with optional token usage.
//
// Parameters:
//   - logFile: absolute path to the log file
//   - sessionID: session identifier
//   - ctxDir: absolute path to the context directory (forwarded to
//     [oversizeContent] so it does not re-resolve)
//   - count: current prompt count
//   - tokens: token usage count
//   - pct: context window usage percentage
//   - windowSize: total context window size
//
// Returns:
//   - string: formatted nudge box, or empty string if silenced
//   - error: propagated from [EmitAndRelay] so callers can honour
//     the log-first principle and skip printing the box when the
//     relay audit entry could not be written.
func EmitCheckpoint(
	logFile, sessionID, ctxDir string,
	count, tokens, pct, windowSize int,
) (string, error) {
	fallback := desc.Text(text.DescKeyCheckContextSizeCheckpointFallback)
	content := message.Load(
		hook.CheckContextSize, hook.VariantCheckpoint,
		nil, fallback,
	)
	if content == "" {
		log.Message(logFile, sessionID,
			fmt.Sprintf(
				desc.Text(text.DescKeyCheckContextSizeSilencedCheckpointLog),
				count,
			),
		)
		return "", nil
	}
	// Append optional token usage and oversize nudge to content
	if tokens > 0 {
		content += token.NewlineLF + TokenUsageLine(tokens, pct, windowSize)
	}
	extra, oversizeErr := oversizeContent(ctxDir)
	if oversizeErr != nil {
		return "", oversizeErr
	}
	if extra != "" {
		content += token.NewlineLF + extra
	}
	box := message.NudgeBox(
		desc.Text(text.DescKeyCheckContextSizeRelayPrefix),
		fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeCheckpointBoxTitle), count),
		content)
	log.Message(logFile, sessionID,
		fmt.Sprintf(
			desc.Text(text.DescKeyCheckContextSizeCheckpointLogFormat),
			count, tokens, pct,
		),
	)
	ref := entity.NewTemplateRef(
		hook.CheckContextSize, hook.VariantCheckpoint, nil,
	)
	checkpointMsg := fmt.Sprintf(
		desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckContextSize,
		fmt.Sprintf(
			desc.Text(text.DescKeyCheckContextSizeCheckpointRelayFormat),
			count,
		),
	)
	if err := EmitAndRelay(checkpointMsg, sessionID, ref); err != nil {
		return "", err
	}
	return box, nil
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
//   - error: propagated from [EmitAndRelay] so callers can honor
//     the log-first principle and skip printing the box when the
//     relay audit entry could not be written.
func EmitWindowWarning(
	logFile, sessionID string,
	count, tokens, pct int,
) (string, error) {
	fallback := fmt.Sprintf(
		desc.Text(text.DescKeyCheckContextSizeWindowFallback),
		pct, coreSession.FormatTokenCount(tokens),
	)
	content := message.Load(hook.CheckContextSize, hook.VariantWindow,
		map[string]any{
			stats.VarPercentage: pct,
			stats.VarTokenCount: coreSession.FormatTokenCount(tokens),
		}, fallback)
	if content == "" {
		log.Message(
			logFile, sessionID,
			fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeSilencedWindowLog),
				count, pct,
			),
		)
		return "", nil
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
	ref := entity.NewTemplateRef(hook.CheckContextSize, hook.VariantWindow,
		map[string]any{
			stats.VarPercentage: pct,
			stats.VarTokenCount: coreSession.FormatTokenCount(tokens),
		},
	)
	windowMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckContextSize,
		fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeWindowRelayFormat), pct))
	if err := EmitAndRelay(windowMsg, sessionID, ref); err != nil {
		return "", err
	}
	return box, nil
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
//   - error: propagated from [EmitAndRelay] so callers can honour the
//     log-first principle. The one-shot "warned" marker is touched
//     only on successful emit, so a failed relay will retry next
//     invocation rather than silently burn the one-shot chance.
func EmitBillingWarning(
	logFile, sessionID string,
	count, tokens, threshold int,
) (string, error) {
	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		return "", dirErr
	}
	// One-shot guard: skip if already warned this session.
	warnedFile := filepath.Join(
		stateDir, stats.ContextSizeBillingWarnedPrefix+sessionID,
	)
	if _, statErr := os.Stat(warnedFile); statErr == nil {
		return "", nil // already fired
	}

	fallback := fmt.Sprintf(desc.Text(text.DescKeyCheckContextSizeBillingFallback),
		coreSession.FormatTokenCount(tokens), coreSession.FormatTokenCount(threshold))
	content := message.Load(hook.CheckContextSize, hook.VariantBilling,
		map[string]any{
			stats.VarTokenCount: coreSession.FormatTokenCount(tokens),
			stats.VarThreshold:  coreSession.FormatTokenCount(threshold),
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
		return "", nil
	}

	box := message.NudgeBox(
		desc.Text(text.DescKeyCheckContextSizeBillingRelayPrefix),
		desc.Text(text.DescKeyCheckContextSizeBillingBoxTitle),
		content)

	log.Message(
		logFile, sessionID, fmt.Sprintf(
			desc.Text(text.DescKeyCheckContextSizeBillingLogFormat),
			count, tokens, threshold),
	)
	ref := entity.NewTemplateRef(
		hook.CheckContextSize, hook.VariantBilling,
		map[string]any{
			stats.VarTokenCount: coreSession.FormatTokenCount(tokens),
			stats.VarThreshold:  coreSession.FormatTokenCount(threshold),
		},
	)
	billingMsg := fmt.Sprintf(desc.Text(text.DescKeyRelayPrefixFormat),
		hook.CheckContextSize,
		fmt.Sprintf(
			desc.Text(text.DescKeyCheckContextSizeBillingRelayFormat),
			coreSession.FormatTokenCount(tokens),
			coreSession.FormatTokenCount(threshold),
		),
	)
	if err := EmitAndRelay(billingMsg, sessionID, ref); err != nil {
		return "", err
	}
	io.TouchFile(warnedFile) // one-shot: mark as fired only on success
	return box, nil
}
