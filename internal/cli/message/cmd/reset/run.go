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
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	errTrigger "github.com/ActiveMemory/ctx/internal/err/trigger"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
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
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}
	info := messages.Lookup(hk, variant)
	if info == nil {
		return errTrigger.Validate(messages.Variants(hk) != nil, hk, variant)
	}

	oPath, pathErr := message.OverridePath(hk, variant)
	if pathErr != nil {
		return pathErr
	}

	if removeErr := os.Remove(oPath); removeErr != nil {
		if os.IsNotExist(removeErr) {
			writeMessage.NoOverride(cmd, hk, variant)
			return nil
		}
		return errTrigger.RemoveOverride(oPath, removeErr)
	}

	hookDir := filepath.Dir(oPath)
	if removeErr := os.Remove(hookDir); removeErr != nil {
		ctxLog.Warn(warn.Remove, hookDir, removeErr)
	}
	messagesDir := filepath.Dir(hookDir)
	if removeErr := os.Remove(messagesDir); removeErr != nil {
		ctxLog.Warn(
			warn.Remove, messagesDir, removeErr,
		)
	}

	writeMessage.OverrideRemoved(cmd, hk, variant)
	return nil
}
