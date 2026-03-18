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

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/load_gate"
	"github.com/ActiveMemory/ctx/internal/config/token"
	token2 "github.com/ActiveMemory/ctx/internal/context/token"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	changescore "github.com/ActiveMemory/ctx/internal/cli/changes/core"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/rc"
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
	if !core.Initialized() {
		return nil
	}

	input := core.ReadInput(stdin)
	if input.SessionID == "" {
		return nil
	}

	if core.Paused(input.SessionID) > 0 {
		return nil
	}

	tmpDir := core.StateDir()
	marker := filepath.Join(tmpDir, load_gate.PrefixCtxLoaded+input.SessionID)

	if _, statErr := os.Stat(marker); statErr == nil {
		return nil // already fired this session
	}

	// Create the marker before emitting — ensures one-shot even if
	// the agent makes multiple parallel tool calls.
	core.TouchFile(marker)

	// Auto-prune stale session state files (best-effort, silent).
	// Runs once per session at startup — fast directory scan.
	core.AutoPrune(load_gate.AutoPruneStaleDays)

	dir := rc.ContextDir()
	var content strings.Builder
	var totalTokens int
	var filesLoaded int
	var perFile []core.FileTokenEntry

	content.WriteString(
		assets.TextDesc(assets.TextDescKeyContextLoadGateHeader) +
			strings.Repeat(
				load_gate.ContextLoadSeparatorChar, load_gate.ContextLoadSeparatorWidth,
			) +
			token.NewlineLF + token.NewlineLF,
	)

	for _, f := range ctx.ReadOrder {
		if f == ctx.Glossary {
			continue
		}

		data, readErr := io.SafeReadFile(dir, f)
		if readErr != nil {
			continue // file missing — skip gracefully
		}

		switch f {
		case ctx.Task:
			// One-liner mention in footer, don't inject content
			continue

		case ctx.Decision, ctx.Learning:
			idx := core.ExtractIndex(string(data))
			if idx == "" {
				idx = assets.TextDesc(assets.TextDescKeyContextLoadGateIndexFallback)
			}
			content.WriteString(fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyContextLoadGateIndexHeader), f, idx))
			tokens := token2.EstimateTokensString(idx)
			totalTokens += tokens
			perFile = append(perFile, core.FileTokenEntry{
				Name:   f + load_gate.ContextLoadIndexSuffix,
				Tokens: tokens,
			})
			filesLoaded++

		default:
			content.WriteString(fmt.Sprintf(
				assets.TextDesc(
					assets.TextDescKeyContextLoadGateFileHeader,
				), f, string(data)))
			tokens := token2.EstimateTokens(data)
			totalTokens += tokens
			perFile = append(perFile, core.FileTokenEntry{Name: f, Tokens: tokens})
			filesLoaded++
		}
	}

	// Best-effort changes summary — never blocks injection
	if refTime, refLabel, refErr := changescore.DetectReferenceTime(""); refErr == nil {
		ctxChanges, _ := changescore.FindContextChanges(refTime)
		codeChanges, _ := changescore.SummarizeCodeChanges(refTime)
		if len(ctxChanges) > 0 || codeChanges.CommitCount > 0 {
			content.WriteString(token.NewlineLF + changescore.RenderChangesForHook(
				refLabel, ctxChanges, codeChanges))
		}
	}

	content.WriteString(
		strings.Repeat(
			load_gate.ContextLoadSeparatorChar, load_gate.ContextLoadSeparatorWidth,
		) + token.NewlineLF)
	content.WriteString(fmt.Sprintf(
		assets.TextDesc(assets.TextDescKeyContextLoadGateFooter),
		filesLoaded, totalTokens))

	core.PrintHookContext(cmd, hook.EventPreToolUse, content.String())

	// Webhook: metadata only — never send file content externally
	webhookMsg := fmt.Sprintf(
		assets.TextDesc(assets.TextDescKeyContextLoadGateWebhook),
		filesLoaded, totalTokens)
	core.Relay(webhookMsg, input.SessionID, nil)

	// Oversize nudge: write the flag for check-context-size to pick up
	core.WriteOversizeFlag(dir, totalTokens, perFile)

	return nil
}
