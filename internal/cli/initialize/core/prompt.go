//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	loop2 "github.com/ActiveMemory/ctx/internal/assets/read/loop"
	"github.com/ActiveMemory/ctx/internal/assets/read/template"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/loop"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	initErr "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// HandlePromptMd creates or merges PROMPT.md with ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - force: If true, overwrite existing prompt section
//   - autoMerge: If true, skip interactive confirmation
//   - ralph: If true, use ralph mode template
//
// Returns:
//   - error: Non-nil if file operations fail
func HandlePromptMd(cmd *cobra.Command, force, autoMerge, ralph bool) error {
	var templateContent []byte
	var readErr error
	if ralph {
		templateContent, readErr = loop2.RalphTemplate(loop.PromptMd)
		if readErr != nil {
			return initErr.ReadTemplate(loop.PromptMd, readErr)
		}
	} else {
		templateContent, readErr = template.Template(loop.PromptMd)
		if readErr != nil {
			return initErr.ReadTemplate(loop.PromptMd, readErr)
		}
	}

	created, mergeErr := CreateOrMerge(cmd, MergeParams{
		Filename:        loop.PromptMd,
		MarkerStart:     marker.PromptMarkerStart,
		TemplateContent: templateContent,
		Force:           force,
		AutoMerge:       autoMerge,
		ConfirmPrompt:   desc.TextDesc(text.DescKeyInitConfirmPrompt),
		UpdateFn:        UpdatePromptSection,
	})
	if mergeErr != nil {
		return mergeErr
	}
	if created {
		mode := ""
		if ralph {
			mode = desc.TextDesc(text.DescKeyInitRalphMode)
		}
		initialize.CreatedWith(cmd, loop.PromptMd, mode)
	}
	return nil
}

// UpdatePromptSection replaces the existing prompt section between markers with
// new content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - existing: Existing file content
//   - newTemplate: New template content
//
// Returns:
//   - error: Non-nil if markers are missing or file operations fail
func UpdatePromptSection(
	cmd *cobra.Command, existing string, newTemplate []byte,
) error {
	return UpdateMarkedSection(
		cmd, loop.PromptMd, existing, newTemplate,
		marker.PromptMarkerStart, marker.PromptMarkerEnd,
		initialize.UpdatedPromptSection,
	)
}
