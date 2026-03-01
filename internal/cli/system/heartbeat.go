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
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// heartbeatCmd returns the "ctx system heartbeat" command.
//
// Sends a heartbeat webhook notification on every prompt, providing
// continuous session-alive visibility with metadata. Unlike other hooks,
// the heartbeat never produces stdout — it only fires a webhook and
// logs locally.
//
// Hook event: UserPromptSubmit
// Output: none (webhook + event log only)
func heartbeatCmd() *cobra.Command {
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
	if !isInitialized() {
		return nil
	}
	input := readInput(stdin)
	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = sessionUnknown
	}
	if paused(sessionID) > 0 {
		return nil
	}

	tmpDir := secureTempDir()
	counterFile := filepath.Join(tmpDir, "heartbeat-"+sessionID)
	mtimeFile := filepath.Join(tmpDir, "heartbeat-mtime-"+sessionID)
	contextDir := rc.ContextDir()
	logFile := filepath.Join(contextDir, "logs", "heartbeat.log")

	// Increment prompt counter.
	count := readCounter(counterFile) + 1
	writeCounter(counterFile, count)

	// Detect context modification since last heartbeat.
	currentMtime := getLatestContextMtime(contextDir)
	lastMtime := readMtime(mtimeFile)
	contextModified := currentMtime > lastMtime
	writeMtime(mtimeFile, currentMtime)

	// Build and send notification.
	ref := notify.NewTemplateRef("heartbeat", "pulse", map[string]any{
		"prompt_count":     count,
		"session_id":       sessionID,
		"context_modified": contextModified,
	})
	msg := fmt.Sprintf("heartbeat: prompt #%d (context_modified=%t)", count, contextModified)
	_ = notify.Send("heartbeat", msg, sessionID, ref)
	eventlog.Append("heartbeat", msg, sessionID, ref)

	logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d context_modified=%t", count, contextModified))

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

