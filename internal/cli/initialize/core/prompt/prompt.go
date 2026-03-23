//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	readLoop "github.com/ActiveMemory/ctx/internal/assets/read/loop"
	"github.com/ActiveMemory/ctx/internal/assets/read/template"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/merge"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/loop"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/entity"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// HandlePromptMd creates or merges PROMPT.md with ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - force: If true, overwrite the existing prompt section
//   - autoMerge: If true, skip interactive confirmation
//   - ralph: If true, use ralph mode template
//
// Returns:
//   - error: Non-nil if file operations fail
func HandlePromptMd(cmd *cobra.Command, force, autoMerge, ralph bool) error {
	var templateContent []byte
	var readErr error
	if ralph {
		templateContent, readErr = readLoop.RalphTemplate(loop.PromptMd)
	} else {
		templateContent, readErr = template.Template(loop.PromptMd)
	}
	if readErr != nil {
		return errInit.ReadTemplate(loop.PromptMd, readErr)
	}

	created, mergeErr := merge.CreateOrMerge(cmd, entity.MergeParams{
		Filename:        loop.PromptMd,
		MarkerStart:     marker.PromptMarkerStart,
		MarkerEnd:       marker.PromptMarkerEnd,
		TemplateContent: templateContent,
		Force:           force,
		AutoMerge:       autoMerge,
		ConfirmPrompt:   desc.Text(text.DescKeyInitConfirmPrompt),
		UpdateTextKey:   text.DescKeyWriteInitUpdatedPromptSection,
	})
	if mergeErr != nil {
		return mergeErr
	}
	if created {
		mode := ""
		if ralph {
			mode = desc.Text(text.DescKeyInitRalphMode)
		}
		initialize.CreatedWith(cmd, loop.PromptMd, mode)
	}
	return nil
}
