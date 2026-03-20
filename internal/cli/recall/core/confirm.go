//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"bufio"
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/export"
	"github.com/spf13/cobra"
)

// ConfirmExport prints the plan summary and prompts for confirmation.
//
// Parameters:
//   - cmd: Cobra command for output.
//   - plan: the export plan to summarize.
//
// Returns:
//   - bool: true if the user confirms.
//   - error: non-nil if reading input fails.
func ConfirmExport(cmd *cobra.Command, plan ExportPlan) (bool, error) {
	export.Summary(cmd, plan.NewCount, plan.RegenCount, plan.SkipCount, plan.LockedCount, false)
	cmd.Print(desc.TextDesc(text.DescKeyConfirmProceed))
	reader := bufio.NewReader(os.Stdin)
	response, readErr := reader.ReadString(token.NewlineLF[0])
	if readErr != nil {
		return false, ctxerr.ReadInput(readErr)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == cli.ConfirmShort || response == cli.ConfirmLong, nil
}
