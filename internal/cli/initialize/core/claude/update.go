//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"github.com/ActiveMemory/ctx/internal/cli/initialize/core"
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/marker"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
	"github.com/spf13/cobra"
)

// updateCtxSection replaces the existing ctx section between markers with
// new content.
//
// Parameters:
//   - cmd: Cobra command for output
//   - existing: Existing file content
//   - newTemplate: New template content
//
// Returns:
//   - error: Non-nil if markers are missing or file operations fail
func updateCtxSection(cmd *cobra.Command, existing string, newTemplate []byte) error {
	return core.UpdateMarkedSection(
		cmd, claude.Md, existing, newTemplate,
		marker.CtxMarkerStart, marker.CtxMarkerEnd,
		initialize.UpdatedCtxSection,
	)
}
