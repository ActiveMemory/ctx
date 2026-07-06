//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package revoke

import (
	"context"

	"github.com/spf13/cobra"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connection/core/config"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/hub"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	writeHub "github.com/ActiveMemory/ctx/internal/write/hub"
)

// Run revokes a client's token on the hub.
//
// The hub address is read from the saved connection config
// (same as `ctx hub status`). Authentication uses the admin
// token, not the stored bearer token, so a fresh client is
// dialed without one.
//
// Parameters:
//   - cmd: cobra command for output
//   - clientID: ID of the client to revoke
//   - adminToken: hub admin token (already resolved from flag
//     or environment by the caller)
//
// Returns:
//   - error: non-nil if config load, dial, or revocation fails
func Run(
	cmd *cobra.Command,
	clientID string,
	adminToken string,
) error {
	cfg, loadErr := connectCfg.Load()
	if loadErr != nil {
		return loadErr
	}

	client, dialErr := hub.NewClient(cfg.HubAddr, "")
	if dialErr != nil {
		return dialErr
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			logWarn.Warn(cfgWarn.CloseHubClient, cerr)
		}
	}()

	if revErr := client.Revoke(
		context.Background(), adminToken, clientID,
	); revErr != nil {
		return revErr
	}

	writeHub.Revoked(cmd, clientID)
	return nil
}
