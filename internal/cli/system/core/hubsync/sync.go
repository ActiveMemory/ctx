//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hubsync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connect/core/config"
	"github.com/ActiveMemory/ctx/internal/cli/connect/core/render"
	"github.com/ActiveMemory/ctx/internal/hub"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// connectFile is the encrypted connection config filename.
const connectFile = ".connect.enc"

// Connected reports whether a hub connection config exists.
//
// Returns:
//   - bool: true if .context/.connect.enc exists
func Connected() bool {
	path := filepath.Join(rc.ContextDir(), connectFile)
	_, statErr := os.Stat(path)
	return statErr == nil
}

// Sync pulls new entries from the hub and writes them to
// .context/shared/. Returns the count of synced entries
// and a formatted status message, or empty string if no
// new entries.
//
// Parameters:
//   - sessionID: current session ID (unused, for future)
//
// Returns:
//   - string: status message or empty if nothing synced
func Sync(_ string) string {
	cfg, loadErr := connectCfg.Load()
	if loadErr != nil {
		return ""
	}

	client, dialErr := hub.NewClient(
		cfg.HubAddr, cfg.Token,
	)
	if dialErr != nil {
		return ""
	}
	defer func() { _ = client.Close() }()

	entries, syncErr := client.Sync(
		context.Background(), cfg.Types, 0,
	)
	if syncErr != nil || len(entries) == 0 {
		return ""
	}

	if writeErr := render.WriteEntries(entries); writeErr != nil {
		return ""
	}

	return fmt.Sprintf(
		"Hub sync: %d shared entries updated",
		len(entries),
	)
}
