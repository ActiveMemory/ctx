//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unpublish

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	memory2 "github.com/ActiveMemory/ctx/internal/config/memory"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/memory"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/publish"
	"github.com/ActiveMemory/ctx/internal/write/sync"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
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

	memoryPath, discoverErr := memory.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		sync.ErrAutoMemoryNotActive(cmd, discoverErr)
		return ctxerr.NotFound()
	}

	data, readErr := io.SafeReadFile(
		filepath.Dir(memoryPath), filepath.Base(memoryPath),
	)
	if readErr != nil {
		return ctxerr.Read(readErr)
	}

	cleaned, found := memory.RemovePublished(string(data))
	if !found {
		publish.UnpublishNotFound(cmd, memory2.MemorySource)
		return nil
	}

	if writeErr := os.WriteFile(
		memoryPath, []byte(cleaned), fs.PermFile,
	); writeErr != nil {
		return ctxerr.Write(writeErr)
	}

	publish.UnpublishDone(cmd, memory2.MemorySource)
	return nil
}
