//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/crypto"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// runResolve reads and prints both sides of a merge conflict.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if no conflict files found or decryption fails
func runResolve(cmd *cobra.Command) error {
	if !rc.ScratchpadEncrypt() {
		return errors.New("resolve is only needed for encrypted scratchpads")
	}

	kp := core.KeyPath()
	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		return ctxerr.LoadKey(loadErr, kp)
	}

	dir := rc.ContextDir()

	ours, errOurs := core.DecryptFile(key, dir, config.FileScratchpadEnc+".ours")
	theirs, errTheirs := core.DecryptFile(key, dir, config.FileScratchpadEnc+".theirs")

	if errOurs != nil && errTheirs != nil {
		return fmt.Errorf("no conflict files found (%s.ours / %s.theirs)",
			config.FileScratchpadEnc, config.FileScratchpadEnc)
	}

	if errOurs == nil {
		cmd.Println("=== OURS ===")
		for i, entry := range ours {
			cmd.Println(fmt.Sprintf("  %d. %s", i+1, core.DisplayEntry(entry)))
		}
	}

	if errTheirs == nil {
		cmd.Println("=== THEIRS ===")
		for i, entry := range theirs {
			cmd.Println(fmt.Sprintf("  %d. %s", i+1, core.DisplayEntry(entry)))
		}
	}

	return nil
}
