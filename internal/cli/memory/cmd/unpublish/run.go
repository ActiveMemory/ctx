//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unpublish

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/memory/core/resolve"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgMem "github.com/ActiveMemory/ctx/internal/config/memory"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/memory"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/write/publish"
)

// Run removes the ctx-managed marker block from MEMORY.md,
// preserving all Claude-owned content outside the markers.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//
// Returns:
//   - error: on discovery, read, or write failure.
func Run(cmd *cobra.Command) error {
	_, projectRoot, err := resolve.ContextAndRoot(cmd)
	if err != nil {
		cmd.SilenceUsage = true
		return err
	}
	memoryPath, discoverErr := resolve.DiscoverSource(cmd, projectRoot)
	if discoverErr != nil {
		return discoverErr
	}
	data, readErr := resolve.ReadSource(memoryPath)
	if readErr != nil {
		return readErr
	}

	cleaned, found := memory.RemovePublished(string(data))
	if !found {
		publish.NotFound(cmd, cfgMem.Source)
		return nil
	}

	if writeErr := io.SafeWriteFile(
		memoryPath, []byte(cleaned), fs.PermFile,
	); writeErr != nil {
		return ctxErr.Write(writeErr)
	}

	publish.Unpublished(cmd, cfgMem.Source)
	return nil
}
