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
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// checkContextSizeCmd returns the "ctx system check-context-size" command.
//
// Counts prompts per session and outputs reminders at adaptive intervals,
// prompting Claude to assess remaining context capacity.
//
// Adaptive frequency:
//
//	Prompts  1-15: silent
//	Prompts 16-30: every 5th prompt
//	Prompts   30+: every 3rd prompt
func checkContextSizeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-context-size",
		Short: "Context size checkpoint hook",
		Long: `Counts prompts per session and emits VERBATIM relay reminders at
adaptive intervals, prompting the user to consider wrapping up.

  Prompts  1-15: silent
  Prompts 16-30: every 5th prompt
  Prompts   30+: every 3rd prompt

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
		sessionID = "unknown"
	}

	tmpDir := secureTempDir()
	counterFile := filepath.Join(tmpDir, "context-check-"+sessionID)
	logFile := filepath.Join(rc.ContextDir(), "logs", "check-context-size.log")

	// Increment counter
	count := readCounter(counterFile) + 1
	writeCounter(counterFile, count)

	// Adaptive frequency
	shouldCheck := false
	if count > 30 {
		if count%3 == 0 {
			shouldCheck = true
		}
	} else if count > 15 {
		if count%5 == 0 {
			shouldCheck = true
		}
	}

	if shouldCheck {
		fallback := "This session is getting deep. Consider wrapping up\n" +
			"soon. If there are unsaved learnings, decisions, or\n" +
			"conventions, now is a good time to persist them."
		content := loadMessage("check-context-size", "checkpoint", nil, fallback)
		if content == "" {
			logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d silenced-by-template", count))
			return nil
		}
		msg := fmt.Sprintf("IMPORTANT: Relay this context checkpoint to the user VERBATIM before answering their question.\n\n"+
			"┌─ Context Checkpoint (prompt #%d) ────────────────\n", count)
		msg += boxLines(content)
		if line := contextDirLine(); line != "" {
			msg += "│ " + line + "\n"
		}
		msg += appendOversizeNudge()
		msg += "└──────────────────────────────────────────────────"
		cmd.Println(msg)
		cmd.Println()
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d CHECKPOINT", count))
		_ = notify.Send("nudge", fmt.Sprintf("check-context-size: Context Checkpoint at prompt #%d", count), sessionID, msg)
		_ = notify.Send("relay", fmt.Sprintf("check-context-size: Context Checkpoint at prompt #%d", count), sessionID, msg)
	} else {
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d silent", count))
	}

	return nil
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
