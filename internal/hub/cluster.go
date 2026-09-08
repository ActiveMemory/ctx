//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/io"
)

// NewCluster creates a Raft cluster node for leader
// election only.
//
// Raft is NOT used for data consensus; entries are
// replicated via sequence-based gRPC sync. Raft only
// determines which node is the current master.
//
// Parameters:
//   - nodeID: unique identifier for this node
//   - bindAddr: address for Raft communication
//   - dataDir: directory for Raft state
//   - peers: other cluster nodes (empty = single node)
//
// Returns:
//   - *Cluster: initialized Raft cluster node
//   - error: non-nil if setup fails
func NewCluster(
	nodeID string,
	bindAddr string,
	dataDir string,
	peers []string,
) (*Cluster, error) {
	raftDir := filepath.Join(dataDir, cfgHub.RaftDir)
	if mkErr := io.SafeMkdirAll(
		raftDir, fs.PermKeyDir,
	); mkErr != nil {
		return nil, mkErr
	}

	cfg := raft.DefaultConfig()
	cfg.LocalID = raft.ServerID(nodeID)
	cfg.LogOutput = os.Stderr

	addr, resolveErr := net.ResolveTCPAddr(
		cfgHub.RaftTransport, bindAddr,
	)
	if resolveErr != nil {
		return nil, resolveErr
	}

	transport, transErr := raft.NewTCPTransport(
		bindAddr, addr, 3,
		10*time.Second, os.Stderr,
	)
	if transErr != nil {
		return nil, transErr
	}

	logStore, logErr := raftboltdb.NewBoltStore(
		filepath.Join(raftDir, cfgHub.RaftLogDB),
	)
	if logErr != nil {
		return nil, logErr
	}

	snapshotStore := raft.NewDiscardSnapshotStore()

	fsm := &leaderFSM{}

	r, raftErr := raft.NewRaft(
		cfg, fsm, logStore, logStore,
		snapshotStore, transport,
	)
	if raftErr != nil {
		return nil, raftErr
	}

	// Bootstrap this node plus any configured peers. A single
	// node bootstraps a one-server cluster and elects itself.
	servers := make([]raft.Server, 0, len(peers)+1)
	servers = append(servers, raft.Server{
		ID:      raft.ServerID(nodeID),
		Address: raft.ServerAddress(bindAddr),
	})
	for _, p := range peers {
		servers = append(servers, raft.Server{
			ID:      raft.ServerID(p),
			Address: raft.ServerAddress(p),
		})
	}

	// ErrCantBootstrap is what a restart against existing Raft
	// state returns: that node is already bootstrapped and the
	// on-disk configuration wins. Any other error leaves a node
	// that never elects anyone, which is precisely the state the
	// Status RPC's leadership fields exist to make visible -- so
	// it surfaces here instead of being discarded.
	bootErr := r.BootstrapCluster(raft.Configuration{
		Servers: servers,
	}).Error()
	if bootErr != nil &&
		!errors.Is(bootErr, raft.ErrCantBootstrap) {
		return nil, bootErr
	}

	return &Cluster{
		raftNode:  r,
		transport: transport,
	}, nil
}

// IsLeader reports whether this node is the Raft leader.
//
// Returns:
//   - bool: true if this node is the current leader
func (c *Cluster) IsLeader() bool {
	return c.raftNode.State() == raft.Leader
}

// LeaderAddr returns the Raft address of the current leader.
//
// Empty while no leader is known: during an election, or once
// quorum is lost.
//
// Returns:
//   - string: leader address, or empty if unknown
func (c *Cluster) LeaderAddr() string {
	addr, _ := c.raftNode.LeaderWithID()
	return string(addr)
}

// Peers reports how many servers other than this node the
// committed Raft configuration holds. A single-node cluster
// reports zero.
//
// Returns:
//   - uint32: server count excluding this node
//   - error: non-nil if the configuration read fails
func (c *Cluster) Peers() (uint32, error) {
	future := c.raftNode.GetConfiguration()
	if cfgErr := future.Error(); cfgErr != nil {
		return 0, cfgErr
	}

	// Counted rather than converted from len() so the wire type
	// needs no int-to-uint32 narrowing.
	var peers uint32
	for range future.Configuration().Servers {
		peers++
	}
	if peers > 0 {
		peers--
	}

	return peers, nil
}

// Stepdown transfers leadership to another node.
//
// Returns:
//   - error: non-nil if leadership transfer fails
func (c *Cluster) Stepdown() error {
	return c.raftNode.LeadershipTransfer().Error()
}

// Shutdown gracefully stops the Raft node.
//
// Returns:
//   - error: non-nil if shutdown fails
func (c *Cluster) Shutdown() error {
	return c.raftNode.Shutdown().Error()
}
