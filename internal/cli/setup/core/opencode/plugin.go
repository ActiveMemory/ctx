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
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errSetup "github.com/ActiveMemory/ctx/internal/err/setup"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deployPlugin writes the embedded plugin to
// .opencode/plugins/ctx.ts. OpenCode auto-loads top-level files
// under .opencode/plugins/; subdirectories are not scanned, so a
// flat single-file deployment is required. Skips if the target
// already exists.
//
// The package.json that v0.8.x and earlier shipped alongside
// index.ts is no longer embedded: the plugin uses a type-only
// import of @opencode-ai/plugin (erased at compile time) and the
// host runtime provides the plugin context, so no runtime
// dependency tree is needed.
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
	)
	target := filepath.Join(
		pluginDir, cfgHook.FileOpenCodePluginDeploy,
	)
	if _, statErr := os.Stat(target); statErr == nil {
		writeSetup.InfoOpenCodeSkipped(cmd, target)
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
	content, ok := files[cfgHook.FileIndexTs]
	if !ok {
		return errSetup.MissingEmbeddedAsset(cfgHook.FileIndexTs)
	}

	if wErr := ctxIo.SafeWriteFile(
		target, content, fs.PermFile,
	); wErr != nil {
		return errFs.FileWrite(target, wErr)
	}
	writeSetup.InfoOpenCodeCreated(cmd, target)

	return nil
}
