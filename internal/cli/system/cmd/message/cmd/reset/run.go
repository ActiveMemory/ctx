//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reset

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/hook"
	writeMessage "github.com/ActiveMemory/ctx/internal/write/message"
)

// Run executes the message reset logic.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hk: Hook name
//   - variant: Template variant name
//
// Returns:
//   - error: Non-nil if the hook/variant is unknown or removal fails
func Run(cmd *cobra.Command, hk, variant string) error {
	info := messages.Lookup(hk, variant)
	if info == nil {
		return core.ValidationError(hk, variant)
	}

	oPath := core.OverridePath(hk, variant)

	if removeErr := os.Remove(oPath); removeErr != nil {
		if os.IsNotExist(removeErr) {
			writeMessage.NoOverride(cmd, hk, variant)
			return nil
		}
		return ctxerr.RemoveOverride(oPath, removeErr)
	}

	hookDir := filepath.Dir(oPath)
	_ = os.Remove(hookDir)
	messagesDir := filepath.Dir(hookDir)
	_ = os.Remove(messagesDir)

	writeMessage.OverrideRemoved(cmd, hk, variant)
	return nil
}
