//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package unknown

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	errCli "github.com/ActiveMemory/ctx/internal/err/cli"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// relay is the event-log + webhook relay leg, indirected through a
// package variable so tests can observe the call without a live
// webhook or an initialized context. Production binds it to
// [nudge.Relay].
var relay = nudge.Relay

// handle is [Handler] with the stdin source injected for testability.
// Bare invocation prints help; an unknown verb emits the verbatim
// relay box (through the write layer), best-effort records the relay
// event when a session is present, suppresses cobra's help dump, and
// returns the unknown-subcommand error.
//
// Parameters:
//   - cmd: the system command (for output and SilenceUsage)
//   - args: leftover args; non-empty means an unknown subcommand
//   - stdin: hook-input source; ReadID is TTY-safe and timeout-guarded
//
// Returns:
//   - error: nil for a bare `ctx system` (help printed); otherwise the
//     unknown-subcommand error from [errCli.UnknownSubcommand].
func handle(cmd *cobra.Command, args []string, stdin *os.File) error {
	if len(args) == 0 {
		// Bare `ctx system`: preserve help + exit 0.
		return cmd.Help()
	}
	verb := args[0]

	prefix := desc.Text(text.DescKeySystemUnknownRelayPrefix)
	title := desc.Text(text.DescKeySystemUnknownBoxTitle)
	body := fmt.Sprintf(desc.Text(text.DescKeySystemUnknownBody), verb)
	writeSetup.Nudge(cmd, message.NudgeBox(prefix, title, body))

	// Best-effort relay leg: only when a hook supplied a session on
	// stdin. ReadID is TTY-safe and timeout-guarded, so a manual typo
	// at a terminal returns IDUnknown without blocking. A relay
	// failure is logged, not returned: the stdout box already reached
	// the agent, and the user's real problem is the unknown verb.
	if sid := coreSession.ReadID(stdin); sid != cfgSession.IDUnknown {
		msg := fmt.Sprintf(
			desc.Text(text.DescKeySystemUnknownRelayMessage), verb,
		)
		ref := entity.NewTemplateRef(
			hook.System, hook.VariantUnknownSubcommand, nil,
		)
		if relayErr := relay(msg, sid, ref); relayErr != nil {
			logWarn.Warn(warn.RelayUnknownSubcommand, relayErr)
		}
	}

	// Suppress cobra's help dump on error — that dump is the very
	// pollution this handler exists to kill.
	cmd.SilenceUsage = true
	return errCli.UnknownSubcommand(verb)
}
