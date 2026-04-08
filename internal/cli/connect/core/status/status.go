//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"context"

	"github.com/spf13/cobra"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connect/core/config"
	"github.com/ActiveMemory/ctx/internal/hub"
	writeConnect "github.com/ActiveMemory/ctx/internal/write/connect"
)

// Run shows hub connection status and entry statistics.
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

	writeConnect.Status(
		cmd, cfg.HubAddr,
		resp.TotalEntries, resp.ConnectedClients,
	)
	return nil
}
