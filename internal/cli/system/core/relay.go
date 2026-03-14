//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/eventlog"
	"github.com/ActiveMemory/ctx/internal/notify"
)

// Relay sends a relay notification and appends the same event to the
// local event log. This is the standard two-sink pattern used by most
// hooks after emitting output.
//
// Parameters:
//   - message: human-readable event description
//   - sessionID: current session identifier
//   - ref: template reference for filtering/aggregation (may be nil)
func Relay(message, sessionID string, ref *notify.TemplateRef) {
	_ = notify.Send(hook.NotifyChannelRelay, message, sessionID, ref)
	eventlog.Append(hook.NotifyChannelRelay, message, sessionID, ref)
}

// NudgeAndRelay sends both a nudge and a relay notification, then
// appends the relay event to the local event log. Used by hooks that
// emit both notification types with the same message.
//
// Parameters:
//   - message: human-readable event description
//   - sessionID: current session identifier
//   - ref: template reference for filtering/aggregation (may be nil)
func NudgeAndRelay(message, sessionID string, ref *notify.TemplateRef) {
	_ = notify.Send(hook.NotifyChannelNudge, message, sessionID, ref)
	Relay(message, sessionID, ref)
}
