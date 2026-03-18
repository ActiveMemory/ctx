//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"github.com/ActiveMemory/ctx/internal/config/pad"
	crypto2 "github.com/ActiveMemory/ctx/internal/err/crypto"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/pad"
	pad2 "github.com/ActiveMemory/ctx/internal/write/pad"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run reads and prints both sides of a merge conflict.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if no conflict files found or decryption fails
func Run(cmd *cobra.Command) error {
	if !rc.ScratchpadEncrypt() {
		return ctxerr.ResolveNotEncrypted()
	}

	kp := core.KeyPath()
	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		return crypto2.LoadKey(loadErr, kp)
	}

	dir := rc.ContextDir()

	ours, errOurs := core.DecryptFile(key, dir, pad.Enc+".ours")
	theirs, errTheirs := core.DecryptFile(key, dir, pad.Enc+".theirs")

	if errOurs != nil && errTheirs != nil {
		return ctxerr.NoConflictFiles(pad.Enc)
	}

	if errOurs == nil {
		pad2.PadResolveSide(cmd, "OURS", displayAll(ours))
	}

	if errTheirs == nil {
		pad2.PadResolveSide(cmd, "THEIRS", displayAll(theirs))
	}

	return nil
}

// displayAll converts entries to their display form.
func displayAll(entries []string) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = core.DisplayEntry(e)
	}
	return out
}
