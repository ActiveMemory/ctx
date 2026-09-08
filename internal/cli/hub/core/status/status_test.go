//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"testing"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/hub"
)

// TestRenderInfo_Standalone pins the case the old renderer got
// wrong in every field: a hub with no Raft node used to report
// its role from the listener count, name itself the leader, and
// report the project count as a peer count.
func TestRenderInfo_Standalone(t *testing.T) {
	info := renderInfo(&hub.StatusResponse{
		TotalEntries:     42,
		ConnectedClients: 3,
		EntriesByProject: map[string]uint64{
			"alpha": 1, "beta": 2,
		},
	})

	if info.Role != cfgHub.RoleStandalone {
		t.Errorf("role = %q, want %q",
			info.Role, cfgHub.RoleStandalone)
	}
	if info.Clustered {
		t.Error("standalone response rendered as clustered")
	}
	if info.Leader != "" {
		t.Errorf("leader = %q, want empty", info.Leader)
	}
	if info.Peers != 0 {
		t.Errorf("peers = %d, want 0 (projects are not peers)",
			info.Peers)
	}
}

// TestRenderInfo_Leader covers the node that answers yes to
// "am I the leader?".
func TestRenderInfo_Leader(t *testing.T) {
	info := renderInfo(&hub.StatusResponse{
		TotalEntries:   42,
		ClusterEnabled: true,
		IsLeader:       true,
		LeaderAddr:     "10.0.0.5:9901",
		ClusterPeers:   2,
	})

	if info.Role != cfgHub.RoleLeader {
		t.Errorf("role = %q, want %q",
			info.Role, cfgHub.RoleLeader)
	}
	if !info.Clustered {
		t.Error("clustered response rendered as standalone")
	}
	if info.Leader != "10.0.0.5:9901" {
		t.Errorf("leader = %q, want the reported address",
			info.Leader)
	}
	if info.Peers != 2 {
		t.Errorf("peers = %d, want 2", info.Peers)
	}
}

// TestRenderInfo_FollowerWithNoListeners is the inversion of the
// old heuristic: a follower with subscribers used to read Active,
// and a leader with none read Follower. The role now comes from
// Raft, so the listener count does not move it.
func TestRenderInfo_FollowerWithNoListeners(t *testing.T) {
	info := renderInfo(&hub.StatusResponse{
		ClusterEnabled: true,
		LeaderAddr:     "10.0.0.5:9901",
		ClusterPeers:   2,
	})

	if info.Role != cfgHub.RoleFollower {
		t.Errorf("role = %q, want %q",
			info.Role, cfgHub.RoleFollower)
	}
}
