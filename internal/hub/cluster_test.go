//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"testing"
	"time"
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

	lis := listenRandom(t)
	addr := lis.Addr().String()
	if closeErr := lis.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}

	cluster, clusterErr := NewCluster(
		clusterNodeID, addr, t.TempDir(), nil,
	)
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

		cluster, clusterErr := NewCluster(
			addr, addr, t.TempDir(), peers,
		)
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
