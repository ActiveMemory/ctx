//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hubsync

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connection/core/config"
	"github.com/ActiveMemory/ctx/internal/cli/connection/core/render"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/hub"
)

// Connected reports whether a hub connection config exists.
//
// ctxDir is supplied by the caller (typically a FullPreamble-gated
// hook) so this function does not re-resolve it; a second resolution
// would be dead code today and would pair an ambiguous (false, err)
// return with the genuine "no hub configured" result.
//
// Returns (false, nil) when the encrypted connect file is absent:
// ordinary "no hub configured" state. A stat failure other than
// not-exist is propagated so the caller can distinguish "no
// connection" from "we could not check."
//
// Parameters:
//   - ctxDir: absolute path to the context directory
//
// Returns:
//   - bool: true if .context/.connect.enc exists
//   - error: non-nil on stat failure other than not-exist
func Connected(ctxDir string) (bool, error) {
	path := filepath.Join(ctxDir, cfgHub.FileConnect)
	_, statErr := os.Stat(path)
	if statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return false, nil
		}
		return false, statErr
	}
	return true, nil
}

// Sync pulls new entries from the hub and writes them to
// .context/hub/. Returns the count of synced entries
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
		desc.Text(text.DescKeyWriteConnectHubSync),
		len(entries),
	)
}
