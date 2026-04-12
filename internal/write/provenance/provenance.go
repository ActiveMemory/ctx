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
	"github.com/ActiveMemory/ctx/internal/config/stats"
)

// Line prints a single provenance line to cmd output with an
// optional context-free percentage suffix.
//
// Parameters:
//   - cmd: Cobra command for output
//   - session: Short session ID
//   - branch: Git branch name
//   - commit: Short commit hash
//   - contextSuffix: Pre-formatted suffix (e.g. " | Context: 45% free")
//     or empty when token data is unavailable
func Line(
	cmd *cobra.Command,
	session, branch, commit, contextSuffix string,
) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteProvenanceLine),
		session, branch, commit, contextSuffix,
	))
}

// ContextSuffix formats the context-free percentage suffix.
// Returns an empty string when freePct is zero or negative.
//
// Parameters:
//   - freePct: Context-free percentage (0-100)
//
// Returns:
//   - string: Pre-formatted suffix including leading separator,
//     or empty string when freePct is out of range
func ContextSuffix(freePct int) string {
	if freePct <= 0 || freePct > stats.PercentMultiplier {
		return ""
	}
	return fmt.Sprintf(
		desc.Text(text.DescKeyWriteProvenanceContext),
		freePct,
	)
}
