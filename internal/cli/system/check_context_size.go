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
	"regexp"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// contextWindowThresholdPct is the percentage of context window usage that
// triggers an independent warning, regardless of prompt count.
const contextWindowThresholdPct = 80

// checkContextSizeCmd returns the "ctx system check-context-size" command.
//
// Counts prompts per session and outputs reminders at adaptive intervals,
// prompting Claude to assess remaining context capacity. Also monitors
// actual context window usage from session JSONL data and fires an
// independent warning when usage exceeds 80%.
//
// Adaptive frequency:
//
//	Prompts  1-15: silent
//	Prompts 16-30: every 5th prompt
//	Prompts   30+: every 3rd prompt
//
// Independent trigger: >80% context window fires regardless of counter.
func checkContextSizeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-context-size",
		Short: "Context size checkpoint hook",
		Long: `Counts prompts per session and emits VERBATIM relay reminders at
adaptive intervals, prompting the user to consider wrapping up.

  Prompts  1-15: silent
  Prompts 16-30: every 5th prompt
  Prompts   30+: every 3rd prompt

Also monitors actual context window token usage from session JSONL data.
Fires an independent warning when context window exceeds 80%, regardless
of prompt count.

Hook event: UserPromptSubmit
Output: VERBATIM relay (when triggered), silent otherwise
Silent when: early in session or between checkpoints`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckContextSize(cmd, os.Stdin)
		},
	}
}

func runCheckContextSize(cmd *cobra.Command, stdin *os.File) error {
	if !isInitialized() {
		return nil
	}
	input := readInput(stdin)
	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = sessionUnknown
	}

	// Pause check — this hook is the designated single emitter
	if turns := paused(sessionID); turns > 0 {
		cmd.Println(pausedMessage(turns))
		return nil
	}

	tmpDir := secureTempDir()
	counterFile := filepath.Join(tmpDir, "context-check-"+sessionID)
	logFile := filepath.Join(rc.ContextDir(), "logs", "check-context-size.log")

	// Increment counter
	count := readCounter(counterFile) + 1
	writeCounter(counterFile, count)

	// Wrap-up suppression: if the user recently ran /ctx-wrap-up,
	// suppress checkpoint nudges to avoid noise during/after the
	// wrap-up ceremony. The marker expires after 2 hours.
	if wrappedUpRecently() {
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d suppressed (wrapped up)", count))
		return nil
	}

	// Adaptive frequency (prompt counter)
	counterTriggered := false
	if count > 30 {
		counterTriggered = count%3 == 0
	} else if count > 15 {
		counterTriggered = count%5 == 0
	}

	// Read actual context window usage from session JSONL
	tokens, _ := readSessionTokenUsage(sessionID)
	windowSize := rc.ContextWindow()
	pct := 0
	if windowSize > 0 && tokens > 0 {
		pct = tokens * 100 / windowSize
	}
	windowTrigger := pct >= contextWindowThresholdPct

	switch {
	case counterTriggered:
		// Checkpoint fires (token line appended when available)
		emitCheckpoint(cmd, logFile, sessionID, count, tokens, pct, windowSize)
	case windowTrigger:
		// >80% context window, no checkpoint due — independent warning
		emitWindowWarning(cmd, logFile, sessionID, count, tokens, pct)
	default:
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d silent", count))
	}

	return nil
}

// emitCheckpoint emits the standard checkpoint box with optional token usage.
func emitCheckpoint(cmd *cobra.Command, logFile, sessionID string, count, tokens, pct, windowSize int) {
	fallback := "This session is getting deep. Consider wrapping up\n" +
		"soon. If there are unsaved learnings, decisions, or\n" +
		"conventions, now is a good time to persist them."
	content := loadMessage("check-context-size", "checkpoint", nil, fallback)
	if content == "" {
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d silenced-by-template", count))
		return
	}
	msg := fmt.Sprintf("IMPORTANT: Relay this context checkpoint to the user VERBATIM before answering their question.\n\n"+
		"┌─ Context Checkpoint (prompt #%d) ────────────────\n", count)
	msg += boxLines(content)
	if tokens > 0 {
		msg += "│ " + tokenUsageLine(tokens, pct, windowSize) + "\n"
	}
	if line := contextDirLine(); line != "" {
		msg += "│ " + line + "\n"
	}
	msg += appendOversizeNudge()
	msg += boxBottom
	cmd.Println(msg)
	cmd.Println()
	logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d CHECKPOINT tokens=%d pct=%d%%", count, tokens, pct))
	ref := notify.NewTemplateRef("check-context-size", "checkpoint", nil)
	checkpointMsg := fmt.Sprintf("check-context-size: Context Checkpoint at prompt #%d", count)
	_ = notify.Send("nudge", checkpointMsg, sessionID, ref)
	_ = notify.Send("relay", checkpointMsg, sessionID, ref)
	eventlog.Append("relay", checkpointMsg, sessionID, ref)
}

