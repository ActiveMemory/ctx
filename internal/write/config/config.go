//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/file"
)

// ProfileStatus prints the active profile status line.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - profile: Detected profile name (dev, base, or empty)
func ProfileStatus(cmd *cobra.Command, profile string) {
	if cmd == nil {
		return
	}
	switch profile {
	case file.ProfileDev:
		cmd.Println(desc.Text(text.DescKeyWriteConfigProfileDev))
	case file.ProfileBase:
		cmd.Println(desc.Text(text.DescKeyWriteConfigProfileBase))
	default:
		cmd.Println(fmt.Sprintf(
			desc.Text(text.DescKeyWriteConfigProfileNone),
			file.CtxRC,
		))
	}
}

// Schema prints the raw JSON schema content.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - data: Raw schema bytes to display
func Schema(cmd *cobra.Command, data []byte) {
	if cmd == nil {
		return
	}
	cmd.Print(string(data))
}

// SwitchConfirm prints the profile switch confirmation message.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - msg: Confirmation message from the switch operation
func SwitchConfirm(cmd *cobra.Command, msg string) {
	if cmd == nil {
		return
	}
	cmd.Println(msg)
}
