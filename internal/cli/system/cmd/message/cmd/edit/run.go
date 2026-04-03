//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package edit

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/assets/read/hook"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/err/fs"
	errTrigger "github.com/ActiveMemory/ctx/internal/err/trigger"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeMessage "github.com/ActiveMemory/ctx/internal/write/message"
)

// Run executes the message edit logic.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hk: Hook name
//   - variant: Template variant name
//
// Returns:
//   - error: Non-nil if the hook/variant is unknown, override exists,
//     or file operations fail
func Run(cmd *cobra.Command, hk, variant string) error {
	info := messages.Lookup(hk, variant)
	if info == nil {
		return errTrigger.Validate(messages.Variants(hk) != nil, hk, variant)
	}

	oPath := message.OverridePath(hk, variant)

	if _, statErr := os.Stat(oPath); statErr == nil {
		return errTrigger.OverrideExists(oPath, hk, variant)
	}

	if info.Category == messages.CategoryCtxSpecific {
		writeMessage.CtxSpecificWarning(cmd)
	}

	data, readErr := hook.Message(hk, variant+file.ExtTxt)
	if readErr != nil {
		return errTrigger.EmbeddedTemplateNotFound(hk, variant)
	}

	dir := filepath.Dir(oPath)
	if mkdirErr := ctxIo.SafeMkdirAll(dir, 0o750); mkdirErr != nil {
		return fs.CreateDir(dir, mkdirErr)
	}

	if writeErr := ctxIo.SafeWriteFile(oPath, data, 0o600); writeErr != nil {
		return errTrigger.WriteOverride(oPath, writeErr)
	}

	writeMessage.OverrideCreated(cmd, oPath)
	writeMessage.EditHint(cmd)
	writeMessage.TemplateVars(cmd, message.FormatTemplateVars(info))

	return nil
}
