//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

// ClusterStatusInfo carries what [ClusterStatus] renders.
//
// Clustered decides the shape: a hub with no Raft node has no
// leader and no peers, and printing either would be inventing
// them. Every field comes from the Status RPC response, not
// from the client's view of the connection.
//
// Fields:
//   - Role: node role label (Leader, Follower, Standalone)
//   - Clustered: a Raft node is attached to the hub
//   - Leader: Raft address of the leader, empty if unknown
//   - Entries: total entry count in the hub's store
//   - Peers: Raft servers other than the one answering
//   - Dropped: cumulative slow-listener disconnects
type ClusterStatusInfo struct {
	Role      string
	Clustered bool
	Leader    string
	Entries   uint64
	Peers     uint32
	Dropped   uint64
}
