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

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/heartbeat"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/log"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the heartbeat hook logic.
//
// Increments a per-session prompt counter, detects context file
// modifications since the last heartbeat, reads token usage, and
// emits a notification plus event log entry. Produces no stdout
// output — the agent never sees this hook.
//
// Parameters:
//   - cmd: Cobra command (unused, heartbeat produces no output)
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(_ *cobra.Command, stdin *os.File) error {
	if !core.Initialized() {
		return nil
	}
	_, sessionID, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	tmpDir := core.StateDir()
	counterFile := filepath.Join(tmpDir, heartbeat.HeartbeatCounterPrefix+sessionID)
	mtimeFile := filepath.Join(tmpDir, heartbeat.HeartbeatMtimePrefix+sessionID)
	contextDir := rc.ContextDir()
	logFile := filepath.Join(contextDir, dir.Logs, heartbeat.HeartbeatLogFile)

	// Increment prompt counter.
	count := core.ReadCounter(counterFile) + 1
	core.WriteCounter(counterFile, count)

	// Detect context modification since the last heartbeat.
	currentMtime := core.GetLatestContextMtime(contextDir)
	lastMtime := core.ReadMtime(mtimeFile)
	contextModified := currentMtime > lastMtime
	core.WriteMtime(mtimeFile, currentMtime)

	// Read token usage for this session.
	info, _ := core.ReadSessionTokenInfo(sessionID)
	tokens := info.Tokens
	window := core.EffectiveContextWindow(info.Model)

	// Build and send notification.
	vars := map[string]any{
		tpl.VarHeartbeatPromptCount:     count,
		tpl.VarHeartbeatSessionID:       sessionID,
		tpl.VarHeartbeatContextModified: contextModified,
	}
	if tokens > 0 {
		pct := tokens * stats.PercentMultiplier / window
		vars[tpl.VarHeartbeatTokens] = tokens
		vars[tpl.VarHeartbeatContextWindow] = window
		vars[tpl.VarHeartbeatUsagePct] = pct
	}
	ref := notify.NewTemplateRef(hook.Heartbeat, hook.VariantPulse, vars)

	var msg string
	if tokens > 0 {
		pct := tokens * stats.PercentMultiplier / window
		msg = fmt.Sprintf(desc.TextDesc(text.DescKeyHeartbeatNotifyTokens),
			count, contextModified, core.FormatTokenCount(tokens), pct)
	} else {
		msg = fmt.Sprintf(desc.TextDesc(text.DescKeyHeartbeatNotifyPlain),
			count, contextModified)
	}
	_ = notify.Send(hook.NotifyChannelHeartbeat, msg, sessionID, ref)
	log.AppendEvent(hook.NotifyChannelHeartbeat, msg, sessionID, ref)

	var logLine string
	if tokens > 0 {
		pct := tokens * stats.PercentMultiplier / window
		logLine = fmt.Sprintf(desc.TextDesc(text.DescKeyHeartbeatLogTokens),
			count, contextModified, core.FormatTokenCount(tokens), pct)
	} else {
		logLine = fmt.Sprintf(desc.TextDesc(text.DescKeyHeartbeatLogPlain),
			count, contextModified)
	}
	core.LogMessage(logFile, sessionID, logLine)

	// No stdout — agent never sees this hook.
	return nil
}
