//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package project

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgProject "github.com/ActiveMemory/ctx/internal/config/project"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// WriteGettingStarted saves an anatomy preamble (what `.context/`
// is and how the project-root contract works), the activation hint,
// next-steps, and workflow-tips text to GETTING_STARTED.md in the
// project root. The file is the human's durable primer after
// running `ctx init`: the preamble names the contract so future
// readers know which directory rule is load-bearing; the activation
// hint comes next because every subsequent `ctx <command>`
// requires CTX_DIR to be declared. Best-effort: failures are
// non-fatal since the activation hint and next-steps were already
// printed to stdout.
//
// Parameters:
//   - cmd:        Cobra command for status output.
//   - contextDir: Absolute path of the just-created .context/
//     directory, used in the activation hint.
func WriteGettingStarted(cmd *cobra.Command, contextDir string) {
	activateHint := fmt.Sprintf(
		desc.Text(text.DescKeyWriteInitActivateHint),
		contextDir,
	)
	content := desc.Text(text.DescKeyWriteInitAnatomyPreamble) +
		token.NewlineLF +
		activateHint +
		token.NewlineLF +
		desc.Text(text.DescKeyWriteInitNextStepsBlock) +
		token.NewlineLF +
		desc.Text(text.DescKeyWriteInitWorkflowTips) +
		token.NewlineLF
	if writeErr := ctxIo.SafeWriteFile(
		cfgProject.GettingStarted, []byte(content), fs.PermFile,
	); writeErr != nil {
		return
	}
	initialize.InfoGettingStartedSaved(cmd, cfgProject.GettingStarted)
}
