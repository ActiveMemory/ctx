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
	errInitialize "github.com/ActiveMemory/ctx/internal/err/initialize"
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
	content, err := makefile.Ctx()

	if err != nil {
		return errInitialize.ReadTemplate(project.MakefileCtx, err)
	}

	if err = os.WriteFile(project.MakefileCtx, content, fs.PermFile); err != nil {
		return errFs.FileWrite(project.MakefileCtx, err)
	}

	initialize.Created(cmd, project.MakefileCtx)

	existing, readErr := os.ReadFile(project.Makefile)
	if readErr != nil {
		minimal := project.MakefileIncludeDirective + token.NewlineLF
		if err := os.WriteFile(
			project.Makefile, []byte(minimal), fs.PermFile,
		); err != nil {
			return errInitialize.CreateMakefile(err)
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
	if err := os.WriteFile(
		project.Makefile, []byte(amended), fs.PermFile,
	); err != nil {
		return errFs.FileAmend(project.Makefile, err)
	}

	initialize.MakefileAppended(cmd, project.MakefileCtx)
	return nil
}
