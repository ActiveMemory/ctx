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

// LatestFile returns the parsed File of the most recent handover
// under handoversDir, or (File{}, false, nil) when no handover
// exists yet.
//
// Parameters:
//   - handoversDir: absolute path to .context/handovers/.
//
// Returns:
//   - File: parsed latest handover file.
//   - bool: true when a handover was found.
//   - error: non-nil on directory-walk failures.
func LatestFile(handoversDir string) (File, bool, error) {
	files, err := listHandovers(handoversDir)
	if err != nil {
		return File{}, false, err
	}
	if len(files) == 0 {
		return File{}, false, nil
	}
	return files[len(files)-1], true, nil
}
