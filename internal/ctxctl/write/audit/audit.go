//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package audit

import (
	"fmt"

	"github.com/spf13/cobra"
)

// None prints the message shown when `ctxctl audit list` runs
// against an empty audit directory.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - msg: the "no reports" message (supplied by the caller;
//     ctxctl owns its user-facing text)
func None(cmd *cobra.Command, msg string) {
	if cmd == nil {
		return
	}
	cmd.Println(msg)
}

// ListItem prints one row of `ctxctl audit list` output.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - format: row format string with four verbs
//     (id, status, commit-range, generated-at)
//   - id: report basename (e.g. "surface")
//   - status: report status (findings | clean | dismissed)
//   - commitRange: commit-range header from the report
//   - generatedAt: pre-formatted timestamp string
func ListItem(
	cmd *cobra.Command,
	format, id, status, commitRange, generatedAt string,
) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		format, id, status, commitRange, generatedAt,
	))
}

// Dismissed prints confirmation for a single dismissal.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - format: confirmation format string with one verb (id)
//   - id: dismissed report id
func Dismissed(cmd *cobra.Command, format, id string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(format, id))
}

// DismissedAll prints confirmation for a bulk dismissal.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - format: confirmation format string with one verb (count)
//   - count: number of reports dismissed
func DismissedAll(cmd *cobra.Command, format string, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(format, count))
}

// Body prints an audit report's body verbatim. Used by
// `ctxctl audit show`.
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
