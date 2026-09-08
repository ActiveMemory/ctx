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
// [ClusterStatus] prints the status dashboard from a
// [ClusterStatusInfo], whose fields all come from the
// hub's Status RPC. A standalone hub prints its role
// (Standalone) and entry count; a clustered one adds
// the leader address -- or an election-in-progress
// note while Raft has no leader -- and the peer count.
// It also prints the cumulative slow-listener
// disconnect count, but only when that count is
// non-zero, so a healthy hub keeps its former output.
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
//	hub.ClusterStatus(cmd, hub.ClusterStatusInfo{
//		Role:      role,
//		Clustered: true,
//		Leader:    leaderAddr,
//		Entries:   entries,
//		Peers:     peers,
//	})
//	hub.PeerAdded(cmd, peerAddr)
//	hub.SteppedDown(cmd)
package hub
