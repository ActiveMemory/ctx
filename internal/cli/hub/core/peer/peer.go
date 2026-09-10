//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package peer

import (
	"context"

	"github.com/spf13/cobra"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connection/core/config"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
	"github.com/ActiveMemory/ctx/internal/hub"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	writeHub "github.com/ActiveMemory/ctx/internal/write/hub"
)

// Run adds or removes a peer in the hub's Raft configuration.
//
// The hub address comes from the saved connection config (same
// as ctx hub status) and authentication uses the admin token,
// so a fresh client is dialed without a bearer token. Only the
// leader can change the configuration; a follower answers
// FailedPrecondition and the message names ctx hub status as
// the way to find the leader.
//
// The address is the peer's --raft-bind address, not its hub
// port: Raft membership is keyed on the Raft transport.
//
// Parameters:
//   - cmd: cobra command for output
//   - action: "add" or "remove"
//   - addr: Raft address of the peer
//   - adminToken: hub admin token, already resolved
//
// Returns:
//   - error: non-nil on a bad action, config load, dial, or a
//     rejected configuration change
func Run(
	cmd *cobra.Command,
	action string,
	addr string,
	adminToken string,
) error {
	if action != cfgHub.ActionAdd &&
		action != cfgHub.ActionRemove {
		return errHub.InvalidPeerAction(action)
	}

	cfg, loadErr := connectCfg.Load()
	if loadErr != nil {
		return loadErr
	}

	client, dialErr := hub.NewClient(cfg.HubAddr, "")
	if dialErr != nil {
		return dialErr
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			logWarn.Warn(cfgWarn.CloseHubClient, cerr)
		}
	}()

	if peerErr := client.Peer(
		context.Background(), adminToken, action, addr,
	); peerErr != nil {
		return peerErr
	}

	if action == cfgHub.ActionAdd {
		writeHub.PeerAdded(cmd, addr)
	} else {
		writeHub.PeerRemoved(cmd, addr)
	}

	return nil
}
