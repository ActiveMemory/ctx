//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
)

// clusterNodeID is the Raft ServerID the test node registers
// under. It is deliberately not an address: LeaderWithID returns
// (address, id), and an ID that cannot be mistaken for the bound
// address is what makes the LeaderAddr assertions meaningful.
const clusterNodeID = "test-node"

// singleNodeCluster starts a one-server Raft node on a free
// loopback port and blocks until it elects itself. A single node
// is quorum on its own, so the election is deterministic; the
// deadline only bounds how long the test waits for the default
// election timeout to fire.
//
// Returns the cluster and the address it bound, which is the
// leader address every caller asserts against.
func singleNodeCluster(t *testing.T) (*Cluster, string) {
	t.Helper()

	return newSingleNode(t, clusterNodeID, freeAddrs(t, 1)[0])
}

// singleNodeClusterAt is singleNodeCluster on a caller-chosen
// address, registered under that address the way
// ctx hub start --raft-bind registers a node.
func singleNodeClusterAt(
	t *testing.T, addr string,
) (*Cluster, string) {
	t.Helper()

	return newSingleNode(t, addr, addr)
}

// newSingleNode boots one self-electing Raft node and blocks
// until it has elected itself.
func newSingleNode(
	t *testing.T, nodeID, addr string,
) (*Cluster, string) {
	t.Helper()

	cluster, clusterErr := NewCluster(ClusterConfig{
		NodeID:   nodeID,
		BindAddr: addr,
		DataDir:  t.TempDir(),
	})
	if clusterErr != nil {
		t.Fatal(clusterErr)
	}
	t.Cleanup(func() {
		if shutErr := cluster.Shutdown(); shutErr != nil {
			t.Log(shutErr)
		}
	})

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if cluster.IsLeader() {
			return cluster, addr
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("raft node never became leader")

	return nil, ""
}

// TestCluster_PeersExcludesSelf pins what the Peers: line counts.
// The committed configuration of a bootstrapped single node holds
// exactly one server — itself — so the peer count is zero, not
// one.
func TestCluster_PeersExcludesSelf(t *testing.T) {
	cluster, _ := singleNodeCluster(t)

	peers, peersErr := cluster.Peers()
	if peersErr != nil {
		t.Fatal(peersErr)
	}
	if peers != 0 {
		t.Errorf("peers = %d, want 0 for a single node", peers)
	}
}

// TestCluster_LeaderAddrIsTheAddress pins the method to its name.
// LeaderWithID returns (address, id); the method used to discard
// the address and return the ID. The node registers under an ID
// that is not an address, so returning the ID again fails here.
func TestCluster_LeaderAddrIsTheAddress(t *testing.T) {
	cluster, addr := singleNodeCluster(t)

	if leader := cluster.LeaderAddr(); leader != addr {
		t.Errorf("leader addr = %q, want %q",
			leader, addr)
	}
}

// TestHubStatus_StandaloneReportsClusterDisabled pins the
// disambiguator. A hub with no Raft node must report
// ClusterEnabled false: without that flag, a standalone hub and a
// clustered node that has lost its leader are indistinguishable,
// since both report IsLeader false and an empty LeaderAddr.
func TestHubStatus_StandaloneReportsClusterDisabled(
	t *testing.T,
) {
	srv := listenTestServer(t)

	resp, statusErr := srv.hubStatus(testCtx())
	if statusErr != nil {
		t.Fatal(statusErr)
	}

	if resp.ClusterEnabled {
		t.Error("cluster reported enabled with no cluster set")
	}
	if resp.IsLeader {
		t.Error("standalone hub claims leadership")
	}
	if resp.LeaderAddr != "" {
		t.Errorf("leader addr = %q, want empty",
			resp.LeaderAddr)
	}
	if resp.ClusterPeers != 0 {
		t.Errorf("peers = %d, want 0", resp.ClusterPeers)
	}
}

// TestHubStatus_ClusterReportsLeadershipState is the contract
// issue #96 asks for: with a Raft node attached, Status answers
// the leadership question instead of leaving the operator to
// guess from a role derived from the listener count.
func TestHubStatus_ClusterReportsLeadershipState(t *testing.T) {
	srv := listenTestServer(t)
	cluster, addr := singleNodeCluster(t)
	srv.SetCluster(cluster)

	resp, statusErr := srv.hubStatus(testCtx())
	if statusErr != nil {
		t.Fatal(statusErr)
	}

	if !resp.ClusterEnabled {
		t.Error("cluster reported disabled with a cluster set")
	}
	if !resp.IsLeader {
		t.Error("elected node does not report leadership")
	}
	if resp.LeaderAddr != addr {
		t.Errorf("leader addr = %q, want %q",
			resp.LeaderAddr, addr)
	}
	if resp.ClusterPeers != 0 {
		t.Errorf("peers = %d, want 0 for a single node",
			resp.ClusterPeers)
	}
}

// freeAddrs returns n loopback addresses with nothing listening
// on them. Raft binds these itself, so the listeners only serve
// to have the kernel pick free ports.
func freeAddrs(t *testing.T, n int) []string {
	t.Helper()

	addrs := make([]string, 0, n)
	for range n {
		lis := listenRandom(t)
		addrs = append(addrs, lis.Addr().String())
		if closeErr := lis.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
	}

	return addrs
}

// startCluster brings up one Raft node per address, each
// bootstrapped with the full server list. Every node registers
// under its own address as both ID and transport address, which
// is the shape ctx hub start --raft-bind gives them.
func startCluster(t *testing.T, addrs []string) []*Cluster {
	t.Helper()

	clusters := make([]*Cluster, 0, len(addrs))
	for i, addr := range addrs {
		peers := make([]string, 0, len(addrs)-1)
		for j, peer := range addrs {
			if i != j {
				peers = append(peers, peer)
			}
		}

		cluster, clusterErr := NewCluster(ClusterConfig{
			NodeID:   addr,
			BindAddr: addr,
			DataDir:  t.TempDir(),
			Peers:    peers,
		})
		if clusterErr != nil {
			t.Fatal(clusterErr)
		}
		t.Cleanup(func() {
			if shutErr := cluster.Shutdown(); shutErr != nil {
				t.Log(shutErr)
			}
		})
		clusters = append(clusters, cluster)
	}

	return clusters
}

// TestCluster_ThreeNodesElectOneLeader is the contract the
// Status fields exist to report, over a real three-node Raft
// cluster: exactly one node leads, all three name the same
// leader, and each counts the other two as peers.
//
// It is also the regression test for the address arithmetic
// this branch removed. Binding Raft to ":port+1" and advertising
// the same wildcard made NewCluster fail with "local bind
// address is not advertisable" on every node, so no ctx cluster
// could start at all.
func TestCluster_ThreeNodesElectOneLeader(t *testing.T) {
	addrs := freeAddrs(t, 3)
	clusters := startCluster(t, addrs)

	leader := waitForLeader(t, clusters)

	for i, cluster := range clusters {
		if got := cluster.LeaderAddr(); got != leader {
			t.Errorf("node %d names leader %q, want %q",
				i, got, leader)
		}

		peers, peersErr := cluster.Peers()
		if peersErr != nil {
			t.Fatal(peersErr)
		}
		if peers != 2 {
			t.Errorf("node %d reports %d peers, want 2",
				i, peers)
		}
	}
}

// waitForLeader blocks until exactly one node reports leadership
// and every node agrees on its address, then returns that
// address.
func waitForLeader(
	t *testing.T, clusters []*Cluster,
) string {
	t.Helper()

	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		leaders := 0
		addr := ""
		agreed := true

		for _, cluster := range clusters {
			if cluster.IsLeader() {
				leaders++
			}
			switch {
			case cluster.LeaderAddr() == "":
				agreed = false
			case addr == "":
				addr = cluster.LeaderAddr()
			case cluster.LeaderAddr() != addr:
				agreed = false
			}
		}

		if leaders == 1 && agreed {
			return addr
		}
		time.Sleep(50 * time.Millisecond)
	}

	t.Fatal("cluster never settled on a single leader")

	return ""
}

