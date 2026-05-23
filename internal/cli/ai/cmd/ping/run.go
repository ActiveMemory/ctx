//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ping

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/ai/core/resolve"
	writeAI "github.com/ActiveMemory/ctx/internal/write/ai"
)

// Run executes `ctx ai ping`: builds the backend
// registry, picks the named or default backend, and
// issues Ping. Fails closed on unconfigured registry or
// unreachable backend.
//
// Parameters:
//   - cobraCmd: cobra command for output.
//   - backendName: --backend value; empty falls back to
//     the registry's default.
//
// Returns:
//   - error: typed err/backend sentinel on failure
//     (no backends configured, ambiguous default,
//     unreachable, etc.).
func Run(cobraCmd *cobra.Command, backendName string) error {
	r, buildErr := resolve.Build()
	if buildErr != nil {
		return buildErr
	}
	b, pickErr := resolve.Pick(r, backendName)
	if pickErr != nil {
		return pickErr
	}
	if pingErr := b.Ping(cobraCmd.Context()); pingErr != nil {
		return pingErr
	}
	writeAI.InfoPingOK(cobraCmd, b.Name())
	return nil
}
