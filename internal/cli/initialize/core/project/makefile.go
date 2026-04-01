//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package project

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/makefile"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/project"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// HandleMakefileCtx deploys Makefile.ctx and amends the user Makefile.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if file operations fail
func HandleMakefileCtx(cmd *cobra.Command) error {
	content, tplErr := makefile.Ctx()

	if tplErr != nil {
		return errInit.ReadTemplate(project.MakefileCtx, tplErr)
	}

	if writeErr := os.WriteFile(
		project.MakefileCtx, content, fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(project.MakefileCtx, writeErr)
	}

	initialize.Created(cmd, project.MakefileCtx)

	existing, readErr := os.ReadFile(project.Makefile)
	if readErr != nil {
		minimal := project.MakefileIncludeDirective + token.NewlineLF
		if writeErr := os.WriteFile(
			project.Makefile, []byte(minimal), fs.PermFile,
		); writeErr != nil {
			return errInit.CreateMakefile(writeErr)
		}
		initialize.MakefileCreated(cmd)
		return nil
	}

	if strings.Contains(string(existing), project.MakefileIncludeDirective) {
		initialize.MakefileIncludes(cmd, project.MakefileCtx)
		return nil
	}

	amended := string(existing)
	if !strings.HasSuffix(amended, token.NewlineLF) {
		amended += token.NewlineLF
	}

	amended += token.NewlineLF + project.MakefileIncludeDirective + token.NewlineLF
	if writeErr := os.WriteFile( //nolint:gosec // path built from trusted project root
		project.Makefile, []byte(amended), fs.PermFile,
	); writeErr != nil {
		return errFs.FileAmend(project.Makefile, writeErr)
	}

	initialize.MakefileAppended(cmd, project.MakefileCtx)
	return nil
}