// joinCluster starts a Raft node that bootstraps nothing and
// waits to be added by a leader, which is what
// ctx hub start --join does.
func joinCluster(t *testing.T, addr string) *Cluster {
	t.Helper()

	cluster, clusterErr := NewCluster(ClusterConfig{
		NodeID:   addr,
		BindAddr: addr,
		DataDir:  t.TempDir(),
		Join:     true,
	})
	if clusterErr != nil {
		t.Fatal(clusterErr)
	}
	t.Cleanup(func() {
		if shutErr := cluster.Shutdown(); shutErr != nil {
			t.Log(shutErr)
		}
	})

	return cluster
}

// waitFor polls cond until it holds or the deadline passes.
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()

	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}

	t.Fatalf("timed out waiting for %s", what)
}

// TestCluster_PeerAddJoinsNode is what ctx hub peer add has to
// mean: a node started in join mode holds no configuration
// until the leader adds it, and then it is a member -- it
// counts as a peer and it names the same leader. Before this
// change the command printed "Added peer" and reached nothing.
func TestCluster_PeerAddJoinsNode(t *testing.T) {
	addrs := freeAddrs(t, 2)
	leader, leaderAddr := singleNodeClusterAt(t, addrs[0])
	joiner := joinCluster(t, addrs[1])

	if addErr := leader.AddPeer(addrs[1]); addErr != nil {
		t.Fatal(addErr)
	}

	waitFor(t, "the joiner to follow the leader", func() bool {
		return joiner.LeaderAddr() == leaderAddr
	})

	peers, peersErr := leader.Peers()
	if peersErr != nil {
		t.Fatal(peersErr)
	}
	if peers != 1 {
		t.Errorf("leader reports %d peers, want 1", peers)
	}
	if joiner.IsLeader() {
		t.Error("the joiner elected itself instead of joining")
	}
}

