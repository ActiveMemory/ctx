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
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the "ctx system heartbeat" subcommand.
//
// Returns:
//   - *cobra.Command: Configured heartbeat subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "heartbeat",
		Short: "Session heartbeat webhook",
		Long: `Sends a heartbeat webhook notification on every prompt, providing
continuous session-alive visibility with metadata (prompt count, session ID,
context modification status).

Unlike other hooks, the heartbeat never produces stdout — the agent never
sees it. It only fires a webhook and writes to the event log.

Hook event: UserPromptSubmit
Output: none (webhook + event log only)
Silent when: not initialized, paused, or no webhook configured`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runHeartbeat(cmd, os.Stdin)
		},
	}
}

func runHeartbeat(_ *cobra.Command, stdin *os.File) error {
	if !core.IsInitialized() {
		return nil
	}
	input := core.ReadInput(stdin)
	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = core.SessionUnknown
	}
	if core.Paused(sessionID) > 0 {
		return nil
	}

	tmpDir := core.StateDir()
	counterFile := filepath.Join(tmpDir, "heartbeat-"+sessionID)
	mtimeFile := filepath.Join(tmpDir, "heartbeat-mtime-"+sessionID)
	contextDir := rc.ContextDir()
	logFile := filepath.Join(contextDir, "logs", "heartbeat.log")

	// Increment prompt counter.
	count := core.ReadCounter(counterFile) + 1
	core.WriteCounter(counterFile, count)

	// Detect context modification since last heartbeat.
	currentMtime := core.GetLatestContextMtime(contextDir)
	lastMtime := readMtime(mtimeFile)
	contextModified := currentMtime > lastMtime
	writeMtime(mtimeFile, currentMtime)

	// Read token usage for this session.
	info, _ := core.ReadSessionTokenInfo(sessionID)
	tokens := info.Tokens
	window := core.EffectiveContextWindow(info.Model)

	// Build and send notification.
	vars := map[string]any{
		"prompt_count":     count,
		"session_id":       sessionID,
		"context_modified": contextModified,
	}
	if tokens > 0 {
		pct := tokens * 100 / window
		vars["tokens"] = tokens
		vars["context_window"] = window
		vars["usage_pct"] = pct
	}
	ref := notify.NewTemplateRef("heartbeat", "pulse", vars)

	var msg string
	if tokens > 0 {
		pct := tokens * 100 / window
		msg = fmt.Sprintf("heartbeat: prompt #%d (context_modified=%t tokens=%s pct=%d%%)",
			count, contextModified, core.FormatTokenCount(tokens), pct)
	} else {
		msg = fmt.Sprintf("heartbeat: prompt #%d (context_modified=%t)", count, contextModified)
	}
	_ = notify.Send("heartbeat", msg, sessionID, ref)
	eventlog.Append("heartbeat", msg, sessionID, ref)

	var logLine string
	if tokens > 0 {
		pct := tokens * 100 / window
		logLine = fmt.Sprintf("prompt#%d context_modified=%t tokens=%s pct=%d%%",
			count, contextModified, core.FormatTokenCount(tokens), pct)
	} else {
		logLine = fmt.Sprintf("prompt#%d context_modified=%t", count, contextModified)
	}
	core.LogMessage(logFile, sessionID, logLine)

	// No stdout — agent never sees this hook.
	return nil
}

// readMtime reads a stored mtime value from a file. Returns 0 if the
// file does not exist or cannot be parsed.
func readMtime(path string) int64 {
	data, readErr := os.ReadFile(path) //nolint:gosec // temp file path
	if readErr != nil {
		return 0
	}
	n, parseErr := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if parseErr != nil {
		return 0
	}
	return n
}

// writeMtime writes a mtime value to a file.
func writeMtime(path string, mtime int64) {
	_ = os.WriteFile(path, []byte(strconv.FormatInt(mtime, 10)), 0o600)
}
