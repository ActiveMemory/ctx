//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sessionevent

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	coreState "github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgEvent "github.com/ActiveMemory/ctx/internal/config/event"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	errSession "github.com/ActiveMemory/ctx/internal/err/session"
	"github.com/ActiveMemory/ctx/internal/log/event"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/notify"
	wSession "github.com/ActiveMemory/ctx/internal/write/session"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
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
	isAgy := caller == cfgHook.ToolAgy ||
		caller == cfgHook.ToolAntigravity ||
		coreSession.Agy()
	if isAgy {
		_ = coreCheck.ReadInput(os.Stdin)
	}
	if coreSession.Agy() {
		defer writeSetup.EmitAgyEmpty(cmd)
	}

	initialized, initErr := coreState.Initialized()
	if initErr != nil {
		logWarn.Warn(warn.StateInitializedProbe, initErr)
		return nil
	}
	if !initialized {
		return nil
	}

	if eventType != cfgEvent.TypeStart && eventType != cfgEvent.TypeEnd {
		return errSession.EventInvalidType(
			cfgEvent.TypeStart, cfgEvent.TypeEnd, eventType)
	}

	msg := fmt.Sprintf(
		desc.Text(text.DescKeyWriteSessionEvent), eventType, caller,
	)
	ref := entity.NewTemplateRef(cfgHook.SessionEvent, eventType,
		map[string]any{cfgEvent.VarCaller: caller})

	// Log-first: the event-log entry IS the authoritative record of
	// the session lifecycle. If it cannot be written, neither the
	// webhook nor the stdout marker should run; both would claim a
	// session event whose audit trail never landed. See
	// docs/security/reporting.md → "Log-First Audit Trail".
	if appendErr := event.Append(
		cfgEvent.CategorySession, msg, "", ref,
	); appendErr != nil {
		return appendErr
	}
	if sendErr := notify.Send(
		cfgEvent.CategorySession, msg, "", ref,
	); sendErr != nil {
		return sendErr
	}

	if isAgy || coreSession.Agy() {
		return nil
	}

	wSession.Event(cmd, eventType, caller)
	return nil
}
