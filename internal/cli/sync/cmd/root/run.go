//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/sync/core/action"
	"github.com/ActiveMemory/ctx/internal/context/load"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/sync"
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
	ctx, err := load.Do("")
	if err != nil {
		var notFoundError *errCtx.NotFoundError
		if errors.As(err, &notFoundError) {
			return errInit.ContextNotInitialized()
		}
		return err
	}

	actions := action.Detect(ctx)

	if len(actions) == 0 {
		sync.InSync(cmd)
		return nil
	}

	sync.Header(cmd, dryRun)

	for i, action := range actions {
		sync.Action(
			cmd, i+1, action.Type, action.Description, action.Suggestion,
		)
	}

	sync.Summary(cmd, len(actions), dryRun)

	return nil
}
