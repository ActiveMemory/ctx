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
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connection/core/config"
	"github.com/ActiveMemory/ctx/internal/cli/connection/core/render"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/hub"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
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

// syncTimeout bounds the session-start pull RPC. It is a var,
// not the bare constant, so tests can shrink it to exercise the
// deadline against a deliberately unresponsive hub without
// waiting the full production interval.
var syncTimeout = time.Duration(cfgHub.HubSyncTimeout) * time.Second

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
		logWarn.Warn(cfgWarn.HubSyncLoadConfig, loadErr)
		return ""
	}

	client, dialErr := hub.NewClient(
		cfg.HubAddr, cfg.Token,
	)
	if dialErr != nil {
		logWarn.Warn(cfgWarn.HubSyncDial, cfg.HubAddr, dialErr)
		return ""
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			logWarn.Warn(cfgWarn.CloseHubClient, cerr)
		}
	}()

	// Bound the pull with a deadline so a hub that accepts the
	// connection but never responds (hung server, black-hole
	// proxy) cannot stall session start: the RPC has no inherent
	// timeout, and this hook must never block. An exceeded
	// deadline surfaces as a warning like any other sync failure;
	// the next session retries.
	ctx, cancel := context.WithTimeout(
		context.Background(), syncTimeout,
	)
	defer cancel()

	entries, syncErr := client.Sync(ctx, cfg.Types, 0)
	if syncErr != nil {
		logWarn.Warn(cfgWarn.HubSyncPull, cfg.HubAddr, syncErr)
		return ""
	}
	if len(entries) == 0 {
		// Genuine empty result: not an error, no warning.
		return ""
	}

	if writeErr := render.WriteEntries(entries); writeErr != nil {
		logWarn.Warn(cfgWarn.HubSyncWrite, len(entries), writeErr)
		return ""
	}

	return fmt.Sprintf(
		desc.Text(text.DescKeyWriteConnectHubSync),
		len(entries),
	)
}
