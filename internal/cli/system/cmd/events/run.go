//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package events

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/recall"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/log"
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
	hook, _ := cmd.Flags().GetString("hook")
	session, _ := cmd.Flags().GetString("session")
	event, _ := cmd.Flags().GetString("event")
	last, _ := cmd.Flags().GetInt("last")
	jsonOut, _ := cmd.Flags().GetBool("json")
	includeAll, _ := cmd.Flags().GetBool("all")

	opts := log.QueryOpts{
		Hook:           hook,
		Session:        session,
		Event:          event,
		Last:           last,
		IncludeRotated: includeAll,
	}

	evts, queryErr := log.Query(opts)
	if queryErr != nil {
		return ctxerr.EventLogRead(queryErr)
	}

	if len(evts) == 0 {
		cmd.Println(desc.TextDesc(text.DescKeyEventsEmpty))
		return nil
	}

	if jsonOut {
		return core.OutputEventsJSON(cmd, evts)
	}
	return core.OutputEventsHuman(cmd, evts)
}
