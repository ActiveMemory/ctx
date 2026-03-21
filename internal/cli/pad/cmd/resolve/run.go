//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
	"github.com/ActiveMemory/ctx/internal/config/pad"
	"github.com/ActiveMemory/ctx/internal/crypto"
	errCrypto "github.com/ActiveMemory/ctx/internal/err/crypto"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/pad"
	"github.com/ActiveMemory/ctx/internal/rc"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
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
		return ctxErr.ResolveNotEncrypted()
	}

	kp := core.KeyPath()
	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		return errCrypto.LoadKey(loadErr, kp)
	}

	dir := rc.ContextDir()

	ours, errOurs := core.DecryptFile(key, dir, pad.EncOurs)
	theirs, errTheirs := core.DecryptFile(key, dir, pad.EncTheirs)

	if errOurs != nil && errTheirs != nil {
		return ctxErr.NoConflictFiles(pad.Enc)
	}

	if errOurs == nil {
		writePad.PadResolveSide(cmd, pad.SideOurs, displayAll(ours))
	}

	if errTheirs == nil {
		writePad.PadResolveSide(cmd, pad.SideTheirs, displayAll(theirs))
	}

	return nil
}
