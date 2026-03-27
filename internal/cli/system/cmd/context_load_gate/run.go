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

	"github.com/ActiveMemory/ctx/internal/cli/change/core/detect"
	"github.com/ActiveMemory/ctx/internal/cli/change/core/render"
	changeCore "github.com/ActiveMemory/ctx/internal/cli/change/core/scan"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/health"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/load"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/entity"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/load_gate"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxToken "github.com/ActiveMemory/ctx/internal/context/token"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
)

// Run executes the context-load-gate hook logic.
//
// Injects project context files into the agent's context window on the
// first tool call of each session. Reads context files in priority order,
// extracts indexes for large files, appends a changes summary, and emits
// a webhook notification with token counts. Writes an oversize flag when
// total injected tokens exceed the configured threshold.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !state.Initialized() {
		return nil
	}

	input := coreSession.ReadInput(stdin)
	if input.SessionID == "" {
		return nil
	}

	if nudge.Paused(input.SessionID) > 0 {
		return nil
	}

	tmpDir := state.Dir()
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

	dir := rc.ContextDir()
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

	for _, f := range ctx.ReadOrder {
		if f == ctx.Glossary {
			continue
		}

		data, readErr := internalIo.SafeReadFile(dir, f)
		if readErr != nil {
			continue // file missing - skip gracefully
		}

		switch f {
		case ctx.Task:
			// One-liner mention in footer, don't inject content
			continue

		case ctx.Decision, ctx.Learning:
			idx := load.ExtractIndex(string(data))
			if idx == "" {
				idx = desc.Text(text.DescKeyContextLoadGateIndexFallback)
			}
			content.WriteString(fmt.Sprintf(
				desc.Text(text.DescKeyContextLoadGateIndexHeader), f, idx))
			tokens := ctxToken.EstimateTokensString(idx)
			totalTokens += tokens
			perFile = append(perFile, entity.FileTokenEntry{
				Name:   f + load_gate.ContextLoadIndexSuffix,
				Tokens: tokens,
			})
			filesLoaded++

		default:
			content.WriteString(fmt.Sprintf(
				desc.Text(
					text.DescKeyContextLoadGateFileHeader,
				), f, string(data)))
			tokens := ctxToken.EstimateTokens(data)
			totalTokens += tokens
			perFile = append(perFile, entity.FileTokenEntry{Name: f, Tokens: tokens})
			filesLoaded++
		}
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
	content.WriteString(fmt.Sprintf(
		desc.Text(text.DescKeyContextLoadGateFooter),
		filesLoaded, totalTokens))

	writeHook.HookContext(
		cmd, coreSession.FormatContext(hook.EventPreToolUse, content.String()),
	)

	// Webhook: metadata only - never send file content externally
	webhookMsg := fmt.Sprintf(
		desc.Text(text.DescKeyContextLoadGateWebhook),
		filesLoaded, totalTokens)
	nudge.Relay(webhookMsg, input.SessionID, nil)

	// Oversize nudge: write the flag for check-context-size to pick up
	load.WriteOversizeFlag(dir, totalTokens, perFile)

	return nil
}
