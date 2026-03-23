//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_context_size

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/counter"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/hook"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/event"
	"github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
)

// Run executes the check-context-size hook logic.
//
// Reads hook input from stdin, tracks per-session prompt counts, and emits
// context checkpoint or window warning messages at adaptive intervals.
// Also fires a one-shot billing warning when token usage exceeds the
// user-configured threshold.
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
	input := hook.ReadInput(stdin)
	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = session.IDUnknown
	}

	// Pause check — this hook is the designated single emitter
	if turns := core.Paused(sessionID); turns > 0 {
		writeHook.Nudge(cmd, core.PausedMessage(turns))
		return nil
	}

	tmpDir := core.StateDir()
	counterFile := filepath.Join(tmpDir, stats.ContextSizeCounterPrefix+sessionID)
	logFile := filepath.Join(rc.ContextDir(), dir.Logs, stats.ContextSizeLogFile)

	// Increment counter
	count := counter.Read(counterFile) + 1
	counter.Write(counterFile, count)

	// Read actual context window usage from session JSONL
	info, _ := hook.ReadSessionTokenInfo(sessionID)
	tokens := info.Tokens
	windowSize := hook.EffectiveContextWindow(info.Model)
	pct := 0
	if windowSize > 0 && tokens > 0 {
		pct = tokens * stats.PercentMultiplier / windowSize
	}

	// Billing threshold: one-shot warning when tokens exceed the
	// user-configured billing_token_warn. Independent of all other
	// triggers — fires even during wrap-up suppression because cost
	// guards are never convenience nudges.
	if billingThreshold := rc.BillingTokenWarn(); billingThreshold > 0 && tokens >= billingThreshold {
		writeHook.NudgeBlock(cmd, nudge.EmitBillingWarning(logFile, sessionID, count, tokens, billingThreshold))
	}

	// Wrap-up suppression: if the user recently ran /ctx-wrap-up,
	// suppress checkpoint and window nudges to avoid noise during/after
	// the wrap-up ceremony. The marker expires after 2 hours.
	// Stats are still recorded so token usage tracking is continuous.
	if core.WrappedUpRecently() {
		core.LogMessage(
			logFile, sessionID,
			fmt.Sprintf(
				desc.Text(text.DescKeyCheckContextSizeSuppressedLogFormat), count),
		)
		hook.WriteSessionStats(sessionID, hook.SessionStats{
			Timestamp:  time.Now().Format(time.RFC3339),
			Prompt:     count,
			Tokens:     tokens,
			Pct:        pct,
			WindowSize: windowSize,
			Model:      info.Model,
			Event:      event.EventSuppressed,
		})
		return nil
	}

	// Adaptive frequency (prompt counter)
	counterTriggered := false
	if count > 30 {
		counterTriggered = count%3 == 0
	} else if count > 15 {
		counterTriggered = count%5 == 0
	}

	windowTrigger := pct >= stats.ContextWindowThresholdPct

	evt := event.EventSilent
	switch {
	case counterTriggered:
		evt = event.EventCheckpoint
		writeHook.NudgeBlock(cmd, nudge.EmitCheckpoint(logFile, sessionID, count, tokens, pct, windowSize))
	case windowTrigger:
		evt = event.EventWindowWarning
		writeHook.NudgeBlock(cmd, nudge.EmitWindowWarning(logFile, sessionID, count, tokens, pct))
	default:
		core.LogMessage(logFile, sessionID,
			fmt.Sprintf(desc.Text(
				text.DescKeyCheckContextSizeSilentLogFormat), count),
		)
	}

	hook.WriteSessionStats(sessionID, hook.SessionStats{
		Timestamp:  time.Now().Format(time.RFC3339),
		Prompt:     count,
		Tokens:     tokens,
		Pct:        pct,
		WindowSize: windowSize,
		Model:      info.Model,
		Event:      evt,
	})

	return nil
}
