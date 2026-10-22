//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agents

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/merge"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/entity"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// HandleMd creates or merges AGENTS.md with ctx content.
//
// AGENTS.md is the tool-neutral instruction file read by Claude Code,
// Codex, Copilot, Cursor, OpenCode, Pi and others, so it carries the
// full ctx instructions; CLAUDE.md only imports it. The ctx block sits
// between marker.AgentsStart and marker.AgentsEnd, the markers
// ctx setup agents|codex|opencode|pi key on, so a later setup run
// leaves the file alone (specs/init-agents-md.md).
//
// Parameters:
//   - cmd: Cobra command for output
//   - force: If true, overwrite the existing ctx section
//   - autoMerge: If true, skip interactive confirmation
//
// Returns:
//   - error: Non-nil if file operations fail
func HandleMd(cmd *cobra.Command, force, autoMerge bool) error {
	templateContent, err := agent.AgentsMd()
	if err != nil {
		return errInit.ReadTemplate(cfgHook.FileAgentsMd, err)
	}

	created, mergeErr := merge.OrCreate(cmd, entity.MergeParams{
		Filename:        cfgHook.FileAgentsMd,
		MarkerStart:     marker.AgentsStart,
		MarkerEnd:       marker.AgentsEnd,
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
		initialize.Created(cmd, cfgHook.FileAgentsMd)
	}

	return nil
}
