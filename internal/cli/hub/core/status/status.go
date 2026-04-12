//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"context"

	"github.com/spf13/cobra"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connection/core/config"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/hub"
	writeHub "github.com/ActiveMemory/ctx/internal/write/hub"
)

// Run shows cluster status via the hub Status RPC.
//
// Parameters:
//   - cmd: cobra command for output
//   - args: unused (cobra signature)
//
// Returns:
//   - error: non-nil if config load or status call fails
func Run(cmd *cobra.Command, _ []string) error {
	cfg, loadErr := connectCfg.Load()
	if loadErr != nil {
		return loadErr
	}

	client, dialErr := hub.NewClient(
		cfg.HubAddr, cfg.Token,
	)
	if dialErr != nil {
		return dialErr
	}
	defer func() { _ = client.Close() }()

	resp, statusErr := client.Status(
		context.Background(),
	)
	if statusErr != nil {
		return statusErr
	}

	role := cfgHub.RoleFollower
	if resp.ConnectedClients > 0 {
		role = cfgHub.RoleActive
	}

	writeHub.ClusterStatus(
		cmd, role, cfg.HubAddr,
		resp.TotalEntries,
		len(resp.EntriesByProject),
	)
	return nil
}
