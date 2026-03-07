//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package contextloadgate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	changescore "github.com/ActiveMemory/ctx/internal/cli/changes/core"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// fileTokenEntry tracks per-file token counts during injection.
type fileTokenEntry struct {
	name   string
	tokens int
}

func runContextLoadGate(cmd *cobra.Command, stdin *os.File) error {
	if !core.IsInitialized() {
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
	marker := filepath.Join(tmpDir, "ctx-loaded-"+input.SessionID)

	if _, statErr := os.Stat(marker); statErr == nil {
		return nil // already fired this session
	}

	// Create marker before emitting — ensures one-shot even if
	// the agent makes multiple parallel tool calls.
	core.TouchFile(marker)

	// Auto-prune stale session state files (best-effort, silent).
	// Runs once per session at startup — fast directory scan.
	core.AutoPrune(7)

	dir := rc.ContextDir()
	var content strings.Builder
	var totalTokens int
	var filesLoaded int
	var perFile []fileTokenEntry

	content.WriteString(
		"PROJECT CONTEXT (auto-loaded by system hook" +
			" — already in your context window)\n" +
			strings.Repeat("=", 80) + "\n\n")

	for _, f := range config.FileReadOrder {
		if f == config.FileGlossary {
			continue
		}

		path := filepath.Join(dir, f)
		data, readErr := os.ReadFile(path) //#nosec G304 — path is within .context/
		if readErr != nil {
			continue // file missing — skip gracefully
		}

		switch f {
		case config.FileTask:
			// One-liner mention in footer, don't inject content
			continue

		case config.FileDecision, config.FileLearning:
			idx := extractIndex(string(data))
			if idx == "" {
				idx = "(no index entries)"
			}
			content.WriteString(fmt.Sprintf(
				"--- %s (index — read full entries by date "+
					"when relevant) ---\n%s\n\n", f, idx))
			tokens := context.EstimateTokensString(idx)
			totalTokens += tokens
			perFile = append(perFile, fileTokenEntry{name: f + " (idx)", tokens: tokens})
			filesLoaded++

		default:
			content.WriteString(fmt.Sprintf(
				"--- %s ---\n%s\n\n", f, string(data)))
			tokens := context.EstimateTokens(data)
			totalTokens += tokens
			perFile = append(perFile, fileTokenEntry{name: f, tokens: tokens})
			filesLoaded++
		}
	}

	// Best-effort changes summary — never blocks injection
	if refTime, refLabel, refErr := changescore.DetectReferenceTime(""); refErr == nil {
		ctxChanges, _ := changescore.FindContextChanges(refTime)
		codeChanges, _ := changescore.SummarizeCodeChanges(refTime)
		if len(ctxChanges) > 0 || codeChanges.CommitCount > 0 {
			content.WriteString(config.NewlineLF + changescore.RenderChangesForHook(
				refLabel, ctxChanges, codeChanges))
		}
	}

	content.WriteString(strings.Repeat("=", 80) + config.NewlineLF)
	content.WriteString(fmt.Sprintf(
		"Context: %d files loaded (~%d tokens). "+
			"Order follows config.FileReadOrder.\n\n"+
			"TASKS.md contains the project's prioritized work items. "+
			"Read it when discussing priorities, picking up work, "+
			"or when the user asks about tasks.\n\n"+
			"For full decision or learning details, read the entry "+
			"in DECISIONS.md or LEARNINGS.md by timestamp.\n",
		filesLoaded, totalTokens))

	core.PrintHookContext(cmd, "PreToolUse", content.String())

	// Webhook: metadata only — never send file content externally
	webhookMsg := fmt.Sprintf(
		"context-load-gate: injected %d files (~%d tokens)",
		filesLoaded, totalTokens)
	_ = notify.Send("relay", webhookMsg, input.SessionID, nil)
	eventlog.Append("relay", webhookMsg, input.SessionID, nil)

	// Oversize nudge: write flag for check-context-size to pick up
	writeOversizeFlag(dir, totalTokens, perFile)

	return nil
}

// writeOversizeFlag writes an injection-oversize flag file when the total
// injected tokens exceed the configured threshold.
func writeOversizeFlag(contextDir string, totalTokens int, perFile []fileTokenEntry) {
	threshold := rc.InjectionTokenWarn()
	if threshold == 0 || totalTokens <= threshold {
		return
	}

	sd := filepath.Join(contextDir, config.DirState)
	_ = os.MkdirAll(sd, 0o750)

	var flag strings.Builder
	flag.WriteString("Context injection oversize warning\n")
	flag.WriteString(strings.Repeat("=", 35) + config.NewlineLF)
	flag.WriteString(fmt.Sprintf("Timestamp: %s\n", time.Now().UTC().Format(time.RFC3339)))
	flag.WriteString(fmt.Sprintf("Injected:  %d tokens (threshold: %d)\n\n", totalTokens, threshold))
	flag.WriteString("Per-file breakdown:\n")
	for _, entry := range perFile {
		flag.WriteString(fmt.Sprintf("  %-22s %5d tokens\n", entry.name, entry.tokens))
	}
	flag.WriteString("\nAction: Run /ctx-consolidate to distill context files.\n")
	flag.WriteString("Files with the most growth are the best candidates.\n")

	_ = os.WriteFile(
		filepath.Join(sd, "injection-oversize"),
		[]byte(flag.String()), 0o600)
}

// extractIndex returns the content between INDEX:START and INDEX:END
// markers, or empty string if markers are not found.
func extractIndex(content string) string {
	start := strings.Index(content, config.IndexStart)
	end := strings.Index(content, config.IndexEnd)
	if start < 0 || end < 0 || end <= start {
		return ""
	}
	startPos := start + len(config.IndexStart)
	return strings.TrimSpace(content[startPos:end])
}
