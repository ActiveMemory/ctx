//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package confirm

import (
	"bufio"
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/export"
	"github.com/ActiveMemory/ctx/internal/write/recall"
	"github.com/spf13/cobra"
)

// Export prints the plan summary and prompts for confirmation.
//
// Parameters:
//   - cmd: Cobra command for output.
//   - plan: the export plan to summarize.
//
// Returns:
//   - bool: true if the user confirms.
//   - error: non-nil if reading input fails.
func Export(cmd *cobra.Command, plan entity.ExportPlan) (bool, error) {
	export.Summary(
		cmd,
		plan.NewCount, plan.RegenCount, plan.SkipCount, plan.LockedCount, false,
	)
	recall.ConfirmPrompt(cmd)
	reader := bufio.NewReader(os.Stdin)
	response, readErr := reader.ReadString(token.NewlineLF[0])
	if readErr != nil {
		return false, errFs.ReadInput(readErr)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == cli.ConfirmShort || response == cli.ConfirmLong, nil
}
