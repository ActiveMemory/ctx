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

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connection/core/config"
	"github.com/ActiveMemory/ctx/internal/cli/connection/core/render"
	"github.com/ActiveMemory/ctx/internal/hub"
	writeConnect "github.com/ActiveMemory/ctx/internal/write/connect"
)

// Run streams entries from the hub in real-time via the
// Listen RPC. Writes each entry to .context/hub/ as
// it arrives. Stops on Ctrl-C.
//
// Parameters:
//   - cmd: cobra command for output
//   - args: unused (cobra signature)
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

	listenErr := client.Listen(
		ctx, cfg.Types, 0,
		func(msg hub.EntryMsg) error {
			writeErr := render.WriteEntries(
				[]hub.EntryMsg{msg},
			)
			if writeErr != nil {
				return writeErr
			}
			writeConnect.EntryReceived(cmd, msg.Type)
			return nil
		},
	)

	// Context cancellation (Ctrl-C) is expected.
	if listenErr != nil && ctx.Err() == nil {
		return listenErr
	}
	return nil
}
