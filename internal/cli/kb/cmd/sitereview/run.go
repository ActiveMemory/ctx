//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sitereview

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Run prints the canonical `/ctx-kb-site-review` skill
// invocation. The CLI surface itself does not perform the
// structural audit; the skill does.
//
// Parameters:
//   - cobraCmd: cobra command for output.
//
// Returns:
//   - error: always nil (the pointer is informational; the
//     skill carries the refusal contract).
func Run(cobraCmd *cobra.Command) error {
	io.SafeFprintf(
		cobraCmd.OutOrStdout(), token.FormatString,
		desc.Text(text.DescKeyWriteKbSiteReviewDrivenHint),
	)
	io.SafeFprintf(
		cobraCmd.OutOrStdout(), token.FormatString,
		desc.Text(text.DescKeyWriteKbSiteReviewContractPointer),
	)
	return nil
}
