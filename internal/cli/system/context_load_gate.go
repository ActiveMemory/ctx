//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// contextLoadGateCmd returns the "ctx system context-load-gate" command.
//
// Auto-injects project context into the agent's context window on the first
// tool use of a session. Uses PreToolUse hook timing so the content arrives
// at the moment of action — when the context window is fresh and attention
// is highest.
//
// Design rationale (v2): instead of telling the agent to read files (which
// the agent can evaluate and skip), the hook reads the files itself and
// injects the content directly via additionalContext. The agent never
// chooses whether to comply — the content is already present.
//
// This moves enforcement from the reasoning layer (soft instruction, subject
// to judgment) to the infrastructure layer (content injection, not subject
// to evaluation). See specs/context-load-gate-v2.md.
//
// Injection strategy per file (follows config.FileReadOrder):
//   - CONSTITUTION, CONVENTIONS, ARCHITECTURE, AGENT_PLAYBOOK: verbatim
//   - DECISIONS, LEARNINGS: index table only (INDEX:START to INDEX:END)
//   - TASKS: one-liner mention in footer (read on demand)
//   - GLOSSARY: skipped (corpus covers terminology)
//
// Webhook payloads contain metadata only (file count, token estimate),
// never file content — to avoid logging sensitive project context in
// external systems.
func contextLoadGateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "context-load-gate",
		Short: "Auto-inject project context on first tool use",
		Long: `Auto-injects project context into the agent's context window.
Fires on the first tool use per session via PreToolUse hook. Subsequent
tool calls in the same session are silent (tracked by session marker file).

Reads context files directly and injects content — no delegation to
bootstrap command, no agent compliance required.
See specs/context-load-gate-v2.md for design rationale.

Hook event: PreToolUse (.*)
Output: JSON HookResponse (additionalContext) on first tool use, silent otherwise
Silent when: marker exists for session_id, or context not initialized`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runContextLoadGate(cmd, os.Stdin)
		},
	}
}

// fileTokenEntry tracks per-file token counts during injection.
type fileTokenEntry struct {
	name   string
	tokens int
}

func runContextLoadGate(cmd *cobra.Command, stdin *os.File) error {
	if !isInitialized() {
		return nil
	}

	input := readInput(stdin)
	if input.SessionID == "" {
		return nil
	}

	if paused(input.SessionID) > 0 {
		return nil
	}

	tmpDir := secureTempDir()
	marker := filepath.Join(tmpDir, "ctx-loaded-"+input.SessionID)

	if _, err := os.Stat(marker); err == nil {
		return nil // already fired this session
	}

	// Create marker before emitting — ensures one-shot even if
	// the agent makes multiple parallel tool calls.
	touchFile(marker)

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

	content.WriteString(strings.Repeat("=", 80) + "\n")
	content.WriteString(fmt.Sprintf(
		"Context: %d files loaded (~%d tokens). "+
			"Order follows config.FileReadOrder.\n\n"+
			"TASKS.md contains the project's prioritized work items. "+
			"Read it when discussing priorities, picking up work, "+
			"or when the user asks about tasks.\n\n"+
			"For full decision or learning details, read the entry "+
			"in DECISIONS.md or LEARNINGS.md by timestamp.\n",
		filesLoaded, totalTokens))

	printHookContext(cmd, "PreToolUse", content.String())

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
// injected tokens exceed the configured threshold. The flag is consumed by
// check-context-size, which appends a nudge to the VERBATIM checkpoint.
func writeOversizeFlag(contextDir string, totalTokens int, perFile []fileTokenEntry) {
	threshold := rc.InjectionTokenWarn()
	if threshold == 0 || totalTokens <= threshold {
		return
	}

	stateDir := filepath.Join(contextDir, config.DirState)
	_ = os.MkdirAll(stateDir, 0o750)

	var flag strings.Builder
	flag.WriteString("Context injection oversize warning\n")
	flag.WriteString(strings.Repeat("=", 35) + "\n")
	flag.WriteString(fmt.Sprintf("Timestamp: %s\n", time.Now().UTC().Format(time.RFC3339)))
	flag.WriteString(fmt.Sprintf("Injected:  %d tokens (threshold: %d)\n\n", totalTokens, threshold))
	flag.WriteString("Per-file breakdown:\n")
	for _, entry := range perFile {
		flag.WriteString(fmt.Sprintf("  %-22s %5d tokens\n", entry.name, entry.tokens))
	}
	flag.WriteString("\nAction: Run /ctx-consolidate to distill context files.\n")
	flag.WriteString("Files with the most growth are the best candidates.\n")

	_ = os.WriteFile(
		filepath.Join(stateDir, "injection-oversize"),
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
