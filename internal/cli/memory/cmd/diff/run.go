//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package diff

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/memory/core/resolve"
	errMemory "github.com/ActiveMemory/ctx/internal/err/memory"
	mem "github.com/ActiveMemory/ctx/internal/memory"
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
	contextDir, projectRoot, err := resolve.ContextAndRoot(cmd)
	if err != nil {
		return err
	}

	sourcePath, discoverErr := mem.DiscoverPath(projectRoot)
	if discoverErr != nil {
		return errMemory.DiscoverFailed(discoverErr)
	}

	diff, diffErr := mem.Diff(contextDir, sourcePath)
	if diffErr != nil {
		return errMemory.DiffFailed(diffErr)
	}

	if diff == "" {
		memory.NoChanges(cmd)
		return nil
	}

	memory.DiffOutput(cmd, diff)
	return nil
}
