//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
)

// runAgent executes the agent command logic.
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
func runAgent(
	cmd *cobra.Command,
	budget int,
	format string,
	cooldown time.Duration,
	session string,
) error {
	if cooldownActive(session, cooldown) {
		return nil
	}

	ctx, err := context.Load("")
	if err != nil {
		var notFoundError *context.NotFoundError
		if errors.As(err, &notFoundError) {
			return fmt.Errorf(
				"no .context/ directory found. Run 'ctx init' first",
			)
		}
		return err
	}

	var outputErr error
	if format == config.FormatJSON {
		outputErr = outputAgentJSON(cmd, ctx, budget)
	} else {
		outputErr = outputAgentMarkdown(cmd, ctx, budget)
	}

	if outputErr == nil {
		touchTombstone(session)
	}

	return outputErr
}
