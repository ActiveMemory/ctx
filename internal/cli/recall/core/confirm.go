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

	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
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
	write.ExportSummary(cmd, plan.NewCount, plan.RegenCount, plan.SkipCount, plan.LockedCount, false)
	cmd.Print(assets.TextDesc(assets.TextDescKeyConfirmProceed))
	reader := bufio.NewReader(os.Stdin)
	response, readErr := reader.ReadString(token.NewlineLF[0])
	if readErr != nil {
		return false, ctxerr.ReadInput(readErr)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == cli.ConfirmShort || response == cli.ConfirmLong, nil
}
