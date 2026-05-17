//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package path

import (
	"path/filepath"

	cfgHandover "github.com/ActiveMemory/ctx/internal/config/handover"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Dir returns the `.context/handovers/` directory.
//
// Returns:
//   - string: full path to .context/handovers/.
//   - error: non-nil when the context directory is not declared.
func Dir() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, cfgHandover.Subdir), nil
}
