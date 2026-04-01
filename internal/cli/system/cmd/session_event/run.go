//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sessionevent

import (
	"fmt"

	"github.com/spf13/cobra"

	coreState "github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	cfgEvent "github.com/ActiveMemory/ctx/internal/config/event"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
	"github.com/ActiveMemory/ctx/internal/log/event"
	"github.com/ActiveMemory/ctx/internal/notify"
	wSession "github.com/ActiveMemory/ctx/internal/write/session"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Run executes the session-event command logic.
//
// Records a session lifecycle event (start or end) to the event log
// and sends a notification. No-op if the context directory is not
// initialized.
//
// Parameters:
//   - cmd: Cobra command for output
//   - eventType: "start" or "end"
//   - caller: identifier of the calling editor (e.g. "vscode")
//
// Returns:
//   - error: Non-nil if eventType is invalid
func Run(cmd *cobra.Command, eventType, caller string) error {
	if !coreState.Initialized() {
		return nil
	}

	if eventType != cfgEvent.TypeStart && eventType != cfgEvent.TypeEnd {
		return errSession.EventInvalidType(
			cfgEvent.TypeStart, cfgEvent.TypeEnd, eventType)
	}

	msg := fmt.Sprintf(desc.Text(text.DescKeyWriteSessionEvent), eventType, caller)
	ref := notify.NewTemplateRef(cfgHook.SessionEvent, eventType,
		map[string]any{"Caller": caller})

	event.Append(cfgEvent.CategorySession, msg, "", ref)
	_ = notify.Send(cfgEvent.CategorySession, msg, "", ref)

	wSession.Event(cmd, eventType, caller)
	return nil
}
