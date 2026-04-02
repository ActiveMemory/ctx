//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unpublish

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgMem "github.com/ActiveMemory/ctx/internal/config/memory"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/memory"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/publish"
	"github.com/ActiveMemory/ctx/internal/write/sync"
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
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	memoryPath, discoverErr := memory.DiscoverPath(projectRoot)
	if discoverErr != nil {
		sync.ErrAutoMemoryNotActive(cmd, discoverErr)
		return ctxErr.NotFound()
	}

	data, readErr := io.SafeReadFile(
		filepath.Dir(memoryPath), filepath.Base(memoryPath),
	)
	if readErr != nil {
		return ctxErr.Read(readErr)
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
