//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"os"
	"path/filepath"
	"strings"

	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/heading"
)

// gistBasenames returns the set of theme-file basenames the root's gists
// link to.
//
// Parameters:
//   - root: the parsed root
//
// Returns:
//   - map[string]bool: set of basenames (e.g. "hooks.md")
func gistBasenames(root Root) map[string]bool {
	names := map[string]bool{}
	for _, t := range root.Themes {
		if t.Link != "" {
			names[filepath.Base(t.Link)] = true
		}
	}
	return names
}

// themeFiles returns the full paths of the .md files in themeDir. A
// missing directory is not an error — it means the root is not yet
// migrated, so there are no theme files.
//
// Parameters:
//   - themeDir: directory to scan
//
// Returns:
//   - []string: full paths of theme files (nil when the dir is absent)
//   - error: a read error other than "does not exist"
func themeFiles(themeDir string) ([]string, error) {
	entries, err := os.ReadDir(themeDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var paths []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), cfgFile.ExtMarkdown) {
			paths = append(paths, filepath.Join(themeDir, e.Name()))
		}
	}
	return paths, nil
}

// id is the block's identity, in the same timestamp+separator+title shape
// as entryID. A convention's empty timestamp collapses this to a
// title-only identity without needing a second id scheme.
//
// Returns:
//   - string: the block's identity
func (b stagedBlock) id() string {
	return b.Timestamp + cfgDisc.IDSeparator + b.Title
}

// conventionBlocks enumerates the curated prose sections of a convention
// staging zone: every line opening with "## " that is neither the
// structural "## Themes" nor inside an HTML comment.
//
// This is a deliberately dumb line scanner. It does not track fenced code
// blocks: fence detection (nested, tilde, indented) is a rabbit hole in a
// parser that drives destructive moves, and the safety net is elsewhere —
// byte conservation makes a mis-cut a reversible mis-grouping rather than
// data loss, and a spurious section is visible to a human at inspect time
// before any plan is approved. A stray "## " inside a fence is fixed by
// rewriting that line, never by the tool.
//
// Parameters:
//   - staging: the root's raw staging zone
//
// Returns:
//   - []stagedBlock: one per section, in file order; nil when none
func conventionBlocks(staging string) []stagedBlock {
	if staging == "" {
		return nil
	}
	spans := htmlCommentSpans(staging)
	offs := lineByteOffsets(staging)
	var blocks []stagedBlock
	for i, line := range strings.Split(staging, token.NewlineLF) {
		if !strings.HasPrefix(line, cfgDisc.SectionLinePrefix) {
			continue
		}
		if strings.TrimSpace(line) == cfgDisc.HeadingThemes {
			continue
		}
		if insideAnySpan(offs[i], spans) {
			continue
		}
		title := strings.TrimSpace(
			strings.TrimPrefix(line, cfgDisc.SectionLinePrefix),
		)
		if title == "" {
			continue
		}
		blocks = append(blocks, stagedBlock{Title: title, StartLine: i})
	}
	return blocks
}

// stagedBlocks enumerates a staging zone's addressable units for kind k:
// timestamped "## [ts] Title" entries for the entry kinds, curated "## "
// prose sections for conventions.
//
// Parameters:
//   - staging: the root's raw staging zone
//   - k: the root kind
//
// Returns:
//   - []stagedBlock: the blocks in file order; nil when none
func stagedBlocks(staging string, k Kind) []stagedBlock {
	if k == KindConvention {
		return conventionBlocks(staging)
	}
	parsed := heading.ParseEntryBlocks(staging)
	if len(parsed) == 0 {
		return nil
	}
	blocks := make([]stagedBlock, 0, len(parsed))
	for _, b := range parsed {
		blocks = append(blocks, stagedBlock{
			Timestamp: b.Entry.Timestamp,
			Title:     b.Entry.Title,
			StartLine: b.StartIndex,
		})
	}
	return blocks
}

// entryIDs returns the identity of every staged entry in content, in
// order. Identity is timestamp+title, not timestamp alone: two entries
// added in the same second share a timestamp but are distinct entries
// (observed in LEARNINGS.md). Conventions carry no timestamp, so theirs
// collapses to title alone.
//
// Parameters:
//   - content: markdown to scan for entry blocks
//   - k: the kind whose entry shape content follows
//
// Returns:
//   - []string: entry identities, nil if none
func entryIDs(content string, k Kind) []string {
	blocks := stagedBlocks(content, k)
	if len(blocks) == 0 {
		return nil
	}
	ids := make([]string, 0, len(blocks))
	for _, b := range blocks {
		ids = append(ids, b.id())
	}
	return ids
}

// entryID is a single staged entry's identity: timestamp joined to title
// by IDSeparator, matching the ids in entryIDs and the keys SplitStaging
// returns. Identity is timestamp+title, not timestamp alone: two entries
// added in the same second share a timestamp but are distinct.
//
// Parameters:
//   - e: a staged entry
//
// Returns:
//   - string: the entry's identity
func entryID(e StagedEntry) string {
	return e.Timestamp + cfgDisc.IDSeparator + e.Title
}
