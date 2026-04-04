//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

import (
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	cfgTrigger "github.com/ActiveMemory/ctx/internal/config/trigger"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/trigger"
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
		return desc.Text(text.DescKeyMCPHooksDisabled), nil
	}

	hooksDir := rc.HooksDir()
	timeout := time.Duration(rc.HookTimeout()) * time.Second

	input := &entity.TriggerInput{
		TriggerType: cfgTrigger.SessionStart,
		Parameters:  map[string]any{},
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	agg, runErr := trigger.RunAll(
		hooksDir, cfgTrigger.SessionStart, input, timeout,
	)
	if runErr != nil {
		return "", runErr
	}

	if agg.Cancelled {
		return agg.Message, nil
	}

	if agg.Context == "" {
		return desc.Text(text.DescKeyMCPSessionStartOK), nil
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
		return desc.Text(text.DescKeyMCPHooksDisabled), nil
	}

	hooksDir := rc.HooksDir()
	timeout := time.Duration(rc.HookTimeout()) * time.Second

	params := map[string]any{}
	if summary != "" {
		params[field.Summary] = summary
	}

	input := &entity.TriggerInput{
		TriggerType: cfgTrigger.SessionEnd,
		Parameters:  params,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	agg, runErr := trigger.RunAll(
		hooksDir, cfgTrigger.SessionEnd, input, timeout,
	)
	if runErr != nil {
		return "", runErr
	}

	if agg.Cancelled {
		return agg.Message, nil
	}

	if agg.Context == "" {
		return desc.Text(text.DescKeyMCPSessionEndOK), nil
	}

	return agg.Context, nil
}
