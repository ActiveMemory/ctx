//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_memory_drift

import (
	"os"
	"path/filepath"

	hook2 "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the check-memory-drift hook logic.
func Run(cmd *cobra.Command, stdin *os.File) error {
	if !core.Initialized() {
		return nil
	}

	input, sessionID, paused := hook2.Preamble(stdin)
	if paused {
		return nil
	}

	// Session tombstone: nudge once per session, per session ID
	tombstone := filepath.Join(
		core.StateDir(), hook.PrefixMemoryDriftThrottle+sessionID,
	)
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

	fallback := desc.Text(text.DescKeyCheckMemoryDriftContent)
	content := core.LoadMessage(
		hook.CheckMemoryDrift, hook.VariantNudge, nil, fallback,
	)
	if content == "" {
		return nil
	}

	core.EmitNudge(cmd, content,
		desc.Text(text.DescKeyCheckMemoryDriftRelayPrefix),
		desc.Text(text.DescKeyCheckMemoryDriftBoxTitle),
		hook.CheckMemoryDrift, hook.VariantNudge,
		desc.Text(text.DescKeyCheckMemoryDriftRelayMessage),
		input.SessionID, nil, tombstone)

	return nil
}
