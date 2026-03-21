//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mark_journal

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// StageChecked prints the result of a --check query. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: journal filename
//   - stage: processing stage name
//   - val: current stage value
func StageChecked(cmd *cobra.Command, filename, stage, val string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyMarkJournalChecked), filename, stage, val))
}

// StageMarked prints the confirmation after marking a stage. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: journal filename
//   - stage: processing stage name
func StageMarked(cmd *cobra.Command, filename, stage string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyMarkJournalMarked), filename, stage))
}
