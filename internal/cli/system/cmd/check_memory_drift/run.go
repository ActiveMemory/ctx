//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_memory_drift

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/memory"
)

// Run executes the check-memory-drift hook logic.
//
// Parameters:
//   - cmd: Cobra command instance
//   - stdin: Standard input for reading hook payload
//
// Returns:
//   - error: Non-nil if the drift check encounters an unrecoverable error
func Run(cmd *cobra.Command, stdin *os.File) error {
	input, sessionID, contextDir, stateDir, ok := coreCheck.FullPreamble(stdin)
	bailSilently := !ok
	if bailSilently {
		return nil
	}

	// Session tombstone: nudge once per session, per session ID
	tombstone := filepath.Join(
		stateDir, hook.PrefixMemoryDriftThrottle+sessionID,
	)
	if _, statErr := os.Stat(tombstone); statErr == nil {
		return nil
	}

	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := memory.DiscoverPath(projectRoot)
	if discoverErr != nil {
		// Auto memory not active - skip silently
		return nil
	}

	if !memory.HasDrift(contextDir, sourcePath) {
		return nil
	}

	return nudge.LoadAndEmit(cmd,
		hook.CheckMemoryDrift, hook.VariantNudge, nil,
		desc.Text(text.DescKeyCheckMemoryDriftContent),
		desc.Text(text.DescKeyCheckMemoryDriftRelayPrefix),
		desc.Text(text.DescKeyCheckMemoryDriftBoxTitle),
		desc.Text(text.DescKeyCheckMemoryDriftRelayMessage),
		input.SessionID, tombstone,
	)
}
