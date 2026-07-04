//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// StatuslineDeployed reports the ctx status line wired into settings.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func StatuslineDeployed(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitStatuslineDeployed), path,
	))
}

// StatuslineBackedUp reports a displaced statusLine entry saved to the
// state directory.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: backup file path
func StatuslineBackedUp(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitStatuslineBackedUp), path,
	))
}

// StatuslineRestored reports a backed-up statusLine entry restored
// after statusline.enabled was set to false.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: backup file path the entry was restored from
func StatuslineRestored(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitStatuslineRestored), path,
	))
}

// StatuslineRemoved reports the ctx statusLine entry dropped after
// statusline.enabled was set to false with no backup to restore.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: settings file path
func StatuslineRemoved(cmd *cobra.Command, path string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitStatuslineRemoved), path,
	))
}
