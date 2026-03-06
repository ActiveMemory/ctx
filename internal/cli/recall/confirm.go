//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"bufio"
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
	"github.com/spf13/cobra"
)

// confirmExport prints the plan summary and prompts for confirmation.
//
// Parameters:
//   - cmd: Cobra command for output.
//   - plan: the export plan to summarize.
//
// Returns:
//   - bool: true if the user confirms.
//   - error: non-nil if reading input fails.
func confirmExport(cmd *cobra.Command, plan exportPlan) (bool, error) {
	write.ExportSummary(cmd, planCounts(plan), false)
	cmd.Print("Proceed? [y/N] ")
	reader := bufio.NewReader(os.Stdin)
	response, readErr := reader.ReadString('\n')
	if readErr != nil {
		return false, ctxerr.ReadInput(readErr)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == config.ConfirmShort || response == config.ConfirmLong, nil
}
