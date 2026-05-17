//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package path

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// KBDir returns the .context/kb/ directory.
//
// Returns:
//   - string: full path to .context/kb/
//   - error: non-nil when the context directory is not declared
func KBDir() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, cfgKB.KBSubdir), nil
}

// KBIndexFile returns the kb landing-page path.
//
// Returns:
//   - string: full path to .context/kb/index.md
//   - error: non-nil when the context directory is not declared
func KBIndexFile() (string, error) {
	kb, err := KBDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(kb, cfgKB.KBIndex), nil
}

// KBTopicsDir returns the .context/kb/topics/ directory.
//
// Returns:
//   - string: full path to .context/kb/topics/
//   - error: non-nil when the context directory is not declared
func KBTopicsDir() (string, error) {
	kb, err := KBDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(kb, cfgKB.TopicsSubdir), nil
}

// KBTopicDir returns the .context/kb/topics/<slug>/ directory
// for a given topic slug.
//
// Parameters:
//   - slug: kebab-case topic slug (e.g. "cursor-hooks").
//
// Returns:
//   - string: full path to .context/kb/topics/<slug>/
//   - error: non-nil when the context directory is not declared
func KBTopicDir(slug string) (string, error) {
	topics, err := KBTopicsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(topics, slug), nil
}

// KBTopicIndexFile returns the index.md path inside a topic
// folder.
//
// Parameters:
//   - slug: kebab-case topic slug.
//
// Returns:
//   - string: full path to .context/kb/topics/<slug>/index.md
//   - error: non-nil when the context directory is not declared
func KBTopicIndexFile(slug string) (string, error) {
	topic, err := KBTopicDir(slug)
	if err != nil {
		return "", err
	}
	return filepath.Join(topic, cfgKB.TopicIndex), nil
}

// IngestDir returns the .context/ingest/ directory.
//
// Returns:
//   - string: full path to .context/ingest/
//   - error: non-nil when the context directory is not declared
func IngestDir() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, cfgKB.IngestSubdir), nil
}

// IngestArtifactFile returns a path under .context/ingest/ for
// a named artifact (e.g. Rules, ModeIngest, SessionLog).
//
// Parameters:
//   - name: filename constant from internal/config/kb.
//
// Returns:
//   - string: full path to .context/ingest/<name>
//   - error: non-nil when the context directory is not declared
func IngestArtifactFile(name string) (string, error) {
	ingest, err := IngestDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ingest, name), nil
}

// CloseoutsDir returns the .context/ingest/closeouts/ directory.
//
// Returns:
//   - string: full path to .context/ingest/closeouts/
//   - error: non-nil when the context directory is not declared
func CloseoutsDir() (string, error) {
	ingest, err := IngestDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ingest, cfgKB.CloseoutsSubdir), nil
}

// ArchiveCloseoutsDir returns the .context/archive/closeouts/
// directory where folded closeouts are moved by the handover
// fold mechanism.
//
// Returns:
//   - string: full path to .context/archive/closeouts/
//   - error: non-nil when the context directory is not declared
func ArchiveCloseoutsDir() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, dir.Archive, cfgKB.ArchiveCloseoutsSubdir), nil
}

// SiteDir returns the .context/site/ root (gitignored).
//
// Returns:
//   - string: full path to .context/site/
//   - error: non-nil when the context directory is not declared
func SiteDir() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, cfgKB.SiteSubdir), nil
}

// SiteKBDir returns the .context/site/kb/ render-output
// directory.
//
// Returns:
//   - string: full path to .context/site/kb/
//   - error: non-nil when the context directory is not declared
func SiteKBDir() (string, error) {
	site, err := SiteDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(site, cfgKB.SiteKBSubdir), nil
}
