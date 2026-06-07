//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dispose

import (
	"github.com/spf13/cobra"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	engine "github.com/ActiveMemory/ctx/internal/dream"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
	writeDream "github.com/ActiveMemory/ctx/internal/write/dream"
)

// Accept loads the proposal with id from the latest run and applies
// its recommended action through the engine, then prints the
// disposition.
//
// Parameters:
//   - cmd: cobra command for output
//   - id: the proposal ID
//   - note: optional human note
//
// Returns:
//   - error: a resolution, not-found, guard, mutation, or ledger
//     failure
func Accept(cmd *cobra.Command, id, note string) error {
	loc, p, loadErr := load(id)
	if loadErr != nil {
		cmd.SilenceUsage = true
		return loadErr
	}
	res, applyErr := engine.Accept(loc.Root, loc.Dreams, p, note)
	if applyErr != nil {
		cmd.SilenceUsage = true
		return applyErr
	}
	writeDream.Disposition(cmd, id, cfgDream.DecisionAccepted, res)
	return nil
}

// Reject loads the proposal with id from the latest run and records a
// rejection (no mutation), then prints the disposition.
//
// Parameters:
//   - cmd: cobra command for output
//   - id: the proposal ID
//   - note: optional human note
//
// Returns:
//   - error: a resolution, not-found, or ledger failure
func Reject(cmd *cobra.Command, id, note string) error {
	loc, p, loadErr := load(id)
	if loadErr != nil {
		cmd.SilenceUsage = true
		return loadErr
	}
	res, applyErr := engine.Reject(loc.Dreams, p, note)
	if applyErr != nil {
		cmd.SilenceUsage = true
		return applyErr
	}
	writeDream.Disposition(cmd, id, cfgDream.DecisionRejected, res)
	return nil
}

// Amend loads the proposal with id from the latest run and applies
// action in place of its recommendation, then prints the disposition.
//
// Parameters:
//   - cmd: cobra command for output
//   - id: the proposal ID
//   - action: the action to apply instead of the recommendation
//   - note: optional human note
//
// Returns:
//   - error: a resolution, not-found, unknown-action, guard,
//     mutation, or ledger failure
func Amend(cmd *cobra.Command, id, action, note string) error {
	if action == "" {
		cmd.SilenceUsage = true
		return errDream.UnknownAction(action, id)
	}
	loc, p, loadErr := load(id)
	if loadErr != nil {
		cmd.SilenceUsage = true
		return loadErr
	}
	res, applyErr := engine.Amend(loc.Root, loc.Dreams, p, action, note)
	if applyErr != nil {
		cmd.SilenceUsage = true
		return applyErr
	}
	writeDream.Disposition(cmd, id, cfgDream.DecisionAmended, res)
	return nil
}
