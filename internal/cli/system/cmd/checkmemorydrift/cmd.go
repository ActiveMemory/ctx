//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkmemorydrift

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the "ctx system check-memory-drift" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-memory-drift subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:    "check-memory-drift",
		Short:  "Memory drift nudge",
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCheckMemoryDrift(cmd, os.Stdin)
		},
	}
}

func runCheckMemoryDrift(cmd *cobra.Command, stdin *os.File) error {
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

	// Session tombstone — nudge once per session per session ID
	tombstone := filepath.Join(core.StateDir(), "memory-drift-nudged-"+sessionID)
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

	msg := "IMPORTANT: Relay this memory drift notice to the user VERBATIM before answering their question.\n\n" +
		"\u250c\u2500 Memory Drift \u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\n" +
		"\u2502 MEMORY.md has changed since last sync.\n" +
		"\u2502 Run: ctx memory sync\n"
	if line := core.ContextDirLine(); line != "" {
		msg += "\u2502 " + line + config.NewlineLF
	}
	msg += "\u2514\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500\u2500"
	cmd.Println(msg)

	core.TouchFile(tombstone)

	return nil
}