// TestCluster_PeerRemoveShrinksCluster pins the other half:
// removing a decommissioned node takes it back out of the
// committed configuration, so it stops counting toward quorum.
func TestCluster_PeerRemoveShrinksCluster(t *testing.T) {
	addrs := freeAddrs(t, 2)
	leader, leaderAddr := singleNodeClusterAt(t, addrs[0])
	joiner := joinCluster(t, addrs[1])

	if addErr := leader.AddPeer(addrs[1]); addErr != nil {
		t.Fatal(addErr)
	}
	waitFor(t, "the joiner to follow the leader", func() bool {
		return joiner.LeaderAddr() == leaderAddr
	})

	if remErr := leader.RemovePeer(addrs[1]); remErr != nil {
		t.Fatal(remErr)
	}

	waitFor(t, "the peer count to drop", func() bool {
		peers, peersErr := leader.Peers()
		return peersErr == nil && peers == 0
	})
}

// TestCluster_StepdownHandsOffLeadership pins ctx hub stepdown:
// the node that was leading is not leading afterwards, and the
// other node is. The command used to print "Leadership
// transferred" without asking anyone.
func TestCluster_StepdownHandsOffLeadership(t *testing.T) {
	addrs := freeAddrs(t, 2)
	clusters := startCluster(t, addrs)
	waitForLeader(t, clusters)

	var leader, other *Cluster
	if clusters[0].IsLeader() {
		leader, other = clusters[0], clusters[1]
	} else {
		leader, other = clusters[1], clusters[0]
	}

	if stepErr := leader.Stepdown(); stepErr != nil {
		t.Fatal(stepErr)
	}

	waitFor(t, "leadership to move", func() bool {
		return !leader.IsLeader() && other.IsLeader()
	})
}

// TestPeer_FollowerIsPrecondition pins the error mapping that
// makes a follower's refusal actionable: raft.ErrNotLeader
// becomes FailedPrecondition with a message pointing at
// ctx hub status, not an opaque Internal.
func TestPeer_FollowerIsPrecondition(t *testing.T) {
	addrs := freeAddrs(t, 2)
	clusters := startCluster(t, addrs)
	waitForLeader(t, clusters)

	follower := clusters[0]
	if follower.IsLeader() {
		follower = clusters[1]
	}

	srv := listenTestServer(t)
	srv.SetCluster(follower)

	_, peerErr := srv.peer(testCtx(), &PeerRequest{
		AdminToken: srv.adminToken,
		Action:     cfgHub.ActionRemove,
		Address:    addrs[0],
	})

	if got := status.Code(peerErr); got != codes.FailedPrecondition {
		t.Errorf("code = %v, want FailedPrecondition", got)
	}
}
