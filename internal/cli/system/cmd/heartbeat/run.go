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

	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/counter"
	coreHeartbeat "github.com/ActiveMemory/ctx/internal/cli/system/core/heartbeat"
	coreLog "github.com/ActiveMemory/ctx/internal/cli/system/core/log"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/time"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/heartbeat"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/log"
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
	if !state.Initialized() {
		return nil
	}
	_, sessionID, paused := coreCheck.Preamble(stdin)
	if paused {
		return nil
	}

	tmpDir := state.Dir()
	counterFile := filepath.Join(
		tmpDir, heartbeat.HeartbeatCounterPrefix+sessionID,
	)
	mtimeFile := filepath.Join(
		tmpDir, heartbeat.HeartbeatMtimePrefix+sessionID,
	)
	contextDir := rc.ContextDir()
	logFile := filepath.Join(contextDir, dir.Logs, heartbeat.HeartbeatLogFile)

	// Increment prompt counter.
	count := counter.Read(counterFile) + 1
	counter.Write(counterFile, count)

	// Detect context modification since the last heartbeat.
	currentMtime := time.GetLatestContextMtime(contextDir)
	lastMtime := coreHeartbeat.ReadMtime(mtimeFile)
	contextModified := currentMtime > lastMtime
	coreHeartbeat.WriteMtime(mtimeFile, currentMtime)

	// Read token usage for this session.
	info, _ := session.ReadSessionTokenInfo(sessionID)
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
	_ = notify.Send(hook.NotifyChannelHeartbeat, msg, sessionID, ref)
	log.AppendEvent(hook.NotifyChannelHeartbeat, msg, sessionID, ref)

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
