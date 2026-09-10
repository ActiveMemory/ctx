//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// ClusterStatus prints the node role and hub statistics.
//
// A standalone hub prints its role and entry count and stops
// there: it has no leader and no peers to report. A clustered
// hub adds the leader address -- or a note that an election is
// in progress when Raft has no leader yet -- and the peer
// count. The dropped-listener line is omitted when the count is
// zero, so a healthy hub keeps its current output.
//
// Parameters:
//   - cmd: Cobra command for output
//   - info: status fields as reported by the Status RPC
func ClusterStatus(
	cmd *cobra.Command, info ClusterStatusInfo,
) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHubRole), info.Role,
	))

	if info.Clustered {
		clusterLines(cmd, info)
	} else {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteHubEntries),
			info.Entries,
		))
	}

	if info.Dropped > 0 {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteHubDroppedListeners),
			info.Dropped,
		))
	}
}

// PeerAdded confirms a peer was added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - addr: peer address that was added
func PeerAdded(cmd *cobra.Command, addr string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHubAddedPeer), addr,
	))
}

// PeerRemoved confirms a peer was removed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - addr: peer address that was removed
func PeerRemoved(cmd *cobra.Command, addr string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHubRemovedPeer), addr,
	))
}

// Revoked confirms a client token was revoked.
//
// Parameters:
//   - cmd: Cobra command for output
//   - clientID: ID of the client that was revoked
func Revoked(cmd *cobra.Command, clientID string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHubRevoked), clientID,
	))
}

// SteppedDown confirms leadership transfer.
//
// Parameters:
//   - cmd: Cobra command for output
func SteppedDown(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteHubLeadershipTransferred))
}
