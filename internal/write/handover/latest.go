//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handover

import (
	"time"
)

// Latest returns the GeneratedAt of the most recent
// handover under handoversDir, or the zero value when no
// handover exists yet. The third return is the path of that
// latest handover (empty when no handover exists).
//
// Parameters:
//   - handoversDir: absolute path to .context/handovers/.
//
// Returns:
//   - time.Time: GeneratedAt of the latest handover, or zero.
//   - string: path of the latest handover, or empty.
//   - error: non-nil on directory-walk failures (not on
//     individual-file parse failures).
func Latest(handoversDir string) (time.Time, string, error) {
	files, err := listHandovers(handoversDir)
	if err != nil {
		return time.Time{}, "", err
	}
	if len(files) == 0 {
		return time.Time{}, "", nil
	}
	latest := files[len(files)-1]
	return latest.Frontmatter.GeneratedAt, latest.Path, nil
}
