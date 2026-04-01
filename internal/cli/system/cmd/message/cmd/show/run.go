//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/hooks/messages"
	readHook "github.com/ActiveMemory/ctx/internal/assets/read/hook"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/config/file"
	errHook "github.com/ActiveMemory/ctx/internal/err/hook"
	"github.com/ActiveMemory/ctx/internal/io"
	writeMessage "github.com/ActiveMemory/ctx/internal/write/message"
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
		return errHook.Validate(messages.Variants(hk) != nil, hk, variant)
	}

	oPath := message.OverridePath(hk, variant)
	if data, readErr := io.SafeReadUserFile(oPath); readErr == nil {
		writeMessage.SourceOverride(cmd, oPath)
		writeMessage.TemplateVars(cmd, message.FormatTemplateVars(info))
		writeMessage.ContentBlock(cmd, data)
		return nil
	}

	data, readErr := readHook.Message(hk, variant+file.ExtTxt)
	if readErr != nil {
		return errHook.EmbeddedTemplateNotFound(hk, variant)
	}

	writeMessage.SourceDefault(cmd)
	writeMessage.TemplateVars(cmd, message.FormatTemplateVars(info))
	writeMessage.ContentBlock(cmd, data)
	return nil
}
