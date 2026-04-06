//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package provenance

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Line prints a single provenance line to cmd output.
//
// Parameters:
//   - cmd: Cobra command for output
//   - session: Short session ID
//   - branch: Git branch name
//   - commit: Short commit hash
func Line(
	cmd *cobra.Command,
	session, branch, commit string,
) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteProvenanceLine),
		session, branch, commit,
	))
}
