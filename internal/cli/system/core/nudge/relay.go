//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/log/event"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Relay sends a relay notification and appends the same event to the
// local event log. This is the standard two-sink pattern used by most
// hooks after emitting output.
//
// Parameters:
//   - msg: human-readable event description
//   - sessionID: current session identifier
//   - ref: template reference for filtering/aggregation (may be nil)
func Relay(msg, sessionID string, ref *notify.TemplateRef) {
	_ = notify.Send(hook.NotifyChannelRelay, msg, sessionID, ref)
	event.Append(hook.NotifyChannelRelay, msg, sessionID, ref)
}

// EmitAndRelay sends both a nudge and a relay notification, then
// appends the relay event to the local event log. Used by hooks that
// emit both notification types with the same message.
//
// Parameters:
//   - msg: human-readable event description
//   - sessionID: current session identifier
//   - ref: template reference for filtering/aggregation (may be nil)
func EmitAndRelay(msg, sessionID string, ref *notify.TemplateRef) {
	_ = notify.Send(hook.NotifyChannelNudge, msg, sessionID, ref)
	Relay(msg, sessionID, ref)
}

// LoadAndEmit loads a hook message template and, if non-empty, emits the
// standard nudge box + relay notification + throttle marker sequence.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hookName: hook name for message lookup and notifications
//   - variant: hook variant for message lookup and template ref
//   - vars: template variables (may be nil)
//   - fallback: fallback text if no template is found
//   - relayPrefix: relay prefix text
//   - boxTitle: nudge box title
//   - relayMessage: human-readable relay suffix
//   - sessionID: current session identifier
//   - markerPath: throttle file to touch (empty string skips)
func LoadAndEmit(
	cmd *cobra.Command,
	hookName, variant string,
	vars map[string]any,
	fallback,
	relayPrefix, boxTitle, relayMessage, sessionID, markerPath string,
) {
	content := message.Load(hookName, variant, vars, fallback)
	if content == "" {
		return
	}
	Emit(cmd, content,
		relayPrefix, boxTitle, hookName, variant,
		relayMessage, sessionID, vars, markerPath,
	)
}

// Emit is the standard hook tail: print nudge box, send
// nudge+relay notifications, and touch the throttle marker.
//
// Parameters:
//   - cmd: Cobra command for output
//   - content: nudge box content (from Load)
//   - relayPrefix: relay prefix text (e.g., "check-backup-age")
//   - boxTitle: nudge box title
//   - hookName: hook name for notifications
//   - variant: hook variant for template ref
//   - relayMessage: human-readable relay suffix
//   - sessionID: current session identifier
//   - vars: template variables for the template ref (may be nil)
//   - markerPath: throttle file to touch (empty string skips)
func Emit(
	cmd *cobra.Command,
	content, relayPrefix, boxTitle,
	hookName, variant, relayMessage, sessionID string,
	vars map[string]any,
	markerPath string,
) {
	writeSetup.Nudge(cmd, message.NudgeBox(relayPrefix, boxTitle, content))
	ref := notify.NewTemplateRef(hookName, variant, vars)
	Relay(hookName+": "+relayMessage, sessionID, ref)
	if markerPath != "" {
		internalIo.TouchFile(markerPath)
	}
}
