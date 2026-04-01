//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"github.com/spf13/cobra"

	readClaude "github.com/ActiveMemory/ctx/internal/assets/read/claude"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/merge"
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/entity"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// HandleMd creates or merges CLAUDE.md with ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - force: If true, overwrite the existing ctx section
//   - autoMerge: If true, skip interactive confirmation
//
// Returns:
//   - error: Non-nil if file operations fail
func HandleMd(cmd *cobra.Command, force, autoMerge bool) error {
	templateContent, err := readClaude.Md()
	if err != nil {
		return errInit.ReadTemplate(claude.Md, err)
	}

	created, mergeErr := merge.OrCreate(cmd, entity.MergeParams{
		Filename:        claude.Md,
		MarkerStart:     marker.CtxStart,
		MarkerEnd:       marker.CtxEnd,
		TemplateContent: templateContent,
		Force:           force,
		AutoMerge:       autoMerge,
		ConfirmPrompt:   desc.Text(text.DescKeyInitConfirmClaude),
		UpdateTextKey:   text.DescKeyWriteInitUpdatedCtxSection,
	})

	if mergeErr != nil {
		return mergeErr
	}

	if created {
		initialize.Created(cmd, claude.Md)
	}

	return nil
}
