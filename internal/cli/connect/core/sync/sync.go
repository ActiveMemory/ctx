//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"context"

	"github.com/spf13/cobra"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connect/core/config"
	"github.com/ActiveMemory/ctx/internal/cli/connect/core/render"
	"github.com/ActiveMemory/ctx/internal/hub"
	writeConnect "github.com/ActiveMemory/ctx/internal/write/connect"
)

// Run syncs entries from the hub to .context/shared/.
//
// Loads connection config, pulls entries since last sync,
// renders them as markdown, and updates sync state.
//
// Parameters:
//   - cmd: cobra command for output
//
// Returns:
//   - error: non-nil if config, sync, or write fails
func Run(cmd *cobra.Command) error {
	cfg, loadErr := connectCfg.Load()
	if loadErr != nil {
		return loadErr
	}

	syncState, stateErr := loadState()
	if stateErr != nil {
		return stateErr
	}

	client, dialErr := hub.NewClient(
		cfg.HubAddr, cfg.Token,
	)
	if dialErr != nil {
		return dialErr
	}
	defer func() { _ = client.Close() }()

	entries, syncErr := client.Sync(
		context.Background(),
		cfg.Types,
		syncState.LastSequence,
	)
	if syncErr != nil {
		return syncErr
	}

	if len(entries) == 0 {
		writeConnect.Synced(cmd, 0)
		return nil
	}

	if writeErr := render.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	syncState.LastSequence = entries[len(entries)-1].Sequence
	if saveErr := saveState(syncState); saveErr != nil {
		return saveErr
	}

	writeConnect.Synced(cmd, len(entries))
	return nil
}
