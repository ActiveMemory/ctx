//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/status/core/out"
	"github.com/ActiveMemory/ctx/internal/context/load"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
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
		if _, ok := errors.AsType[*errCtx.NotFoundError](err); ok {
			return errInit.ContextNotInitialized()
		}
		return err
	}

	if jsonOutput {
		return out.PersistStatusJSON(cmd, ctx, verbose)
	}

	return out.PersistStatusText(cmd, ctx, verbose)
}
