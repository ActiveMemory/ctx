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

	"github.com/spf13/cobra"

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
		cmd.Println("IMPORTANT: Relay this context checkpoint to the user VERBATIM before answering their question.")
		cmd.Println()
		cmd.Println(fmt.Sprintf("┌─ Context Checkpoint (prompt #%d) ────────────────", count))
		cmd.Println("│ This session is getting deep. Consider wrapping up")
		cmd.Println("│ soon. If there are unsaved learnings, decisions, or")
		cmd.Println("│ conventions, now is a good time to persist them.")
		if line := contextDirLine(); line != "" {
			cmd.Println("│ " + line)
		}
		cmd.Println("└──────────────────────────────────────────────────")
		cmd.Println()
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d CHECKPOINT", count))
		_ = notify.Send("nudge", fmt.Sprintf("check-context-size: Context Checkpoint at prompt #%d", count), sessionID)
		_ = notify.Send("relay", fmt.Sprintf("check-context-size: Context Checkpoint at prompt #%d", count), sessionID)
	} else {
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d silent", count))
	}

	return nil
}
