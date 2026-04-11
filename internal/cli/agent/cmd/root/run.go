//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"errors"
	"time"

	"github.com/spf13/cobra"

	coreBudget "github.com/ActiveMemory/ctx/internal/cli/agent/core/budget"
	coreCooldown "github.com/ActiveMemory/ctx/internal/cli/agent/core/cooldown"
	"github.com/ActiveMemory/ctx/internal/config/fmt"
	"github.com/ActiveMemory/ctx/internal/context/load"
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
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
//   - steeringBodies: pre-loaded steering file bodies (may be nil)
//   - skillBody: pre-loaded skill content (empty to omit)
//
// Returns:
//   - error: Non-nil if context loading fails or .context/ is not found
func Run(
	cmd *cobra.Command,
	budget int,
	format string,
	cooldown time.Duration,
	session string,
	steeringBodies []string,
	skillBody string,
	sharedBodies []string,
) error {
	if coreCooldown.Active(session, cooldown) {
		return nil
	}

	ctx, err := load.Do("")
	if err != nil {
		if _, ok := errors.AsType[*errCtx.NotFoundError](err); ok {
			return errInit.NotInitialized()
		}
		return err
	}

	var outputErr error
	if format == fmt.FormatJSON {
		outputErr = coreBudget.OutputAgentJSON(
			cmd, ctx, budget,
			steeringBodies, skillBody,
			sharedBodies,
		)
	} else {
		outputErr = coreBudget.OutputAgentMarkdown(
			cmd, ctx, budget,
			steeringBodies, skillBody,
			sharedBodies,
		)
	}

	if outputErr == nil {
		coreCooldown.TouchTombstone(session)
	}

	return outputErr
}
