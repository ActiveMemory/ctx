//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package peer

import (
	"github.com/spf13/cobra"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
	writeHub "github.com/ActiveMemory/ctx/internal/write/hub"
)

// Run handles peer add/remove subcommands.
//
// Parameters:
//   - cmd: cobra command for output
//   - args: [action, address] where action is add or remove
//
// Returns:
//   - error: non-nil if action is invalid
func Run(cmd *cobra.Command, args []string) error {
	action := args[0]
	addr := args[1]

	switch action {
	case cfgHub.ActionAdd:
		writeHub.PeerAdded(cmd, addr)
	case cfgHub.ActionRemove:
		writeHub.PeerRemoved(cmd, addr)
	default:
		return errHub.InvalidPeerAction(action)
	}
	return nil
}
