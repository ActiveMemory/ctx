//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	padCrypto "github.com/ActiveMemory/ctx/internal/cli/pad/core/crypto"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/pad"
	"github.com/ActiveMemory/ctx/internal/crypto"
	errCrypto "github.com/ActiveMemory/ctx/internal/err/crypto"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
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
		return errPad.ResolveNotEncrypted()
	}

	kp := store.KeyPath()
	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		return errCrypto.LoadKey(loadErr, kp)
	}

	dir := rc.ContextDir()

	ours, errOurs := padCrypto.DecryptFile(key, dir, pad.EncOurs)
	theirs, errTheirs := padCrypto.DecryptFile(key, dir, pad.EncTheirs)

	if errOurs != nil && errTheirs != nil {
		return errPad.NoConflictFiles(pad.Enc)
	}

	if errOurs == nil {
		writePad.ResolveSide(cmd, pad.SideOurs, displayAll(ours))
	}

	if errTheirs == nil {
		writePad.ResolveSide(cmd, pad.SideTheirs, displayAll(theirs))
	}

	return nil
}
