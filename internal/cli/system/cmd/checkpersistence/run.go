//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkpersistence

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// persistenceState holds the counter state for persistence nudging.
type persistenceState struct {
	Count     int
	LastNudge int
	LastMtime int64
}

func readPersistenceState(path string) (persistenceState, bool) {
	data, readErr := os.ReadFile(path) //nolint:gosec // state dir path
	if readErr != nil {
		return persistenceState{}, false
	}

	var ps persistenceState
	for _, line := range strings.Split(strings.TrimSpace(string(data)), config.NewlineLF) {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case "count":
			n, parseErr := strconv.Atoi(parts[1])
			if parseErr == nil {
				ps.Count = n
			}
		case "last_nudge":
			n, parseErr := strconv.Atoi(parts[1])
			if parseErr == nil {
				ps.LastNudge = n
			}
		case "last_mtime":
			n, parseErr := strconv.ParseInt(parts[1], 10, 64)
			if parseErr == nil {
				ps.LastMtime = n
			}
		}
	}
	return ps, true
}

func writePersistenceState(path string, s persistenceState) {
	content := fmt.Sprintf("count=%d\nlast_nudge=%d\nlast_mtime=%d\n",
		s.Count, s.LastNudge, s.LastMtime)
	_ = os.WriteFile(path, []byte(content), 0o600)
}

func runCheckPersistence(cmd *cobra.Command, stdin *os.File) error {
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
	stateFile := filepath.Join(tmpDir, "persistence-nudge-"+sessionID)
	contextDir := rc.ContextDir()
	logFile := filepath.Join(contextDir, "logs", "check-persistence.log")

	// Initialize state if needed
	ps, exists := readPersistenceState(stateFile)
	if !exists {
		initialMtime := core.GetLatestContextMtime(contextDir)
		ps = persistenceState{
			Count:     1,
			LastNudge: 0,
			LastMtime: initialMtime,
		}
		writePersistenceState(stateFile, ps)
		core.LogMessage(logFile, sessionID, fmt.Sprintf("init count=1 mtime=%d", initialMtime))
		return nil
	}

	ps.Count++
	currentMtime := core.GetLatestContextMtime(contextDir)

	// If context files were modified since last check, reset the nudge counter
	if currentMtime > ps.LastMtime {
		ps.LastNudge = ps.Count
		ps.LastMtime = currentMtime
		writePersistenceState(stateFile, ps)
		core.LogMessage(logFile, sessionID, fmt.Sprintf("prompt#%d context-modified, reset nudge counter", ps.Count))
		return nil
	}

	sinceNudge := ps.Count - ps.LastNudge

	// Determine if we should nudge
	shouldNudge := false
	if ps.Count >= 11 && ps.Count <= 25 && sinceNudge >= 20 {
		shouldNudge = true
	} else if ps.Count > 25 && sinceNudge >= 15 {
		shouldNudge = true
	}

	if shouldNudge {
		fallback := fmt.Sprintf("No context files updated in %d+ prompts.\n", sinceNudge) +
			"Have you discovered learnings, made decisions,\n" +
			"established conventions, or completed tasks\n" +
			"worth persisting?\n" +
			config.NewlineLF +
			"Run /ctx-wrap-up to capture session context."
		content := core.LoadMessage("check-persistence", "nudge",
			map[string]any{
				"PromptCount":       ps.Count,
				"PromptsSinceNudge": sinceNudge,
			}, fallback)
		if content == "" {
			core.LogMessage(logFile, sessionID, fmt.Sprintf("prompt#%d silenced-by-template", ps.Count))
			writePersistenceState(stateFile, ps)
			return nil
		}
		msg := fmt.Sprintf("IMPORTANT: Relay this persistence checkpoint to the user VERBATIM before answering their question.\n\n"+
			"┌─ Persistence Checkpoint (prompt #%d) ───────────\n", ps.Count)
		msg += core.BoxLines(content)
		if line := core.ContextDirLine(); line != "" {
			msg += "│ " + line + config.NewlineLF
		}
		msg += config.NudgeBoxBottom
		cmd.Println(msg)
		cmd.Println()
		core.LogMessage(logFile, sessionID, fmt.Sprintf("prompt#%d NUDGE since_nudge=%d", ps.Count, sinceNudge))
		ref := notify.NewTemplateRef("check-persistence", "nudge",
			map[string]any{"PromptCount": ps.Count, "PromptsSinceNudge": sinceNudge})
		_ = notify.Send("nudge", fmt.Sprintf("check-persistence: Persistence Checkpoint at prompt #%d", ps.Count), sessionID, ref)
		_ = notify.Send("relay", fmt.Sprintf("check-persistence: No context updated in %d+ prompts", sinceNudge), sessionID, ref)
		eventlog.Append("relay", fmt.Sprintf("check-persistence: No context updated in %d+ prompts", sinceNudge), sessionID, ref)
		ps.LastNudge = ps.Count
	} else {
		core.LogMessage(logFile, sessionID, fmt.Sprintf("prompt#%d silent since_nudge=%d", ps.Count, sinceNudge))
	}

	writePersistenceState(stateFile, ps)
	return nil
}
