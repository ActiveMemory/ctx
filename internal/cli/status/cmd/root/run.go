//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/status/core"
	"github.com/ActiveMemory/ctx/internal/context"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
)

// Run executes the status command logic.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - jsonOutput: If true, output as JSON
//   - verbose: If true, include file content previews
//
// Returns:
//   - error: Non-nil if context loading fails
func Run(cmd *cobra.Command, jsonOutput, verbose bool) error {
	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return ctxerr.ContextNotInitialized()
		}
		return err
	}

	if jsonOutput {
		return core.OutputStatusJSON(cmd, ctx, verbose)
	}

	return core.OutputStatusText(cmd, ctx, verbose)
}
