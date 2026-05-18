//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ask

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errKbCli "github.com/ActiveMemory/ctx/internal/err/kb/cli"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Run validates a non-empty question and prints the canonical
// `/ctx-kb-ask` skill invocation. The CLI surface itself does
// not perform the Q&A pass; the skill does.
//
// Parameters:
//   - cobraCmd: cobra command for output.
//   - question: trimmed question text (empty is rejected).
//   - args: raw positional args, joined for the printed
//     skill invocation.
//
// Returns:
//   - error: refusal on empty question.
func Run(cobraCmd *cobra.Command, question string, args []string) error {
	if question == "" {
		cobraCmd.SilenceUsage = true
		return errKbCli.ErrAskNoQuestion
	}
	io.SafeFprintf(
		cobraCmd.OutOrStdout(), token.FormatString,
		desc.Text(text.DescKeyWriteKbAskDrivenHint),
	)
	io.SafeFprintf(
		cobraCmd.OutOrStdout(),
		desc.Text(text.DescKeyWriteKbAskInvokeFormat),
		strings.Join(args, token.Space),
	)
	io.SafeFprintf(
		cobraCmd.OutOrStdout(), token.FormatString,
		desc.Text(text.DescKeyWriteKbAskContractPointer),
	)
	return nil
}
