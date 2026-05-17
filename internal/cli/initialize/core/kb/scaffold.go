//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package kb

import (
	"path/filepath"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHandover "github.com/ActiveMemory/ctx/internal/config/handover"
	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	errInitKB "github.com/ActiveMemory/ctx/internal/err/initialize/kb"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Scaffold extracts the embedded KB / ingest / handover
// templates into contextDir. Existing files are skipped
// (preserves curated content).
//
// Parameters:
//   - contextDir: absolute path to .context/.
//
// Returns:
//   - error: wrapped errors for directory create / file copy
//     failures. Per-file existence is non-fatal (skipped with
//     a noop).
func Scaffold(contextDir string) error {
	dirs := []string{
		filepath.Join(contextDir, cfgKB.KBSubdir),
		filepath.Join(contextDir, cfgKB.KBSubdir, cfgKB.TopicsSubdir),
		filepath.Join(contextDir, cfgKB.IngestSubdir),
		filepath.Join(
			contextDir, cfgKB.IngestSubdir, cfgKB.CloseoutsSubdir,
		),
		filepath.Join(
			contextDir, cfgKB.IngestSubdir, cfgKB.SchemasSubdir,
		),
		filepath.Join(contextDir, cfgHandover.Subdir),
	}
	for _, d := range dirs {
		if mkErr := io.SafeMkdirAll(d, cfgFs.PermExec); mkErr != nil {
			return errInitKB.Mkdir(d, mkErr)
		}
	}

	// Place .gitkeep stubs so the empty dirs round-trip
	// through git. Skipped silently when any file in the dir
	// already exists.
	gitkeepDirs := []string{
		filepath.Join(contextDir, cfgKB.KBSubdir, cfgKB.TopicsSubdir),
		filepath.Join(
			contextDir, cfgKB.IngestSubdir, cfgKB.CloseoutsSubdir,
		),
		filepath.Join(contextDir, cfgHandover.Subdir),
	}
	for _, d := range gitkeepDirs {
		if gkErr := PlaceGitkeep(d); gkErr != nil {
			return gkErr
		}
	}

	// Copy embedded ingest-side templates.
	ingestErr := CopyEmbedTree(
		cfgKB.AssetTemplatesIngest,
		filepath.Join(contextDir, cfgKB.IngestSubdir),
	)
	if ingestErr != nil {
		return errInitKB.CopyIngestTemplates(ingestErr)
	}

	// Copy embedded schemas.
	schemasErr := CopyEmbedTree(
		cfgKB.AssetTemplatesSchemas,
		filepath.Join(
			contextDir, cfgKB.IngestSubdir, cfgKB.SchemasSubdir,
		),
	)
	if schemasErr != nil {
		return errInitKB.CopySchemas(schemasErr)
	}

	// Copy embedded kb landing.
	landingErr := CopyEmbedFile(
		cfgKB.AssetTemplateKBIndex,
		filepath.Join(
			contextDir, cfgKB.KBSubdir, cfgKB.KBIndex,
		),
	)
	if landingErr != nil {
		return errInitKB.CopyLanding(landingErr)
	}

	return nil
}
