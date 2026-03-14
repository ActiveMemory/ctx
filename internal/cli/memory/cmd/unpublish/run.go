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
	"github.com/spf13/cobra"

	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/validation"
	"github.com/ActiveMemory/ctx/internal/write"
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
		write.ErrAutoMemoryNotActive(cmd, discoverErr)
		return ctxerr.MemoryNotFound()
	}

	data, readErr := validation.SafeReadFile(
		filepath.Dir(memoryPath), filepath.Base(memoryPath),
	)
	if readErr != nil {
		return ctxerr.ReadMemory(readErr)
	}

	cleaned, found := memory.RemovePublished(string(data))
	if !found {
		write.UnpublishNotFound(cmd, memory2.MemorySource)
		return nil
	}

	if writeErr := os.WriteFile(
		memoryPath, []byte(cleaned), fs.PermFile,
	); writeErr != nil {
		return ctxerr.WriteMemory(writeErr)
	}

	write.UnpublishDone(cmd, memory2.MemorySource)
	return nil
}
