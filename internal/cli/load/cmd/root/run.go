//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/load/core"
	"github.com/ActiveMemory/ctx/internal/context"
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
	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return fmt.Errorf("no .context/ directory found. Run 'ctx init' first")
		}
		return err
	}

	if raw {
		return core.OutputRaw(cmd, ctx)
	}

	return core.OutputAssembled(cmd, ctx, budget)
}
