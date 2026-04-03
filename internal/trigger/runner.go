//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/token"
	errTrigger "github.com/ActiveMemory/ctx/internal/err/trigger"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// RunAll executes all enabled hooks for the given type in alphabetical
// order. It passes input as JSON via stdin and reads HookOutput as JSON
// from stdout.
//
// Behaviour per hook:
//   - cancel:true in output → halt, set Cancelled/Message, return
//   - non-empty context → append to AggregatedOutput.Context
//   - non-zero exit → log error, record in Errors, continue
//   - invalid JSON stdout → log warning, record in Errors, continue
//   - timeout exceeded → kill process, log warning, continue
//
// Returns an empty AggregatedOutput (not nil) when no hooks exist.
//
// Parameters:
//   - hooksDir: root hooks directory (e.g. .context/hooks)
//   - hookType: lifecycle event category
//   - input: JSON object sent to each hook via stdin
//   - timeout: per-hook execution timeout; zero uses DefaultTimeout
//
// Returns:
//   - *AggregatedOutput: aggregated results from all hooks
//   - error: non-nil only on discovery failure
func RunAll(
	hooksDir string,
	hookType HookType,
	input *HookInput,
	timeout time.Duration,
) (*AggregatedOutput, error) {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	agg := &AggregatedOutput{}

	all, discoverErr := Discover(hooksDir)
	if discoverErr != nil {
		return nil, errTrigger.DiscoverFailed(discoverErr)
	}

	hooks := all[hookType]
	if len(hooks) == 0 {
		return agg, nil
	}

	inputJSON, marshalErr := json.Marshal(input)
	if marshalErr != nil {
		return nil, errTrigger.MarshalInput(marshalErr)
	}

	for _, h := range hooks {
		if !h.Enabled {
			continue
		}

		out, runErr := runOne(h, inputJSON, timeout)
		if runErr != nil {
			ctxLog.Warn("hook %s: %v", h.Path, runErr)
			agg.Errors = append(agg.Errors, fmt.Sprintf("%s: %s", h.Path, runErr))
			continue
		}

		if out.Cancel {
			agg.Cancelled = true
			agg.Message = out.Message
			return agg, nil
		}

		if out.Context != "" {
			if agg.Context != "" {
				agg.Context += token.NewlineLF
			}
			agg.Context += out.Context
		}
	}

	return agg, nil
}

// DefaultTimeout is the per-hook execution timeout when none is
// specified by the caller.
const DefaultTimeout = 10 * time.Second
