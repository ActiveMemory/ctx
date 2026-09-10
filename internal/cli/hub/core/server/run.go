//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/hub"
	writeServe "github.com/ActiveMemory/ctx/internal/write/serve"
)

// ParsePeers splits a comma-separated peer string into
// a slice. Returns nil for empty input.
//
// Parameters:
//   - s: comma-separated peer addresses
//
// Returns:
//   - []string: peer addresses, or nil if empty
func ParsePeers(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, token.Comma)
}

// DefaultPort returns the default hub listen port.
//
// Returns:
//   - int: default port number (9900)
func DefaultPort() int { return defaultPort }

// Run starts the ctx Hub gRPC server.
//
// On first run, generates an admin token and prints it.
// On subsequent runs, loads the existing token.
// If dataDir is empty, uses ~/.ctx/hub-data/.
//
// A raftBind address starts the Raft node that makes
// leadership queryable through the Status RPC; peers are the
// raftBind addresses of the other nodes, so every node
// bootstraps the same configuration. Without raftBind the hub
// runs standalone, and asking for peers without it is an error
// rather than a hub that quietly is not in a cluster.
//
// A joining node bootstraps nothing and waits for a leader to
// add it with ctx hub peer add. That is the only way into an
// existing cluster: a node that bootstrapped its own
// configuration would be a second cluster of one.
//
// Parameters:
//   - cmd: cobra command for output
//   - port: TCP port to listen on
//   - dataDir: hub data directory (empty = default)
//   - raftBind: address this node advertises to peers
//   - peers: peer Raft addresses (may be nil)
//   - join: wait to be added instead of bootstrapping
//
// Returns:
//   - error: non-nil if setup or server startup fails
func Run(
	cmd *cobra.Command,
	port int,
	dataDir string,
	raftBind string,
	peers []string,
	join bool,
) error {
	dataDir, resolveErr := resolveDataDir(dataDir)
	if resolveErr != nil {
		return resolveErr
	}

	store, storeErr := hub.NewStore(dataDir)
	if storeErr != nil {
		return storeErr
	}

	adminToken, tokenErr := loadOrCreateAdmin(
		cmd, dataDir,
	)
	if tokenErr != nil {
		return tokenErr
	}

	srv := hub.NewServer(store, adminToken)

	// Start the Raft node when the operator named an address
	// for it. The node registers under that address as both its
	// ID and its transport address, which is the shape peers are
	// given, so every node bootstraps the same configuration.
	if raftBind != "" || len(peers) > 0 || join {
		if clusterErr := startCluster(
			srv, dataDir, raftBind, peers, join,
		); clusterErr != nil {
			return clusterErr
		}
	}

	addr := fmt.Sprintf(cfgHub.FmtPort, port)
	lis, lisErr := net.Listen(cfgHub.RaftTransport, addr)
	if lisErr != nil {
		return lisErr
	}

	writeServe.HubStarted(cmd, lis.Addr())

	return srv.Serve(lis)
}
