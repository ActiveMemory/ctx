//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package paths

import (
	"path/filepath"

	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Resolve derives the project root from the cwd-anchored context
// directory and computes the dreams/ and ideas/ paths under it.
//
// Returns:
//   - Resolved: the project root, dreams/, and ideas/ paths
//   - error: the ContextDir resolver failure, propagated unchanged
func Resolve() (Resolved, error) {
	ctxDir, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		return Resolved{}, ctxErr
	}
	root := filepath.Dir(ctxDir)
	return Resolved{
		Root:   root,
		Dreams: filepath.Join(root, cfgDir.Dreams),
		Ideas:  filepath.Join(root, cfgDir.Ideas),
	}, nil
}
