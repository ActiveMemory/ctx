//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// Watching prints the initial "watching" status line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func Watching(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWatchWatching))
}

// DryRun prints the dry-run notice.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func DryRun(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWatchDryRun))
}

// StopHint prints the Ctrl+C stop hint.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func StopHint(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWatchStopHint))
}

// CloseLogError prints a log file close error.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - err: the close error.
func CloseLogError(cmd *cobra.Command, err error) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWatchCloseLogError), err),
	)
}

// DryRunPreview prints a dry-run preview of an update that would be applied.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - updateType: the context update type.
//   - content: the update content.
func DryRunPreview(cmd *cobra.Command, updateType, content string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			desc.Text(text.DescKeyWatchDryRunPreview),
			updateType, content,
		),
	)
}

// ApplyFailed prints a failure message for an update that could not be applied.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - updateType: the context update type.
//   - err: the apply error.
func ApplyFailed(cmd *cobra.Command, updateType string, err error) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			desc.Text(text.DescKeyWatchApplyFailed), updateType, err,
		),
	)
}

// ApplySuccess prints a success message for an applied update.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - updateType: the context update type.
//   - content: the update content.
func ApplySuccess(cmd *cobra.Command, updateType, content string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWatchApplySuccess), updateType, content),
	)
}
