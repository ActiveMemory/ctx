//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"

	"github.com/ActiveMemory/ctx/internal/cli/status/core/out"
	"github.com/ActiveMemory/ctx/internal/context/load"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/spf13/cobra"

	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
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
	ctx, err := load.Do("")
	if err != nil {
		var notFoundError *errCtx.NotFoundError
		if errors.As(err, &notFoundError) {
			return errInit.ContextNotInitialized()
		}
		return err
	}

	if jsonOutput {
		return out.PersistStatusJSON(cmd, ctx, verbose)
	}

	return out.PersistStatusText(cmd, ctx, verbose)
}
