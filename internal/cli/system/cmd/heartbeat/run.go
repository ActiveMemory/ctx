//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package heartbeat

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/counter"
	coreHeartbeat "github.com/ActiveMemory/ctx/internal/cli/system/core/heartbeat"
	coreLog "github.com/ActiveMemory/ctx/internal/cli/system/core/log"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/time"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/heartbeat"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/log/event"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the heartbeat hook logic.
//
// Increments a per-session prompt counter, detects context file
// modifications since the last heartbeat, reads token usage, and
// emits a notification plus event log entry. Produces no stdout
// output; the agent never sees this hook.
//
// Parameters:
//   - cmd: Cobra command (unused, heartbeat produces no output)
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(_ *cobra.Command, stdin *os.File) error {
	initialized, initErr := state.Initialized()
	if initErr != nil {
		logWarn.Warn(warn.StateInitializedProbe, initErr)
		return nil
	}
	if !initialized {
		return nil
	}
	_, sessionID, paused := coreCheck.Preamble(stdin)
	if paused {
		return nil
	}

	tmpDir, dirErr := state.Dir()
	if dirErr != nil {
		logWarn.Warn(warn.StateDirProbe, dirErr)
		return nil
	}
	counterFile := filepath.Join(
		tmpDir, heartbeat.CounterPrefix+sessionID,
	)
	mtimeFile := filepath.Join(
		tmpDir, heartbeat.MtimePrefix+sessionID,
	)
	// Unreachable under normal flow: state.Initialized() above already
	// proved ContextDir succeeds. Kept defensive so a future ContextDir
	// failure surfaces instead of the heartbeat silently going dark.
	contextDir, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		logWarn.Warn(warn.ContextDirResolve, ctxErr)
		return nil
	}
	logFile := filepath.Join(contextDir, dir.Logs, heartbeat.LogFile)

	// Increment prompt counter.
	count := counter.Read(counterFile) + 1
	counter.Write(counterFile, count)

	// Detect context modification since the last heartbeat.
	currentMtime := time.GetLatestMtime(contextDir)
	lastMtime := coreHeartbeat.ReadMtime(mtimeFile)
	contextModified := currentMtime > lastMtime
	coreHeartbeat.WriteMtime(mtimeFile, currentMtime)

	// Read token usage for this session.
	info, _ := session.ReadTokenInfo(sessionID)
	tokens := info.Tokens
	window := session.EffectiveContextWindow(info.Model)

	// Build and send notification.
	vars := map[string]any{
		heartbeat.VarPromptCount:     count,
		heartbeat.VarSessionID:       sessionID,
		heartbeat.VarContextModified: contextModified,
	}
	if tokens > 0 {
		pct := tokens * stats.PercentMultiplier / window
		vars[heartbeat.VarTokens] = tokens
		vars[heartbeat.VarContextWindow] = window
		vars[heartbeat.VarUsagePct] = pct
	}
	ref := notify.NewTemplateRef(hook.Heartbeat, hook.VariantPulse, vars)

	var msg string
	if tokens > 0 {
		pct := tokens * stats.PercentMultiplier / window
		msg = fmt.Sprintf(desc.Text(text.DescKeyHeartbeatNotifyTokens),
			count, contextModified, session.FormatTokenCount(tokens), pct)
	} else {
		msg = fmt.Sprintf(desc.Text(text.DescKeyHeartbeatNotifyPlain),
			count, contextModified)
	}
	// Log-first: if the event log cannot be written, neither the
	// webhook nor the operational log line should claim the
	// heartbeat happened. See docs/security/reporting.md →
	// "Log-First Audit Trail".
	appendErr := event.Append(
		hook.NotifyChannelHeartbeat, msg, sessionID, ref,
	)
	if appendErr != nil {
		return appendErr
	}
	sendErr := notify.Send(
		hook.NotifyChannelHeartbeat, msg, sessionID, ref,
	)
	if sendErr != nil {
		return sendErr
	}

	var logLine string
	if tokens > 0 {
		pct := tokens * stats.PercentMultiplier / window
		logLine = fmt.Sprintf(desc.Text(text.DescKeyHeartbeatLogTokens),
			count, contextModified, session.FormatTokenCount(tokens), pct)
	} else {
		logLine = fmt.Sprintf(desc.Text(text.DescKeyHeartbeatLogPlain),
			count, contextModified)
	}
	coreLog.Message(logFile, sessionID, logLine)

	// No stdout - agent never sees this hook.
	return nil
}
