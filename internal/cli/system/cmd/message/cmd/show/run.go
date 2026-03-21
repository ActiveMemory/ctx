//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	hook "github.com/ActiveMemory/ctx/internal/assets/read/hook"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/hook"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Run executes the message show logic.
//
// Parameters:
//   - cmd: Cobra command for output
//   - hk: Hook name
//   - variant: Template variant name
//
// Returns:
//   - error: Non-nil if the hook/variant is unknown or template is missing
func Run(cmd *cobra.Command, hk, variant string) error {
	info := messages.Lookup(hk, variant)
	if info == nil {
		return core.ValidationError(hk, variant)
	}

	oPath := core.OverridePath(hk, variant)
	if data, readErr := io.SafeReadUserFile(oPath); readErr == nil {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyMessageSourceOverride), oPath))
		core.PrintTemplateVars(cmd, info)
		cmd.Println()
		cmd.Print(string(data))
		if len(data) > 0 && data[len(data)-1] != '\n' {
			cmd.Println()
		}
		return nil
	}

	data, readErr := hook.Message(hk, variant+file.ExtTxt)
	if readErr != nil {
		return ctxerr.EmbeddedTemplateNotFound(hk, variant)
	}

	cmd.Println(desc.Text(text.DescKeyMessageSourceDefault))
	core.PrintTemplateVars(cmd, info)
	cmd.Println()
	cmd.Print(string(data))
	if len(data) > 0 && data[len(data)-1] != '\n' {
		cmd.Println()
	}
	return nil
}
