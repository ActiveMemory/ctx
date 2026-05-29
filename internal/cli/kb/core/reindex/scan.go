//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

import (
	"os"
	"path/filepath"
	"strings"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	cfgToken "github.com/ActiveMemory/ctx/internal/config/token"
)

// collectTopicDirs records, in indexed (keyed by slash-separated
// path relative to root), every directory at or below root/rel that
// directly contains a topic index.md. The root itself (rel == "")
// is never recorded: a topics/index.md is not a topic page.
//
// Parameters:
//   - root: absolute topics-root path.
//   - rel: slash-separated path under root currently being scanned
//     ("" for the root itself).
//   - indexed: accumulator of topic-index-bearing relative paths.
//
// Returns:
//   - error: a directory-read failure.
func collectTopicDirs(root, rel string, indexed map[string]bool) error {
	entries, readErr := os.ReadDir(
		filepath.Join(root, filepath.FromSlash(rel)),
	)
	if readErr != nil {
		return readErr
	}
	for _, e := range entries {
		if e.IsDir() {
			child := e.Name()
			if rel != "" {
				child = rel + cfgToken.Slash + e.Name()
			}
			if walkErr := collectTopicDirs(
				root, child, indexed,
			); walkErr != nil {
				return walkErr
			}
			continue
		}
		if e.Name() == cfgKB.TopicIndex && rel != "" {
			indexed[rel] = true
		}
	}
	return nil
}

// topicLeaves returns the indexed directories that are topic pages
// rather than group-landings: those with no other indexed directory
// nested beneath them.
//
// Parameters:
//   - indexed: all topic-index-bearing relative paths.
//
// Returns:
//   - []string: leaf topic slugs (unsorted).
func topicLeaves(indexed map[string]bool) []string {
	var leaves []string
	for rel := range indexed {
		if !hasNestedTopic(rel, indexed) {
			leaves = append(leaves, rel)
		}
	}
	return leaves
}

// hasNestedTopic reports whether any indexed directory is a strict
// descendant of rel (making rel a group-landing, not a topic).
//
// Parameters:
//   - rel: candidate slash-separated topic path.
//   - indexed: all topic-index-bearing relative paths.
//
// Returns:
//   - bool: true when rel has a nested topic beneath it.
func hasNestedTopic(rel string, indexed map[string]bool) bool {
	prefix := rel + cfgToken.Slash
	for other := range indexed {
		if other != rel && strings.HasPrefix(other, prefix) {
			return true
		}
	}
	return false
}
