//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hub provides terminal output for the hub
// cluster management commands (ctx hub).
//
// # Cluster Status
//
// [ClusterStatus] prints the full cluster dashboard:
// the current node role (Leader or Follower), the
// leader address, total entry count, and peer count.
//
// # Peer Management
//
// [PeerAdded] confirms a peer was added to the cluster
// and prints the peer address. [PeerRemoved] confirms
// a peer was removed with its address.
//
// # Leadership
//
// [SteppedDown] confirms that leadership was
// transferred to another node. This is printed after
// a successful step-down operation.
//
// # Message Categories
//
//   - Info: cluster status, peer changes, leadership
//     transfer confirmations
//
// # Usage
//
//	hub.ClusterStatus(cmd, role, leader, entries, peers)
//	hub.PeerAdded(cmd, peerAddr)
//	hub.SteppedDown(cmd)
package hub
