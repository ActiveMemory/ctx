//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ingest

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errKbCli "github.com/ActiveMemory/ctx/internal/err/kb/cli"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Run validates non-empty source args and prints the
// canonical `/ctx-kb-ingest` skill invocation. The CLI
// surface itself does not perform the editorial pass; the
// skill does.
//
// Parameters:
//   - cobraCmd: cobra command for output.
//   - args: positional source arguments (paths, URLs, or
//     inline gestures). Empty is rejected.
//
// Returns:
//   - error: refusal when no sources were supplied.
func Run(cobraCmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		cobraCmd.SilenceUsage = true
		return errKbCli.ErrIngestNoSources
	}
	io.SafeFprintf(
		cobraCmd.OutOrStdout(), token.FormatString,
		desc.Text(text.DescKeyWriteKbIngestDrivenHint),
	)
	io.SafeFprintf(
		cobraCmd.OutOrStdout(),
		desc.Text(text.DescKeyWriteKbIngestInvokeFormat),
		strings.Join(args, token.Space),
	)
	io.SafeFprintf(
		cobraCmd.OutOrStdout(), token.FormatString,
		desc.Text(text.DescKeyWriteKbIngestFallbackHint),
	)
	return nil
}
