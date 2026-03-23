//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"github.com/spf13/cobra"

	readPrompt "github.com/ActiveMemory/ctx/internal/assets/read/prompt"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/tpl"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// CreatePromptTemplates creates prompt template files in .context/prompts/.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: The .context/ directory path
//   - force: If true, overwrite existing files
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func CreatePromptTemplates(
	cmd *cobra.Command, contextDir string, force bool,
) error {
	return tpl.DeployTemplates(cmd, contextDir, force,
		entity.DeployParams{
			SubDir:     dir.Prompts,
			ListErrKey: text.DescKeyErrPromptListPromptTemplates,
			ReadErrKey: text.DescKeyErrPromptReadPromptTemplate,
		},
		readPrompt.TemplateList, readPrompt.Template,
	)
}
