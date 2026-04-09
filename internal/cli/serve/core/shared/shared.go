//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package shared

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

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
	return strings.Split(s, ",")
}

// DefaultPort returns the default hub listen port.
//
// Returns:
//   - int: default port number (9900)
func DefaultPort() int { return defaultPort }

// Run starts the shared context hub gRPC server.
//
// On first run, generates an admin token and prints it.
// On subsequent runs, loads the existing token.
// If dataDir is empty, uses ~/.ctx/hub-data/.
// If peers is non-empty, starts Raft cluster for HA.
//
// Parameters:
//   - cmd: cobra command for output
//   - port: TCP port to listen on
//   - dataDir: hub data directory (empty = default)
//   - peers: peer addresses for cluster mode (may be nil)
//
// Returns:
//   - error: non-nil if setup or server startup fails
func Run(
	cmd *cobra.Command,
	port int,
	dataDir string,
	peers []string,
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

	// Start Raft cluster if peers are configured.
	if len(peers) > 0 {
		bindAddr := fmt.Sprintf(":%d", port+1)
		cluster, clusterErr := hub.NewCluster(
			fmt.Sprintf(":%d", port),
			bindAddr, dataDir, peers,
		)
		if clusterErr != nil {
			return clusterErr
		}
		srv.SetCluster(cluster)
	}

	addr := fmt.Sprintf(":%d", port)
	lis, lisErr := net.Listen("tcp", addr)
	if lisErr != nil {
		return lisErr
	}

	writeServe.HubStarted(cmd, lis.Addr())

	return srv.Serve(lis)
}
