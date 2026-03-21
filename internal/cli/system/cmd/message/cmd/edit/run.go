//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package edit

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	hook "github.com/ActiveMemory/ctx/internal/assets/read/hook"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/err/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/hook"
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
		return core.ValidationError(hk, variant)
	}

	oPath := core.OverridePath(hk, variant)

	if _, statErr := os.Stat(oPath); statErr == nil {
		return ctxerr.OverrideExists(oPath, hk, variant)
	}

	if info.Category == messages.CategoryCtxSpecific {
		cmd.Println(desc.Text(text.DescKeyMessageCtxSpecificWarning))
		cmd.Println()
	}

	data, readErr := hook.Message(hk, variant+file.ExtTxt)
	if readErr != nil {
		return ctxerr.EmbeddedTemplateNotFound(hk, variant)
	}

	dir := filepath.Dir(oPath)
	if mkdirErr := os.MkdirAll(dir, 0o750); mkdirErr != nil {
		return fs.CreateDir(dir, mkdirErr)
	}

	if writeErr := os.WriteFile(oPath, data, 0o600); writeErr != nil {
		return ctxerr.WriteOverride(oPath, writeErr)
	}

	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyMessageOverrideCreated), oPath))
	cmd.Println(desc.Text(text.DescKeyMessageEditHint))
	core.PrintTemplateVars(cmd, info)

	return nil
}
