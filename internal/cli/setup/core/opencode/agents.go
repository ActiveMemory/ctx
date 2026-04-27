//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deployAgents creates AGENTS.md in the project root using the
// shared agent template. Skips if the file already exists.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if file read or write fails
func deployAgents(cmd *cobra.Command) error {
	target := cfgHook.FileAgentsMd

	if _, statErr := os.Stat(target); statErr == nil {
		writeSetup.InfoOpenCodeSkipped(cmd, target)
		return nil
	}

	content, readErr := agent.AgentsMd()
	if readErr != nil {
		return readErr
	}

	if wErr := io.SafeWriteFile(
		target, content, fs.PermFile,
	); wErr != nil {
		return wErr
	}
	writeSetup.InfoOpenCodeCreated(cmd, target)
	return nil
}
