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

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/eventlog"
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
	if !core.IsInitialized() {
		return nil
	}
	_, sessionID, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	tmpDir := core.StateDir()
	counterFile := filepath.Join(tmpDir, config.HeartbeatCounterPrefix+sessionID)
	mtimeFile := filepath.Join(tmpDir, config.HeartbeatMtimePrefix+sessionID)
	contextDir := rc.ContextDir()
	logFile := filepath.Join(contextDir, config.LogsDir, config.HeartbeatLogFile)

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
		config.TplVarHeartbeatPromptCount:     count,
		config.TplVarHeartbeatSessionID:       sessionID,
		config.TplVarHeartbeatContextModified: contextModified,
	}
	if tokens > 0 {
		pct := tokens * config.PercentMultiplier / window
		vars[config.TplVarHeartbeatTokens] = tokens
		vars[config.TplVarHeartbeatContextWindow] = window
		vars[config.TplVarHeartbeatUsagePct] = pct
	}
	ref := notify.NewTemplateRef(config.HookHeartbeat, config.VariantPulse, vars)

	var msg string
	if tokens > 0 {
		pct := tokens * config.PercentMultiplier / window
		msg = fmt.Sprintf(assets.TextDesc(assets.TextDescKeyHeartbeatNotifyTokens),
			count, contextModified, core.FormatTokenCount(tokens), pct)
	} else {
		msg = fmt.Sprintf(assets.TextDesc(assets.TextDescKeyHeartbeatNotifyPlain),
			count, contextModified)
	}
	_ = notify.Send(config.NotifyChannelHeartbeat, msg, sessionID, ref)
	eventlog.Append(config.NotifyChannelHeartbeat, msg, sessionID, ref)

	var logLine string
	if tokens > 0 {
		pct := tokens * config.PercentMultiplier / window
		logLine = fmt.Sprintf(assets.TextDesc(assets.TextDescKeyHeartbeatLogTokens),
			count, contextModified, core.FormatTokenCount(tokens), pct)
	} else {
		logLine = fmt.Sprintf(assets.TextDesc(assets.TextDescKeyHeartbeatLogPlain),
			count, contextModified)
	}
	core.LogMessage(logFile, sessionID, logLine)

	// No stdout — agent never sees this hook.
	return nil
}
