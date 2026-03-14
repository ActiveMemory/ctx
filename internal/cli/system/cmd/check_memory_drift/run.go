//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_memory_drift

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/notify"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/spf13/cobra"
)

// Run executes the check-memory-drift hook logic.
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !core.Initialized() {
		return nil
	}

	input, sessionID, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	// Session tombstone: nudge once per session, per session ID
	tombstone := filepath.Join(core.StateDir(), hook.PrefixMemoryDriftThrottle+sessionID)
	if _, statErr := os.Stat(tombstone); statErr == nil {
		return nil
	}

	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := memory.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		// Auto memory not active — skip silently
		return nil
	}

	if !memory.HasDrift(contextDir, sourcePath) {
		return nil
	}

	fallback := assets.TextDesc(assets.TextDescKeyCheckMemoryDriftContent)
	content := core.LoadMessage(hook.CheckMemoryDrift, hook.VariantNudge, nil, fallback)
	if content == "" {
		return nil
	}

	cmd.Println(core.NudgeBox(
		assets.TextDesc(assets.TextDescKeyCheckMemoryDriftRelayPrefix),
		assets.TextDesc(assets.TextDescKeyCheckMemoryDriftBoxTitle),
		content))

	ref := notify.NewTemplateRef(hook.CheckMemoryDrift, hook.VariantNudge, nil)
	core.NudgeAndRelay(
		hook.CheckMemoryDrift+": "+assets.TextDesc(assets.TextDescKeyCheckMemoryDriftRelayMessage),
		input.SessionID, ref,
	)

	core.TouchFile(tombstone)

	return nil
}
