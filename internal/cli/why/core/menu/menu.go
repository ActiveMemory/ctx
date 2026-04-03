//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package menu

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/why/core/data"
	"github.com/ActiveMemory/ctx/internal/cli/why/core/show"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errCli "github.com/ActiveMemory/ctx/internal/err/cli"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/why"
)

// Show presents a numbered menu and reads user selection
// from stdin.
//
// Parameters:
//   - cmd: Cobra command for output and context
//
// Returns:
//   - error: Non-nil on read failure or invalid selection
func Show(cmd *cobra.Command) error {
	why.Banner(cmd)
	why.Separator(cmd)
	for i, doc := range data.DocOrder {
		why.MenuItem(cmd, i+1, doc.Label)
	}
	why.MenuPrompt(cmd)

	reader := bufio.NewReader(os.Stdin)
	input, readErr := reader.ReadString(token.NewlineLF[0])
	if readErr != nil {
		return errFs.ReadInput(readErr)
	}

	input = strings.TrimSpace(input)
	choice, parseErr := strconv.Atoi(input)
	if parseErr != nil ||
		choice < 1 || choice > len(data.DocOrder) {
		return errCli.InvalidSelection(
			input, len(data.DocOrder),
		)
	}

	why.Separator(cmd)
	return show.Doc(cmd, data.DocOrder[choice-1].Alias)
}
