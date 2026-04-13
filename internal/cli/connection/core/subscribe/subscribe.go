//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package subscribe

import (
	"github.com/spf13/cobra"

	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connection/core/config"
	writeConnect "github.com/ActiveMemory/ctx/internal/write/connect"
)

// Run updates the subscription types in the connection config.
//
// Parameters:
//   - cmd: cobra command for output
//   - args: entry types to subscribe to (cobra args)
//
// Returns:
//   - error: non-nil if config load or save fails
func Run(cmd *cobra.Command, args []string) error {
	types := args
	cfg, loadErr := connectCfg.Load()
	if loadErr != nil {
		return loadErr
	}

	cfg.Types = types
	if saveErr := connectCfg.Save(cfg); saveErr != nil {
		return saveErr
	}

	writeConnect.Subscribed(cmd, types)
	return nil
}
