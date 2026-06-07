//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"io/fs"
	"path/filepath"
	"strings"

	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// ScanIdeas walks ideasDir for markdown idea files and returns a map
// of each file's path (relative to projectRoot) to its content hash.
// The dream's own dreams/ notebook and the ideas/done/ archive are
// excluded, as are non-markdown binaries. The result feeds
// DeltaSelect to drive the discipline clock.
//
// Parameters:
//   - projectRoot: absolute path to the project root (keys are
//     relative to this)
//   - ideasDir: absolute path to the ideas/ directory
//
// Returns:
//   - map[string]string: relative idea path → content hash
//   - error: non-nil on a walk or read failure (a missing ideasDir
//     yields an empty map, not an error)
func ScanIdeas(projectRoot, ideasDir string) (map[string]string, error) {
	result := make(map[string]string)
	walkErr := filepath.WalkDir(
		ideasDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if d.Name() == cfgDir.Done && path != ideasDir {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(d.Name(), cfgDream.IdeaGlob) {
				return nil
			}
			data, readErr := ctxIo.SafeReadUserFile(path)
			if readErr != nil {
				return readErr
			}
			rel, relErr := filepath.Rel(projectRoot, path)
			if relErr != nil {
				return relErr
			}
			result[rel] = Hash(data)
			return nil
		},
	)
	if walkErr != nil {
		if pathMissing(walkErr) {
			return result, nil
		}
		return nil, errDream.ScanIdeas(ideasDir, walkErr)
	}
	return result, nil
}
