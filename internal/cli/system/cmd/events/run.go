//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package events

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	errRcall "github.com/ActiveMemory/ctx/internal/err/recall"
	"github.com/ActiveMemory/ctx/internal/log"
	writeEvents "github.com/ActiveMemory/ctx/internal/write/events"
)

// Run executes the events subcommand, querying and displaying event log
// entries filtered by hook, session, event type, and count.
//
// Parameters:
//   - cmd: Cobra command for flag access and output
//
// Returns:
//   - error: Non-nil on event log read failure
func Run(cmd *cobra.Command) error {
	hook, _ := cmd.Flags().GetString(cFlag.Hook)
	session, _ := cmd.Flags().GetString(cFlag.Session)
	event, _ := cmd.Flags().GetString(cFlag.Event)
	last, _ := cmd.Flags().GetInt(cFlag.Last)
	jsonOut, _ := cmd.Flags().GetBool(cFlag.JSON)
	includeAll, _ := cmd.Flags().GetBool(cFlag.All)

	opts := log.QueryOpts{
		Hook:           hook,
		Session:        session,
		Event:          event,
		Last:           last,
		IncludeRotated: includeAll,
	}

	evts, queryErr := log.Query(opts)
	if queryErr != nil {
		return errRcall.EventLogRead(queryErr)
	}

	if len(evts) == 0 {
		writeEvents.Empty(cmd)
		return nil
	}

	if jsonOut {
		writeEvents.JSON(cmd, core.FormatEventsJSON(evts))
	} else {
		writeEvents.Human(cmd, core.FormatEventsHuman(evts))
	}
	return nil
}
