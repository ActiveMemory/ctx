//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ClusterStatus prints cluster role and stats.
//
// Parameters:
//   - cmd: Cobra command for output
//   - role: current node role (Leader/Follower)
//   - leader: leader address
//   - entries: total entry count
//   - peers: number of peers
func ClusterStatus(
	cmd *cobra.Command,
	role, leader string,
	entries uint64,
	peers int,
) {
	cmd.Println("Role:", role)
	cmd.Println("Leader:", leader)
	cmd.Println(fmt.Sprintf(
		"Entries: %d  Peers: %d", entries, peers,
	))
}

// PeerAdded confirms a peer was added.
//
// Parameters:
//   - cmd: Cobra command for output
//   - addr: peer address that was added
func PeerAdded(cmd *cobra.Command, addr string) {
	cmd.Println("Added peer", addr)
}

// PeerRemoved confirms a peer was removed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - addr: peer address that was removed
func PeerRemoved(cmd *cobra.Command, addr string) {
	cmd.Println("Removed peer", addr)
}

// SteppedDown confirms leadership transfer.
//
// Parameters:
//   - cmd: Cobra command for output
func SteppedDown(cmd *cobra.Command) {
	cmd.Println("Leadership transferred")
}
