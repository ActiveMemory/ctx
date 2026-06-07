//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"time"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
)

// runDirName reports whether name parses as a per-run timestamp
// directory under the RunTimeLayout, so notebook artifacts like the
// state file, ledger, lock, and failmark are not mistaken for runs.
//
// Parameters:
//   - name: a directory base name under dreams/
//
// Returns:
//   - bool: true when name parses under cfgDream.RunTimeLayout
func runDirName(name string) bool {
	_, parseErr := time.Parse(cfgDream.RunTimeLayout, name)
	return parseErr == nil
}
