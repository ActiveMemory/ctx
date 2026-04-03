//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

import (
	"time"

	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trigger"
)

// Trigger result messages.
const (
	// msgHooksDisabled is returned when triggers are not enabled.
	msgHooksDisabled = "Hooks disabled."
	// msgSessionStartOK is returned when start triggers produce no
	// additional context.
	msgSessionStartOK = "Session start hooks executed. " +
		"No additional context."
	// msgSessionEndOK is returned when end triggers produce no
	// additional context.
	msgSessionEndOK = "Session end hooks executed."
	// paramSummary is the parameter key for session summary.
	paramSummary = "summary"
)

// SessionStartHooks executes session-start triggers and returns
// aggregated context.
//
// Returns success with empty context when no triggers exist or
// triggers are disabled.
//
// Returns:
//   - string: aggregated context from trigger outputs
//   - error: trigger discovery or execution error
func (h *Handler) SessionStartHooks() (string, error) {
	if !rc.HooksEnabled() {
		return msgHooksDisabled, nil
	}

	hooksDir := rc.HooksDir()
	timeout := time.Duration(rc.HookTimeout()) * time.Second

	input := &entity.TriggerInput{
		TriggerType: string(entity.TriggerSessionStart),
		Parameters:  map[string]any{},
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	agg, runErr := trigger.RunAll(
		hooksDir, trigger.SessionStart, input, timeout,
	)
	if runErr != nil {
		return "", runErr
	}

	if agg.Cancelled {
		return agg.Message, nil
	}

	if agg.Context == "" {
		return msgSessionStartOK, nil
	}

	return agg.Context, nil
}

// SessionEndHooks executes session-end triggers with the given summary
// in the trigger input parameters.
//
// Returns success with empty context when no triggers exist or
// triggers are disabled.
//
// Parameters:
//   - summary: optional session summary passed to triggers via parameters
//
// Returns:
//   - string: aggregated context from trigger outputs
//   - error: trigger discovery or execution error
func (h *Handler) SessionEndHooks(summary string) (string, error) {
	if !rc.HooksEnabled() {
		return msgHooksDisabled, nil
	}

	hooksDir := rc.HooksDir()
	timeout := time.Duration(rc.HookTimeout()) * time.Second

	params := map[string]any{}
	if summary != "" {
		params[paramSummary] = summary
	}

	input := &entity.TriggerInput{
		TriggerType: string(entity.TriggerSessionEnd),
		Parameters:  params,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	agg, runErr := trigger.RunAll(
		hooksDir, trigger.SessionEnd, input, timeout,
	)
	if runErr != nil {
		return "", runErr
	}

	if agg.Cancelled {
		return agg.Message, nil
	}

	if agg.Context == "" {
		return msgSessionEndOK, nil
	}

	return agg.Context, nil
}
