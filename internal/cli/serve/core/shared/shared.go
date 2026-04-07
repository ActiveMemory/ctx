//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package shared

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/hub"
	"github.com/ActiveMemory/ctx/internal/io"
	writeServe "github.com/ActiveMemory/ctx/internal/write/serve"
)

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
//
// Parameters:
//   - cmd: cobra command for output
//   - port: TCP port to listen on
//   - dataDir: hub data directory (empty = default)
//
// Returns:
//   - error: non-nil if setup or server startup fails
func Run(
	cmd *cobra.Command, port int, dataDir string,
) error {
	if dataDir == "" {
		defaultDir, dirErr := hubDir()
		if dirErr != nil {
			return dirErr
		}
		dataDir = defaultDir
	} else {
		if mkErr := io.SafeMkdirAll(
			dataDir, dataDirPerm,
		); mkErr != nil {
			return mkErr
		}
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

	addr := fmt.Sprintf(":%d", port)
	lis, lisErr := net.Listen("tcp", addr)
	if lisErr != nil {
		return lisErr
	}

	writeServe.HubStarted(cmd, lis.Addr())

	return srv.Serve(lis)
}
