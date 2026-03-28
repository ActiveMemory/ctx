//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/load/core"
	"github.com/ActiveMemory/ctx/internal/context/load"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	writeLoad "github.com/ActiveMemory/ctx/internal/write/load"
)

// Run executes the load command logic.
//
// Loads context from .context/ and outputs it in either raw or assembled
// format based on the flags.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - budget: Token budget for assembled output
//   - raw: If true, output raw file contents without assembly
//
// Returns:
//   - error: Non-nil if context loading fails or .context/ is not found
func Run(cmd *cobra.Command, budget int, raw bool) error {
	ctx, err := load.Do("")
	if err != nil {
		var notFoundError *errCtx.NotFoundError
		if errors.As(err, &notFoundError) {
			return errInit.NotInit()
		}
		return err
	}

	files := core.SortByReadOrder(ctx.Files)

	if raw {
		return writeLoad.Raw(cmd, files)
	}

	return writeLoad.Assembled(
		cmd, files, budget, ctx.TotalTokens, core.FileNameToTitle,
	)
}
