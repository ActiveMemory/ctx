//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	errKbCli "github.com/ActiveMemory/ctx/internal/err/kb/cli"
	"github.com/ActiveMemory/ctx/internal/io"
)

// ListTopics returns every subdirectory of topicsDir that
// contains an index.md file. Slugs are returned sorted.
//
// Parameters:
//   - topicsDir: absolute path to .context/kb/topics/.
//
// Returns:
//   - []string: sorted topic slugs (slashes preserved for
//     vendor-namespaced topology).
//   - error: wrapped enumeration failure.
func ListTopics(topicsDir string) ([]string, error) {
	entries, readErr := os.ReadDir(topicsDir)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errKbCli.ReadTopicsDir(readErr)
	}
	var slugs []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		// Test if <topicsDir>/<name>/index.md exists.
		idx := filepath.Join(topicsDir, e.Name(), cfgKB.TopicIndex)
		if _, statErr := io.SafeStat(idx); statErr == nil {
			slugs = append(slugs, e.Name())
		}
	}
	sort.Strings(slugs)
	return slugs, nil
}
