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
	cfgObsidian "github.com/ActiveMemory/ctx/internal/config/obsidian"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errJournal "github.com/ActiveMemory/ctx/internal/err/journal"
	"github.com/ActiveMemory/ctx/internal/io"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	writeObsidian "github.com/ActiveMemory/ctx/internal/write/obsidian"
)

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
		filepath.Join(output, cfgObsidian.DirEntries),
		filepath.Join(output, cfgObsidian.DirConfig),
		filepath.Join(output, dir.JournTopics),
		filepath.Join(output, dir.JournalFiles),
		filepath.Join(output, dir.JournalTypes),
	}
	for _, d := range dirs {
		if mkErr := io.SafeMkdirAll(d, fs.PermExec); mkErr != nil {
			return errFs.Mkdir(d, mkErr)
		}
	}

	// Write .obsidian/app.json
	appConfigPath := filepath.Join(
		output, cfgObsidian.DirConfig, cfgObsidian.AppConfigFile,
	)
	if wErr := io.SafeWriteFile(
		appConfigPath, []byte(cfgObsidian.AppConfig), fs.PermFile,
	); wErr != nil {
		return errFs.FileWrite(appConfigPath, wErr)
	}

	// Write README
	readmePath := filepath.Join(output, file.Readme)
	if wErr := io.SafeWriteFile(
		readmePath,
		[]byte(fmt.Sprintf(tpl.ObsidianReadme, journalDir)),
		fs.PermFile,
	); wErr != nil {
		return errFs.FileWrite(readmePath, wErr)
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
		dst := filepath.Join(output, cfgObsidian.DirEntries, entry.Filename)

		content, readErr := io.SafeReadUserFile(filepath.Clean(src))
		if readErr != nil {
			writeErr.WarnFile(cmd, entry.Filename, readErr)
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
			entry, topicIndex, cfgObsidian.MaxRelated,
		)

		if wErr := io.SafeWriteFile(
			dst, []byte(transformed), fs.PermFile,
		); wErr != nil {
			writeErr.WarnFile(cmd, entry.Filename, wErr)
			continue
		}
	}

	// Write topic MOC and pages
	if len(topics) > 0 {
		topicsDir := filepath.Join(output, dir.JournTopics)
		mocPath := filepath.Join(output, cfgObsidian.MOCTopics)
		if wErr := io.SafeWriteFile(
			mocPath, []byte(moc.ObsidianTopics(topics)),
			fs.PermFile,
		); wErr != nil {
			return errFs.FileWrite(mocPath, wErr)
		}

		for _, t := range topics {
			if !t.Popular {
				continue
			}
			pagePath := filepath.Join(topicsDir, t.Name+file.ExtMarkdown)
			if wErr := io.SafeWriteFile(
				pagePath, []byte(moc.GenerateObsidianTopicPage(t)),
				fs.PermFile,
			); wErr != nil {
				writeErr.WarnFile(cmd, pagePath, wErr)
			}
		}
	}

	// Write key files MOC and pages
	if len(keyFiles) > 0 {
		filesDir := filepath.Join(output, dir.JournalFiles)
		mocPath := filepath.Join(output, cfgObsidian.MOCFiles)
		if wErr := io.SafeWriteFile(
			mocPath, []byte(moc.ObsidianFiles(keyFiles)),
			fs.PermFile,
		); wErr != nil {
			return errFs.FileWrite(mocPath, wErr)
		}

		for _, kf := range keyFiles {
			if !kf.Popular {
				continue
			}
			slug := format.KeyFileSlug(kf.Path)
			pagePath := filepath.Join(filesDir, slug+file.ExtMarkdown)
			if wErr := io.SafeWriteFile(
				pagePath, []byte(moc.GenerateObsidianFilePage(kf)),
				fs.PermFile,
			); wErr != nil {
				writeErr.WarnFile(cmd, pagePath, wErr)
			}
		}
	}

	// Write types MOC and pages
	if len(sessionTypes) > 0 {
		typesDir := filepath.Join(output, dir.JournalTypes)
		mocPath := filepath.Join(output, cfgObsidian.MOCTypes)
		if wErr := io.SafeWriteFile(
			mocPath, []byte(moc.ObsidianTypes(sessionTypes)),
			fs.PermFile,
		); wErr != nil {
			return errFs.FileWrite(mocPath, wErr)
		}

		for _, st := range sessionTypes {
			pagePath := filepath.Join(typesDir, st.Name+file.ExtMarkdown)
			if wErr := io.SafeWriteFile(
				pagePath,
				[]byte(moc.GenerateObsidianTypePage(st)), fs.PermFile,
			); wErr != nil {
				writeErr.WarnFile(cmd, pagePath, wErr)
			}
		}
	}

	// Write Home.md
	homePath := filepath.Join(output, cfgObsidian.MOCHome)
	if wErr := io.SafeWriteFile(
		homePath,
		[]byte(moc.Home(
			regularEntries,
			len(topics) > 0, len(keyFiles) > 0, len(sessionTypes) > 0,
		)),
		fs.PermFile,
	); wErr != nil {
		return errFs.FileWrite(homePath, wErr)
	}

	writeObsidian.InfoGenerated(cmd, len(entries), output)

	return nil
}
