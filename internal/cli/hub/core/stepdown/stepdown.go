//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stepdown

import (
	"context"

	"github.com/spf13/cobra"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connection/core/config"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/hub"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	writeHub "github.com/ActiveMemory/ctx/internal/write/hub"
)

// Run asks the hub to hand leadership to another node.
//
// The hub address comes from the saved connection config (same
// as ctx hub status) and authentication uses the admin token,
// so a fresh client is dialed without a bearer token. Only a
// leader can transfer leadership: asking a follower fails with
// FailedPrecondition rather than printing a confirmation for
// something that did not happen.
//
// The confirmation prints only after the transfer returns; run
// ctx hub status to see which node won.
//
// Parameters:
//   - cmd: cobra command for output
//   - adminToken: hub admin token, already resolved
//
// Returns:
//   - error: non-nil if config load, dial, or transfer fails
func Run(
	cmd *cobra.Command, adminToken string,
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

	if stepErr := client.Stepdown(
		context.Background(), adminToken,
	); stepErr != nil {
		return stepErr
	}

	writeHub.SteppedDown(cmd)

	return nil
}
