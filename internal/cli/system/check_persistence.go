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

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// checkPersistenceCmd returns the "ctx system check-persistence" command.
//
// Counts prompts since the last .context/ file modification and nudges
// the agent to persist learnings, decisions, or task updates.
//
// Nudge frequency:
//
//	Prompts  1-10: silent (too early)
//	Prompts 11-25: nudge once at prompt 20 since last modification
//	Prompts   25+: every 15th prompt since last modification
func checkPersistenceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-persistence",
		Short: "Persistence nudge hook",
		Long: `Tracks prompts since the last .context/ file modification and nudges
the agent to persist learnings, decisions, or task updates.

  Prompts  1-10: silent (too early)
  Prompts 11-25: nudge once at prompt 20 since last modification
  Prompts   25+: every 15th prompt since last modification

Hook event: UserPromptSubmit
Output: agent directive (when triggered), silent otherwise
Silent when: context files were recently modified`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckPersistence(cmd, os.Stdin)
		},
	}
}

// persistenceState holds the counter state for persistence nudging.
type persistenceState struct {
	Count     int
	LastNudge int
	LastMtime int64
}

func readPersistenceState(path string) (persistenceState, bool) {
	data, err := os.ReadFile(path) //nolint:gosec // temp file path
	if err != nil {
		return persistenceState{}, false
	}

	var state persistenceState
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case "count":
			n, err := strconv.Atoi(parts[1])
			if err == nil {
				state.Count = n
			}
		case "last_nudge":
			n, err := strconv.Atoi(parts[1])
			if err == nil {
				state.LastNudge = n
			}
		case "last_mtime":
			n, err := strconv.ParseInt(parts[1], 10, 64)
			if err == nil {
				state.LastMtime = n
			}
		}
	}
	return state, true
}

func writePersistenceState(path string, s persistenceState) {
	content := fmt.Sprintf("count=%d\nlast_nudge=%d\nlast_mtime=%d\n",
		s.Count, s.LastNudge, s.LastMtime)
	_ = os.WriteFile(path, []byte(content), 0o600)
}

// getLatestContextMtime returns the most recent mtime of any .context/*.md file.
func getLatestContextMtime(contextDir string) int64 {
	entries, err := os.ReadDir(contextDir)
	if err != nil {
		return 0
	}

	var latest int64
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		mtime := info.ModTime().Unix()
		if mtime > latest {
			latest = mtime
		}
	}
	return latest
}

func runCheckPersistence(cmd *cobra.Command, stdin *os.File) error {
	if !isInitialized() {
		return nil
	}
	input := readInput(stdin)
	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = "unknown"
	}

	tmpDir := secureTempDir()
	stateFile := filepath.Join(tmpDir, "persistence-nudge-"+sessionID)
	contextDir := rc.ContextDir()
	logFile := filepath.Join(contextDir, "logs", "check-persistence.log")

	// Initialize state if needed
	state, exists := readPersistenceState(stateFile)
	if !exists {
		initialMtime := getLatestContextMtime(contextDir)
		state = persistenceState{
			Count:     1,
			LastNudge: 0,
			LastMtime: initialMtime,
		}
		writePersistenceState(stateFile, state)
		logMessage(logFile, sessionID, fmt.Sprintf("init count=1 mtime=%d", initialMtime))
		return nil
	}

	state.Count++
	currentMtime := getLatestContextMtime(contextDir)

	// If context files were modified since last check, reset the nudge counter
	if currentMtime > state.LastMtime {
		state.LastNudge = state.Count
		state.LastMtime = currentMtime
		writePersistenceState(stateFile, state)
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d context-modified, reset nudge counter", state.Count))
		return nil
	}

	sinceNudge := state.Count - state.LastNudge

	// Determine if we should nudge
	shouldNudge := false
	if state.Count >= 11 && state.Count <= 25 && sinceNudge >= 20 {
		shouldNudge = true
	} else if state.Count > 25 && sinceNudge >= 15 {
		shouldNudge = true
	}

	if shouldNudge {
		fallback := fmt.Sprintf("No context files updated in %d+ prompts.\n", sinceNudge) +
			"Have you discovered learnings, made decisions,\n" +
			"established conventions, or completed tasks\n" +
			"worth persisting?\n" +
			"\n" +
			"Run /ctx-wrap-up to capture session context."
		content := loadMessage("check-persistence", "nudge",
			map[string]any{
				"PromptCount":       state.Count,
				"PromptsSinceNudge": sinceNudge,
			}, fallback)
		if content == "" {
			logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d silenced-by-template", state.Count))
			writePersistenceState(stateFile, state)
			return nil
		}
		msg := fmt.Sprintf("IMPORTANT: Relay this persistence checkpoint to the user VERBATIM before answering their question.\n\n"+
			"┌─ Persistence Checkpoint (prompt #%d) ───────────\n", state.Count)
		msg += boxLines(content)
		if line := contextDirLine(); line != "" {
			msg += "│ " + line + "\n"
		}
		msg += config.NudgeBoxBottom
		cmd.Println(msg)
		cmd.Println()
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d NUDGE since_nudge=%d", state.Count, sinceNudge))
		_ = notify.Send("nudge", fmt.Sprintf("check-persistence: Persistence Checkpoint at prompt #%d", state.Count), sessionID, msg)
		_ = notify.Send("relay", fmt.Sprintf("check-persistence: No context updated in %d+ prompts", sinceNudge), sessionID, msg)
		state.LastNudge = state.Count
	} else {
		logMessage(logFile, sessionID, fmt.Sprintf("prompt#%d silent since_nudge=%d", state.Count, sinceNudge))
	}

	writePersistenceState(stateFile, state)
	return nil
}
