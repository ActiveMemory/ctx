//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/fmt"
	"github.com/ActiveMemory/ctx/internal/context/load"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/agent/core"
	errctx "github.com/ActiveMemory/ctx/internal/err/context"
)

// Run executes the agent command logic.
//
// When a session and cooldown are provided, it checks a tombstone file
// to suppress repeated output within the cooldown window. On the first
// invocation (or after cooldown expires), it loads context from .context/
// and outputs a context packet in the specified format.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - budget: Token budget to include in the output
//   - format: Output format, "json" for JSON, or any other value
//     for Markdown
//   - cooldown: duration to suppress repeated output (0 to disable)
//   - session: session identifier for tombstone isolation (empty to
//     disable cooldown)
//
// Returns:
//   - error: Non-nil if context loading fails or .context/ is not found
func Run(
	cmd *cobra.Command,
	budget int,
	format string,
	cooldown time.Duration,
	session string,
) error {
	if core.CooldownActive(session, cooldown) {
		return nil
	}

	ctx, err := load.Do("")
	if err != nil {
		var notFoundError *errctx.NotFoundError
		if errors.As(err, &notFoundError) {
			return ctxerr.NotInitialized()
		}
		return err
	}

	var outputErr error
	if format == fmt.FormatJSON {
		outputErr = core.OutputAgentJSON(cmd, ctx, budget)
	} else {
		outputErr = core.OutputAgentMarkdown(cmd, ctx, budget)
	}

	if outputErr == nil {
		core.TouchTombstone(session)
	}

	return outputErr
}
