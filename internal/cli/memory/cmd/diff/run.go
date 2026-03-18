//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package diff

import (
	"path/filepath"

	ctxerr "github.com/ActiveMemory/ctx/internal/err/memory"
	"github.com/ActiveMemory/ctx/internal/write/memory"
	"github.com/spf13/cobra"

	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
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
		return ctxerr.DiscoverFailed(discoverErr)
	}

	diff, diffErr := mem.Diff(contextDir, sourcePath)
	if diffErr != nil {
		return ctxerr.DiffFailed(diffErr)
	}

	if diff == "" {
		memory.NoChanges(cmd)
		return nil
	}

	cmd.Print(diff)
	return nil
}
