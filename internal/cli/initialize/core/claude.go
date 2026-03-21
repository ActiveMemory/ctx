//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/spf13/cobra"

	readClaude "github.com/ActiveMemory/ctx/internal/assets/read/claude"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// HandleClaudeMd creates or merges CLAUDE.md with ctx content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - force: If true, overwrite the existing ctx section
//   - autoMerge: If true, skip interactive confirmation
//
// Returns:
//   - error: Non-nil if file operations fail
func HandleClaudeMd(cmd *cobra.Command, force, autoMerge bool) error {
	templateContent, err := readClaude.Md()
	if err != nil {
		return ctxErr.ReadTemplate(claude.Md, err)
	}

	created, mergeErr := CreateOrMerge(cmd, MergeParams{
		Filename:        claude.Md,
		MarkerStart:     marker.CtxMarkerStart,
		TemplateContent: templateContent,
		Force:           force,
		AutoMerge:       autoMerge,
		ConfirmPrompt:   desc.TextDesc(text.DescKeyInitConfirmClaude),
		UpdateFn:        UpdateCtxSection,
	})

	if mergeErr != nil {
		return mergeErr
	}

	if created {
		initialize.Created(cmd, claude.Md)
	}

	return nil
}

// UpdateCtxSection replaces the existing ctx section between markers with
// new content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - existing: Existing file content
//   - newTemplate: New template content
//
// Returns:
//   - error: Non-nil if markers are missing or file operations fail
func UpdateCtxSection(cmd *cobra.Command, existing string, newTemplate []byte) error {
	return UpdateMarkedSection(
		cmd, claude.Md, existing, newTemplate,
		marker.CtxMarkerStart, marker.CtxMarkerEnd,
		initialize.UpdatedCtxSection,
	)
}