// emitWindowWarning emits an independent context window warning (>80%).
func emitWindowWarning(cmd *cobra.Command, logFile, sessionID string, count, tokens, pct int) {
	fallback := fmt.Sprintf("⚠ Context window is %d%% full (~%s tokens).\n"+
		"The session will lose older context soon. Consider wrapping up\n"+
		"or starting a fresh session with /ctx-wrap-up.", pct, formatTokenCount(tokens))
	content := loadMessage("check-context-size", "window",
		map[string]any{"Percentage": pct, "TokenCount": formatTokenCount(tokens)}, fallback)
	if content == "" {
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d window-silenced pct=%d%%", count, pct))
		return
	}
	msg := "IMPORTANT: Relay this context window warning to the user VERBATIM before answering their question.\n\n" +
		"┌─ Context Window Warning ─────────────────────────\n"
	msg += boxLines(content)
	if line := contextDirLine(); line != "" {
		msg += "│ " + line + "\n"
	}
	msg += boxBottom
	cmd.Println(msg)
	cmd.Println()
	logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d WINDOW-WARNING tokens=%d pct=%d%%", count, tokens, pct))
	ref := notify.NewTemplateRef("check-context-size", "window",
		map[string]any{"Percentage": pct, "TokenCount": formatTokenCount(tokens)})
	windowMsg := fmt.Sprintf("check-context-size: Context window at %d%%", pct)
	_ = notify.Send("nudge", windowMsg, sessionID, ref)
	_ = notify.Send("relay", windowMsg, sessionID, ref)
	eventlog.Append("relay", windowMsg, sessionID, ref)
}

// tokenUsageLine formats a context window usage line for display inside
// checkpoint boxes.
//
// Under 80%: ⏱ Context window: ~52k tokens (~26% of 200k)
// At/over 80%: ⚠ Context window: ~164k tokens (~82% of 200k) — running low
func tokenUsageLine(tokens, pct, windowSize int) string {
	icon := "⏱"
	suffix := ""
	if pct >= contextWindowThresholdPct {
		icon = "⚠"
		suffix = " — running low"
	}
	return fmt.Sprintf("%s Context window: ~%s tokens (~%d%% of %s)%s",
		icon, formatTokenCount(tokens), pct, formatWindowSize(windowSize), suffix)
}

// appendOversizeNudge checks for an injection-oversize flag file and returns
// box-formatted nudge lines if present. Deletes the flag after reading (one-shot).
// Returns empty string if no flag exists or the template is silenced.
func appendOversizeNudge() string {
	flagPath := filepath.Join(rc.ContextDir(), config.DirState, "injection-oversize")
	data, readErr := os.ReadFile(flagPath) //nolint:gosec // project-local state path
	if readErr != nil {
		return ""
	}

	tokenCount := extractOversizeTokens(data)
	fallback := fmt.Sprintf("⚠ Context injection is large (~%d tokens).\n"+
		"Run /ctx-consolidate to distill your context files.", tokenCount)
	content := loadMessage("check-context-size", "oversize",
		map[string]any{"TokenCount": tokenCount}, fallback)
	if content == "" {
		_ = os.Remove(flagPath) // silenced, still consume the flag
		return ""
	}

	_ = os.Remove(flagPath) // one-shot: consumed
	return boxLines(content)
}

// oversizeTokenRe matches "Injected:  NNNNN tokens" in the flag file.
var oversizeTokenRe = regexp.MustCompile(`Injected:\s+(\d+)\s+tokens`)

// extractOversizeTokens parses the token count from an injection-oversize flag file.
// Returns 0 if the format is unexpected.
func extractOversizeTokens(data []byte) int {
	m := oversizeTokenRe.FindSubmatch(data)
	if m == nil {
		return 0
	}
	n, err := strconv.Atoi(string(m[1]))
	if err != nil {
		return 0
	}
	return n
}
