//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"

	"github.com/ActiveMemory/ctx/internal/write/sync"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/sync/core"
	"github.com/ActiveMemory/ctx/internal/context"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
)

// Run executes the sync command logic.
//
// Loads context, detects discrepancies between codebase and documentation,
// and displays suggested actions. In dry-run mode, only shows what would
// be suggested without prompting for changes.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - dryRun: If true, only show suggestions without prompting for changes
//
// Returns:
//   - error: Non-nil if context loading fails or .context/ is not found
func Run(cmd *cobra.Command, dryRun bool) error {
	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return ctxerr.ContextNotInitialized()
		}
		return err
	}

	actions := core.DetectSyncActions(ctx)

	if len(actions) == 0 {
		sync.CtxSyncInSync(cmd)
		return nil
	}

	sync.CtxSyncHeader(cmd, dryRun)

	for i, action := range actions {
		sync.CtxSyncAction(cmd, i+1, action.Type, action.Description, action.Suggestion)
	}

	sync.CtxSyncSummary(cmd, len(actions), dryRun)

	return nil
}
