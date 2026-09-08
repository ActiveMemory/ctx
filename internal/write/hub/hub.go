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

// ClusterStatus prints cluster role and stats. The dropped-listener
// line is omitted when the count is zero so a healthy hub keeps its
// current output.
//
// Parameters:
//   - cmd: Cobra command for output
//   - role: current node role (Leader/Follower)
//   - leader: leader address
//   - entries: total entry count
//   - peers: number of peers
//   - dropped: cumulative slow-listener disconnects
func ClusterStatus(
	cmd *cobra.Command,
	role, leader string,
	entries uint64,
	peers int,
	dropped uint64,
) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHubRole), role,
	))
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHubLeader), leader,
	))
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteHubClusterStats),
		entries, peers,
	))
	if dropped > 0 {
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteHubDroppedListeners),
			dropped,
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
