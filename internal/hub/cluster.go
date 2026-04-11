//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb/v2"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Cluster wraps a Raft node for leader election only.
//
// Raft is NOT used for data consensus — entries are
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
	raftDir := filepath.Join(dataDir, "raft")
	if mkErr := io.SafeMkdirAll(
		raftDir, fs.PermKeyDir,
	); mkErr != nil {
		return nil, mkErr
	}

	cfg := raft.DefaultConfig()
	cfg.LocalID = raft.ServerID(nodeID)
	cfg.LogOutput = os.Stderr

	addr, resolveErr := net.ResolveTCPAddr(
		"tcp", bindAddr,
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
		filepath.Join(raftDir, "log.db"),
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

	// Bootstrap if single node or first startup.
	if len(peers) == 0 {
		config := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raft.ServerID(nodeID),
					Address: raft.ServerAddress(bindAddr),
				},
			},
		}
		r.BootstrapCluster(config)
	} else {
		servers := make(
			[]raft.Server, 0, len(peers)+1,
		)
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
		config := raft.Configuration{
			Servers: servers,
		}
		r.BootstrapCluster(config)
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

// LeaderAddr returns the address of the current leader.
//
// Returns:
//   - string: leader address, or empty if unknown
func (c *Cluster) LeaderAddr() string {
	_, id := c.raftNode.LeaderWithID()
	return string(id)
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
