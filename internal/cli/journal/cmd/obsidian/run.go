//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package obsidian

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/tpl"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/consolidate"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/format"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/frontmatter"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/moc"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/parse"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/reduce"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/section"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/turn"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/wikilink"
	"github.com/ActiveMemory/ctx/internal/cli/journal/core/wrap"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/obsidian"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errJournal "github.com/ActiveMemory/ctx/internal/err/journal"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/err"
	writeObsidian "github.com/ActiveMemory/ctx/internal/write/obsidian"
)

// Run generates an Obsidian vault from journal entries.
//
// Pipeline:
//  1. Scan entries (reuse core.ScanJournalEntries)
//  2. Create output dirs (entries/, topics/, files/, types/, .obsidian/)
//  3. Write .obsidian/app.json
//  4. Transform and write entries (normalize, convert links, transform
//     frontmatter, add related footer)
//  5. Build indices (reuse core.BuildTopicIndex etc.)
//  6. Generate and write MOC pages
//  7. Generate and write Home.md
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - output: Output directory for the vault
//
// Returns:
//   - error: Non-nil if generation fails
func Run(cmd *cobra.Command, output string) error {
	return BuildVault(
		cmd, filepath.Join(rc.ContextDir(), dir.Journal), output,
	)
}

// BuildVault generates an Obsidian vault from journal entries in
// journalDir and writes the output to the output directory.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - journalDir: Path to the source journal directory
//   - output: Output directory for the vault
//
// Returns:
//   - error: Non-nil if generation fails
func BuildVault(cmd *cobra.Command, journalDir, output string) error {
	if _, statErr := os.Stat(journalDir); os.IsNotExist(statErr) {
		return errJournal.NoDir(journalDir)
	}

	entries, scanErr := parse.ScanJournalEntries(journalDir)
	if scanErr != nil {
		return errJournal.Scan(scanErr)
	}

	if len(entries) == 0 {
		return errJournal.NoEntries(journalDir)
	}

	// Create output directory structure
	dirs := []string{
		output,
		filepath.Join(output, obsidian.DirEntries),
		filepath.Join(output, obsidian.DirConfig),
		filepath.Join(output, dir.JournTopics),
		filepath.Join(output, dir.JournalFiles),
		filepath.Join(output, dir.JournalTypes),
	}
	for _, d := range dirs {
		if mkErr := os.MkdirAll(d, fs.PermExec); mkErr != nil {
			return errFs.Mkdir(d, mkErr)
		}
	}

	// Write .obsidian/app.json
	appConfigPath := filepath.Join(
		output, obsidian.DirConfig, obsidian.AppConfigFile,
	)
	if writeErr := os.WriteFile(
		appConfigPath, []byte(obsidian.AppConfig), fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(appConfigPath, writeErr)
	}

	// Write README
	readmePath := filepath.Join(output, file.Readme)
	if writeErr := os.WriteFile(
		readmePath,
		[]byte(fmt.Sprintf(tpl.ObsidianReadme, journalDir)),
		fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(readmePath, writeErr)
	}

	// Build indices for MOC pages and related footer
	regularEntries := moc.FilterRegularEntries(entries)

	topicEntries := moc.FilterEntriesWithTopics(entries)
	topics := section.BuildTopicIndex(topicEntries)

	keyFileEntries := moc.FilterEntriesWithKeyFiles(entries)
	keyFiles := section.BuildKeyFileIndex(keyFileEntries)

	typeEntries := moc.FilterEntriesWithType(entries)
	sessionTypes := section.BuildTypeIndex(typeEntries)

	// Build topic lookup for related footer
	topicIndex := moc.BuildTopicLookup(topicEntries)

	// Transform and write entries
	for _, entry := range entries {
		src := entry.Path
		dst := filepath.Join(output, obsidian.DirEntries, entry.Filename)

		content, readErr := os.ReadFile(filepath.Clean(src))
		if readErr != nil {
			err.WarnFile(cmd, entry.Filename, readErr)
			continue
		}

		// Normalize content (read-only - do NOT write back to source)
		normalized := wrap.Content(
			turn.MergeConsecutive(
				consolidate.ToolRuns(
					reduce.CleanToolOutputJSON(
						reduce.StripSystemReminders(string(content)),
					),
				),
			),
		)

		// Transform for Obsidian
		sourcePath := filepath.Join(
			dir.Context, dir.Journal, entry.Filename,
		)
		transformed := frontmatter.Transform(normalized, sourcePath)
		transformed = wikilink.ConvertMarkdownLinks(transformed)
		transformed += moc.GenerateRelatedFooter(
			entry, topicIndex, obsidian.MaxRelated,
		)

		if writeErr := os.WriteFile( //nolint:gosec // dst built from output dir + filename
			dst, []byte(transformed), fs.PermFile,
		); writeErr != nil {
			err.WarnFile(cmd, entry.Filename, writeErr)
			continue
		}
	}

	// Write topic MOC and pages
	if len(topics) > 0 {
		topicsDir := filepath.Join(output, dir.JournTopics)
		mocPath := filepath.Join(output, obsidian.MOCTopics)
		if writeErr := os.WriteFile(
			mocPath, []byte(moc.ObsidianTopics(topics)),
			fs.PermFile,
		); writeErr != nil {
			return errFs.FileWrite(mocPath, writeErr)
		}

		for _, t := range topics {
			if !t.Popular {
				continue
			}
			pagePath := filepath.Join(topicsDir, t.Name+file.ExtMarkdown)
			if writeErr := os.WriteFile(
				pagePath, []byte(moc.GenerateObsidianTopicPage(t)),
				fs.PermFile,
			); writeErr != nil {
				err.WarnFile(cmd, pagePath, writeErr)
			}
		}
	}

	// Write key files MOC and pages
	if len(keyFiles) > 0 {
		filesDir := filepath.Join(output, dir.JournalFiles)
		mocPath := filepath.Join(output, obsidian.MOCFiles)
		if writeErr := os.WriteFile(
			mocPath, []byte(moc.ObsidianFiles(keyFiles)),
			fs.PermFile,
		); writeErr != nil {
			return errFs.FileWrite(mocPath, writeErr)
		}

		for _, kf := range keyFiles {
			if !kf.Popular {
				continue
			}
			slug := format.KeyFileSlug(kf.Path)
			pagePath := filepath.Join(filesDir, slug+file.ExtMarkdown)
			if writeErr := os.WriteFile(
				pagePath, []byte(moc.GenerateObsidianFilePage(kf)),
				fs.PermFile,
			); writeErr != nil {
				err.WarnFile(cmd, pagePath, writeErr)
			}
		}
	}

	// Write types MOC and pages
	if len(sessionTypes) > 0 {
		typesDir := filepath.Join(output, dir.JournalTypes)
		mocPath := filepath.Join(output, obsidian.MOCTypes)
		if writeErr := os.WriteFile(
			mocPath, []byte(moc.ObsidianTypes(sessionTypes)),
			fs.PermFile,
		); writeErr != nil {
			return errFs.FileWrite(mocPath, writeErr)
		}

		for _, st := range sessionTypes {
			pagePath := filepath.Join(typesDir, st.Name+file.ExtMarkdown)
			if writeErr := os.WriteFile(
				pagePath,
				[]byte(moc.GenerateObsidianTypePage(st)), fs.PermFile,
			); writeErr != nil {
				err.WarnFile(cmd, pagePath, writeErr)
			}
		}
	}

	// Write Home.md
	homePath := filepath.Join(output, obsidian.MOCHome)
	if writeErr := os.WriteFile(
		homePath,
		[]byte(moc.Home(
			regularEntries,
			len(topics) > 0, len(keyFiles) > 0, len(sessionTypes) > 0,
		)),
		fs.PermFile,
	); writeErr != nil {
		return errFs.FileWrite(homePath, writeErr)
	}

	writeObsidian.InfoGenerated(cmd, len(entries), output)

	return nil
}
