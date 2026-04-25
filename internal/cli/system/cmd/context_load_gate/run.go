//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context_load_gate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/change/core/detect"
	"github.com/ActiveMemory/ctx/internal/cli/change/core/render"
	changeCore "github.com/ActiveMemory/ctx/internal/cli/change/core/scan"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/health"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/load"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/load_gate"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	ctxToken "github.com/ActiveMemory/ctx/internal/context/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Run executes the context-load-gate hook logic.
//
// Injects CONSTITUTION and distilled AGENT_PLAYBOOK_GATE into the
// agent's context window on the first tool call of each session.
// Appends a changes summary, emits a webhook notification with token
// counts, and writes an oversize flag when total injected tokens
// exceed the configured threshold. Full context files are loaded
// on-demand via CLAUDE.md instructions.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	initialized, initErr := state.Initialized()
	if initErr != nil {
		logWarn.Warn(warn.StateInitializedProbe, initErr)
		return nil
	}
	if !initialized {
		return nil
	}

	input := coreSession.ReadInput(stdin)
	if input.SessionID == "" {
		return nil
	}

	if nudge.Paused(input.SessionID) > 0 {
		return nil
	}

	tmpDir, dirErr := state.Dir()
	if dirErr != nil {
		logWarn.Warn(warn.StateDirProbe, dirErr)
		return nil
	}
	marker := filepath.Join(tmpDir, load_gate.PrefixCtxLoaded+input.SessionID)

	if _, statErr := os.Stat(marker); statErr == nil {
		return nil // already fired this session
	}

	// Create the marker before emitting - ensures one-shot even if
	// the agent makes multiple parallel tool calls.
	internalIo.TouchFile(marker)

	// Auto-prune stale session state files (best-effort, silent).
	// Runs once per session at startup - fast directory scan.
	health.AutoPrune(load_gate.AutoPruneStaleDays)

	// Unreachable under normal flow: state.Initialized() above already
	// proved ContextDir succeeds. Kept defensive so a future ContextDir
	// failure mode surfaces loudly instead of silently hanging the gate.
	dir, err := rc.ContextDir()
	if err != nil {
		logWarn.Warn(warn.ContextDirResolve, err)
		return nil
	}
	var content strings.Builder
	var totalTokens int
	var filesLoaded int
	var perFile []entity.FileTokenEntry

	content.WriteString(
		desc.Text(text.DescKeyContextLoadGateHeader) +
			strings.Repeat(
				load_gate.ContextLoadSeparatorChar, load_gate.ContextLoadSeparatorWidth,
			) +
			token.NewlineLF + token.NewlineLF,
	)

	// Inject only hard rules and distilled directives. Full context
	// files are loaded on-demand via CLAUDE.md instructions.
	gateFiles := []string{ctx.Constitution, ctx.AgentPlaybookGate}

	for _, f := range gateFiles {
		data, readErr := internalIo.SafeReadFile(dir, f)
		if readErr != nil {
			continue // file missing - skip gracefully
		}

		internalIo.SafeFprintf(&content, desc.Text(
			text.DescKeyContextLoadGateFileHeader,
		), f, string(data))
		tokens := ctxToken.Estimate(data)
		totalTokens += tokens
		perFile = append(perFile, entity.FileTokenEntry{Name: f, Tokens: tokens})
		filesLoaded++
	}

	// Best-effort changes summary - never blocks injection
	if refTime, refLabel, refErr := detect.ReferenceTime(""); refErr == nil {
		ctxChanges, _ := changeCore.FindContextChanges(refTime)
		codeChanges, _ := changeCore.SummarizeCodeChanges(refTime)
		if len(ctxChanges) > 0 || codeChanges.CommitCount > 0 {
			content.WriteString(token.NewlineLF + render.ChangesForHook(
				refLabel, ctxChanges, codeChanges))
		}
	}

	content.WriteString(
		strings.Repeat(
			load_gate.ContextLoadSeparatorChar, load_gate.ContextLoadSeparatorWidth,
		) + token.NewlineLF)
	internalIo.SafeFprintf(&content, desc.Text(text.DescKeyContextLoadGateFooter),
		filesLoaded, totalTokens)

	writeSetup.Context(
		cmd, coreSession.FormatContext(hook.EventPreToolUse, content.String()),
	)

	// Webhook: metadata only - never send file content externally.
	// Log-first: Relay writes the event log then sends the webhook;
	// if either fails, the oversize flag is NOT written; we do not
	// want an oversize nudge to fire for a gate event we never
	// recorded.
	webhookMsg := fmt.Sprintf(
		desc.Text(text.DescKeyContextLoadGateWebhook),
		filesLoaded, totalTokens)
	if relayErr := nudge.Relay(
		webhookMsg, input.SessionID, nil,
	); relayErr != nil {
		return relayErr
	}

	// Oversize nudge: write the flag for check-context-size to pick up
	load.WriteOversizeFlag(dir, totalTokens, perFile)

	return nil
}
