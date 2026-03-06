//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

// InfoPathConversionExists reports that a path conversion target already
// exists at the destination. Used during init to show which template files
// were skipped.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - rootDir: project root directory for path resolution.
//   - oldPath: original template-relative path.
//   - newPath: destination-relative path joined with rootDir.
func InfoPathConversionExists(
	cmd *cobra.Command, rootDir, oldPath, newPath string,
) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplPathExists, oldPath, filepath.Join(rootDir, newPath))
}

// InfoExistsWritingAsAlternative reports that a file already exists and the
// content is being written to an alternative filename instead.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - path: the original target path that already exists.
//   - alternative: the fallback path where content was written.
func InfoExistsWritingAsAlternative(
	cmd *cobra.Command, path, alternative string,
) {
	if cmd == nil {
		return
	}
	sprintf(cmd, tplExistsWritingAsAlternative, path, alternative)
}
