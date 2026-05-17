//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ground

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	kbPath "github.com/ActiveMemory/ctx/internal/cli/kb/core/path"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errKbCli "github.com/ActiveMemory/ctx/internal/err/kb/cli"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Run validates the grounding-sources file and prints the
// canonical `/ctx-kb-ground` skill invocation. The CLI
// surface itself does not perform the re-grounding pass;
// the skill does.
//
// Parameters:
//   - cobraCmd: cobra command for output.
//
// Returns:
//   - error: refusal when `grounding-sources.md` is missing
//     or empty.
func Run(cobraCmd *cobra.Command) error {
	gPath, pathErr := kbPath.IngestArtifactFile(
		cfgKB.GroundingSources,
	)
	if pathErr != nil {
		return pathErr
	}
	info, statErr := io.SafeStat(gPath)
	if errors.Is(statErr, os.ErrNotExist) {
		cobraCmd.SilenceUsage = true
		return errKbCli.GroundingMissing(gPath)
	}
	if statErr != nil {
		return statErr
	}
	if info.Size() == 0 {
		cobraCmd.SilenceUsage = true
		return errKbCli.GroundingEmpty(gPath)
	}
	io.SafeFprintf(
		cobraCmd.OutOrStdout(), token.FormatString,
		desc.Text(text.DescKeyWriteKbCliGroundDrivenHint),
	)
	io.SafeFprintf(
		cobraCmd.OutOrStdout(), token.FormatString,
		desc.Text(text.DescKeyWriteKbCliGroundContractPointer),
	)
	return nil
}
