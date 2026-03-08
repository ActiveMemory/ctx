//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_memory_drift

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write"
)

// Run executes the check-memory-drift hook logic.
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !core.IsInitialized() {
		return nil
	}

	_, sessionID, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}

	// Session tombstone: nudge once per session, per session ID
	tombstone := filepath.Join(core.StateDir(), config.MemoryDriftThrottlePrefix+sessionID)
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

	content := fmt.Sprintf(assets.TextDesc(
		assets.TextDescKeyCheckMemoryDriftContent), config.NewlineLF,
	)
	write.HookNudge(cmd, core.NudgeBox(
		assets.TextDesc(assets.TextDescKeyCheckMemoryDriftRelayPrefix),
		assets.TextDesc(assets.TextDescKeyCheckMemoryDriftBoxTitle),
		content))

	core.TouchFile(tombstone)

	return nil
}
