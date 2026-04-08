//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package listen

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connect/core/config"
	"github.com/ActiveMemory/ctx/internal/cli/connect/core/render"
	"github.com/ActiveMemory/ctx/internal/hub"
	writeConnect "github.com/ActiveMemory/ctx/internal/write/connect"
)

// Run streams entries from the hub in real-time.
//
// Connects to the hub, starts a Sync to get recent entries,
// then blocks waiting for new entries. Writes each entry to
// .context/shared/ as it arrives. Stops on Ctrl-C.
//
// Parameters:
//   - cmd: cobra command for output
//
// Returns:
//   - error: non-nil if config, connection, or write fails
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

	ctx, stop := signal.NotifyContext(
		context.Background(), os.Interrupt,
	)
	defer stop()

	writeConnect.Listening(cmd)

	entries, syncErr := client.Sync(
		ctx, cfg.Types, 0,
	)
	if syncErr != nil {
		return syncErr
	}
	if len(entries) > 0 {
		if writeErr := render.WriteEntries(
			entries,
		); writeErr != nil {
			return writeErr
		}
		writeConnect.Synced(cmd, len(entries))
	}

	<-ctx.Done()
	return nil
}
