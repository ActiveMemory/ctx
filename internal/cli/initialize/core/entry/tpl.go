//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/entry"
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core/tpl"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// CreateTemplates creates entry template files in .context/templates/.
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: The .context/ directory path
//   - force: If true, overwrite existing files
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func CreateTemplates(
	cmd *cobra.Command, contextDir string, force bool,
) error {
	return tpl.DeployTemplates(cmd, contextDir, force,
		entity.DeployParams{
			SubDir:     dir.Templates,
			ListErrKey: text.DescKeyErrPromptListEntryTemplates,
			ReadErrKey: text.DescKeyErrPromptReadEntryTemplate,
		},
		entry.List, entry.ForName,
	)
}
