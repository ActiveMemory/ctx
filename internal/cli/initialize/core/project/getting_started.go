//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgProject "github.com/ActiveMemory/ctx/internal/config/project"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// WriteGettingStarted saves the next-steps and workflow-tips text to
// GETTING_STARTED.md in the project root. Best-effort: failures are
// non-fatal since the same content was already printed to stdout.
//
// Parameters:
//   - cmd: Cobra command for status output
func WriteGettingStarted(cmd *cobra.Command) {
	content := desc.Text(text.DescKeyWriteInitNextStepsBlock) +
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
