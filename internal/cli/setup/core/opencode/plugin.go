//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deployPlugin creates .opencode/plugins/ctx/ with the embedded
// plugin files (index.ts and package.json). Skips if index.ts
// already exists.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func deployPlugin(cmd *cobra.Command) error {
	pluginDir := filepath.Join(
		cfgHook.DirOpenCode,
		cfgHook.DirOpenCodePlugins,
		mcpServer.Name,
	)

	indexPath := filepath.Join(
		pluginDir, cfgHook.FileIndexTs,
	)
	if _, statErr := os.Stat(indexPath); statErr == nil {
		writeSetup.InfoOpenCodeSkipped(cmd, pluginDir)
		return nil
	}

	if mkErr := ctxIo.SafeMkdirAll(
		pluginDir, fs.PermExec,
	); mkErr != nil {
		return errFs.Mkdir(pluginDir, mkErr)
	}

	files, readErr := agent.OpenCodePlugin()
	if readErr != nil {
		return readErr
	}

	for name, content := range files {
		target := filepath.Join(pluginDir, name)
		if wErr := ctxIo.SafeWriteFile(
			target, content, fs.PermFile,
		); wErr != nil {
			return errFs.FileWrite(target, wErr)
		}
		writeSetup.InfoOpenCodeCreated(cmd, target)
	}

	return nil
}
