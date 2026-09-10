//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/hub"
	writeHub "github.com/ActiveMemory/ctx/internal/write/hub"
)

// renderInfo maps a Status response onto the rendered fields.
//
// The role, the leader and the peer count are the hub's answers,
// not the client's: a hub with no Raft node reports Standalone
// and has neither a leader nor peers to name.
//
// Parameters:
//   - resp: response from the hub Status RPC
//
// Returns:
//   - writeHub.ClusterStatusInfo: fields to render
func renderInfo(
	resp *hub.StatusResponse,
) writeHub.ClusterStatusInfo {
	info := writeHub.ClusterStatusInfo{
		Role:    cfgHub.RoleStandalone,
		Entries: resp.TotalEntries,
		Dropped: resp.DroppedListeners,
	}

	if !resp.ClusterEnabled {
		return info
	}

	info.Clustered = true
	info.Role = cfgHub.RoleFollower
	if resp.IsLeader {
		info.Role = cfgHub.RoleLeader
	}
	info.Leader = resp.LeaderAddr
	info.Peers = resp.ClusterPeers

	return info
}
