//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package publish

import (
	"path/filepath"

	"github.com/spf13/cobra"

	errMemory "github.com/ActiveMemory/ctx/internal/err/memory"
	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/publish"
	"github.com/ActiveMemory/ctx/internal/write/sync"
)

// Run selects high-value context, formats it, and writes a marked block
// into MEMORY.md. In dry-run mode it reports what would be published.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//   - budget: maximum line count for the published block.
//   - dryRun: when true, show the plan without writing.
//
// Returns:
//   - error: on discovery, selection, or publish failure.
func Run(cmd *cobra.Command, budget int, dryRun bool) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	memoryPath, discoverErr := mem.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		sync.ErrAutoMemoryNotActive(cmd, discoverErr)
		return errMemory.NotFound()
	}

	result, selectErr := mem.SelectContent(contextDir, budget)
	if selectErr != nil {
		return errMemory.SelectContentFailed(selectErr)
	}

	publish.PublishPlan(cmd, budget,
		len(result.Tasks), len(result.Decisions),
		len(result.Conventions), len(result.Learnings),
		result.TotalLines,
	)

	if dryRun {
		publish.PublishDryRun(cmd)
		return nil
	}

	if _, publishErr := mem.Publish(
		contextDir, memoryPath, budget,
	); publishErr != nil {
		return errMemory.PublishFailed(publishErr)
	}

	publish.PublishDone(cmd)

	return nil
}
