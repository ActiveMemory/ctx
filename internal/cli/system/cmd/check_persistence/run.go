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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/log"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	coreNudge "github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/persistence"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/time"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/nudge"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
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
	if !state.Initialized() {
		return nil
	}
	_, sessionID, paused := coreCheck.Preamble(stdin)
	if paused {
		return nil
	}

	tmpDir := state.Dir()
	stateFile := filepath.Join(tmpDir, nudge.PersistenceNudgePrefix+sessionID)
	contextDir := rc.ContextDir()
	logFile := filepath.Join(contextDir, dir.Logs, nudge.PersistenceLogFile)

	// Initialize state if needed
	ps, exists := persistence.ReadPersistenceState(stateFile)
	if !exists {
		initialMtime := time.GetLatestContextMtime(contextDir)
		ps = persistence.State{
			Count:     1,
			LastNudge: 0,
			LastMtime: initialMtime,
		}
		persistence.WritePersistenceState(stateFile, ps)
		log.Message(logFile, sessionID, fmt.Sprintf(
			desc.Text(text.DescKeyCheckPersistenceInitLogFormat), initialMtime),
		)
		return nil
	}

	ps.Count++
	currentMtime := time.GetLatestContextMtime(contextDir)

	// If context files were modified since the last check, reset the nudge counter
	if currentMtime > ps.LastMtime {
		ps.LastNudge = ps.Count
		ps.LastMtime = currentMtime
		persistence.WritePersistenceState(stateFile, ps)
		log.Message(logFile, sessionID, fmt.Sprintf(
			desc.Text(text.DescKeyCheckPersistenceModifiedLogFormat), ps.Count),
		)
		return nil
	}

	sinceNudge := ps.Count - ps.LastNudge

	// Gate persistence nudges behind minimum context window usage.
	// Below the threshold, prompt count is a poor proxy for session depth.
	pct := coreSession.LatestSessionPct(sessionID)
	if pct > 0 && pct < stats.ContextCheckpointMinPct {
		log.Message(logFile, sessionID, fmt.Sprintf(
			desc.Text(text.DescKeyCheckPersistenceSuppressedLogFormat),
			pct, stats.ContextCheckpointMinPct, ps.Count))
		persistence.WritePersistenceState(stateFile, ps)
		return nil
	}

	if persistence.PersistenceNudgeNeeded(ps.Count, sinceNudge) {
		fallback := fmt.Sprintf(
			desc.Text(text.DescKeyCheckPersistenceFallback), sinceNudge,
		)
		content := message.LoadMessage(hook.CheckPersistence, hook.VariantNudge,
			map[string]any{
				nudge.VarPromptCount:       ps.Count,
				nudge.VarPromptsSinceNudge: sinceNudge,
			}, fallback)
		if content == "" {
			log.Message(logFile, sessionID, fmt.Sprintf(
				desc.Text(text.DescKeyCheckPersistenceSilencedLogFormat), ps.Count),
			)
			persistence.WritePersistenceState(stateFile, ps)
			return nil
		}

		boxTitle := desc.Text(text.DescKeyCheckPersistenceBoxTitle)
		relayPrefix := desc.Text(text.DescKeyCheckPersistenceRelayPrefix)

		writeHook.NudgeBlock(cmd,
			message.NudgeBox(
				relayPrefix, fmt.Sprintf(
					desc.Text(text.DescKeyCheckPersistenceBoxTitleFormat),
					boxTitle, ps.Count),
				content,
			),
		)
		log.Message(logFile, sessionID, fmt.Sprintf(
			desc.Text(
				text.DescKeyCheckPersistenceNudgeLogFormat), ps.Count, sinceNudge,
		),
		)
		ref := notify.NewTemplateRef(hook.CheckPersistence, hook.VariantNudge,
			map[string]any{
				nudge.VarPromptCount:       ps.Count,
				nudge.VarPromptsSinceNudge: sinceNudge,
			},
		)
		_ = notify.Send(
			hook.NotifyChannelNudge,
			fmt.Sprintf(
				desc.Text(text.DescKeyRelayPrefixFormat),
				hook.CheckPersistence,
				fmt.Sprintf(
					desc.Text(text.DescKeyCheckPersistenceCheckpointFormat),
					ps.Count,
				),
			),
			sessionID, ref,
		)
		coreNudge.Relay(
			fmt.Sprintf(
				desc.Text(text.DescKeyRelayPrefixFormat),
				hook.CheckPersistence,
				fmt.Sprintf(
					desc.Text(text.DescKeyCheckPersistenceRelayFormat), sinceNudge,
				),
			),
			sessionID, ref,
		)
		ps.LastNudge = ps.Count
	} else {
		log.Message(
			logFile, sessionID,
			fmt.Sprintf(
				desc.Text(text.DescKeyCheckPersistenceSilentLogFormat),
				ps.Count, sinceNudge,
			),
		)
	}

	persistence.WritePersistenceState(stateFile, ps)
	return nil
}
