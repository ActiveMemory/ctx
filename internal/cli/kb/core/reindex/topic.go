//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

import (
	"errors"
	"os"
	"sort"

	errKbCli "github.com/ActiveMemory/ctx/internal/err/kb/cli"
	"github.com/ActiveMemory/ctx/internal/io"
)

// ListTopics returns every topic slug under topicsDir: the
// slash-separated path, relative to the topics root, of each
// directory that holds a topic index.md. The scan is recursive, so
// a grouped layout (topics/<group>/<slug>/index.md) enumerates as
// "<group>/<slug>". A directory whose index.md sits above nested
// topics is a group-landing (orientation) page, not a topic, and is
// excluded. Slugs are returned sorted.
//
// Parameters:
//   - topicsDir: absolute path to .context/kb/topics/.
//
// Returns:
//   - []string: sorted topic slugs (slashes preserved for grouped /
//     vendor-namespaced topology).
//   - error: wrapped enumeration failure.
func ListTopics(topicsDir string) ([]string, error) {
	if _, statErr := io.SafeStat(topicsDir); statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errKbCli.ReadTopicsDir(statErr)
	}

	indexed := make(map[string]bool)
	if walkErr := collectTopicDirs(topicsDir, "", indexed); walkErr != nil {
		return nil, errKbCli.ReadTopicsDir(walkErr)
	}

	slugs := topicLeaves(indexed)
	sort.Strings(slugs)
	return slugs, nil
}
