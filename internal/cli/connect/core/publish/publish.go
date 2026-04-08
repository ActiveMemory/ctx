//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package publish

import (
	"context"

	"github.com/spf13/cobra"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connect/core/config"
	"github.com/ActiveMemory/ctx/internal/hub"
	writeConnect "github.com/ActiveMemory/ctx/internal/write/connect"
)

// Run publishes local entries to the hub.
//
// Currently publishes entries passed as arguments.
// Future: read from local context files with --new flag.
//
// Parameters:
//   - cmd: cobra command for output
//   - entries: entries to publish
//
// Returns:
//   - error: non-nil if config load or publish fails
func Run(
	cmd *cobra.Command, entries []hub.PublishEntry,
) error {
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

	_, pubErr := client.Publish(
		context.Background(), entries,
	)
	if pubErr != nil {
		return pubErr
	}

	writeConnect.Published(cmd, len(entries))
	return nil
}
