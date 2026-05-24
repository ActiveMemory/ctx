//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// None prints the message shown when `ctx audit list` runs
// against an empty audit directory.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func None(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWriteAuditNone))
}

// ListItem prints one row of `ctx audit list` output.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - id: report basename (e.g. "surface")
//   - status: report status (findings | clean | dismissed)
//   - commitRange: commit-range header from the report
//   - generatedAt: pre-formatted timestamp string
func ListItem(
	cmd *cobra.Command,
	id, status, commitRange, generatedAt string,
) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteAuditListItem),
		id, status, commitRange, generatedAt,
	))
}

// Dismissed prints confirmation for a single dismissal.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - id: dismissed report id
func Dismissed(cmd *cobra.Command, id string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteAuditDismissed), id,
	))
}

// DismissedAll prints confirmation for a bulk dismissal.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of reports dismissed
func DismissedAll(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteAuditDismissedAll), count,
	))
}

// Body prints an audit report's body verbatim. Used by
// `ctx audit show`.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - body: report body, already free of frontmatter
func Body(cmd *cobra.Command, body string) {
	if cmd == nil {
		return
	}
	cmd.Println(body)
}
