//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_persistence

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/nudge"
	"github.com/ActiveMemory/ctx/internal/config/tpl"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the check-persistence hook logic.
//
// Tracks how many prompts have passed without any context file updates
// and emits a persistence nudge when the threshold is reached. State is
// stored per-session in the .context/state/ directory.
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
	_, sessionID, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	tmpDir := core.StateDir()
	stateFile := filepath.Join(tmpDir, nudge.PersistenceNudgePrefix+sessionID)
	contextDir := rc.ContextDir()
	logFile := filepath.Join(contextDir, dir.Logs, nudge.PersistenceLogFile)

	// Initialize state if needed
	ps, exists := core.ReadPersistenceState(stateFile)
	if !exists {
		initialMtime := core.GetLatestContextMtime(contextDir)
		ps = core.PersistenceState{
			Count:     1,
			LastNudge: 0,
			LastMtime: initialMtime,
		}
		core.WritePersistenceState(stateFile, ps)
		core.LogMessage(logFile, sessionID, fmt.Sprintf(desc.TextDesc(text.DescKeyCheckPersistenceInitLogFormat), initialMtime))
		return nil
	}

	ps.Count++
	currentMtime := core.GetLatestContextMtime(contextDir)

	// If context files were modified since last check, reset the nudge counter
	if currentMtime > ps.LastMtime {
		ps.LastNudge = ps.Count
		ps.LastMtime = currentMtime
		core.WritePersistenceState(stateFile, ps)
		core.LogMessage(logFile, sessionID, fmt.Sprintf(desc.TextDesc(text.DescKeyCheckPersistenceModifiedLogFormat), ps.Count))
		return nil
	}

	sinceNudge := ps.Count - ps.LastNudge

	if core.PersistenceNudgeNeeded(ps.Count, sinceNudge) {
		fallback := fmt.Sprintf(desc.TextDesc(text.DescKeyCheckPersistenceFallback), sinceNudge)
		content := core.LoadMessage(hook.CheckPersistence, hook.VariantNudge,
			map[string]any{
				tpl.VarPromptCount:       ps.Count,
				tpl.VarPromptsSinceNudge: sinceNudge,
			}, fallback)
		if content == "" {
			core.LogMessage(logFile, sessionID, fmt.Sprintf(desc.TextDesc(text.DescKeyCheckPersistenceSilencedLogFormat), ps.Count))
			core.WritePersistenceState(stateFile, ps)
			return nil
		}

		boxTitle := desc.TextDesc(text.DescKeyCheckPersistenceBoxTitle)
		relayPrefix := desc.TextDesc(text.DescKeyCheckPersistenceRelayPrefix)

		cmd.Println(core.NudgeBox(relayPrefix, fmt.Sprintf(desc.TextDesc(text.DescKeyCheckPersistenceBoxTitleFormat), boxTitle, ps.Count), content))
		cmd.Println()
		core.LogMessage(logFile, sessionID, fmt.Sprintf("prompt#%d NUDGE since_nudge=%d", ps.Count, sinceNudge))
		ref := notify.NewTemplateRef(hook.CheckPersistence, hook.VariantNudge,
			map[string]any{tpl.VarPromptCount: ps.Count, tpl.VarPromptsSinceNudge: sinceNudge})
		_ = notify.Send(hook.NotifyChannelNudge, hook.CheckPersistence+": "+fmt.Sprintf(desc.TextDesc(text.DescKeyCheckPersistenceCheckpointFormat), ps.Count), sessionID, ref)
		core.Relay(hook.CheckPersistence+": "+fmt.Sprintf(desc.TextDesc(text.DescKeyCheckPersistenceRelayFormat), sinceNudge), sessionID, ref)
		ps.LastNudge = ps.Count
	} else {
		core.LogMessage(logFile, sessionID, fmt.Sprintf(desc.TextDesc(text.DescKeyCheckPersistenceSilentLogFormat), ps.Count, sinceNudge))
	}

	core.WritePersistenceState(stateFile, ps)
	return nil
}
