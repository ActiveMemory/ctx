//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package diff

import (
	"path/filepath"

	"github.com/spf13/cobra"

	ctxErr "github.com/ActiveMemory/ctx/internal/err/memory"
	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/memory"
)

// Run computes and prints a line-based diff between the mirror and
// current MEMORY.md.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//
// Returns:
//   - error: on discovery or diff failure.
func Run(cmd *cobra.Command) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := mem.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		return ctxErr.DiscoverFailed(discoverErr)
	}

	diff, diffErr := mem.Diff(contextDir, sourcePath)
	if diffErr != nil {
		return ctxErr.DiffFailed(diffErr)
	}

	if diff == "" {
		memory.NoChanges(cmd)
		return nil
	}

	cmd.Print(diff)
	return nil
}
