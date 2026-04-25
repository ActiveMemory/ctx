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
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/log/event"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Relay appends a relay event to the local event log and, only on
// success, sends the relay webhook notification. The order is
// deliberate: the log is the authoritative record; the webhook is a
// side effect that must not claim an event happened unless it was
// first recorded. See docs/security/reporting.md →
// "Log-First Audit Trail".
//
// Parameters:
//   - msg: human-readable event description
//   - sessionID: current session identifier
//   - ref: template reference for filtering/aggregation (may be nil)
//
// Returns:
//   - error: non-nil when the event-log append fails (webhook is
//     skipped in that case) or when the webhook send itself fails
//     (log was written). Callers propagate to surface real failures
//     rather than pretend the notification succeeded.
func Relay(msg, sessionID string, ref *entity.TemplateRef) error {
	if appendErr := event.Append(
		hook.NotifyChannelRelay, msg, sessionID, ref,
	); appendErr != nil {
		return appendErr
	}
	return notify.Send(hook.NotifyChannelRelay, msg, sessionID, ref)
}

// EmitAndRelay sends both a nudge and a relay notification, then
// appends the relay event to the local event log.
//
// The nudge webhook has no corresponding event-log channel today,
// so log-first ordering cannot apply to the nudge leg; this is a
// known gap. A future refactor may add a nudge channel to
// [event.Append]; until then the nudge webhook can fire even if the
// later relay log fails. The relay leg itself follows [Relay]'s
// log-first ordering.
//
// Parameters:
//   - msg: human-readable event description
//   - sessionID: current session identifier
//   - ref: template reference for filtering/aggregation (may be nil)
//
// Returns:
//   - error: non-nil when the nudge send, the relay log, or the
//     relay webhook fails. A nudge failure short-circuits the relay
//     so we do not send half a story.
func EmitAndRelay(msg, sessionID string, ref *entity.TemplateRef) error {
	if sendErr := notify.Send(
		hook.NotifyChannelNudge, msg, sessionID, ref,
	); sendErr != nil {
		return sendErr
	}
	return Relay(msg, sessionID, ref)
}

// LoadAndEmit loads a hook message template and, if non-empty, emits
// the standard nudge box + relay notification + throttle marker
// sequence.
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
//
// Returns:
//   - error: propagated from [Emit] when the message is non-empty;
//     nil when the template resolved to an empty message (intentional
//     silence).
func LoadAndEmit(
	cmd *cobra.Command,
	hookName, variant string,
	vars map[string]any,
	fallback,
	relayPrefix, boxTitle, relayMessage, sessionID, markerPath string,
) error {
	content := message.Load(hookName, variant, vars, fallback)
	if content == "" {
		return nil
	}
	return Emit(cmd, content,
		relayPrefix, boxTitle, hookName, variant,
		relayMessage, sessionID, vars, markerPath,
	)
}

// Emit is the standard hook tail: print nudge box, send nudge+relay
// notifications, and touch the throttle marker.
//
// The throttle marker is only touched on a successful relay: marking
// a hook as recently-emitted when the emit actually failed would
// suppress retries on a real problem.
//
// Parameters:
//   - cmd: Cobra command for output
//   - content: nudge box content (from Load)
//   - relayPrefix: relay prefix text (e.g., "check-ceremony")
//   - boxTitle: nudge box title
//   - hookName: hook name for notifications
//   - variant: hook variant for template ref
//   - relayMessage: human-readable relay suffix
//   - sessionID: current session identifier
//   - vars: template variables for the template ref (may be nil)
//   - markerPath: throttle file to touch (empty string skips)
//
// Returns:
//   - error: propagated from [Relay] (log or webhook failure).
func Emit(
	cmd *cobra.Command,
	content, relayPrefix, boxTitle,
	hookName, variant, relayMessage, sessionID string,
	vars map[string]any,
	markerPath string,
) error {
	writeSetup.Nudge(cmd, message.NudgeBox(relayPrefix, boxTitle, content))
	ref := entity.NewTemplateRef(hookName, variant, vars)
	if relayErr := Relay(
		hookName+token.ColonSpace+relayMessage, sessionID, ref,
	); relayErr != nil {
		return relayErr
	}
	if markerPath != "" {
		internalIo.TouchFile(markerPath)
	}
	return nil
}
